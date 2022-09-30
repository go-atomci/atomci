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

package api

import (
	"github.com/astaxie/beego"
	"github.com/go-atomci/atomci/internal/core/pipelinemgr"
	"github.com/go-atomci/atomci/internal/core/publish"
	"github.com/go-atomci/atomci/internal/middleware/log"
	"github.com/go-atomci/atomci/pkg/notification"
)

// PipelineController ...
type PipelineController struct {
	BaseController
}

// GetStepInfo ...
func (p *PipelineController) GetStepInfo() {
	projectID, _ := p.GetInt64FromPath(":project_id")
	publishID, _ := p.GetInt64FromPath(":publish_id")
	stageID, _ := p.GetInt64FromPath(":stage_id")
	stepName := p.GetStringFromPath(":step_name")

	pm := pipelinemgr.NewPipelineManager()
	rsp, err := pm.GetStepInfo(projectID, publishID, stageID, stepName)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("get pipeline step info: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, rsp, "")
	p.ServeJSON()
}

// RunStep ..
func (p *PipelineController) RunStep() {
	creator := p.User
	projectID, _ := p.GetInt64FromPath(":project_id")
	publishID, _ := p.GetInt64FromPath(":publish_id")
	stageID, _ := p.GetInt64FromPath(":stage_id")
	stepName := p.GetStringFromPath(":step_name")

	pm := pipelinemgr.NewPipelineManager()
	var err error
	var publishStatus, runID int64
	var message, jobName string
	switch stepName {
	case "manual":
		request := &pipelinemgr.ManualStepReq{}
		p.DecodeJSONReq(&request)
		message = request.Message
		publishStatus, err = pm.RunManualStep(publishID, stageID, request)
	case "build":
		request := &pipelinemgr.BuildStepReq{}
		p.DecodeJSONReq(&request)
		publishStatus, runID, jobName, err = pm.RunBuildStep(projectID, publishID, stageID, creator, stepName, request)
	case "deploy":
		request := &pipelinemgr.DeployStepReq{}
		p.DecodeJSONReq(&request)
		publishStatus, runID, jobName, err = pm.RunDeployStep(projectID, publishID, stageID, creator, stepName, request)
	default:
		log.Log.Error("unknow step_name: %s", stepName)
	}
	publishmgr := publish.NewPublishManager()
	updateErr := publishmgr.UpdatePublish(publishID, stageID, publishStatus, runID, creator, message, jobName)
	if err != nil || updateErr != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("Run Publish error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, nil, "")
	p.ServeJSON()
}

// RunStepCallback ..
func (p *PipelineController) RunStepCallback() {
	creator := p.User
	publishID, _ := p.GetInt64FromPath(":publish_id")
	stageID, _ := p.GetInt64FromPath(":stage_id")
	stepName := p.GetStringFromPath(":step_name")

	pm := pipelinemgr.NewPipelineManager()
	var err error
	var publishStatus int64
	var message string
	switch stepName {
	case "build", "deploy":
		request := &pipelinemgr.BuildStepCallbackReq{}
		p.DecodeJSONReq(&request)
		publishStatus, err = pm.RunBuildDeployCallBackStep(request)
	default:
		log.Log.Error("callback occur erro: unknow step_name: %s", stepName)
	}
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("RunStep callback error: %s", err.Error())
		return
	}
	var runID int64
	var jobName string
	publishmgr := publish.NewPublishManager()
	updateErr := publishmgr.UpdatePublish(publishID, stageID, publishStatus, runID, creator, message, jobName)
	if updateErr != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("RunStep callback, update publish Order occur error: %s", err.Error())
		return
	}
	publishInfo, _ := publishmgr.GetPublishInfo(publishID)

	dingEnable := beego.AppConfig.DefaultBool("notification::dingEnable", false)
	mailEnable := beego.AppConfig.DefaultBool("notification::mailEnable", false)
	dingURL := beego.AppConfig.String("notification::ding")

	smtpHost := beego.AppConfig.String("notification::smtpHost")
	smtpAccount := beego.AppConfig.String("notification::smtpAccount")
	smtpPassword := beego.AppConfig.String("notification::smtpPassword")
	smtpPort, _ := beego.AppConfig.Int("notification::smtpPort")

	pushOptions := notification.PushNotification{
		// message
		Status:      publishStatus,
		PublishName: publishInfo.Name,
		StageName:   publishInfo.StageName,
		StepName:    publishInfo.Step,
		// dingding
		DingURL:    dingURL,
		DingEnable: dingEnable,
		// email
		EmailEnable:   mailEnable,
		EmailHost:     smtpHost,
		EmailPort:     smtpPort,
		EmailUser:     smtpAccount,
		EmailPassword: smtpPassword,
	}
	// publishID
	go notification.Send(pushOptions)

	p.Data["json"] = NewResult(true, nil, "")
	p.ServeJSON()
}

