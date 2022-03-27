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

package pipelinemgr

import (
	"fmt"

	appmgr "github.com/go-atomci/atomci/internal/core/apps"
	"github.com/go-atomci/atomci/internal/core/settings"
	"github.com/go-atomci/atomci/internal/dao"
	"github.com/go-atomci/atomci/internal/middleware/log"
	"github.com/go-atomci/atomci/internal/models"
	"github.com/go-atomci/atomci/utils/query"
)

// PipelineManager ...
type PipelineManager struct {
	model           *dao.PipelineStageModel
	modelProject    *dao.ProjectModel
	modelPublish    *dao.PublishModel
	modelArrange    *dao.AppArrangeModel
	modelPublishJob *dao.PublishJobModel
	modelK8s        *dao.K8sClusterModel
	appHandler      *appmgr.AppManager
	// TODO: modelApp, modelAppArrnage change to appHandler
	modelApp        *dao.ScmAppModel
	modelAppArrange *dao.AppArrangeModel
	settingsHandler *settings.SettingManager
}

// NewPipelineManager ...
func NewPipelineManager() *PipelineManager {
	return &PipelineManager{
		model:           dao.NewPipelineStageModel(),
		modelProject:    dao.NewProjectModel(),
		modelPublish:    dao.NewPublishModel(),
		modelArrange:    dao.NewAppArrangeModel(),
		modelPublishJob: dao.NewPublishJobModel(),
		modelK8s:        dao.NewK8sClusterModel(),
		modelApp:        dao.NewScmAppModel(),
		modelAppArrange: dao.NewAppArrangeModel(),
		appHandler:      appmgr.NewAppManager(),
		settingsHandler: settings.NewSettingManager(),
	}
}

// GetPipelineInstanceJSONByID ..
func (pm *PipelineManager) GetPipelineInstanceJSONByID(pipelineInstanceID int64) (PipelineConfig, error) {
	pipelineInstance, err := pm.model.GetPipelineInstanceConfigByID(pipelineInstanceID)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}
	pipelineConfig := PipelineConfig{}
	return pipelineConfig.Struct(pipelineInstance.Config)
}

// GetPipelineInstanceEnvStageByID ..
func (pm *PipelineManager) GetPipelineInstanceEnvStageByID(pipelineInstanceID, stageID int64) (*PipelineStageStruct, error) {
	pipelineInstanceConfigJSON, err := pm.GetPipelineInstanceJSONByID(pipelineInstanceID)

	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("get pipeline instance %v, occur error: %s", pipelineInstanceID, err))
	}

	envStageJSON := &PipelineStageStruct{}
	for _, item := range pipelineInstanceConfigJSON {
		if item.StageID == stageID {
			envStageJSON = item
			break
		}
	}
	if envStageJSON.StageID == 0 {
		return nil, fmt.Errorf("can not get env stage based on lastpipelineinstance id: %v", pipelineInstanceID)
	}
	return envStageJSON, nil
}

// GetFlowComponents ..
func (pm *PipelineManager) GetFlowComponents() ([]*models.FlowComponent, error) {
	return pm.model.GetFlowComponents()
}

// GetTaskTmpls ..
func (pm *PipelineManager) GetTaskTmplByName(name string) (*models.TaskTmpl, error) {
	return pm.model.GetTaskTmplByName(name)
}

// GetTaskTmpls ..
func (pm *PipelineManager) GetTaskTmpls() ([]TaskTmplResp, error) {
	stepList, err := pm.model.GetTaskTmpls()
	if err != nil {
		return nil, err
	}
	return stepRespFormat(stepList), nil
}

func stepRespFormat(stepList []*models.TaskTmpl) []TaskTmplResp {
	stepsRsp := []TaskTmplResp{}
	for _, item := range stepList {
		subTaskStruct, err := StepStruct(item.SubTask)
		if err != nil {
			log.Log.Error("when parse flow setp id: %v sub task error: %s", item.ID, err.Error())
			subTaskStruct = []subTask{}
		}
		stepsRsp = append(stepsRsp, TaskTmplResp{
			Name:        item.Name,
			Description: item.Description,
			Creator:     item.Creator,
			TypeDisplay: item.TypeDisplay,
			Type:        item.Type,
			SubTask:     subTaskStruct,
			ComponentID: item.ComponentID,
			CreateAt:    item.CreateAt,
			UpdateAt:    item.UpdateAt,
			ID:          item.ID,
		})
	}
	return stepsRsp
}

// GetTaskTmplsByPagination ..
func (pm *PipelineManager) GetTaskTmplsByPagination(filter *query.FilterQuery) (*query.QueryResult, error) {
	res := &query.QueryResult{Item: []TaskTmplResp{}}
	stepList, count, err := pm.model.GetTaskTmplsByPagination(filter)
	if err != nil {
		return res, err
	}

	res.Item = stepRespFormat(stepList)
	if err = query.FillPageInfo(res, filter.PageIndex, filter.PageSize, int(count)); err != nil {
		return nil, err
	}
	return res, err
}

