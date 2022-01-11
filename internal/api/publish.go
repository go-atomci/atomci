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
	"github.com/go-atomci/atomci/internal/core/publish"
	"github.com/go-atomci/atomci/internal/middleware/log"
	"github.com/go-atomci/atomci/internal/models"
)

// PublishController ...
type PublishController struct {
	BaseController
}

// Create publish
func (p *PublishController) Create() {
	user := p.User
	projectID, _ := p.GetInt64FromPath(":project_id")
	req := &publish.PublishReq{}
	p.DecodeJSONReq(req)
	pm := publish.NewPublishManager()
	err := pm.CreatePublish(user, projectID, req)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("Create Publish error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, nil, "")
	p.ServeJSON()
}

// PublishList ...
func (p *PublishController) PublishList() {
	projectID, _ := p.GetInt64FromPath(":project_id")
	filterQuery := models.ProejctReleaseFilterQuery{}
	p.DecodeJSONReq(&filterQuery)
	pm := publish.NewPublishManager()
	rsp, err := pm.PublishList(projectID, &filterQuery)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("Get publish list error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, rsp, "")
	p.ServeJSON()
}

// Delete publish base publish_id
func (p *PublishController) Delete() {
	publishID, _ := p.GetInt64FromPath(":publish_id")
	pm := publish.NewPublishManager()
	// TODO: permission control added later
	rsp := pm.DeletePublish(publishID)
	p.Data["json"] = NewResult(true, rsp, "")
	p.ServeJSON()
}

// GetPublish ...
func (p *PublishController) GetPublish() {
	publishID, _ := p.GetInt64FromPath(":publish_id")
	pm := publish.NewPublishManager()
	result, err := pm.GetPublishInfo(publishID)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("Update Publish error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, result, "")
	p.ServeJSON()
}

// ClosePublish ..
func (p *PublishController) ClosePublish() {
	pm := publish.NewPublishManager()
	publishID, _ := p.GetInt64FromPath(":publish_id")
	req := &publish.PublishUpdate{}
	p.DecodeJSONReq(req)

	var err error
	if req.Name != "" && req.VersionNo != "" {
		err = pm.UpdatePublishBaseInfo(publishID, req)
	} else {
		err = pm.ClosePublish(publishID)
	}
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("update Publish id: %v error: %s", publishID, err.Error())
		return
	}
	p.Data["json"] = NewResult(true, nil, "")
	p.ServeJSON()
}

// DeletePublish ..
func (p *PublishController) DeletePublish() {
	pm := publish.NewPublishManager()
	publishID, _ := p.GetInt64FromPath(":publish_id")
	err := pm.DeletePublish(publishID)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("Delete Publish error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, nil, "")
	p.ServeJSON()
}

// CanAddedApps ..
func (p *PublishController) CanAddedApps() {
	pm := publish.NewPublishManager()
	publishID, _ := p.GetInt64FromPath(":publish_id")
	rsp, err := pm.GetCanAddedApps(publishID)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("get can added app occur error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, rsp, "")
	p.ServeJSON()
}

// AddPublishApp ..
func (p *PublishController) AddPublishApp() {
	req := &publish.PublishAddApps{}
	p.DecodeJSONReq(req)
	publishID, _ := p.GetInt64FromPath(":publish_id")
	pm := publish.NewPublishManager()
	err := pm.AddPublishApps(publishID, req)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("add publish App occur error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, nil, "")
	p.ServeJSON()
}

// DeletePublishApp ..
func (p *PublishController) DeletePublishApp() {
	publishAppID, _ := p.GetInt64FromPath(":publish_app_id")
	pm := publish.NewPublishManager()
	err := pm.DeletePublishApp(publishAppID)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("add publish App occur error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, nil, "")
	p.ServeJSON()
}

// GetBackTo ..
func (p *PublishController) GetBackTo() {
	projectID, _ := p.GetInt64FromPath(":project_id")
	publishID, _ := p.GetInt64FromPath(":publish_id")
	stageID, _ := p.GetInt64FromPath(":stage_id")
	pm := publish.NewPublishManager()
	result, err := pm.GetBackTo(projectID, publishID, stageID)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("Get Publish BackTo list error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, result, "")
	p.ServeJSON()
}

// TriggerBackTo ..
func (p *PublishController) TriggerBackTo() {
	projectID, _ := p.GetInt64FromPath(":project_id")
	publishID, _ := p.GetInt64FromPath(":publish_id")
	stageID, _ := p.GetInt64FromPath(":stage_id")
	currentUser := p.User
	req := &publish.TriggerBackToReq{}
	p.DecodeJSONReq(req)
	pm := publish.NewPublishManager()
	err := pm.TriggerBackTo(projectID, publishID, stageID, req, currentUser)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("Get Publish BackTo list error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, nil, "")
	p.ServeJSON()
}

// GetNextStage ..
func (p *PublishController) GetNextStage() {
	projectID, _ := p.GetInt64FromPath(":project_id")
	publishID, _ := p.GetInt64FromPath(":publish_id")
	envID, _ := p.GetInt64FromPath(":stage_id")
	pm := publish.NewPublishManager()
	result, err := pm.GetNextStage(projectID, publishID, envID)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("Get Publish NextStage list error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, result, "")
	p.ServeJSON()
}

// TriggerNextStage ..
func (p *PublishController) TriggerNextStage() {
	projectID, _ := p.GetInt64FromPath(":project_id")
	publishID, _ := p.GetInt64FromPath(":publish_id")
	stageID, _ := p.GetInt64FromPath(":stage_id")
	currentUser := p.User
	req := &publish.TriggerBackToReq{}
	p.DecodeJSONReq(req)
	pm := publish.NewPublishManager()
	err := pm.TriggerNextStage(projectID, publishID, stageID, req, currentUser)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("Get Publish NextStage list error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, nil, "")
	p.ServeJSON()
}

// GetOpertaionLogByPagination ..
func (p *PublishController) GetOpertaionLogByPagination() {
	publishID, _ := p.GetInt64FromPath(":publish_id")
	filterQuery := p.GetFilterQuery()
	pm := publish.NewPublishManager()
	result, err := pm.GetPublishOperationLog(publishID, filterQuery)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("Get Publish Operation Log error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, result, "")
	p.ServeJSON()
}
