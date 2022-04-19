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
	"fmt"

	"github.com/go-atomci/atomci/constant"
	"github.com/go-atomci/atomci/internal/core/kuberes"
	"github.com/go-atomci/atomci/internal/core/project"
	"github.com/go-atomci/atomci/internal/middleware/log"
	"github.com/go-atomci/atomci/internal/models"
)

// ProjectController ...
type ProjectController struct {
	BaseController
}

// Create project
func (p *ProjectController) Create() {
	user := p.UserModel
	groupName := p.UserGroup()
	if groupName == "" {
		groupName = "system"
	}
	req := &project.ProjectReq{}
	p.DecodeJSONReq(req)
	pm := project.NewProjectManager()

	result, err := pm.CreateProject(user.User, groupName, req)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("Create Project error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, result, "")
	p.ServeJSON()
}

func (p *ProjectController) GetAppserviceList() {
	cluster := p.GetStringFromPath(":cluster")
	projectID, _ := p.GetInt64FromPath(":project_id")
	filterQuery := p.GetFilterQuery()

	// TODO: envID 0 change to real
	ar, err := kuberes.NewAppRes(cluster, 0, projectID)
	if err != nil {
		log.Log.Error(err.Error())
		p.HandleInternalServerError(err.Error())
		return
	}
	res, err := ar.GetAppListByPagination(filterQuery, projectID, cluster)
	if err != nil {
		log.Log.Error(err.Error())
		p.HandleInternalServerError(err.Error())
		return
	}
	p.Data["json"] = NewResult(true, res, "")
	p.ServeJSON()
}

func (p *ProjectController) AppInspect() {
	cluster := p.GetStringFromPath(":cluster")
	namespace := p.GetStringFromPath(":namespace")
	appname := p.GetStringFromPath(":app")

	// TODO: change envID, projectID to real
	ar, err := kuberes.NewAppRes(cluster, 0, 0)
	if err != nil {
		log.Log.Error(err.Error())
		p.HandleInternalServerError(err.Error())
		return
	}
	result, err := ar.GetAppDetail(namespace, appname)
	if err != nil {
		log.Log.Error("Get application information failed: " + err.Error())
		p.HandleInternalServerError(err.Error())
		return
	}
	p.Data["json"] = NewResult(true, result, "")
	p.ServeJSON()
}

func (p *ProjectController) AppDelete() {
	cluster := p.GetStringFromPath(":cluster")
	namespace := p.GetStringFromPath(":namespace")
	appname := p.GetStringFromPath(":app")

	ar, err := kuberes.NewAppRes(cluster, 0, 0)
	if err != nil {
		log.Log.Error("Delete application failed for: "+err.Error(), "cluster: "+cluster+",", "namespace: "+namespace+",", "name: "+appname, "!")
		p.HandleInternalServerError(err.Error())
		return
	}
	err = ar.DeleteApp(namespace, appname)
	if err != nil {
		log.Log.Error("Delete application failed for: "+err.Error(), "cluster: "+cluster+",", "namespace: "+namespace+",", "name: "+appname, "!")
		p.HandleInternalServerError(err.Error())
		return
	}

	log.Log.Info("Delete application successfully!", "cluster: "+cluster+",", "namespace: "+namespace+",", "name: "+appname, "!")
	p.Data["json"] = NewResult(true, nil, "")
	p.ServeJSON()
}

func (p *ProjectController) PodLog() {
	cluster := p.GetStringFromPath(":cluster")
	namespace := p.GetStringFromPath(":namespace")
	appname := p.GetStringFromPath(":app")

	podName := p.GetString("podName")
	containerName := p.GetString("containerName")
	// TODO: envID, projectID change to real
	ar, err := kuberes.NewAppRes(cluster, 0, 0)
	if err != nil {
		log.Log.Error(err.Error())
		p.HandleInternalServerError(err.Error())
		return
	}

	result, err := ar.GetAppPodLog(namespace, appname, podName, containerName)
	if err != nil {
		log.Log.Error(err.Error())
		p.HandleInternalServerError(err.Error())
		return
	}
	p.Data["json"] = NewResult(true, result, "")
	p.ServeJSON()
}

func (p *ProjectController) AppEvent() {
	cluster := p.GetStringFromPath(":cluster")
	namespace := p.GetStringFromPath(":namespace")
	appname := p.GetStringFromPath(":app")

	// TODO: change envID, projectID to real
	ar, err := kuberes.NewAppRes(cluster, 0, 0)
	if err != nil {
		log.Log.Error(err.Error())
		p.HandleInternalServerError(err.Error())
		return
	}

	result, err := ar.GetAppEvent(namespace, appname)
	if err != nil {
		log.Log.Error(err.Error())
		p.HandleInternalServerError(err.Error())
		return
	}
	p.Data["json"] = NewResult(true, result, "")
	p.ServeJSON()
}

func (p *ProjectController) AppRestart() {
	cluster := p.GetStringFromPath(":cluster")
	namespace := p.GetStringFromPath(":namespace")
	appname := p.GetStringFromPath(":app")
	ar, err := kuberes.NewAppRes(cluster, 0, 0)
	if err != nil {
		log.Log.Error("Restart application failed for: "+err.Error(), "cluster: "+cluster+",", "namespace: "+namespace+",", "name: "+appname, "!")
		p.HandleInternalServerError(err.Error())
		return
	}
	err = ar.Restart(namespace, appname)
	if err != nil {
		log.Log.Error("Restart application failed for: "+err.Error(), "cluster: "+cluster+",", "namespace: "+namespace+",", "name: "+appname, "!")
		p.HandleInternalServerError(err.Error())
		return
	}

	log.Log.Info("Restart application successfully!", "cluster: "+cluster+",", "namespace: "+namespace+",", "name: "+appname, "!")
	p.Data["json"] = NewResult(true, nil, "")
	p.ServeJSON()
}

func (p *ProjectController) AppScale() {
	cluster := p.GetStringFromPath(":cluster")
	appname := p.GetStringFromPath(":app")
	namespace := p.GetStringFromPath(":namespace")

	scale, err := p.GetInt("scaleBy")
	if err != nil {
		p.HandleInternalServerError(err.Error())
		return
	}

	if !(scale >= constant.ReplicasMin && scale <= constant.ReplicasMax) {
		err = fmt.Errorf(
			"replicas error: replicas must be an integer and in the range of %v to %v",
			constant.ReplicasMin, constant.ReplicasMax)
		log.Log.Error("error occur: %v", err.Error())
		p.HandleInternalServerError(err.Error())
		return
	}
	ar, err := kuberes.NewAppRes(cluster, 0, 0)
	if err != nil {
		log.Log.Error("scale application failed for: "+err.Error(), "cluster: "+cluster+",", "namespace: "+namespace+",", "name: "+appname, "!")
		p.HandleInternalServerError(err.Error())
		return
	}
	if err := ar.ScaleApp(namespace, appname, scale); err != nil {
		log.Log.Error("scale application failed for: "+err.Error(), "cluster: "+cluster+",", "namespace: "+namespace+",", "name: "+appname, "!")
		p.HandleInternalServerError(err.Error())
		return
	}
	log.Log.Info("scale application succefully,", "cluster: "+cluster+",", "namespace: "+namespace+",", "name: "+appname, "!")
	p.Data["json"] = NewResult(true, nil, "")
	p.ServeJSON()
}

// ProjectList ...
func (p *ProjectController) ProjectList() {
	filter := models.ProejctFilterQuery{}
	p.DecodeJSONReq(&filter)
	projectIDs, err := p.Projects()
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("Base on permission, filter project error: %s", err.Error())
		return
	}
	pm := project.NewProjectManager()
	rsp, err := pm.ProjectList(projectIDs, &filter)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("Get project list error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, rsp, "")
	p.ServeJSON()
}

// Update project base info
func (p *ProjectController) Update() {
	user := p.User
	req := &project.ProjectUpdateReq{}
	p.DecodeJSONReq(req)
	pm := project.NewProjectManager()
	projectID, _ := p.GetInt64FromPath(":project_id")
	err := pm.UpdateProject(user, projectID, req)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("Update Project error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, nil, "")
	p.ServeJSON()
}

// CheckProjetCreator check whether user is creator or have admin permissions
func (p *ProjectController) CheckProjetCreator() {
	flag := false
	// var err error
	if p.IsSysAdmin() {
		flag = true
	} else {
		groupFlag := p.IsGroupAdmin()
		if groupFlag == 1 {
			flag = true
		}
	}
	p.Data["json"] = NewResult(true, flag, "")
	p.ServeJSON()
}

// Delete project base project_id
func (p *ProjectController) Delete() {
	projectID, _ := p.GetInt64FromPath(":project_id")
	pm := project.NewProjectManager()
	groupAdminFlag := p.IsGroupAdmin()
	flag := false
	if groupAdminFlag == 1 {
		flag = true
	}
	if !flag {
		log.Log.Error("when check project user flag is %v", flag)
		p.HandleInternalServerError("仅允许项目owner及管理员更新基础信息")
		return
	}
	rsp := pm.DeleteProject(projectID)
	if rsp != nil {
		p.HandleInternalServerError(rsp.Error())
		return
	}
	p.Data["json"] = NewResult(true, rsp, "")
	p.ServeJSON()
}

// GetProject ...
func (p *ProjectController) GetProject() {
	projectID, _ := p.GetInt64FromPath(":project_id")
	pm := project.NewProjectManager()
	result, err := pm.GetProjectInfo(projectID)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("Get Project error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, result, "")
	p.ServeJSON()
}

/* -----  Project Setup  -----  */

// GetProjectMembers ..
func (p *ProjectController) GetProjectMembers() {
	projectID, _ := p.GetInt64FromPath(":project_id")
	pm := project.NewProjectManager()
	result, err := pm.GetProjectMembers(projectID)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("Get Project Pipeline error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, result, "")
	p.ServeJSON()
}

// AddProjectMember ..
func (p *ProjectController) AddProjectMember() {
	projectID, _ := p.GetInt64FromPath(":project_id")
	req := &project.ProjectNumberReq{}
	groupName := p.UserGroup()
	if groupName == "" {
		groupName = "system"
	}
	p.DecodeJSONReq(req)
	pm := project.NewProjectManager()
	if err := pm.AddProjectMembers(projectID, req, groupName); err != nil {
		p.HandleInternalServerError(err.Error())
		return
	}
	p.Data["json"] = NewResult(true, nil, "")
	p.ServeJSON()
}

// DeleteProjectMember ..
func (p *ProjectController) DeleteProjectMember() {
	groupName := p.UserGroup()
	if groupName == "" {
		groupName = "system"
	}
	numberID, _ := p.GetInt64FromPath(":id")
	pm := project.NewProjectManager()
	if err := pm.DeleteProjectMember(numberID, groupName); err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("Delete Project number error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, nil, "")
	p.ServeJSON()
}

/*    */

// GetProjectEnvs ..
func (p *ProjectController) GetProjectEnvs() {
	projectID, _ := p.GetInt64FromPath(":project_id")
	pm := project.NewProjectManager()
	rsp, err := pm.GetProjectEnvs(projectID)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("Get flow stages occur error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, rsp, "")
	p.ServeJSON()
}

// GetProjectEnvsByPagination ..
func (p *ProjectController) GetProjectEnvsByPagination() {
	projectID, _ := p.GetInt64FromPath(":project_id")
	filterQuery := p.GetFilterQuery()
	pm := project.NewProjectManager()
	rsp, err := pm.GetProjectEnvsByPagination(filterQuery, projectID)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("Get flow stages occur error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, rsp, "")
	p.ServeJSON()
}

// CreateProjectEnv ..
func (p *ProjectController) CreateProjectEnv() {
	projectID, _ := p.GetInt64FromPath(":project_id")
	request := project.ProjectEnvReq{}
	creator := p.User
	p.DecodeJSONReq(&request)
	pm := project.NewProjectManager()
	err := pm.CreateProjectEnv(&request, creator, projectID)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("Create project env occur error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, nil, "")
	p.ServeJSON()
}

// UpdateProjectEnv ..
func (p *ProjectController) UpdateProjectEnv() {
	stageID, _ := p.GetInt64FromPath(":env_id")
	request := project.ProjectEnvReq{}
	p.DecodeJSONReq(&request)
	pm := project.NewProjectManager()
	err := pm.UpdateProjectEnv(&request, stageID)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("update project env occur error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, nil, "")
	p.ServeJSON()
}

// DeleteProjectEnv ..
func (p *ProjectController) DeleteProjectEnv() {
	envID, _ := p.GetInt64FromPath(":env_id")
	pm := project.NewProjectManager()
	err := pm.DeleteProjectEnv(envID)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("delete project env occur error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, nil, "")
	p.ServeJSON()
}

// GetProjectPipelines ..
func (p *ProjectController) GetProjectPipelines() {
	projectID, _ := p.GetInt64FromPath(":project_id")
	pm := project.NewProjectManager()
	result, err := pm.GetProjectPipelines(projectID)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("Get Project Pipeline error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, result, "")
	p.ServeJSON()
}

// GetPipelinesByPagination ..
func (p *ProjectController) GetPipelinesByPagination() {
	projectID, _ := p.GetInt64FromPath(":project_id")
	filterQuery := p.GetFilterQuery()
	pm := project.NewProjectManager()
	rsp, err := pm.GetPipelinesByPagination(filterQuery, projectID)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("Get pipelines occur error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, rsp, "")
	p.ServeJSON()
}

// UpdatePipelineConfig ...
func (p *ProjectController) UpdatePipelineConfig() {
	projectID, _ := p.GetInt64FromPath(":project_id")
	pipelineID, _ := p.GetInt64FromPath(":id")
	request := project.PipelineReq{}
	currentUser := p.User
	p.DecodeJSONReq(&request)
	mgr := project.NewProjectManager()
	err := mgr.UpdateProjectPipelineConfig(&request, currentUser, projectID, pipelineID)
	if err != nil {
		p.ServeError(err)
		return
	}
	p.ServeResult(NewResult(true, nil, ""))
}

// CreatePipeline ..
func (p *ProjectController) CreatePipeline() {
	req := project.PipelineReq{}
	currentUser := p.User
	p.DecodeJSONReq(&req)
	pm := project.NewProjectManager()
	id, err := pm.CreateProjectPipeline(&req, currentUser)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		return
	}
	p.Data["json"] = NewResult(true, id, "")
	p.ServeJSON()
}

// GetProjectPipeline ..
func (p *ProjectController) GetProjectPipeline() {
	// projectID, _ := p.GetInt64FromPath(":project_id")
	pipelineID, _ := p.GetInt64FromPath(":id")
	mgr := project.NewProjectManager()
	setResult, err := mgr.GetPipelineConfig(pipelineID)
	if err != nil {
		p.ServeError(err)
		return
	}
	p.ServeResult(NewResult(true, setResult, ""))
}

// DeleteProjectPipeline ..
func (p *ProjectController) DeleteProjectPipeline() {
	pipelineBindID, _ := p.GetInt64FromPath(":id")
	pm := project.NewProjectManager()
	if err := pm.DeleteProjectPipeline(pipelineBindID); err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("Delete Project Pipeline error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, nil, "")
	p.ServeJSON()
}

/* project app controller part -- */

// CreateApp for project
func (p *ProjectController) CreateApp() {
	projectID, _ := p.GetInt64FromPath(":project_id")
	req := &project.ProjectAppReq{}
	p.DecodeJSONReq(&req)
	pm := project.NewProjectManager()
	result := pm.CreateProjectApp(projectID, req, p.User)
	if result != nil {
		p.HandleInternalServerError(result.Error())
		log.Log.Error("add project app error: %s", result.Error())
		return
	}
	p.Data["json"] = NewResult(true, result, "")
	p.ServeJSON()
}

// GetApps ..
func (p *ProjectController) GetApps() {
	projectID, _ := p.GetInt64FromPath(":project_id")
	pm := project.NewProjectManager()
	result, err := pm.GetProjectApps(projectID)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("get project app list error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, result, "")
	p.ServeJSON()
}

// GetAppsByPagination ..
func (p *ProjectController) GetAppsByPagination() {
	projectID, _ := p.GetInt64FromPath(":project_id")
	filterQuery := models.ProejctAppFilterQuery{}
	p.DecodeJSONReq(&filterQuery)
	pm := project.NewProjectManager()
	result, err := pm.GetProjectAppsByPagination(projectID, &filterQuery)
	if err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("get project app list error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, result, "")
	p.ServeJSON()
}

// DeleteProjectApp for project
func (p *ProjectController) DeleteProjectApp() {
	projectAppID, _ := p.GetInt64FromPath(":project_app_id")
	pm := project.NewProjectManager()
	result := pm.DeleteProjectApp(projectAppID)
	if result != nil {
		p.HandleInternalServerError(result.Error())
		log.Log.Error("delete project app error: %s", result.Error())
		return
	}
	p.Data["json"] = NewResult(true, result, "")
	p.ServeJSON()
}

// UpdateProjectApp ..
func (p *ProjectController) UpdateProjectApp() {
	projectID, _ := p.GetInt64FromPath(":project_id")
	projectAppID, _ := p.GetInt64FromPath(":project_app_id")
	req := &project.ProjectAppUpdateReq{}
	p.DecodeJSONReq(req)
	pm := project.NewProjectManager()
	if err := pm.UpdateProjectApp(projectID, projectAppID, req); err != nil {
		p.HandleInternalServerError(err.Error())
		log.Log.Error("update project app error: %s", err.Error())
		return
	}
	p.Data["json"] = NewResult(true, nil, "")
	p.ServeJSON()
}