// CreateTaskTmpl ..
func (pm *PipelineManager) CreateTaskTmpl(request *TaskTmplReq, creator string) error {
	componentModel, err := pm.model.GetFlowComponentByType(request.Type)
	if err != nil {
		log.Log.Error("when crate flow step, get component by type: %s", err.Error())
		return fmt.Errorf("请选择有效的节点类型后重试")
	}
	subTaskStr, err := request.String()
	if err != nil {
		log.Log.Error("flow step req sub tasks to string error: %v", err.Error())
		return fmt.Errorf("节点的子任务解析错误，请联系管理员")
	}
	newTaskTmpl := &models.TaskTmpl{
		Name:        request.Name,
		Description: request.Description,
		Type:        request.Type,
		TypeDisplay: componentModel.Name,
		ComponentID: componentModel.ID,
		Creator:     creator,
		SubTask:     subTaskStr,
	}
	return pm.model.CreateTaskTmpl(newTaskTmpl)
}

// UpdateTaskTmpl ..
func (pm *PipelineManager) UpdateTaskTmpl(request *TaskTmplReq, stepID int64) error {
	stepModel, err := pm.model.GetTaskTmplByID(stepID)
	if err != nil {
		return err
	}
	if len(request.Name) > 64 {
		return fmt.Errorf("步骤名称不允许超过64个字符，当前长度：%v", len(request.Name))
	}
	if request.Name != "" {
		stepModel.Name = request.Name
	}
	if request.Description != "" {
		stepModel.Description = request.Description
	}
	if request.Type != "" {
		stepModel.Type = request.Type
		componentModel, err := pm.model.GetFlowComponentByType(request.Type)
		if err != nil {
			log.Log.Error("when crate flow step, get component by type: %s", err.Error())
			return fmt.Errorf("请选择有效的节点类型后重试")
		}
		stepModel.ComponentID = componentModel.ID
		stepModel.TypeDisplay = componentModel.Name
	}

	if len(request.SubTask) > 0 {
		subTaskStr, err := request.String()
		if err != nil {
			log.Log.Error("flow step req sub tasks to string error: %v", err.Error())
			return fmt.Errorf("节点的子任务解析错误，请联系管理员")
		}
		stepModel.SubTask = subTaskStr
	}

	return pm.model.UpdateTaskTmpl(stepModel)
}

// DeleteTaskTmpl ..
func (pm *PipelineManager) DeleteTaskTmpl(stepID int64) error {
	// TODO: add delete flow step delete, verify step is referenced or not.
	return pm.model.DeleteTaskTmpl(stepID)
}

// GetFirstStepByPipelineInstanceIDAndStageID ..
func (pm *PipelineManager) GetFirstStepByPipelineInstanceIDAndStageID(instanceID, stageID int64) (string, string, error) {
	pipelineConfigJSON, err := pm.GetPipelineInstanceJSONByID(instanceID)
	if err != nil {
		log.Log.Error("when GetFirstStepByPipelineInstanceIDAndStageID, get GetPipelineInstanceJSONByID occur error: %s", err.Error())
		return "", "", err
	}
	envStageJSON := &PipelineStageStruct{}
	for _, item := range pipelineConfigJSON {
		if item.StageID == stageID {
			envStageJSON = item
			break
		}
	}
	if envStageJSON.StageID == 0 {
		log.Log.Error("when based on instance id: %v, did not get validate envStage %v ", instanceID, stageID)
		return "", "", fmt.Errorf("when based on instance id: %v, did not get validate envStage %v ", instanceID, stageID)
	}

	var stepName, stepType string
	steps := envStageJSON.Steps
	if len(steps) > 0 {
		firstStep := steps[0]
		if firstStep.Index == 1 {
			TaskTmpl, err := pm.model.GetTaskTmplByID(firstStep.StepID)
			if err == nil {
				stepName = TaskTmpl.Name
				stepType = TaskTmpl.Type
			} else {
				return "", "", fmt.Errorf("当前阶段(ID: %v) 任务: %v", stageID, err.Error())
			}
		} else {
			return "", "", fmt.Errorf("当前阶段(ID: %v)未能获取到有效的第一个任务", stageID)
		}
	} else {
		return "", "", fmt.Errorf("当前阶段(ID: %v)没有找到任务定义", stageID)
	}
	return stepName, stepType, nil
}

// GetNextStepTypeByStageID ..
func (pm *PipelineManager) GetNextStepTypeByStageID(instanceID, stageID int64, stepIndex int) (string, error) {

	stageInstanceJSON, err := pm.GetPipelineInstanceEnvStageByID(instanceID, stageID)
	if err != nil {
		log.Log.Error("when GetNextStepTypeByStageID, get GetPipelineInstanceEnvStageByID by instanceid occur error: %s", err.Error())
		return "", err
	}

	steps := stageInstanceJSON.Steps
	if len(steps) == 0 {
		log.Log.Error("current statgeInstance (ID: %d) did not found the defined of steps", stageInstanceJSON.ID)
		return "", fmt.Errorf("未能获取到该阶段 实例id: %v 的任务定义，请联系管理员重试", stageInstanceJSON.ID)
	}

	var nextStep string
	nextStepIndex := stepIndex + 1
	if len(steps) < nextStepIndex {
		return "None", nil
	}
	for _, step := range steps {
		if step.Index == nextStepIndex {
			TaskTmpl, err := pm.model.GetTaskTmplByID(step.StepID)
			if err != nil {
				log.Log.Error("when GetNextStepTypeByStageID, get flow step by id: %v occur error:%s", step.StepID, err.Error())
				return "", fmt.Errorf("基于id: %v未能获取到任务类型，请联系管理员后重试", step.StepID)
			}
			nextStep = TaskTmpl.Type
			break
		}
	}
	return nextStep, nil
}
