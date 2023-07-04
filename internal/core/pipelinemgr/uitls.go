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
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-atomci/atomci/constant"
	"github.com/go-atomci/atomci/internal/core/apps"
	"github.com/go-atomci/atomci/internal/core/kuberes"
	"github.com/go-atomci/atomci/internal/core/settings"
	"github.com/go-atomci/atomci/internal/dao"
	"github.com/go-atomci/atomci/internal/middleware/log"
	"github.com/go-atomci/atomci/internal/models"
	"github.com/go-atomci/atomci/utils"

	"github.com/drone/go-scm/scm"
	"github.com/go-atomci/workflow"
	"github.com/go-atomci/workflow/jenkins"
	"github.com/go-atomci/workflow/jenkins/templates"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

var (
	atomciServer = beego.AppConfig.String("atomci::url")
)

func (pm *PipelineManager) getManualStepInfo(instanceID, stageID int64, stepIndex int) (*StepRsp, error) {
	operationLogs, err := pm.modelPublish.GetOperationLogByInstanceIDAndStageIDStepType(instanceID, stageID, stepIndex)
	if err != nil {
		log.Log.Error("when get manual step info, get operation log params instanceID: %v, stageID: %v, stepIndex: %v \noccur errror: %s", instanceID, stageID, stepIndex, err.Error())
		return nil, err
	}

	currentStepRsp := &StepRsp{}
	// Get Current Step Operation
	if len(operationLogs) > 0 {
		latestOperation := operationLogs[0]
		currentStepRsp = &StepRsp{
			Name:    latestOperation.Step,
			Creator: latestOperation.Creator,
			Message: latestOperation.Message,
		}
	}
	return currentStepRsp, nil
}

// Pipeline Operation: manual step, getStepInfoByInstanceID
func (pm *PipelineManager) getStepInfoByInstanceID(publishID int64) (*ManualStepResp, error) {
	// base publishID get lastPipelineInstanceID/ step_id
	publishModel, err := pm.modelPublish.GetPublishByID(publishID)
	if err != nil {
		log.Log.Error("when get manual step info, get publish occur error: %s", err.Error())
		return nil, err
	}
	instanceID, stageID, stepIndex := publishModel.LastPipelineInstanceID, publishModel.StageID, publishModel.StepIndex

	rsp := ManualStepResp{}
	// Get Current Step Operation
	StepRsp, err := pm.getManualStepInfo(instanceID, stageID, stepIndex)
	if err != nil {
		return nil, err
	}
	rsp.CurrenStep = StepRsp

	// Get Pervious Step Operation
	if stepIndex == 1 {
		rsp.PreviousStep = nil
		return &rsp, nil
	}
	previousStep, err := pm.getManualStepInfo(instanceID, stageID, stepIndex-1)
	if err != nil {
		log.Log.Error("when get manual step info, get previous step occur error: %s", err.Error())
		return nil, err
	}
	rsp.PreviousStep = previousStep
	return &rsp, nil
}

func (pm *PipelineManager) getUserToken(user string) (string, error) {
	userInfo, err := dao.NewUserRolesModel().GetUserByName(user)
	if err != nil {
		log.Log.Error("when get %v token by name occur error: %s", user, err.Error())
		return "", err
	}
	return userInfo.Token, nil
}

// generate compileEnv based on project app compileEnvID
func (pm *PipelineManager) generateCompileEnvParams(apps []*RunBuildAppReq) []compileEnv {
	compileParams := []compileEnv{}
	for _, item := range apps {
		projectApp, err := pm.modelProject.GetProjectApp(item.ProjectAppID)
		if err != nil {
			logs.Warn("project app error: %s", err.Error())
			continue
		}
		scmApp, err := pm.modelApp.GetScmAppByID(projectApp.ScmID)
		if err != nil {
			logs.Warn("get scm app error: %s", err.Error())
			continue
		}

		if scmApp.CompileEnvID == 0 {
			log.Log.Debug("app: %v didnot setup complie env, use default docker runtime", scmApp.Name)
			continue
		}
		compileItem, err := pm.settingsHandler.GetCompileEnvByID(scmApp.CompileEnvID)
		if err != nil {
			logs.Warn("get compile env by id:%v error: %s", scmApp.CompileEnvID, err.Error())
		}
		if compileItem.Name == constant.DefaultContainerName {
			log.Log.Warn("app: %v setup complie env to %v, skip this compileItem generate", constant.DefaultContainerName, scmApp.Name)
			continue
		}
		compileEnvItem := compileEnv{
			Image:      compileItem.Image,
			Args:       compileItem.Args,
			Command:    compileItem.Command,
			WorkingDir: "/home/jenkins/agent",
			Name:       strings.ToLower(scmApp.Name),
		}
		compileParams = append(compileParams, compileEnvItem)
	}
	return compileParams
}

func (pm *PipelineManager) getSysDefaultCompileEnv(compileEnvName string) (jenkins.ContainerEnv, error) {
	compileEnv, err := pm.settingsHandler.GetCompileEnvByName(compileEnvName)
	if err != nil {
		return jenkins.ContainerEnv{}, err
	}

	return jenkins.ContainerEnv{
		Name:       compileEnv.Name,
		Image:      compileEnv.Image,
		WorkingDir: "/home/jenkins/agent",
		// TODO: command / args is valid
		CommandArr: commandAndArgSplit(compileEnv.Command),
		ArgsArr:    commandAndArgSplit(compileEnv.Args),
	}, nil
}

