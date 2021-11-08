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
	"strings"

	"github.com/go-atomci/atomci/middleware/log"
	"github.com/go-atomci/atomci/models"

	"github.com/go-atomci/workflow"
	"github.com/go-atomci/workflow/jenkins"

	"github.com/astaxie/beego/logs"
)

// NewWorkFlowProvide new workflow provide
func NewWorkFlowProvide(driver, addr, user, token, jobName string, flowProcessor jenkins.FlowProcessor) (workflow.WorkFlow, error) {
	var err error
	var workFlowProvider workflow.WorkFlow
	switch {
	case driver == workflow.DriverJenkins.String():
		workFlowProvider, err = jenkins.NewJenkinsClient(
			jenkins.URL(addr),
			jenkins.JenkinsUser(user),
			jenkins.JenkinsToken(token),
			jenkins.JenkinsJob(jobName),
			jenkins.Processor(flowProcessor),
		)
		if err != nil {
			log.Log.Error("%v: ", err)
			return nil, err
		}
		return workFlowProvider, nil
	}
	log.Log.Error("work flow system not configured")
	return nil, fmt.Errorf("work flow system not configured")
}

// GetStepInfo ...
func (pm *PipelineManager) GetStepInfo(projectID int64, publishID int64, stageID int64, stepName string) (interface{}, error) {
	if err := pm.verifyProjectPublish(projectID, publishID); err != nil {
		return nil, fmt.Errorf("请选择有效的项目/流水线后重试：%s", err.Error())
	}
	switch stepName {
	case "manual":
		return pm.getStepInfoByInstanceID(publishID)
	case "build":
		return pm.getPublishStepPreBranchList(projectID, publishID, stageID)
	case "deploy":
		return pm.getDeployStepAppImages(publishID)
	default:
		log.Log.Error("unknow step_name: %s", stepName)
		return nil, fmt.Errorf(fmt.Sprintf("unknown args step_name: %s", stepName))
	}
}

// RunManualStep .. return publish status, error
func (pm *PipelineManager) RunManualStep(publishID, stageID int64, request *ManualStepReq) (int64, error) {
	if err := pm.verifyProjectPublish(0, publishID); err != nil {
		return models.Skipped, fmt.Errorf("请选择有效的流水线后重试：%s", err.Error())
	}
	if _, err := pm.modelPublish.GetPublishByID(publishID); err != nil {
		return models.Failed, err
	}
	switch request.Status {
	case "success":
		return models.Success, nil
	case "failed":
		return models.Failed, nil
	default:
		log.Log.Error("request status is unexception, status: %v", request.Status)
		return models.Skipped, nil
	}
}

// RunBuildStep publish-order build operation
func (pm *PipelineManager) RunBuildStep(projectID, publishID, stageID int64, creator, stepName string, params *BuildStepReq) (int64, int64, string, error) {
	if err := pm.verifyProjectPublish(projectID, publishID); err != nil {
		return models.Failed, 0, "", fmt.Errorf("请选择有效的项目/流水线后重试：%s", err.Error())
	}

	publish, _ := pm.modelPublish.GetPublishByID(publishID)
	envStageJSON, err := pm.GetPipelineInstanceEnvStageByID(publish.LastPipelineInstanceID, stageID)
	if err != nil {
		return models.Failed, 0, "", fmt.Errorf("can not get env stage based on lastpipelineinstance id: %v", publish.LastPipelineInstanceID)
	}

	log.Log.Debug("run build step params: %+v", params)
	switch params.ActionName {
	case "trigger":
		if len(params.Apps) == 0 {
			return models.Failed, 0, "", fmt.Errorf("至少包含一个代码仓库 才允许触发构建")
		}

		runningJobVerify, jobString := pm.ifHasRunningBuildJob(projectID, stageID, publishID)
		if runningJobVerify {
			return models.Skipped, 0, "", fmt.Errorf(fmt.Sprintf("此阶段的流水线存在构建中的任务, 任务ID: %s", jobString))
		}

		// Create Publish job
		runID, jobName, err := pm.CreateBuildJob(creator, projectID, publishID, envStageJSON, params.Apps, params.EnvVars)
		if err != nil {
			return models.Failed, 0, "", err
		}
		logs.Info("create build job success, job run id: %v", runID)
		return models.Running, runID, jobName, nil
	case "terminate":
		if err := pm.publishTerminatePublish(projectID, publishID, stageID, "build"); err != nil {
			if strings.Contains(err.Error(), "操作拒绝") {
				return models.Skipped, 0, "", err
			}
			return models.TerminateFailed, 0, "", err
		}
		return models.TerminateSuccess, 0, "", nil
	default:
		log.Log.Error("step_name: %s did not defined ActionName, params: %v", stepName, params)
		return models.UnKnown, 0, "", fmt.Errorf(fmt.Sprintf("step_name: %s did not defined ActionName, params: %v", stepName, params))
	}
}