/*  -----  For frontend   ----------   */

// GetPublishStats ..
func (p *PipelineController) GetPublishStats() {
	projectID, _ := p.GetInt64FromPath(":project_id")
	request := &pipelinemgr.PublishStatsReq{}
	p.DecodeJSONReq(&request)
	log.Log.Debug("publish job stats params: %+v", request)
	pm := pipelinemgr.NewPipelineManager()
	rsp, err := pm.GetPublishStats(projectID, request)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("Get publish stats occur error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, rsp, "")
	p.ServeJSON()
}

// GetJenkinsConfig ..
func (p *PipelineController) GetJenkinsConfig() {
	stageID, _ := p.GetInt64FromPath(":stage_id")
	pm := pipelinemgr.NewPipelineManager()
	rsp, err := pm.GetJenkinsConfig(stageID)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("Get jenkins config occur error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, rsp, "")
	p.ServeJSON()
}

// GetFlowComponents ..
func (p *PipelineController) GetFlowComponents() {
	pm := pipelinemgr.NewPipelineManager()
	rsp, err := pm.GetFlowComponents()
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("Get flow components occur error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, rsp, "")
	p.ServeJSON()
}

// GetTaskTmpls ..
func (p *PipelineController) GetTaskTmpls() {
	pm := pipelinemgr.NewPipelineManager()
	rsp, err := pm.GetTaskTmpls()
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("get all task templates occur error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, rsp, "")
	p.ServeJSON()
}

// GetTaskTmplsByPagination ..
func (p *PipelineController) GetTaskTmplsByPagination() {
	filterQuery := p.GetFilterQuery()
	pm := pipelinemgr.NewPipelineManager()
	rsp, err := pm.GetTaskTmplsByPagination(filterQuery)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("get all task templates occur error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, rsp, "")
	p.ServeJSON()
}

// UpdateTaskTmpl ..
func (p *PipelineController) UpdateTaskTmpl() {
	stepID, _ := p.GetInt64FromPath(":step_id")
	request := pipelinemgr.TaskTmplReq{}
	p.DecodeJSONReq(&request)
	pm := pipelinemgr.NewPipelineManager()
	err := pm.UpdateTaskTmpl(&request, stepID)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("update flow step occur error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, nil, "")
	p.ServeJSON()
}

// CreateTaskTmpl ..
func (p *PipelineController) CreateTaskTmpl() {
	request := pipelinemgr.TaskTmplReq{}
	creator := p.User
	p.DecodeJSONReq(&request)
	pm := pipelinemgr.NewPipelineManager()
	err := pm.CreateTaskTmpl(&request, creator)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("create flow step occur error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, nil, "")
	p.ServeJSON()
}

// DeleteTaskTmpl ..
func (p *PipelineController) DeleteTaskTmpl() {
	stepID, _ := p.GetInt64FromPath(":step_id")
	pm := pipelinemgr.NewPipelineManager()
	err := pm.DeleteTaskTmpl(stepID)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("delete flow step occur error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, nil, "")
	p.ServeJSON()
}