// CreateBuildJob return publishjob run id, error
func (pm *PipelineManager) CreateBuildJob(creator string, projectID, publishID int64, envStageJSON *PipelineStageStruct, apps []*RunBuildAppReq, customeEnvVars []EnvItem) (int64, string, error) {
	// Prerequisites -jenkins
	CIInfo, err := pm.GetCIConfig(envStageJSON.StageID)
	if err != nil {
		log.Log.Error("getCIConfig occur error: %s", err.Error())
		return 0, "", err
	}
	if len(CIInfo) != 5 {
		log.Log.Error("get ci config len is not 5, ciinfo: %+v", CIInfo)
		return 0, "", fmt.Errorf("get ci config len is not 5, ciinfo: %+v", CIInfo)
	}
	addr, user, token := CIInfo[0], CIInfo[1], CIInfo[2]

	jenkinsClient, err := NewWorkFlowProvide(workflow.DriverJenkins.String(), addr, user, token, "", nil)
	if err != nil {
		return 0, "", err
	}
	if _, err := jenkinsClient.Ping(); err != nil {
		return 0, "", fmt.Errorf("jenkins is unhealthy, error: %s", err.Error())
	}

	publishItem, err := pm.modelPublish.GetPublishByID(publishID)
	if err != nil {
		log.Log.Error("get publish order occur error: %s", err.Error())
		return 0, "", err
	}

	stepSubTasks := []*subTask{}
	stepIndex := publishItem.StepIndex
	compileParams := pm.generateCompileEnvParams(apps)

	for _, item := range envStageJSON.Steps {
		if item.Index == stepIndex && item.Type == constant.StepBuild {
			log.Log.Debug("step index: %v, item.type: %v", stepIndex, constant.StepBuild)
			// stepSubTasks = item.SubTasks
			// step sub tasks defined
			stepSubTasks = item.SubTask
			if len(stepSubTasks) == 0 {
				logs.Warn("sub tasks redefined")
				stepSubTasks = []*subTask{
					{
						Index: 1,
						Name:  "代码检出",
						Type:  "checkout",
					},
					{
						Index:  2,
						Name:   "编译",
						Type:   "compile",
						Params: compileParams,
					},
					{
						Index: 3,
						Name:  "制作镜像",
						Type:  "build-image",
					},
				}
			} else {
				// type type is compile, add params rely on app language
				for _, subTask := range stepSubTasks {
					if subTask.Type == constant.StepSubTaskCompile {
						subTask.Params = compileParams
						log.Log.Debug("add params for subTask type: %v, index: %v", subTask.Type, subTask.Index)
					}
				}
			}
			break
		}
	}

	if len(stepSubTasks) == 0 {
		log.Log.Error("this build jod did not have sub tasks, or rquest invalid maybe current step index is not %V", constant.StepBuild)
		return 0, "", fmt.Errorf("this build jod did not have sub tasks, or rquest invalid")
	}

	// Aggregate the app parms for build based on request params
	appsAllParams, _ := pm.aggregateAppsParamsForBuild(apps, envStageJSON)

	deployInfo, _, err := pm.getDeployInfo(envStageJSON.StageID)
	if err != nil {
		log.Log.Error("getDeployInfo occur error: %s", err.Error())
		return 0, "", err
	}
	if len(deployInfo) != 4 {
		log.Log.Error("deploy info is validate, len: %v", len(deployInfo))
	}
	// Create publishJob publishJobApps
	appsParamsForJob := []*AppParamsForCreatePublishJob{}
	for _, param := range appsAllParams {
		paramForJob := &AppParamsForCreatePublishJob{
			ProjectAppID: param.ProjectAppID,
			Branch:       param.Branch,
			Path:         param.Path,
			// TODO: image version get based on image tag rule
			ImageVersion: "",
		}
		appsParamsForJob = append(appsParamsForJob, paramForJob)
	}

	publishJobID, err := pm.CreatePublishJob(projectID, publishID, envStageJSON.StageID, creator, "build", appsParamsForJob)
	if err != nil {
		log.Log.Error("when create build job, create publish job error: %s", err.Error())
		return 0, "", err
	}
	jobName := fmt.Sprintf("atomci_%v_%v_%v", projectID, publishID, envStageJSON.StageID)

	jenkinsJNLPTemplate, err := pm.getSysDefaultCompileEnv(constant.DefaultContainerName)
	if err != nil {
		log.Log.Error("when create build job, get sys default %v compile env error: %s", constant.DefaultContainerName, err.Error())
		return 0, "", err
	}

	jenkinsCheckoutTemplate, err := pm.getSysDefaultCompileEnv(constant.CheckoutContainerName)
	if err != nil {
		log.Log.Error("when create build job, get sys default checkout compile env  error: %s", err.Error())
		return 0, "", err
	}
	jenkinsKanikoTemplate, err := pm.getSysDefaultCompileEnv(constant.BuildImageContainerName)
	if err != nil {
		log.Log.Error("when create build job, get sys default kaniko compile env  error: %s", err.Error())
		return 0, "", err
	}

	// default container template
	containerTemplates := []jenkins.ContainerEnv{
		jenkinsJNLPTemplate,
		jenkinsCheckoutTemplate,
		jenkinsKanikoTemplate,
	}
	// TaskTmplItem.SubTask
	taskPipelineXMLStrArr := []string{}
	for _, subTask := range stepSubTasks {
		taskPipelineXMLStr := ""
		switch subTask.Type {
		case constant.StepSubTaskCheckout:
			//
			appCheckoutItems, err := pm.renderAppCheckoutItemsForBuild(projectID, envStageJSON.StageID, publishJobID, appsAllParams)
			if err != nil {
				return 0, "", err
			}
			items := map[string]interface{}{"CheckoutItems": appCheckoutItems}
			taskPipelineXMLStr, err = jenkins.GeneratePipelineXMLStr(templates.Checkout, items)
			if err != nil {
				return 0, "", err
			}
		case constant.StepSubTaskCompile:
			for _, compileItem := range subTask.Params {
				log.Log.Debug("sub task image: %v", compileItem.Name)
				compileContainerItem := jenkins.ContainerEnv{
					Name:       compileItem.Name,
					Image:      compileItem.Image,
					CommandArr: commandAndArgSplit(compileItem.Command),
					ArgsArr:    commandAndArgSplit(compileItem.Args),
					WorkingDir: compileItem.WorkingDir,
				}
				containerTemplates = append(containerTemplates, compileContainerItem)
			}

			appBuildItems, err := pm.renderAppBuildItemsForBuild(projectID, envStageJSON.StageID, publishJobID, appsAllParams, CIInfo)
			if err != nil {
				return 0, "", err
			}
			items := map[string]interface{}{"BuildItems": appBuildItems}
			taskPipelineXMLStr, err = jenkins.GeneratePipelineXMLStr(templates.Compile, items)
			if err != nil {
				return 0, "", err
			}

		case constant.StepSubTaskBuildImage:
			//
			appImageItems, err := pm.renderAppImageitemsForBuild(projectID, publishID, envStageJSON.StageID, publishJobID, appsAllParams, CIInfo, deployInfo)
			if err != nil {
				return 0, "", err
			}
			items := map[string]interface{}{"ImageItems": appImageItems}
			taskPipelineXMLStr, err = jenkins.GeneratePipelineXMLStr(templates.BuildImage, items)
			if err != nil {
				return 0, "", err
			}

		default:
			logs.Info("%v sub task type did not matched, taskPipelineXmlStr is empty value", subTask.Type)
		}
		taskPipelineXMLStrArr = append(taskPipelineXMLStrArr, taskPipelineXMLStr)
	}

	pipelineStagesStr := strings.Join(taskPipelineXMLStrArr, " ")

	if len(apps) == 0 {
		log.Log.Error("project app len is 0, invalidate")
		return 0, "", fmt.Errorf("project app len is 0, invalidate")
	}
	oneAppReq := apps[0]
	projectApp, err := pm.modelProject.GetProjectApp(oneAppReq.ProjectAppID)
	if err != nil {
		log.Log.Error("when crate build job, get project app error: %s", err.Error())
		return 0, "", err
	}
	scmApp, err := pm.modelApp.GetScmAppByID(projectApp.ScmID)
	if err != nil {
		log.Log.Error("when crate build job, get scm app error: %s", err.Error())
		return 0, "", err
	}

	scmIntegrateResp, err := pm.settingsHandler.GetSCMIntegrateSettinByID(scmApp.RepoID)
	if err != nil {
		log.Log.Error("when crate build job, get scm integrate setting error: %s", err.Error())
		return 0, "", err
	}

	baseURL := strings.Replace(scmIntegrateResp.URL, "http://", "", -1)
	baseURL = strings.Replace(baseURL, "https://", "", -1)
	if strings.HasSuffix(baseURL, "/") {
		baseURL = strings.Replace(baseURL, "/", "", -1)
	}
	repoConfStr := fmt.Sprintf("{\"%s\":[\"%s\",\"%s\"]}", baseURL, scmIntegrateResp.User, scmIntegrateResp.Token)

	adminToken, err := pm.getUserToken("admin")
	if err != nil {
		log.Log.Error("get admin token occur error: %v", err.Error())
		return 0, "", fmt.Errorf("网络错误，请重试")
	}

	// TODO: Input correct env values
	envVars := []jenkins.EnvItem{
		{Key: "JENKINS_SLAVE_WORKSPACE", Value: CIInfo[3]},
		{Key: "ACCESS_TOKEN", Value: adminToken},
		{Key: "REPO_CNF", Value: repoConfStr},
		{Key: "DOCKER_AUTH", Value: deployInfo[2]},
		{Key: "REGISTRY_ADDR", Value: deployInfo[1]},
		{Key: "DOCKER_CONFIG", Value: "/kaniko/.docker"},
	}

	for _, env := range customeEnvVars {
		jenkinsEnvItem := jenkins.EnvItem{
			Key:   env.Key,
			Value: env.Value,
		}
		envVars = append(envVars, jenkinsEnvItem)
	}

	callBackURL := fmt.Sprintf("%s/atomci/api/v1/pipelines/%d/publishes/%d/stages/%d/steps/%s/callback", atomciServer, projectID, publishID, envStageJSON.StageID, "build")
	callBackRequestBody := fmt.Sprintf("{\"publish_job_id\": %d}", publishJobID)

	// k8sDeployInfo, err := pm.getDeployInfo(stageJSON.StageID)
	// k8sDeployInfo: []string{harbor.HarborName, harbor.HarborAddr, flowStage.ArrangeEnv, harbor.HarborUser, harbor.HarborPassword}
	// if err != nil {
	// return 0, "", err
	// }

	flowProcessor := &jenkins.CIContext{
		EnvVars:            envVars,
		ContainerTemplates: containerTemplates,
		Stages:             pipelineStagesStr,
		CommonContext: jenkins.CommonContext{
			Namespace:        CIInfo[4],
			PodName:          strings.ReplaceAll(jobName, "_", "-"),
			DefaultContainer: constant.DefaultContainerName,
		},
		CallBack: jenkins.CallbackRequest{
			Token: adminToken,
			URL:   callBackURL,
			Body:  callBackRequestBody,
		},
	}

	workerflowClient, err := NewWorkFlowProvide(workflow.DriverJenkins.String(), addr, user, token, jobName, flowProcessor)
	if err != nil {
		log.Log.Error("when new workflow provide error: %s", err.Error())
		return 0, "", err
	}
	runID, err := workerflowClient.Build()
	if err != nil {
		// TODO: deleted publishjob item already created
		return 0, "", err
	}
	// Update runID/status to publishjob
	err = pm.UpdatePublishJob(publishJobID, runID)
	if err != nil {
		return 0, "", err
	}
	return runID, jobName, nil
}

