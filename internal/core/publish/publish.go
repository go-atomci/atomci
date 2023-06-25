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
	"time"

	"github.com/go-atomci/atomci/internal/core/pipelinemgr"
	"github.com/go-atomci/atomci/internal/core/project"
	"github.com/go-atomci/atomci/internal/dao"
	"github.com/go-atomci/atomci/internal/middleware/log"
	"github.com/go-atomci/atomci/internal/models"
	"github.com/go-atomci/atomci/utils/query"

	"github.com/astaxie/beego/logs"
)

// PublishManager ...
type PublishManager struct {
	model           *dao.PublishModel
	pipelineModel   *dao.PipelineStageModel
	gitAppModel     *dao.ScmAppModel
	projectModel    *dao.ProjectModel
	k8sModel        *dao.K8sClusterModel
	pipelineHandler *pipelinemgr.PipelineManager
	projectHandler  *project.ProjectManager
}

// NewPublishManager ...
func NewPublishManager() *PublishManager {
	return &PublishManager{
		model:           dao.NewPublishModel(),
		pipelineModel:   dao.NewPipelineStageModel(),
		projectModel:    dao.NewProjectModel(),
		gitAppModel:     dao.NewScmAppModel(),
		k8sModel:        dao.NewK8sClusterModel(),
		pipelineHandler: pipelinemgr.NewPipelineManager(),
		projectHandler:  project.NewProjectManager(),
	}
}

// CreatePublish ...
func (pm *PublishManager) CreatePublish(user string, projectID int64, p *PublishReq) error {
	if err := pm.publishCreateParamVerify(p); err != nil {
		log.Log.Error("args verify failed: %v", err.Error())
		return fmt.Errorf("app args verify failed, please check app url/branch info, err: %v", err.Error())
	}
	firstStageID, firstStageName, step, stepType, err := pm.projectHandler.GetStageStepInfo(p.BindPipelineID)
	if err != nil {
		log.Log.Error("get stage step info failed, msg: %s", err)
		return err
	}
	timeNow, _ := time.Parse("2006-01-02 15:04:05", time.Now().Local().Format("2006-01-02 15:04:05"))

	publishModel := models.Publish{
		Addons:     models.NewAddons(),
		StartAt:    timeNow,
		Name:       p.Name,
		ProjectID:  projectID,
		StageID:    firstStageID,
		StageName:  firstStageName,
		Step:       step,
		StepType:   stepType,
		StepIndex:  1,
		Status:     models.Pending,
		PipelineID: p.BindPipelineID,
		Creator:    user,
		VersionNo:  p.VersionNo,
	}
	publishID, err := pm.model.CreatePublishifNotExist(&publishModel)
	log.Log.Debug("create publish success ID: %v", publishID)
	if err != nil {
		log.Log.Error("create publish failed, msg: %s", err)
		return err
	}
	// TODO: add transaction operations
	// create publish app
	if err := pm.createPublishApps(p.Apps, publishID); err != nil {
		return err
	}

	if pipelineInstanceID, err := pm.createPipelineInstance(p.BindPipelineID, publishID, user); err == nil {
		publishModel, err := pm.model.GetPublishByID(publishID)
		if err != nil {
			log.Log.Error("get publishby id error: %s", err.Error())
			return err
		}
		publishModel.LastPipelineInstanceID = pipelineInstanceID
		return pm.model.UpdatePublish(publishModel)
	}
	log.Log.Error("create pipeline instance occur error: %s", err.Error())
	return nil
}

// PublishList ...
func (pm *PublishManager) PublishList(projectID int64, filter *models.ProejctReleaseFilterQuery) (*query.QueryResult, error) {
	log.Log.Debug("publish filter params: %+v", filter)
	publishes, modelDatas, err := pm.model.GetPublishesList(projectID, filter)
	// add operations' label
	for _, item := range modelDatas {
		item.Operations = pm.getPublishItemCanEnableOperations(item)
		item.Previous, item.NextStep = pm.getPublishPreviousAndNextStepbyPublishModel(item)
	}
	publishes.Item = modelDatas
	return publishes, err
}

