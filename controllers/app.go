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

package controllers

import (
	"github.com/go-atomci/atomci/core/apps"
	"github.com/go-atomci/atomci/core/kuberes"
	"github.com/go-atomci/atomci/middleware/log"
	"github.com/go-atomci/atomci/utils/errors"
)

// AppController ...
type AppController struct {
	BaseController
}

// GetArrange ...
func (a *AppController) GetArrange() {
	appID, err := a.GetInt64FromPath(":app_id")
	if err != nil {
		a.ServeError(errors.NewBadRequest().SetMessage("invalid app id"))
		return
	}
	envID, err := a.GetInt64FromPath(":env_id")
	if err != nil {
		a.ServeError(errors.NewBadRequest().SetMessage("invalid env id"))
		return
	}
	mgr := apps.NewAppManager()
	arrange, err := mgr.GetArrange(appID, envID)
	if err != nil {
		a.ServeError(err)
		return
	}
	a.ServeResult(NewResult(true, arrange, ""))
}

// SetArrange ...
func (a *AppController) SetArrange() {
	projectAppID, err := a.GetInt64FromPath(":app_id")
	if err != nil {
		a.ServeError(errors.NewBadRequest().SetMessage("invalid project app id"))
		return
	}
	arrangeEnvID, err := a.GetInt64FromPath(":env_id")
	if err != nil {
		a.ServeError(errors.NewBadRequest().SetMessage("invalid project env id"))
		return
	}
	request := apps.AppArrangeReq{}
	a.DecodeJSONReq(&request)

	native := &kuberes.NativeTemplate{
		Template: request.Config,
	}

	if err := native.Validate(); err != nil {
		a.ServeError(errors.NewBadRequest().SetMessage("yaml parse error: %s", err.Error()))
		return
	}

	mgr := apps.NewAppManager()
	err = mgr.SetArrange(projectAppID, arrangeEnvID, &request)
	if err != nil {
		a.ServeError(err)
		return
	}
	a.ServeResult(NewResult(true, nil, ""))
}

func (a *AppController) ParseArrangeYaml() {
	request := apps.AppArrangConfig{}
	a.DecodeJSONReq(&request)

	native := &kuberes.NativeTemplate{
		Template: request.Config,
	}
	rsp, err := native.GetContainerImages()
	if err != nil {
		log.Log.Debug("get container images error: %s")
	}
	a.ServeResult(NewResult(true, rsp, ""))
}

/* -- repo server start -- */

// GetRepos ..
func (a *AppController) GetRepos() {
	projectID, err := a.GetInt64FromQuery("project_id")
	if err != nil {
		a.HandleInternalServerError(err.Error())
		log.Log.Error("parse project id error: %s", err.Error())
		return
	}
	if projectID == 0 {
		projectID = 1
	}
	log.Log.Debug("args projectID: %v", projectID)
	mgr := apps.NewAppManager()
	rsp, err := mgr.GetRepos(projectID)
	if err != nil {
		a.HandleInternalServerError(err.Error())
		log.Log.Error("get repos error: %s", err.Error())
		return
	}
	a.ServeResult(NewResult(true, rsp, ""))
}

// GetGitProjectsByRepoID ..
func (a *AppController) GetGitProjectsByRepoID() {
	// TODO: change url query
	projectID, _ := a.GetInt64FromQuery("project_id")
	if projectID == 0 {
		projectID = 1
	}
	log.Log.Debug("args projectID: %v", projectID)
	repoID, _ := a.GetInt64FromPath(":repo_id")
	request := apps.SetupRepo{}
	a.DecodeJSONReq(&request)
	mgr := apps.NewAppManager()
	rsp, err := mgr.SetRepoAndGetProjects(projectID, repoID, &request)
	if err != nil {
		a.HandleInternalServerError(err.Error())
		log.Log.Error("get repo's projects error: %s", err.Error())
		return
	}
	a.ServeResult(NewResult(true, rsp, ""))
}

/* -- repo server end -- */

// GetAppBranches ..
func (a *AppController) GetAppBranches() {
	AppID, err := a.GetInt64FromPath(":app_id")
	if err != nil {
		a.ServeError(errors.NewBadRequest().SetMessage("invalid app id"))
		return
	}
	filterQuery := a.GetFilterQuery()
	mgr := apps.NewAppManager()
	rsp, err := mgr.AppBranches(AppID, filterQuery)
	if err != nil {
		a.HandleInternalServerError(err.Error())
		log.Log.Error("Get app list error: %s", err.Error())
		return
	}
	a.Data["json"] = NewResult(true, rsp, "")
	a.ServeJSON()
}

// SyncAppBranches ..
func (a *AppController) SyncAppBranches() {
	AppID, err := a.GetInt64FromPath(":app_id")
	if err != nil {
		a.ServeError(errors.NewBadRequest().SetMessage("invalid app id"))
		return
	}
	mgr := apps.NewAppManager()
	if err := mgr.SyncAppBranches(AppID); err != nil {
		a.HandleInternalServerError(err.Error())
		log.Log.Error("sync app branches error: %s", err.Error())
		return
	}
	a.Data["json"] = NewResult(true, nil, "")
	a.ServeJSON()
}