// CreateDeployJob return publishjob run id, error
func (pm *PipelineManager) CreateDeployJob(creator string, projectID, publishID int64, stageJSON *PipelineStageStruct, apps []*RunDeployAppReq) (int64, string, error) {
	// Aggregate the app parms for deploy based on request params
	appsAllParams, _ := pm.aggregateAppsParamsForDeploy(publishID, stageJSON.StageID, apps, stageJSON)

	CIInfo, err := pm.GetCIConfig(stageJSON.StageID)
	if err != nil {
		log.Log.Error("getCIConfig occur error: %s", err.Error())
		return 0, "", err
	}
	addr, user, token := CIInfo[0], CIInfo[1], CIInfo[2]

	jenkinsClient, err := NewWorkFlowProvide(workflow.DriverJenkins.String(), addr, user, token, "", nil)
	if err != nil {
		return 0, "", err
	}
	if _, err := jenkinsClient.Ping(); err != nil {
		return 0, "", fmt.Errorf("jenkins is unhealthy, error: %s", err.Error())
	}

	// Create publishJob publishJobApps
	appsParamsForJob := []*AppParamsForCreatePublishJob{}
	for _, param := range appsAllParams {
		paramForJob := &AppParamsForCreatePublishJob{
			ProjectAppID: param.ProjectAppID,
			Path:         param.Path,
			ImageAddr:    param.ImageAddr,
		}
		appsParamsForJob = append(appsParamsForJob, paramForJob)
	}

	jobName := fmt.Sprintf("atomci_%v_%v", projectID, stageJSON.StageID)

	// deploy app, combine app arrange to temmplateStr
	templateStr, err := pm.renderTemplateStr(apps, publishID, stageJSON.StageID)
	if err != nil {
		return 0, "", err
	}

	envModel, err := pm.modelProject.GetProjectEnvByID(stageJSON.StageID)
	if err != nil {
		log.Log.Error("when create deploy job, get project env by id occur error: %s", err.Error())
		return 0, "", err
	}

	clusterModel, err := pm.settingsHandler.GetIntegrateSettingByID(envModel.Cluster)
	if err != nil {
		log.Log.Error("when create deploy job, get cluster by id %v occur error: %s", envModel.Cluster, err.Error())
		return 0, "", err
	}

	err = kuberes.TriggerApplicationCreate(clusterModel.Name, envModel.Namespace, templateStr, projectID, stageJSON.StageID, true)
	if err != nil {
		log.Log.Error("when crate deploy job, trigger application create occur error: %s", err.Error())
		return 0, "", err
	}

	appsParamsHealth := []*AppParamsForHealthCheck{}
	for _, param := range appsAllParams {
		item := &AppParamsForHealthCheck{
			Name:     param.Name,
			ID:       param.ProjectAppID,
			FullName: param.FullName,
			Type:     "app",
		}
		appsParamsHealth = append(appsParamsHealth, item)
	}
	publishJobID, err := pm.CreatePublishJob(projectID, publishID, stageJSON.StageID, creator, "deploy", appsParamsForJob)
	if err != nil {
		return 0, "", err
	}

	healthCheckItems, err := pm.renderHealthCheckCommand(projectID, stageJSON.StageID, publishJobID, appsParamsHealth, stageJSON)
	if err != nil {
		return 0, "", err
	}

	callBackURL := fmt.Sprintf("%s/atomci/api/v1/pipelines/%d/publishes/%d/stages/%d/steps/%s/callback", atomciServer, projectID, publishID, stageJSON.StageID, "deploy")
	callBackRequestBody := fmt.Sprintf("{\"publish_job_id\": %d}", publishJobID)
	// TODO: Use pm.getAppConfig() get all config
	adminToken, err := pm.getUserToken("admin")
	if err != nil {
		log.Log.Error("get admin token occur error: %v", err.Error())
		return 0, "", fmt.Errorf("网络错误，请重试")
	}
	userToken, err := pm.getUserToken(creator)
	if err != nil {
		log.Log.Error("get %v token occur error: %v", creator, err.Error())
		return 0, "", fmt.Errorf("网络错误，请重试")
	}

	envVars := []jenkins.EnvItem{
		{Key: "JENKINS_SLAVE_WORKSPACE", Value: CIInfo[3]},
		{Key: "ATOMCI_SERVER", Value: atomciServer},
		{Key: "ACCESS_TOKEN", Value: adminToken},
		{Key: "USER_TOKEN", Value: userToken},
	}

	jenkinsJNLPTemplate, err := pm.getSysDefaultCompileEnv(constant.DefaultContainerName)
	if err != nil {
		log.Log.Error("when create deploy job, get sys default %v compile env error: %s", constant.DefaultContainerName, err.Error())
		return 0, "", err
	}

	jenkinsCheckoutTemplate, err := pm.getSysDefaultCompileEnv(constant.CheckoutContainerName)
	if err != nil {
		log.Log.Error("when create deploy job, get sys default checkout compile env  error: %s", err.Error())
		return 0, "", err
	}

	// default container template
	containerTemplates := []jenkins.ContainerEnv{
		jenkinsJNLPTemplate,
		jenkinsCheckoutTemplate,
	}

	flowProcessor := &jenkins.DeployContext{
		HealthCheckItems:   healthCheckItems,
		EnvVars:            envVars,
		ContainerTemplates: containerTemplates,
		CallBack: jenkins.CallbackRequest{
			Token: adminToken,
			URL:   callBackURL,
			Body:  callBackRequestBody,
		},
		CommonContext: jenkins.CommonContext{
			Namespace:        CIInfo[4],
			PodName:          strings.ReplaceAll(jobName, "_", "-"),
			DefaultContainer: constant.DefaultContainerName,
		},
	}

	workerflowClient, err := NewWorkFlowProvide(workflow.DriverJenkins.String(), addr, user, token, jobName, flowProcessor)
	if err != nil {
		return 0, "", err
	}
	runID, err := workerflowClient.Build()
	if err != nil {
		return 0, "", err
	}
	// Update runID/status to publishjob
	err = pm.UpdatePublishJob(publishJobID, runID)
	if err != nil {
		return 0, "", err
	}
	return runID, jobName, nil
}
func (pm *PipelineManager) renderTemplateStr(apps []*RunDeployAppReq, publishID, envID int64) (string, error) {
	var templateStr string
	for _, item := range apps {
		arrange, err := pm.appHandler.GetRealArrange(item.ProjectAppID, envID)
		if err != nil {
			log.Log.Error("get app id: %v  env id: %v real arrange, occur error: %s", item.ProjectAppID, envID, err.Error())
			continue
		}

		// replace template str
		arrangeConfig := arrange.Config
		publishApp, err := pm.modelPublish.GetPublishAppByPublishIDAndAppID(publishID, item.ProjectAppID)
		if err != nil {
			logs.Warn("when get publish app by publishid/appid occur error:%s, did not update app arrange image info", err.Error())
			continue
		}

		newImageAddr, originImage, err := pm.generateImageAddr(arrange.ID, item.ProjectAppID, publishApp.BranchName)
		if err != nil {
			continue
		}
		arrangeConfig = strings.Replace(arrangeConfig, originImage, newImageAddr, -1)
		if templateStr == "" {
			templateStr = arrangeConfig
		} else {
			templateStr = templateStr + "\n---\n" + arrangeConfig
		}
	}
	return templateStr, nil
}