// RunBuildDeployCallBackStep publish-order build callback operation
func (pm *PipelineManager) RunBuildDeployCallBackStep(request *BuildStepCallbackReq) (int64, error) {
	// update publish job id status
	if err := pm.UpdatePublishJobStatus(request.PublishJobID, "SUCCESS"); err != nil {
		if strings.Contains(err.Error(), "already was end status") {
			return models.Skipped, nil
		}
		log.Log.Error("build callback, update publish job status occur error: %s", err.Error())
		return models.Skipped, err
	}
	return models.Success, nil
}

// RunDeployStep publish-order deploy operation
func (pm *PipelineManager) RunDeployStep(projectID, publishID, stageID int64, creator, stepName string, params *DeployStepReq) (int64, int64, string, error) {
	if err := pm.verifyProjectPublish(projectID, publishID); err != nil {
		return models.Failed, 0, "", fmt.Errorf("请选择有效的项目/流水线后重试：%s", err.Error())
	}
	publish, _ := pm.modelPublish.GetPublishByID(publishID)
	envStageJSON, err := pm.GetPipelineInstanceEnvStageByID(publish.LastPipelineInstanceID, stageID)
	if err != nil {
		return models.Failed, 0, "", fmt.Errorf("can not get env stage based on lastpipelineinstance id: %v", publish.LastPipelineInstanceID)
	}

	log.Log.Debug("run deploy step params: %+v", params)
	switch params.ActionName {
	case "trigger":
		if len(params.Apps) == 0 {
			return models.Failed, 0, "", fmt.Errorf("至少包含一个应用，才允许触发部署")
		}
		runningJobVerify, jobString := pm.ifHasRunningJob(projectID, stageID)
		if runningJobVerify {
			return models.Skipped, 0, "", fmt.Errorf(fmt.Sprintf("此阶段的流水线存在部署中的任务, 任务ID: %s", jobString))
		}

		projectAppsint := []int64{}
		for _, i := range params.Apps {
			projectAppsint = append(projectAppsint, i.ProjectAppID)
		}
		err = pm.checkApparrange(projectID, projectAppsint, envStageJSON)
		if err != nil {
			return models.Failed, 0, "", fmt.Errorf(fmt.Sprintf("checkAppArrange occur error: %s", err))
		}

		// Create Publish job
		runID, jobName, err := pm.CreateDeployJob(creator, projectID, publishID, envStageJSON, params.Apps)
		if err != nil {
			return models.Failed, 0, "", err
		}
		logs.Info("create build job success, job run id: %v", runID)
		return models.Running, runID, jobName, nil
	case "terminate":
		if err := pm.publishTerminatePublish(projectID, publishID, stageID, "deploy"); err != nil {
			if strings.Contains(err.Error(), "操作拒绝") {
				return models.Skipped, 0, "", err
			}
			return models.TerminateFailed, 0, "", err
		}
		return models.TerminateSuccess, 0, "", nil
	default:
		log.Log.Error("step_name: %s did not defined ActionName, params: %v", stepName, params)
		return models.UnKnown, 0, "", fmt.Errorf(fmt.Sprintf("step_name: %s did not defined ActionName, params: %v", stepName, params))
	}
}

// AutoTriggerNextStep ..
func (pm *PipelineManager) AutoTriggerNextStep(publish *models.Publish, nextStepType string) (int64, int64, string, error) {
	switch nextStepType {
	case "build":
		return models.Pending, 0, "", nil
	case "deploy":
		params, err := pm.generateAutoDeployStep(publish.ID)
		if err != nil {
			return models.Failed, 0, "", err
		}
		return pm.RunDeployStep(publish.ProjectID, publish.ID, publish.StageID, "admin", "deploy", params)
	default:
		log.Log.Error("stepType: %s is not exception", nextStepType)
		return models.Pending, 0, "", nil
	}
}

// GetJenkinsConfig ..
func (pm *PipelineManager) GetJenkinsConfig(stageID int64) (*JenkinsConfigRsp, error) {
	jenkinsConfigSlice, err := pm.GetCIConfig(stageID)
	if err != nil {
		return nil, err
	}
	return &JenkinsConfigRsp{
		Jenkins: jenkinsConfigSlice[0],
	}, nil
}