// GetPublishInfo ...
func (pm *PublishManager) GetPublishInfo(publishID int64) (*PublishInfoResp, error) {
	publishBase, err := pm.model.GetPublishByID(publishID)
	publishResp := &PublishInfoResp{}
	if err != nil {
		log.Log.Error("get publish error: %s, publish ID: %v", err, publishID)
		return nil, err
	}
	publishResp.Publish = publishBase

	// get pipeline Name by pipelineID
	if modelPipeline, errt := pm.projectModel.GetProjectPipelineByID(publishBase.PipelineID); errt == nil {
		publishResp.PipelineName = modelPipeline.Name
	} else {
		logs.Warn("when GetPipelineByID occur error: %s", errt.Error())
	}

	infoApps, err := pm.getPublishInfoApps(publishID)
	if err != nil {
		return nil, err
	}
	publishResp.Apps = infoApps

	// get operations
	publishResp.Operations = pm.getPublishItemCanEnableOperations(publishBase)

	// steps
	if steps, err := pm.getPublishStepsByPipelineInstanceID(publishBase); err == nil {
		publishResp.Steps = steps
	}
	return publishResp, err
}

// UpdatePublish ..
func (pm *PublishManager) UpdatePublish(publishID, stageID, status, runID int64, creator, message, jobName string) error {
	// status is -1, mean skip update status
	if status == models.Skipped {
		return nil
	}

	publishItem, err := pm.model.GetPublishByID(publishID)
	if err != nil {
		return err
	}

	// create operation log
	createOperationLogReq := &CreateOperationLogReq{
		Creator:            creator,
		StageName:          publishItem.StageName,
		StepName:           publishItem.Step,
		Message:            message,
		PipelineInstanceID: publishItem.LastPipelineInstanceID,
		StepIndex:          publishItem.StepIndex,
		Status:             status,
		PublishID:          publishItem.ID,
		StageID:            publishItem.StageID,
		RunID:              runID,
		JobName:            jobName,
	}
	if err := pm.createPublishOperationLogItem(createOperationLogReq); err != nil {
		log.Log.Error("when update publish order status, create publish OperationLog occur error: %s", err.Error())
	}

	nextStepIndex := publishItem.StepIndex
	nextStepType := publishItem.StepType
	nextStepName := publishItem.Step
	if status == models.Success {
		lastStage, lastStep, err := pm.pipelineHandler.CheckCurrentStepWhertherLastStageLastStep(publishID, stageID)
		if err != nil {
			log.Log.Error("when updatePublish, check current step Wherther last stage last step occur error: %s", err.Error())
		}
		if lastStep {
			if lastStage {
				status = models.END
			}
		} else {
			// when status is success, and not lastStep
			nextStepIndex = nextStepIndex + 1
			stepType, nstepName, err := pm.pipelineHandler.GetNextStepType(publishID, nextStepIndex)
			if err != nil {
				log.Log.Error("after trigger pipeline operation, get nextStepType failed: %v", err.Error())
			}
			nextStepType = stepType
			nextStepName = nstepName
			status = models.Pending

			// check driver type: auto/ manual
			status, runID, jobName, err = pm.autoDriverCheckAndTrigger(publishItem, nextStepType)
			if err != nil {
				log.Log.Error("when updatePublish, autoDriverCheckAndTrigger, occur error: %s", err.Error())
			}
			if status != models.Pending {
				// create operation log
				operationLog := &CreateOperationLogReq{
					Creator:            "system",
					StageName:          publishItem.StageName,
					StepName:           nextStepName,
					Message:            "",
					Type:               "自动流转",
					PipelineInstanceID: publishItem.LastPipelineInstanceID,
					StepIndex:          nextStepIndex,
					Status:             status,
					PublishID:          publishItem.ID,
					StageID:            publishItem.StageID,
					RunID:              runID,
					JobName:            jobName,
				}
				if err := pm.createPublishOperationLogItem(operationLog); err != nil {
					log.Log.Error("when update publish order status, create publish OperationLog occur error: %s", err.Error())
				}
			}
		}
	}
	log.Log.Debug("==>nextStepType: %v， nextStepName: %v, nextStepIndex: %v", nextStepType, nextStepName, nextStepIndex)
	return pm.updatePublishModel(publishItem, stageID, status, nextStepIndex, nextStepType, nextStepName)
}