func (pm *PipelineManager) generateImageAddr(arrangeID, projectAppID int64, branch string) (string, string, error) {
	imageMapping, err := pm.modelAppArrange.GetAppImageMappingByArrangeIDAndProjectAppID(arrangeID, projectAppID)
	if err != nil {
		log.Log.Error("get imagemapping error: %s", err.Error())
		return "", "", err
	}
	newImageAddr := imageMapping.Image
	switch imageMapping.ImageTagType {
	case models.SystemDefaultTag:
		// branch get from RunBuildAppReq.Branch
		imageTag, err := pm.GetAppCodeCommitByBranch(projectAppID, branch)
		if err != nil {
			logs.Error("when get app code commit by branch error: %s, did not update app arrange image info", err.Error())
			return "", "", err
		}

		imageStr, _ := removeImageUrlTag(imageMapping.Image)
		newImageAddr = fmt.Sprintf("%s:%s", imageStr, imageTag)
	case models.LatestTag:
		imageStr, _ := removeImageUrlTag(imageMapping.Image)
		newImageAddr = fmt.Sprintf("%s:%s", imageStr, "latest")
	case models.OriginTag:
		log.Log.Debug("image tag use from yaml, no need replace")
	}
	return newImageAddr, imageMapping.Image, nil
}

func removeImageUrlTag(imageUrl string) (string, error) {
	// 10.10.0.8/abc
	// 10.10.0.8/abc:tag
	// 10.10.0.8/abc:tag:deg 错误

	// 10.10.0.8:9980/abc
	// 10.10.0.8:9980/abc:tag
	// 10.10.0.8:9980/abc:tag:def 错误

	if !strings.Contains(imageUrl, ":") {
		return imageUrl, nil
	}
	originImageSplit := strings.Split(imageUrl, ":")
	if len(originImageSplit) > 2 {
		real := imageUrl[0:strings.LastIndex(imageUrl, ":")]
		return real, nil
	}
	if strings.Contains(originImageSplit[1], "/") {
		return imageUrl, nil
	}
	real := imageUrl[0:strings.LastIndex(imageUrl, ":")]
	return real, nil

	return "", nil
}

func (pm *PipelineManager) GetAppCodeCommitByBranch(appID int64, branchName string) (string, error) {
	projectApp, err := pm.modelProject.GetProjectApp(appID)
	if err != nil {
		log.Log.Error("when get app code commit, get project ap by id: %v error:%s", appID, err.Error())
		return "", err
	}

	scmApp, err := pm.modelApp.GetScmAppByID(projectApp.ScmID)
	if err != nil {
		log.Log.Error("when get app code commit, get scm ap by id: %v error:%s", appID, err.Error())
		return "", err
	}

	scmIntegrateResp, err := pm.settingsHandler.GetSCMIntegrateSettinByID(scmApp.RepoID)
	if err != nil {
		return "", err
	}
	client, err := apps.NewScmProvider(scmIntegrateResp.Type, scmIntegrateResp.URL, scmIntegrateResp.Token)
	if err != nil {
		return "", err
	}
	opt := scm.CommitListOptions{
		Ref:  branchName,
		Page: 1,
		Size: 10,
	}

	got, _, err := client.Git.ListCommits(context.Background(), scmApp.FullName, opt)
	if err != nil {
		return "", err
	}

	if len(got) > 0 {
		return branchName + "-" + got[0].Sha[0:7], nil
	} else {
		logs.Warn("branch: %v did not include any commit", branchName)
		return "", fmt.Errorf("应用:%v 分支:%v 未包含任何提交, 请通过“我的应用”-“应用详情”-“同步远程分支”后重新选择", scmApp.Name, branchName)
	}
}

