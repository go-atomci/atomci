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

package projects

import (
	"encoding/json"
	"fmt"

	"github.com/go-atomci/atomci/core/pipelinemgr"
	"github.com/go-atomci/atomci/middleware/log"
	"github.com/go-atomci/atomci/models"
	"github.com/go-atomci/atomci/utils/query"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

type PipelineReq struct {
	Description string      `json:"description,omitempty"`
	Name        string      `json:"name,omitempty"`
	ProjectID   int64       `json:"project_id,omitempty"`
	IsDefault   bool        `json:"is_default,omitempty"`
	Config      interface{} `json:"config,omitempty"`
}

type ProjectPipelineRespone struct {
	Description string                     `json:"description,omitempty"`
	Name        string                     `json:"name,omitempty"`
	ID          int64                      `json:"id,omitempty"`
	Creator     string                     `json:"creator,omitempty"`
	Config      pipelinemgr.PipelineConfig `json:"config,omitempty"`
}

type ProjectEnvReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Cluster     int64  `json:"cluster"`
	Namespace   string `json:"namespace"`
	ArrangeEnv  string `json:"arrange_env"`
	CIServer    int64  `json:"ci_server"`
	Harbor      int64  `json:"harbor"`
}

func (s *PipelineReq) String() (string, error) {
	bytes, err := json.Marshal(s.Config)
	return string(bytes), err
}

// Struct ...
func (config *ProjectPipelineRespone) Struct(sc string) ([]*pipelinemgr.PipelineStageStruct, error) {
	stages := []*pipelinemgr.PipelineStageStruct{}
	err := json.Unmarshal([]byte(sc), &stages)
	return stages, err
}

// UpdateProjectPipelineConfig ...
func (pm *ProjectManager) UpdateProjectPipelineConfig(request *PipelineReq, creator string, projectID, pipelineID int64) error {
	if pipelineID == 0 {
		return fmt.Errorf("参数错误(pipelineID: 0)请联系管理员重试")
	}

	pipelineModel, err := pm.model.GetProjectPipelineByID(pipelineID)
	if err != nil {
		return fmt.Errorf("when update pipeline base info, get pipeline occur error: %s", err.Error())
	}

	if len(request.Description) > 0 {
		pipelineModel.Description = request.Description
	}
	if len(request.Name) > 0 {
		pipelineModel.Name = request.Name
	}

	if request.IsDefault {
		pipelineItem, err := pm.model.GetDefaultPipeline(request.ProjectID)
		if err == nil && pipelineItem.ID != pipelineModel.ID {
			return fmt.Errorf("已经存在默认流程: %v，请先将其改为非默认再重试", pipelineItem.Name)
		}
	}

	pipelineModel.IsDefault = request.IsDefault

	// TODO: config verify  请确保设置的每个阶段均已经添加任务节点
	configString, err := request.String()
	if err != nil {
		return err
	}

	if len(configString) > 0 {
		pipelineModel.Config = configString
	}

	return pm.model.UpdateProjectPipeline(pipelineModel)
}

// GetPipelineConfig ..
func (pm *ProjectManager) GetPipelineConfig(pipelineID int64) (ProjectPipelineRespone, error) {
	pipelineInfo, err := pm.model.GetProjectPipelineByID(pipelineID)
	if err != nil {
		return ProjectPipelineRespone{}, fmt.Errorf("%s", err.Error())
	}
	configJSON := pipelinemgr.PipelineConfig{}
	resp := ProjectPipelineRespone{
		Name:        pipelineInfo.Name,
		Description: pipelineInfo.Description,
		Creator:     pipelineInfo.Creator,
	}
	log.Log.Debug("pipeline config: %v", pipelineInfo.Config)
	if len(pipelineInfo.Config) > 0 {
		configJSON, err = configJSON.Struct(pipelineInfo.Config)
		if err != nil {
			errMsg := fmt.Errorf("pipeline config json parse error: %s", err.Error())
			log.Log.Error(err.Error())
			return resp, errMsg
		}
	}
	resp.Config = configJSON
	return resp, nil

}

