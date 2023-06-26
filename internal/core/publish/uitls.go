/*
Copyright 2021 The AtomCI Group Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package publish

import (
	"fmt"
	"strings"

	"github.com/go-atomci/atomci/internal/core/pipelinemgr"
	"github.com/go-atomci/atomci/internal/middleware/log"
	"github.com/go-atomci/atomci/internal/models"

	"github.com/astaxie/beego/logs"
)

func getOperationCanEnableByStepStatus(stepType string, status int64, operations *models.PublishOperation) *models.PublishOperation {
	switch stepType {
	case models.StepManual:
		operations.Manual = true
	case models.StepBuild:
		operations.Build = true
		if status == models.Running {
			operations.BackTo = false
			operations.Build = false
			operations.Terminate = true
		}
	case models.StepDeploy:
		operations.Deploy = true
		if status == models.Running {
			operations.BackTo = false
			operations.Deploy = false
			operations.Terminate = true
		}
	case "None":
		operations.NextStage = true
	default:
		logs.Info("step unmatched, show default enabled operations")
	}
	return operations
}

func (pm *PublishManager) getPublishItemCanEnableOperations(item *models.Publish) *models.PublishOperation {
	operations := &models.PublishOperation{
		BackTo: true,
	}

	switch item.Status {
	case models.END, models.Closed:
		operations.BackTo = false
	case models.Running, models.Pending:
		operations = getOperationCanEnableByStepStatus(item.StepType, item.Status, operations)
	case models.Failed, models.UnKnown:
		operations = getOperationCanEnableByStepStatus(item.StepType, item.Status, operations)
	case models.Success:
		nextStep, err := pm.pipelineHandler.GetNextStepTypeByStageID(item.LastPipelineInstanceID, item.StageID, item.StepIndex)
		if err != nil {
			log.Log.Error("when get publish list, get next step occur error: %s, use default enable publish operations", err)
			return nil
		}
		operations = getOperationCanEnableByStepStatus(nextStep, item.Status, operations)
	case models.TerminateSuccess, models.TerminateFailed:
		operations = getOperationCanEnableByStepStatus(item.StepType, item.Status, operations)
	}
	return operations
}

func (pm *PublishManager) getPublishStepsByPipelineInstanceID(publish *models.Publish) ([]*PublishStep, error) {
	currentEnvStage, err := pm.pipelineHandler.GetPipelineInstanceEnvStageByID(publish.LastPipelineInstanceID, publish.StageID)
	if err != nil {
		log.Log.Error("when get publish steps, get PipelineInstanceEnvStage by pipline instance id: %v occur error: %s", publish.LastPipelineInstanceID, err.Error())
		return nil, err
	}
	currentIndex := publish.StepIndex
	steps := []*PublishStep{}
	for _, step := range currentEnvStage.Steps {
		item := &PublishStep{
			Index: step.Index,
			Type:  step.Type,
			Name:  step.Name,
		}
		steps = append(steps, item)
	}

	for _, s := range steps {
		if s.Index < currentIndex {
			s.Status = models.Success
		} else {
			if s.Index > currentIndex {
				s.Status = models.Pending
			} else {
				s.Status = publish.Status
			}
		}
	}
	return steps, nil
}

func (pm *PublishManager) getCurrentEnvStageJSON(pipelineConfig pipelinemgr.PipelineConfig, envStageID int64, currentEnvStage *pipelinemgr.PipelineStageStruct) *pipelinemgr.PipelineStageStruct {
	for _, envStage := range pipelineConfig {
		if envStage.StageID == envStageID {
			currentEnvStage = envStage
			break
		}
	}
	return currentEnvStage
}

// Base on publish model, get previousStep and nextStep, return step name
func (pm *PublishManager) getPublishPreviousAndNextStepbyPublishModel(publish *models.Publish) (string, string) {
	envStage, err := pm.pipelineHandler.GetPipelineInstanceEnvStageByID(publish.LastPipelineInstanceID, publish.StageID)
	if err != nil {
		log.Log.Error("when get publish steps, get PipelineInstanceEnvStage occur error: %s, reset revious/next step to empty", err.Error())
		return "", ""
	}
	var previousStep, nextStep string
	currentIndex := publish.StepIndex
	for _, index := range envStage.Steps {
		if index.Index == currentIndex-1 {
			previousStep = index.Name
			continue
		}
		if index.Index == currentIndex+1 {
			nextStep = index.Name
		}
	}
	return previousStep, nextStep
}

func (pm *PublishManager) createPublishApps(publishApps []*PubllishReqApp, publishID int64) error {
	for _, app := range publishApps {
		projectApp, _appErr := pm.projectModel.GetProjectApp(app.AppID)
		if _appErr != nil {
			log.Log.Error("get project app failed, msg: %s", _appErr)
			return _appErr
		}

		publishAppModel := models.PublishApp{
			Addons:         models.NewAddons(),
			BranchName:     app.BranchName,
			CompileCommand: app.CompileCommand,
			ProjectAppID:   projectApp.ID,
			PublishID:      publishID,
		}
		_, err := pm.model.CreatePublishAppIfNotExist(&publishAppModel)
		if err != nil {
			log.Log.Error("publish id: %v create publish app failed, msg: %s", publishID, err)
			return err
		}
	}
	return nil
}

func (pm *PublishManager) getPublishInfoApps(publishID int64) ([]*PublishInfoApp, error) {
	modelApps, err := pm.model.GetPublishAppsByID(publishID)
	if err != nil {
		return nil, err
	}
	infoApps := []*PublishInfoApp{}
	for _, app := range modelApps {
		projectApp, err := pm.projectModel.GetProjectApp(app.ProjectAppID)
		if err != nil {
			logs.Warn("publish app is not found, by project app id: %v, error: %s", app.ProjectAppID, err.Error())
			continue
		}
		scpApp, err := pm.gitAppModel.GetScmAppByID(projectApp.ScmID)
		if err != nil {
			logs.Warn("get scm by id %v error: %s", projectApp.ScmID, err.Error())
			continue
		}
		infoApp := &PublishInfoApp{
			PublishApp: app,
			Language:   scpApp.Language,
			Name:       scpApp.Name,
			// use hard code defined app type
			Type: "app",
		}
		infoApps = append(infoApps, infoApp)
	}
	return infoApps, nil
}

func (pm *PublishManager) createPipelineInstance(pipelineID, publishID int64, user string) (int64, error) {

	projectPipeline, err := pm.projectModel.GetProjectPipelineByID(pipelineID)
	if err != nil {
		log.Log.Error("when get project pipeline by id occur error: %s", err.Error())
		return 0, err
	}

	pipelineInstance := &models.PipelineInstance{
		PipelineID:  pipelineID,
		PublishID:   publishID,
		Creator:     user,
		Name:        projectPipeline.Name,
		Description: projectPipeline.Description,
		Config:      projectPipeline.Config,
		ProjectID:   projectPipeline.ProjectID,
		IsDefault:   projectPipeline.IsDefault,
	}
	instanceID, err := pm.pipelineModel.CreatePipelineInstance(pipelineInstance)
	if err != nil {
		log.Log.Error("when create publish, create pipeline instance occur errror: %s", err.Error())
		return 0, fmt.Errorf("网络错误，请重试")
	}
	return instanceID, nil
}

func (pm *PublishManager) publishOrderBaseVerify(publishID, stageID int64) (*models.Publish, error) {
	modelPublish, err := pm.model.GetPublishByID(publishID)
	if err != nil {
		log.Log.Error("when get back-to/next-stage to list, get publish by id occur error: %d", publishID)
		return nil, fmt.Errorf("when get back-to/next-stage to list, get publish by id occur error: %d", publishID)
	}
	if modelPublish.StageID != stageID {
		return nil, fmt.Errorf("发布单id: %d, 当前不在此阶段, 不允许执行该操作", publishID)
	}

	if modelPublish.Status == models.Running {
		return nil, fmt.Errorf("Publish-Order id: %d 's current step is running, operation reject", modelPublish.ID)
	}

	if modelPublish.StageID != stageID {
		return nil, fmt.Errorf("publish-order id: %d not in this stage, operation reject", modelPublish.ID)
	}
	return modelPublish, nil
}

func (pm *PublishManager) getCurrentAndRequestStage(lastPipelineInstanceID, stageID, reqStageID int64) (*pipelinemgr.PipelineStageStruct, *pipelinemgr.PipelineStageStruct, error) {
	pipelineInstanceConfig, err := pm.pipelineHandler.GetPipelineInstanceJSONByID(lastPipelineInstanceID)
	if err != nil {
		return nil, nil, fmt.Errorf("when get pipeline instance json, id: %d, occur error: %s", lastPipelineInstanceID, err)
	}

	currentEnvStage := &pipelinemgr.PipelineStageStruct{}
	currentEnvStage = pm.getCurrentEnvStageJSON(pipelineInstanceConfig, stageID, currentEnvStage)

	reqEnvStage := &pipelinemgr.PipelineStageStruct{}
	reqEnvStage = pm.getCurrentEnvStageJSON(pipelineInstanceConfig, reqStageID, currentEnvStage)

	return currentEnvStage, reqEnvStage, nil
}

// For: Back-to, Next-stage
func (pm *PublishManager) updatePublishOrderStatus(publish *models.Publish, lastPipelineID, reqStageID int64, reqStage *pipelinemgr.PipelineStageStruct, creator, stepLabel, message string) error {
	step, stepType, err := pm.pipelineHandler.GetFirstStepByPipelineInstanceIDAndStageID(lastPipelineID, reqStageID)
	if err != nil {
		log.Log.Error("when first Step/stepType by pipelineID: %v, stageID: %v, occur error: %s", lastPipelineID, reqStageID, err.Error())
		return err
	}
	log.Log.Debug("update publish pipeline status reqStage envID: %v", reqStage.StageID)
	publish.StageID = reqStage.StageID
	publish.StageName = reqStage.Name
	publish.Step = step
	publish.StepType = stepType
	publish.Status = models.Pending
	//
	publish.StepIndex = 1

	err = pm.model.UpdatePublish(publish)
	if err != nil {
		return err
	}
	createOperationLogReq := &CreateOperationLogReq{
		Creator:            creator,
		StageName:          reqStage.Name,
		StepName:           stepLabel,
		Message:            message,
		PipelineInstanceID: reqStage.PipelineInstanceID,
		StageInstanceID:    reqStage.ID,
		StepIndex:          publish.StepIndex,
		Status:             models.Success,
		PublishID:          publish.ID,
		StageID:            reqStage.StageID,
	}
	if err := pm.createPublishOperationLogItem(createOperationLogReq); err != nil {
		log.Log.Error("when back-to/Next-stage publish, create publish operation log occur error:%s", err.Error())
	}
	return nil
}

//
func (pm *PublishManager) createPublishOperationLogItem(co *CreateOperationLogReq) error {
	operationLog := &models.PublishOperationLog{
		Creator:            co.Creator,
		Type:               co.Type,
		Stage:              co.StageName,
		StageID:            co.StageID,
		Step:               co.StepName,
		Message:            co.Message,
		Status:             co.Status,
		PublishID:          co.PublishID,
		PipelineInstanceID: co.PipelineInstanceID,
		StageInstanceID:    co.StageInstanceID,
		StepIndex:          co.StepIndex,
		RunID:              co.RunID,
		JobName:            co.JobName,
	}
	return pm.model.CreatePublishOperation(operationLog)
}

func (pm *PublishManager) publishCreateParamVerify(req *PublishReq) error {
	// name
	if len(req.Name) > 64 {
		return fmt.Errorf("流水线描述过长，不允许超过64个字符")
	}
	// VersionNo
	if len(req.VersionNo) > 64 {
		return fmt.Errorf("流水线名称不允许超过64个字符")
	}
	// App
	if len(req.Apps) == 0 {
		return fmt.Errorf("请至少勾选一个代码库后，重试")
	}
	errs := []string{}
	for _, app := range req.Apps {
		if app.BranchName == "" {
			errs = []string{"请确认分支选择"}
		}
		_, err := pm.pipelineHandler.GetAppCodeCommitByBranch(app.AppID, app.BranchName)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%v", strings.Join(errs, ","))
	}
	return nil
}