// Pipeline Operation:: publish step, get branch list for publish
func (pm *PipelineManager) getPublishStepPreBranchList(projectID, publishID, stageID int64) (*BuildStepResp, error) {
	targetBranch := []string{"master"}
	publishApps, err := pm.modelPublish.GetPublishAppsByID(publishID)

	publishStepResp := []*PublishStepResp{}
	for _, app := range publishApps {
		projectApp, _ := pm.modelProject.GetProjectApp(app.ProjectAppID)
		scmApp, err := pm.modelApp.GetScmAppByID(projectApp.ScmID)
		if err != nil {
			log.Log.Error("get scm app by id %v error: %s", projectApp.ScmID, err.Error())
			continue
		}
		branchHistoryList, _ := pm.modelApp.GetAppBranches(app.ProjectAppID)
		branchItems := []string{}
		for _, branch := range branchHistoryList {
			branchItems = append(branchItems, branch.BranchName)
		}
		if len(branchItems) == 0 {
			branchItems = []string{"master"}
		}
		appInfo := &PublishStepResp{
			BranchName:        app.BranchName,
			AppName:           scmApp.Name,
			Language:          scmApp.Language,
			ProjectAppID:      app.ProjectAppID,
			BuildPath:         scmApp.BuildPath,
			Type:              "app",
			TargetBranch:      targetBranch,
			CompileCommand:    app.CompileCommand,
			BranchHistoryList: branchItems,
		}
		publishStepResp = append(publishStepResp, appInfo)
	}
	publish, _ := pm.modelPublish.GetPublishByID(publishID)
	rsp := &BuildStepResp{
		VersionNo:   publish.VersionNo,
		VersionName: publish.Name,
		Apps:        publishStepResp,
	}
	return rsp, err

}

// Pipeline Operation:: publish step, terminate publish
func (pm *PipelineManager) publishTerminatePublish(projectID, publishID, stageID int64, jobType string) error {
	publishOrder, err := pm.modelPublish.GetPublishByID(publishID)
	if err != nil {
		return err
	}
	if publishOrder.Status == models.TerminateSuccess {
		return fmt.Errorf("publish Order already terminated, operation reject")
	}

	if !utils.IntContains([]int64{models.Running, models.TerminateFailed}, publishOrder.Status) {
		return fmt.Errorf("publish Order current status is not allowed terminate, operation reject")
	}

	CIInfo, err := pm.GetCIConfig(stageID)
	if err != nil {
		log.Log.Error("getCIConfig occur error: %s", err.Error())
		return err
	}
	addr, user, token := CIInfo[0], CIInfo[1], CIInfo[2]

	var jobName string
	switch jobType {
	case "build":
		jobName = fmt.Sprintf("atomci_%v_%v_%v", projectID, publishID, stageID)
	case "deploy":
		jobName = fmt.Sprintf("atomci_%v_%v", projectID, stageID)
	default:
		log.Log.Error("jobType: %s is noexception", jobType)
		return fmt.Errorf("不支持此任务类型: %v 的终止", jobType)
	}
	workerflowClient, err := NewWorkFlowProvide(workflow.DriverJenkins.String(), addr, user, token, jobName, nil)
	if err != nil {
		return err
	}
	latestPublishJob, err := pm.modelPublishJob.GetLastPublishJobByPublishID(publishID)
	if err != nil {
		return err
	}

	if err := workerflowClient.Abort(latestPublishJob.RunID); err != nil {
		return err
	}
	return pm.updatePublishJob(latestPublishJob, models.StatusAbort)
}

// Pipeline Operation:: deploy step getDeployStepAppImages
func (pm *PipelineManager) getDeployStepAppImages(publishID int64) ([]*DeployStepAppRsp, error) {
	publishApps, err := pm.modelPublish.GetPublishAppsByID(publishID)
	if err != nil {
		log.Log.Error("when getDeployStepAppImages, get publishAppbyID occur error: %s", err.Error())
		return nil, err
	}
	rsp := []*DeployStepAppRsp{}
	for _, app := range publishApps {
		projectApp, err := pm.modelProject.GetProjectApp(app.ProjectAppID)
		if err != nil {
			logs.Warn("project app id: %v not exist, err: %s", app.ProjectAppID, err.Error())
			continue
		}
		scmApp, err := pm.modelApp.GetScmAppByID(projectApp.ScmID)
		if err != nil {
			logs.Warn("scm app id: %v not exist, err: %s", projectApp.ScmID, err.Error())
			continue
		}
		item := &DeployStepAppRsp{
			ProjectAppID: projectApp.ID,
			Name:         scmApp.Name,
			Type:         "app",
		}
		rsp = append(rsp, item)
	}
	return rsp, nil
}

func (pm *PipelineManager) verifyProjectPublish(projectID, publishID int64) error {
	if projectID > 0 {
		project, err := pm.modelProject.GetProjectByID(projectID)
		if err != nil {
			log.Log.Error("when verifyProjectPublish, getProjectByID occur error: %s", err.Error())
			return err
		}
		if project.Status == models.ProjectEnd {
			return fmt.Errorf(fmt.Sprintf("项目: %v 已经结束，请联系管理员开启项目后重试", project.Name))
		}
	}
	if publishID > 0 {
		if _, err := pm.modelPublish.GetPublishByID(publishID); err != nil {
			log.Log.Error("when verifyProjectPublish, getProjectByID occur error: %s", err.Error())
			return err
		}
	}
	return nil
}

// check current stage whether have running job.
func (pm *PipelineManager) ifHasRunningBuildJob(projectID, stageID, publishID int64) (bool, string) {
	// Query publishjob app based on projectID & stageID
	runningJobs, err := pm.modelPublishJob.GetCurrentRunningBuildJob(projectID, stageID, publishID, []string{models.StatusRunning, models.StatusInit}, "build")
	if err != nil {
		log.Log.Error("when trigger build, current running job verify occur error: %s, reset verify is true", err)
		return false, ""
	}
	if len(runningJobs) > 0 {
		log.Log.Error("when trigger build, current already have running job, verify is false")
		jobIDs := []string{}
		for _, job := range runningJobs {
			jobIDs = append(jobIDs, strconv.FormatInt(job.ID, 10))
		}
		jobString := strings.Join(jobIDs, ",")
		return true, jobString
	}
	return false, ""
}