// GetStageStepInfo For pipline instance (stage/step) init
func (pm *ProjectManager) GetStageStepInfo(pipelineID int64) (int64, string, string, string, error) {
	pipelineResp, err := pm.GetPipelineConfig(pipelineID)
	if err != nil {
		return 0, "", "", "", err
	}
	if len(pipelineResp.Config) == 0 {
		return 0, "", "", "", fmt.Errorf("流水线 %s 未设置有效阶段，请联系管理员，确认后重试", pipelineResp.Name)
	}
	var (
		firstStageID                       int64
		firstStageName, stepName, stepType string
	)
	// TODO: sort
	for _, stage := range pipelineResp.Config {
		if stage.Index == 1 {
			flowStage, err := pm.model.GetProjectEnvByID(stage.StageID)
			if err != nil {
				log.Log.Error("when get project env, occur error: %s", err.Error())
				return 0, "", "", "", fmt.Errorf("网络错误，请重试")
			}
			firstStageID = flowStage.ID
			firstStageName = flowStage.Name

			if len(stage.Steps) == 0 {
				log.Log.Error("when pipeline %v env: %v steps len is 0, invalidate", pipelineID, stage.StageID)
				return 0, "", "", "", fmt.Errorf("网络错误，请重试")
			}
			stepName = stage.Steps[0].Name
			stepType = stage.Steps[0].Type

			break
		}
	}
	log.Log.Debug("get stage step : firststageID %v, name %s, step: %s", firstStageID, firstStageName, stepName)
	return firstStageID, firstStageName, stepName, stepType, err
}

// GetPipelinesByPagination ..
func (pm *ProjectManager) GetPipelinesByPagination(filter *query.FilterQuery, projectID int64) (*query.QueryResult, error) {
	return pm.model.GetPipelinesByPagination(filter, projectID)
}

func (pm *ProjectManager) CreateProjectPipeline(request *PipelineReq, creator string) (int64, error) {
	if request.ProjectID == 0 {
		return 0, fmt.Errorf("project id %v info is invalidation", 0)
	}

	if request.IsDefault {
		if pipelineItem, err := pm.model.GetDefaultPipeline(request.ProjectID); err == nil {
			return 0, fmt.Errorf("已经存在默认流程: %v，请先将其改为非默认再重试", pipelineItem.Name)
		}
	}

	pipeline := &models.ProjectPipeline{
		Name:        request.Name,
		Description: request.Description,
		Creator:     creator,
		ProjectID:   request.ProjectID,
		IsDefault:   request.IsDefault,
	}

	return pm.model.CreatePipeline(pipeline)
}

// GetProjectEnvs ..
func (pm *ProjectManager) GetProjectEnvs(projectID int64) ([]*models.ProjectEnv, error) {
	return pm.model.GetProjectEnvs(projectID)
}

// GetProjectEnvsByPagination ..
func (pm *ProjectManager) GetProjectEnvsByPagination(filter *query.FilterQuery, projectID int64) (*query.QueryResult, error) {
	return pm.model.GetProjectEnvsByPagination(filter, projectID)
}

// UpdateProjectEnv ..
func (pm *ProjectManager) UpdateProjectEnv(request *ProjectEnvReq, stepID int64) error {
	stageModel, err := pm.model.GetProjectEnvByID(stepID)
	if err != nil {
		return err
	}
	if request.Name != "" {
		stageModel.Name = request.Name
	}

	if request.Description != "" {
		stageModel.Description = request.Description
	}
	if request.Cluster != 0 {
		stageModel.Cluster = request.Cluster
	}
	if request.Namespace != "" {
		stageModel.Namespace = request.Namespace
	}

	if request.CIServer != 0 {
		stageModel.CIServer = request.CIServer
	}
	if request.Harbor != 0 {
		stageModel.Harbor = request.Harbor
	}

	return pm.model.UpdateProjectEnv(stageModel)
}

// CreateProjectEnv ..
func (pm *ProjectManager) CreateProjectEnv(request *ProjectEnvReq, creator string, projectID int64) error {
	// TODO: verify cluster, nammespace
	if request.Cluster == 0 {
		return fmt.Errorf("你请选择集群")
	}
	if request.ArrangeEnv == "" {
		return fmt.Errorf("你请选择环境标识后，再重试")
	}

	// TODO: verify projectID is validate
	if projectID == 0 {
		return fmt.Errorf("无效的 project id: %v", projectID)
	}
	existStage, err := pm.model.GetProjectEnvBycIDAndEnvTag(request.ArrangeEnv, projectID)
	if err == nil {
		return fmt.Errorf("环境标识必须唯一，%v 环境已经使用此标识 %s，请你更新后重试", existStage.Name, request.ArrangeEnv)
	}
	if err != nil && err != orm.ErrNoRows {
		logs.Warn("when create flow stage, GetProjectEnvBycIDAndArrangeEnv check occur error:%s", err.Error())
	}
	newProjectEnv := &models.ProjectEnv{
		ProjectID:   projectID,
		Name:        request.Name,
		Description: request.Description,
		Cluster:     request.Cluster,
		Namespace:   request.Namespace,
		CIServer:    request.CIServer,
		Harbor:      request.Harbor,
		ArrangeEnv:  request.ArrangeEnv,
		Creator:     creator,
	}
	return pm.model.CreateProjectEnv(newProjectEnv)
}

// DeleteProjectEnv ..
func (pm *ProjectManager) DeleteProjectEnv(stageID int64) error {
	// TODO: when delete env, verify env id is referenced or not.
	return pm.model.DeleteProjectEnv(stageID)
}