// auto driver check
func (pm *PublishManager) autoDriverCheckAndTrigger(publishItem *models.Publish, nextStepType string) (int64, int64, string, error) {
	// check driver type: auto/ manual
	stageInstanceJSON, err := pm.pipelineHandler.GetPipelineInstanceEnvStageByID(publishItem.LastPipelineInstanceID, publishItem.StageID)
	if err != nil {
		stageInstanceJSON, err = pm.pipelineHandler.GetPipelineInstanceEnvStageByID(publishItem.LastPipelineInstanceID, publishItem.StageID)
		if err != nil {
			log.Log.Error("when update publish order check pipeline instance occur error: %s", err.Error())
			return models.Failed, 0, "", fmt.Errorf("系统错误，请联系管理员后重试")
		}
	}
	for _, step := range stageInstanceJSON.Steps {
		if step.Type == publishItem.StepType {
			switch step.Driver {
			case "auto":
				log.Log.Debug("step's Driver is auto, start autoTrigger check..")
				if nextStepType == "manual" {
					break
				}
				status, runID, jobName, err := pm.pipelineHandler.AutoTriggerNextStep(publishItem, nextStepType)
				if err != nil {
					log.Log.Error("Auto trigger next step failed, msg: %s", err.Error())
					return models.Failed, 0, "", err
				}
				log.Log.Debug("autoTrigger status: %v, runID: %v, jobName: %s", status, runID, jobName)
				return status, runID, jobName, nil
			case "manual":
				logs.Info("nextStep is manual type, no need trigger next step")
			default:
				logs.Warn("current step Driver type %s is not except", step.Driver)
			}
			break
		}
	}
	return models.Pending, 0, "", nil
}

// updatePublishModel ..
func (pm *PublishManager) updatePublishModel(publishItem *models.Publish, stageID, status int64, nextStepIndex int, nextStepType, nextStepName string) error {
	publishItem.Status = status
	publishItem.StepIndex = nextStepIndex
	publishItem.StepType = nextStepType
	publishItem.Step = nextStepName
	if stageID > 0 {
		publishItem.StageID = stageID
		modelStage, err := pm.projectModel.GetProjectEnvByID(stageID)
		if err != nil {
			log.Log.Error("when update publish, get project env by id: %v, occur errror: %s", stageID, err.Error())
			return err
		}
		publishItem.StageName = modelStage.Name
	}

	return pm.model.UpdatePublish(publishItem)
}

// ClosePublish ..
func (pm *PublishManager) ClosePublish(publishID int64) error {
	publish, err := pm.model.GetPublishByID(publishID)
	if err != nil {
		log.Log.Error("close publish error: %s", err.Error())
		return err
	}
	if publish.Status == models.Running {
		return fmt.Errorf("流水线的状态为：执行中，禁止归档")
	}
	timeNow, _ := time.Parse("2006-01-02 15:04:05", time.Now().Local().Format("2006-01-02 15:04:05"))
	publish.Status = models.Closed
	publish.EndAt = &timeNow
	return pm.model.UpdatePublish(publish)
}

// UpdatePublishBaseInfo ..
func (pm *PublishManager) UpdatePublishBaseInfo(publishID int64, req *PublishUpdate) error {
	publish, err := pm.model.GetPublishByID(publishID)
	if err != nil {
		log.Log.Error("when update publish , get publish error: %s", err.Error())
		return err
	}
	if publish.Status == models.Closed {
		return fmt.Errorf("流水线的状态为：已归档，禁止更新")
	}
	publish.Name = req.Name
	publish.VersionNo = req.VersionNo
	return pm.model.UpdatePublish(publish)
}

// DeletePublish ...
func (pm *PublishManager) DeletePublish(publishID int64) error {
	publish, err := pm.model.GetPublishByID(publishID)
	if err != nil {
		log.Log.Error("delete publish error: %s", err.Error())
		return err
	}
	if publish.Status != models.Running {
		publish.MarkDeleted()
		return pm.model.UpdatePublish(publish)
	}
	return fmt.Errorf("流水线的状态为：执行中，禁止删除操作，请稍后重试")
}

// GetBackTo ...
func (pm *PublishManager) GetBackTo(projectID, publishID, stageID int64) ([]*CirculationRsp, error) {
	_, err := pm.publishOrderBaseVerify(publishID, stageID)
	if err != nil {
		return nil, err
	}
	publishItem, _ := pm.model.GetPublishByID(publishID)

	pipelineInstanceConfig, err := pm.pipelineHandler.GetPipelineInstanceJSONByID(publishItem.LastPipelineInstanceID)
	if err != nil {
		log.Log.Error("when get back to list, get pipeline stage occur error: %s", err)
		return nil, err
	}

	currentEnvStage := &pipelinemgr.PipelineStageStruct{}
	currentEnvStage = pm.getCurrentEnvStageJSON(pipelineInstanceConfig, publishItem.StageID, currentEnvStage)

	index := currentEnvStage.Index
	resp := []*CirculationRsp{}

	// index is 1, return current statge directly
	if index == 1 {
		envStage, err := pm.projectModel.GetProjectEnvByID(currentEnvStage.StageID)
		if err != nil {
			return nil, err
		}
		resp = append(resp, &CirculationRsp{
			ID:   envStage.ID,
			Name: envStage.Name,
		})
		return resp, nil
	}

	for _, stage := range pipelineInstanceConfig {
		if stage.Index <= index {
			envStage, err := pm.projectModel.GetProjectEnvByID(stage.StageID)
			if err != nil {
				return nil, err
			}
			resp = append(resp, &CirculationRsp{
				ID:   envStage.ID,
				Name: envStage.Name,
			})
		}
	}
	return resp, nil
}