// check current stage whether have running job.
func (pm *PipelineManager) ifHasRunningJob(projectID, envID int64) (bool, string) {
	// Query publishjob app based on projectID & envID
	runningJobs, err := pm.modelPublishJob.GetCurrentRunningJob(projectID, envID, []string{models.StatusRunning, models.StatusInit}, "deploy")
	if err != nil {
		log.Log.Error("when trigger publish, current running job verify occur error: %s, reset verify is true", err)
		return false, ""
	}
	if len(runningJobs) > 0 {
		log.Log.Error("when trigger publish, current already have running job, verify is false")
		jobIDs := []string{}
		for _, job := range runningJobs {
			jobIDs = append(jobIDs, strconv.FormatInt(job.ID, 10))
		}
		jobString := strings.Join(jobIDs, ",")
		return true, jobString
	}
	return false, ""
}

// verify apps whether already setup arrange in current stage,
func (pm *PipelineManager) checkApparrange(projectID int64, apps []int64, stage *PipelineStageStruct) error {
	appIDs := []int64{}
	appIDs = append(appIDs, apps...)
	flowStageID := stage.StageID
	envModel, err := pm.modelProject.GetProjectEnvByID(flowStageID)
	if err != nil {
		return err
	}
	arrangeEnvID := envModel.ID
	modelApps, err := pm.modelProject.GetProjectAppsByIDs(projectID, appIDs)
	if err != nil {
		log.Log.Error("get git apps occur error: %s", err)
		return err
	}

	nilArranged := []string{}
	for _, modelApp := range modelApps {
		_, err := pm.appHandler.GetRealArrange(modelApp.ID, arrangeEnvID)
		if err != nil {
			log.Log.Error("get project app id: %v arrnage occur error: %s", modelApp.ID, err)
			scmApp, err := pm.modelApp.GetScmAppByID(modelApp.ScmID)
			if err != nil {
				log.Log.Warn("get scm app error: %s", err.Error())
				continue
			}
			nilArranged = append(nilArranged, scmApp.Name)
		}
	}
	if len(nilArranged) > 0 {
		nilArrangedName := strings.Join(nilArranged, ",")
		log.Log.Error("apps: %s on stage %v(%v) did not setup arrange correctly", nilArrangedName, stage.Name, envModel.Name)
		return fmt.Errorf("请保存应用: %v %v 的『应用编排』后重试", nilArrangedName, envModel.Name)
	}
	return nil
}

// when terminate publish-order, update publish job'status to abort
func (pm *PipelineManager) updatePublishJob(publishJob *models.PublishJob, status string) error {
	publishJob.Status = status
	if err := pm.modelPublishJob.UpdatePublishJob(publishJob); err != nil {
		return err
	}
	return nil
}

/*  Generate Commands For Jenkins Default Pipeline  */

func (pm *PipelineManager) aggregateAppsParamsForBuild(apps []*RunBuildAppReq, stageJSON *PipelineStageStruct) ([]*RunBuildAllParms, error) {
	allParms := []*RunBuildAllParms{}
	for _, app := range apps {
		projectApp, err := pm.modelProject.GetProjectApp(app.ProjectAppID)
		if err != nil {
			log.Log.Error("get proejct modelapp occur error: %s", err)
		}
		scmApp, err := pm.modelApp.GetScmAppByID(projectApp.ScmID)
		if err != nil {
			logs.Warn("get scm app error: %s", err.Error())
			continue
		}
		releaseBranch := "None"
		allParm := &RunBuildAllParms{
			ProjectID:      projectApp.ProjectID,
			ScmApp:         scmApp,
			RunBuildAppReq: app,
			Release:        releaseBranch,
		}
		allParms = append(allParms, allParm)

	}
	return allParms, nil
}

func (pm *PipelineManager) aggregateAppsParamsForDeploy(publishID, stageID int64, apps []*RunDeployAppReq, stageJSON *PipelineStageStruct) ([]*RunDeployAllParms, error) {

	allParms := []*RunDeployAllParms{}
	for _, app := range apps {
		projectApp, err := pm.modelProject.GetProjectApp(app.ProjectAppID)
		if err != nil {
			log.Log.Error("get gitmodelapp occur error: %s", err)
		}

		scmApp, err := pm.modelApp.GetScmAppByID(projectApp.ScmID)
		if err != nil {
			log.Log.Error("get scmapp error: %s", err.Error())
		}
		arrange, err := pm.appHandler.GetRealArrange(app.ProjectAppID, stageID)
		if err != nil {
			log.Log.Error("get app id: %v  env id: %v real arrange, occur error: %s", app.ProjectAppID, stageID, err.Error())
			continue
		}

		publishApp, err := pm.modelPublish.GetPublishAppByPublishIDAndAppID(publishID, app.ProjectAppID)
		if err != nil {
			logs.Warn("when get publish app by publishid/appid occur error:%s, did not update app arrange image info", err.Error())
			continue
		}

		newImageAddr, _, err := pm.generateImageAddr(arrange.ID, app.ProjectAppID, publishApp.BranchName)
		if err != nil {
			continue
		}
		log.Log.Debug("imageAddr: %s", newImageAddr)
		allParm := &RunDeployAllParms{
			ProjectID:       projectApp.ProjectID,
			ScmApp:          scmApp,
			RunDeployAppReq: app,
			ImageAddr:       newImageAddr,
		}
		allParms = append(allParms, allParm)
	}
	return allParms, nil
}

/*  auto Trigger part start */

// AutoTrigger Deploy
func (pm *PipelineManager) generateAutoDeployStep(publishID int64) (*DeployStepReq, error) {
	// TODO: 应该基于 publishID 索引到最新的一个 publishjob, then get according publishjobapps, 而不是直接根据 publish app 来部署；
	publishApps, err := pm.modelPublish.GetPublishAppsByID(publishID)
	if err != nil {
		log.Log.Error("when AutoTriggerNextStep, GetPublishAppsByID occur error: %s", err.Error())
		return nil, err
	}
	apps := []*RunDeployAppReq{}
	for _, app := range publishApps {
		app := &RunDeployAppReq{
			ProjectAppID: app.ProjectAppID,
		}
		apps = append(apps, app)
	}
	params := &DeployStepReq{
		ActionName: "trigger",
		Apps:       apps,
	}
	return params, nil
}

/*  auto Trigger part end */

func (pm *PipelineManager) generateBaseInfo(projectID, stageID, publishJobID int64) string {
	scriptBaseInfo := fmt.Sprintf(" --project-id %d --stage-id %d --publish-job-id %d ", projectID, stageID, publishJobID)
	return scriptBaseInfo
}
func (pm *PipelineManager) generateAppPth(stageID, projectID int64, workSpace string, appArgs *RunBuildAllParms) string {
	appPath := strings.Join([]string{workSpace, strconv.Itoa(int(projectID)), strconv.Itoa(int(stageID)), appArgs.Name, appArgs.Branch, appArgs.BuildPath}, "/")
	return strings.ReplaceAll(appPath, "//", "/")
}

// Rendering parameters for app checkout items's command
func (pm *PipelineManager) renderAppCheckoutItemsForBuild(projectID, stageID, publishJobID int64, allParms []*RunBuildAllParms) ([]jenkins.StepItem, error) {
	appCheckoutItems := []jenkins.StepItem{}

	buildBaseInfo := pm.generateBaseInfo(projectID, stageID, publishJobID)
	for _, app := range allParms {
		// TODO: if GitAPP type is not app, how to deal with this, skip ??
		item := jenkins.StepItem{}
		item.Name = app.Name

		// TODO: app build vcsType use git
		appInfoStr := fmt.Sprintf(" --scm-app-id %d --app-name %s --app-language %s --branch-url %s --vcs-type %s --build-path %s ", app.ProjectAppID, app.Name, app.Language, app.Path, "git", app.BuildPath)
		appParms := fmt.Sprintf(" --branch-name %s ", app.Branch)
		Command := fmt.Sprintf("sh 'checkout.py %s %s %s'", buildBaseInfo, appInfoStr, appParms)
		item.Command = Command
		item.ContainerName = constant.CheckoutContainerName
		appCheckoutItems = append(appCheckoutItems, item)
	}

	return appCheckoutItems, nil
}

// Rendering parameters for app build items's command
func (pm *PipelineManager) renderAppBuildItemsForBuild(projectID, stageID, publishJobID int64, allParms []*RunBuildAllParms, ciConfig []string) ([]*jenkins.StepItem, error) {
	appBuildItems := []*jenkins.StepItem{}

	for _, app := range allParms {
		item := &jenkins.StepItem{}
		item.Name = app.Name
		// Default containername is constant.DefaultContainerName(jnlp)
		item.ContainerName = constant.DefaultContainerName
		command := fmt.Sprintf("sh 'echo app:%v language:%v, did not defined compile command, skip compile'", app.Name, app.Language)
		customCompileCommand := app.RunBuildAppReq.CompileCommand

		appPath := pm.generateAppPth(stageID, projectID, ciConfig[3], app)
		appRootPath := appPath
		if app.CompileEnvID == 0 {
			command = fmt.Sprintf("sh 'echo app:%v language:%v, did not setup compile env,skip compile...'", app.Name, app.Language)
		} else if len(customCompileCommand) > 0 {
			item.ContainerName = strings.ToLower(app.Name)
			command = fmt.Sprintf("sh 'cd %v; %v'", appRootPath, customCompileCommand)
		}
		item.Command = command
		appBuildItems = append(appBuildItems, item)
	}

	return appBuildItems, nil
}

// Rendering parameters for app images items's command
func (pm *PipelineManager) renderAppImageitemsForBuild(projectID, publishID, stageID, publishJobID int64, allParms []*RunBuildAllParms, ciConfig []string, deployInfo []string) ([]*jenkins.StepItem, error) {
	appImageItems := []*jenkins.StepItem{}

	if len(ciConfig) != 5 {
		log.Log.Error("ciConfig is invalide, real len: %v", len(ciConfig))
	}

	if len(deployInfo) == 0 {
		log.Log.Error("deployinfo is invalide, real len: %v", len(deployInfo))
		return nil, fmt.Errorf("ciConfig is invalide, real len: %v", len(deployInfo))
	}
	for _, app := range allParms {
		item := &jenkins.StepItem{}
		item.Name = app.Name

		appPath := pm.generateAppPth(stageID, projectID, ciConfig[3], app)

		arrange, err := pm.appHandler.GetRealArrange(app.ProjectAppID, stageID)
		if err != nil {
			log.Log.Error("get app id: %v  env id: %v real arrange, occur error: %s", app.ProjectAppID, stageID, err.Error())
			continue
		}

		imageURL, _, err := pm.generateImageAddr(arrange.ID, app.ProjectAppID, app.Branch)
		if err != nil {
			continue
		}
		dockerfile := app.Dockerfile
		if dockerfile == "" {
			dockerfile = "Dockerfile"
		}
		var insecure = ""
		if isHttps, _ := strconv.ParseBool(deployInfo[3]); !isHttps {
			insecure = "--insecure --skip-tls-verify --insecure-pull"
		}
		Command := fmt.Sprintf("sh \"cd %v; export DOCKER_CONFIG=$DOCKER_CONFIG; /kaniko/executor -f %v -c ./  -d %v %s \"", appPath, dockerfile, imageURL, insecure)
		item.Command = Command
		item.ContainerName = constant.BuildImageContainerName
		appImageItems = append(appImageItems, item)
	}

	return appImageItems, nil
}

// Rendering parameters for healthcheck command
func (pm *PipelineManager) renderHealthCheckCommand(projectID, stageID, publishJobID int64, allParms []*AppParamsForHealthCheck, stageJSON *PipelineStageStruct) ([]*jenkins.StepItem, error) {
	healthCheckItems := []*jenkins.StepItem{}
	buildBaseInfo := pm.generateBaseInfo(projectID, stageID, publishJobID)

	envStage, err := pm.modelProject.GetProjectEnvByID(stageJSON.StageID)
	if err != nil {
		log.Log.Error("when render healthcheck command, get flow stage byid: %v, occur error: %s", stageJSON.StageID, err.Error())
		return nil, fmt.Errorf("未能获取到指定的阶段")
	}
	for _, app := range allParms {
		// skip GitAPP's type is not app
		if app.Type != "app" {
			logs.Info("app name: %s type: %s, is not app type, skip health check", app.Name, app.Type)
			continue
		}
		item := &jenkins.StepItem{}
		item.Name = app.Name
		// TODO: refactor(WIP) healthcheck stage, use checkout container name tmp.
		item.ContainerName = constant.CheckoutContainerName

		appArrange, err := pm.appHandler.GetRealArrange(app.ID, stageID)
		if err != nil {
			log.Log.Warn("get app id: %v, env id: %v occur error: %s", err.Error())
			continue
		}

		native := &kuberes.NativeTemplate{
			Template: appArrange.Config,
		}
		appResItems, err := native.GetAppResourceNames()
		if err != nil {
			log.Log.Warn("parse app arrange occur error: %s", err.Error())
			continue
		}

		settingKubernetesItem, err := pm.settingsHandler.GetIntegrateSettingByID(envStage.Cluster)
		if err != nil {
			log.Log.Error("integrate setting cluster by id: %v error: %s", envStage.Cluster, err.Error())
			return nil, fmt.Errorf("integrate setting cluster by id: %v error: %s", envStage.Cluster, err.Error())
		}

		for _, appRes := range appResItems {
			svcName := appRes.Name
			svcInfo := fmt.Sprintf(" --cluster %s --namespace %s --app-name %s --service-name %s", settingKubernetesItem.Name, envStage.Namespace, app.Name, svcName)
			item.Command = fmt.Sprintf("sh 'healthcheck.py %s %s'", buildBaseInfo, svcInfo)
			healthCheckItems = append(healthCheckItems, item)
		}
	}
	return healthCheckItems, nil
}