// TriggerBackTo ..
func (pm *PublishManager) TriggerBackTo(projectID, publishID, stageID int64, req *TriggerBackToReq, currentUser string) error {
	modelPublish, err := pm.publishOrderBaseVerify(publishID, stageID)
	if err != nil {
		return err
	}
	currentStage, reqStage, err := pm.getCurrentAndRequestStage(modelPublish.LastPipelineInstanceID, stageID, req.StageID)
	if err != nil {
		log.Log.Error("trigger back to, get stage instance occur error: %s", err.Error())
		return fmt.Errorf("获取阶段实例失败，请联系管理员重试")
	}
	if currentStage.Index < reqStage.Index {
		return fmt.Errorf("backTo operation can only be returned to the current or previous stage, publish-Order id: %d", modelPublish.ID)
	}
	// TODO: confirm
	pipelineInstanceID, err := pm.createPipelineInstance(modelPublish.PipelineID, publishID, currentUser)
	if err != nil {
		log.Log.Error("create pipeline instance occur error: %v", err.Error())
		return fmt.Errorf("创建流水线实例失败，无法回退，请联系管理员重试")
	}
	modelPublish.LastPipelineInstanceID = pipelineInstanceID
	return pm.updatePublishOrderStatus(modelPublish, pipelineInstanceID, req.StageID, reqStage, currentUser, "back-to", req.Message)
}

// GetNextStage ...
func (pm *PublishManager) GetNextStage(projectID, publishID, envID int64) ([]*CirculationRsp, error) {
	_, err := pm.publishOrderBaseVerify(publishID, envID)
	if err != nil {
		return nil, err
	}
	publishItem, _ := pm.model.GetPublishByID(publishID)
	modelStageInstance, err := pm.pipelineHandler.GetPipelineInstanceEnvStageByID(publishItem.LastPipelineInstanceID, envID)

	if err != nil {
		log.Log.Error("when get next stage list, get pipeline stage instance occur error: %s", err)
		return nil, err
	}
	index := modelStageInstance.Index
	resp := []*CirculationRsp{}
	stages, err := pm.pipelineHandler.GetPipelineInstanceJSONByID(publishItem.LastPipelineInstanceID)
	if err != nil {
		log.Log.Error("when get next stage, get GetPipelineInstanceJSONByID occur error: %s, reset revious/next step to empty", err.Error())
		return nil, err
	}

	stageLength := len(stages)
	if index == int64(stageLength) {
		return nil, fmt.Errorf("current publish  pipeline id: %d, env id: %v already is pipeline last stage,  operation reject", publishID, modelStageInstance.ID)
	}
	nextIndex := index + 1
	for _, stage := range stages {
		if stage.Index == nextIndex {
			flowStage, err := pm.projectModel.GetProjectEnvByID(stage.StageID)
			if err != nil {
				return nil, err
			}
			resp = append(resp, &CirculationRsp{
				ID:   stage.StageID,
				Name: flowStage.Name,
			})
			break
		}
	}
	return resp, nil
}

// TriggerNextStage ..
func (pm *PublishManager) TriggerNextStage(projectID, publishID, envID int64, req *TriggerBackToReq, currentUser string) error {
	modelPublish, err := pm.publishOrderBaseVerify(publishID, envID)
	if err != nil {
		return err
	}
	currentStage, reqStage, err := pm.getCurrentAndRequestStage(modelPublish.LastPipelineInstanceID, envID, req.StageID)
	if err != nil {
		return fmt.Errorf("when trigger next stage current envid: %v, req envid: %v error: %s", envID, req.StageID, err.Error())
	}
	if currentStage.Index > reqStage.Index {
		return fmt.Errorf("NextStage operation can only be returned to the next stage, publish-Order id: %d", modelPublish.ID)
	}
	return pm.updatePublishOrderStatus(modelPublish, modelPublish.LastPipelineInstanceID, req.StageID, reqStage, currentUser, "next-stage", "")
}

// GetPublishOperationLog ..
func (pm *PublishManager) GetPublishOperationLog(publishID int64, filter *query.FilterQuery) (*query.QueryResult, error) {
	return pm.model.GetOperationLogsByPublishID(publishID, filter)
}