// GetCIConfig ..
func (pm *PipelineManager) GetCIConfig(stageID int64) ([]string, error) {
	projectEnv, err := pm.modelProject.GetProjectEnvByID(stageID)
	if err != nil {
		log.Log.Error("when getCIConfig, GetProjectEnvByID %v occur error: %s", stageID, err.Error())
		return nil, fmt.Errorf("未能找到到 id: %v 的配置，请联系管理员后重试", stageID)
	}
	CIServer := projectEnv.CIServer
	log.Log.Debug("current CIServer integrate_setting id: %v", CIServer)
	settingItem, err := pm.settingsHandler.GetIntegrateSettingByID(CIServer)
	if err != nil {
		log.Log.Error("when get ci config, get integrate setting by id: %v error: %s", CIServer, err.Error())
		return nil, err
	}
	if settingItem.Type != "jenkins" {
		return []string{}, fmt.Errorf("settings type is: %s, current ci server only support jenkins", settingItem.Type)
	}
	var url, user, token, namespace, workSpace string
	if jenkinsConfig, ok := settingItem.Config.(*settings.JenkinsConfig); ok {
		url = jenkinsConfig.URL
		user = jenkinsConfig.User
		token = jenkinsConfig.Token
		namespace = jenkinsConfig.Namespace
		if namespace == "" {
			namespace = "devops"
		}
		workSpace = jenkinsConfig.WorkSpace
	} else {
		log.Log.Error("parse jenkins config error")
		return []string{}, fmt.Errorf("parse jenkins config error")
	}
	log.Log.Debug("jenkins user: %v, url: %v, token: %v, workspace: %v", user, url, token, workSpace)
	if url == "" || user == "" || token == "" || workSpace == "" {
		return nil, fmt.Errorf("请联系管理员确认 系统管理-服务集成 %v 的配置, 当前配置为: url: %v, user: %v, token: %v, workSpace: %v", settingItem.Name, url, user, token, workSpace)
	}
	return []string{url, user, token, workSpace, namespace}, nil
}

// getDeployInfo cluster,registry auth info,arrangeEnv
func (pm *PipelineManager) getDeployInfo(stageID int64) ([]string, int64, error) {
	envStage, err := pm.modelProject.GetProjectEnvByID(stageID)
	if err != nil {
		log.Log.Error("when get deploy info, get project env by id:%v, errror: %v", stageID, err.Error())
		return nil, 0, err
	}

	settingKubernetesItem, err := pm.settingsHandler.GetIntegrateSettingByID(envStage.Cluster)
	if err != nil {
		log.Log.Error("integrate setting cluster by id: %v error: %s", envStage.Cluster, err.Error())
		return []string{}, 0, fmt.Errorf("integrate setting cluster by id: %v error: %s", envStage.Cluster, err.Error())
	}
	if settingKubernetesItem.Type != "kubernetes" {
		return []string{}, 0, fmt.Errorf("settings type is: %s, current deploy server only support kubernetes", settingKubernetesItem.Type)
	}

	settingRegistryItem, err := pm.settingsHandler.GetIntegrateSettingByID(envStage.Registry)
	if err != nil {
		log.Log.Error("integrate setting registry by id: %v error: %s", envStage.Registry, err.Error())
		return []string{}, 0, fmt.Errorf("integrate setting registry by id: %v error: %s", envStage.Registry, err.Error())
	}
	if settingRegistryItem.Type != "registry" {
		return []string{}, 0, fmt.Errorf("settings type is: %s, current deploy server only support kubernetes", settingRegistryItem.Type)
	}

	var registryAddr, registryAuth string
	var isHttps bool
	if registryConf, ok := settingRegistryItem.Config.(*settings.RegistryConfig); ok {
		registryAddr = registryConf.URL
		registryAuth = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", registryConf.User, registryConf.Password)))
		isHttps = registryConf.IsHttps
	} else {
		log.Log.Error("parse kubernetes config error")
		return []string{}, 0, fmt.Errorf("parse jenkins config error")
	}
	return []string{settingKubernetesItem.Name, registryAddr, registryAuth, strconv.FormatBool(isHttps)}, envStage.ID, nil
}

func (pm *PipelineManager) publishStepVerify(publishID int64, step string) (bool, error) {
	publish, _ := pm.modelPublish.GetPublishByID(publishID)
	err := fmt.Errorf("当前步骤是: %s, 不允许此操作", publish.Step)
	switch step {
	case models.StepVerify:
		if publish.StepType == step && publish.Status == models.Success || publish.Status == models.Failed {
			return false, err
		}
	}
	return true, nil
}

// CheckCurrentStepWhertherLastStageLastStep ..
func (pm *PipelineManager) CheckCurrentStepWhertherLastStageLastStep(publishID, stageID int64) (bool, bool, error) {
	publishItem, _ := pm.modelPublish.GetPublishByID(publishID)
	pipelineInstanceID, stepIndex := publishItem.LastPipelineInstanceID, publishItem.StepIndex

	var lastStage, lastStep bool
	currentStageInstanceJSON, err := pm.GetPipelineInstanceEnvStageByID(pipelineInstanceID, stageID)
	if err != nil {
		return false, false, fmt.Errorf("网络异常，请重试")
	}
	steps := currentStageInstanceJSON.Steps
	lastStepIndex := steps[len(steps)-1:][0].Index
	if lastStepIndex == stepIndex {
		lastStep = true
	}
	pipelineStagesJSON, err := pm.GetPipelineInstanceJSONByID(pipelineInstanceID)
	if err != nil {
		return false, false, err
	}

	lastPipeLinearr := pipelineStagesJSON[len(pipelineStagesJSON)-1:]
	lastPipeLineJSON := lastPipeLinearr[0]

	log.Log.Debug("lastPipelineJSON: %+v", lastPipeLineJSON)
	if lastPipeLineJSON.StageID == stageID {
		lastStage = true
	}
	return lastStage, lastStep, nil
}

// GetNextStepType ..
func (pm *PipelineManager) GetNextStepType(publishID int64, stepIndex int) (string, string, error) {
	publishItem, _ := pm.modelPublish.GetPublishByID(publishID)
	stageJSON, err := pm.GetPipelineInstanceEnvStageByID(publishItem.LastPipelineInstanceID, publishItem.StageID)
	if err != nil {
		return "", "", fmt.Errorf("网络异常，请重试")
	}

	stepsJSON := stageJSON.Steps
	var stepType, stepName string
	for _, step := range stepsJSON {
		if step.Index == stepIndex {
			stepModel, err := pm.model.GetTaskTmplByID(step.StepID)
			if err != nil {
				log.Log.Error("when get nextStepType, get TaskTmpl byID occur error: %s", err.Error())
				return "", "", fmt.Errorf("网络异常，请重试")
			}
			stepType = stepModel.Type
			stepName = stepModel.Name
		}
	}
	if stepType == "" {
		log.Log.Error("when get nextStepType,nomatch step byIndexID: %v, stageID: %v", stepIndex, publishItem.StageID)
		return "", "", fmt.Errorf("网络异常，请重试")
	}
	return stepType, stepName, nil
}

func commandAndArgSplit(itemStr string) (itemArr []string) {
	itemStr = strings.TrimSpace(itemStr)
	if itemStr == "" {
		return
	}
	itemArr = strings.Split(itemStr, " ")
	return
}
