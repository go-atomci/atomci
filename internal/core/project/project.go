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

package project

import (
	"fmt"

	"strconv"
	"time"

	"github.com/go-atomci/atomci/constant"
	"github.com/go-atomci/atomci/internal/dao"
	"github.com/go-atomci/atomci/internal/middleware/log"
	"github.com/go-atomci/atomci/internal/models"
	"github.com/go-atomci/atomci/utils"
	"github.com/go-atomci/atomci/utils/errors"
	"github.com/go-atomci/atomci/utils/query"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// ProjectManager ...
type ProjectManager struct {
	model          *dao.ProjectModel
	pipelineModel  *dao.PipelineStageModel
	scmAppModel    *dao.ScmAppModel
	k8sModel       *dao.K8sClusterModel
	userrolesModel *dao.UserRolesModel
	publishModel   *dao.PublishModel
	settingModel   *dao.SysSettingModel
}

// NewProjectManager ...
func NewProjectManager() *ProjectManager {
	return &ProjectManager{
		model:          dao.NewProjectModel(),
		pipelineModel:  dao.NewPipelineStageModel(),
		settingModel:   dao.NewSysSettingModel(),
		scmAppModel:    dao.NewScmAppModel(),
		k8sModel:       dao.NewK8sClusterModel(),
		userrolesModel: dao.NewUserRolesModel(),
		publishModel:   dao.NewPublishModel(),
	}
}

// CreateProject ...
func (pm *ProjectManager) CreateProject(user, groupName string, p *ProjectReq) (*models.ProjectResponse, error) {

	projectModel := models.Project{
		Addons:      models.NewAddons(),
		Name:        p.Name,
		Description: p.Description,
		Owner:       user,
		Creator:     user,
		Status:      models.ProjectRuning,
	}
	projectID, err := pm.model.CreateProjectifNotExist(&projectModel)
	if err != nil {
		return nil, err
	}

	// Delete project label
	deleteProjectLabel := false
	// add owner to project numbers as developer role
	modelRole, err := dao.GetGroupRoleByName(groupName, constant.SystemMemberRole)
	if err != nil {
		log.Log.Error("when crate project, get role by group,name occur error: %s", err.Error())
		deleteProjectLabel = true

	}
	if !deleteProjectLabel {
		number := &ProjectNumberReq{
			User:   user,
			RoleID: modelRole.ID,
		}
		if err := pm.AddProjectMembers(projectID, number, groupName); err != nil {
			log.Log.Error("after create project, add projectNumber occur error: %s", err.Error())
			deleteProjectLabel = true
		}
	}
	if deleteProjectLabel {
		if err := pm.model.DeleteProject(projectID); err != nil {
			log.Log.Error("add project Number failed, delete project occur error: %s", err.Error())
			return nil, fmt.Errorf("网络异常，请稍后重试")
		}
	}
	projectResp := pm.GetProjectResp(projectID)
	return projectResp, nil
}

// GetProjectResp ..
func (pm *ProjectManager) GetProjectResp(projectID int64) *models.ProjectResponse {
	project, err := pm.model.GetProjectByID(projectID)
	if err != nil {
		log.Log.Error("when get project resp, get project by id occur error: %s", err.Error())
		return nil
	}
	users, err := pm.model.GetProjectUsers(projectID)
	var numbers int
	var membersName []string
	if err != nil {
		log.Log.Error("when GetProjectUsers occur error: %s", err.Error())
	} else {
		numbers = len(users)

		for _, user := range users {
			membersName = append(membersName, user.User)
		}
	}

	projectResp := &models.ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		CreateAt:    project.CreateAt,
		UpdateAt:    project.UpdateAt,
		StartAt:     project.CreateAt,
		EndAt:       project.EndAt,
		Status:      project.Status,
		Creator:     project.Creator,
		Owner:       project.Owner,
		Members:     numbers,
		MembersName: membersName,
	}
	return projectResp
}

// UpdateProject ...
func (pm *ProjectManager) UpdateProject(user string, projectID int64, p *ProjectUpdateReq) error {
	modelProject, err := pm.model.GetProjectByID(projectID)
	if err != nil {
		return err
	}
	timeNow, _ := time.Parse("2006-01-02 15:04:05", time.Now().Local().Format("2006-01-02 15:04:05"))
	if p.Status == models.ProjectRuning {
		modelProject.Status = p.Status
		modelProject.StartAt = timeNow
		return pm.model.UpdateProject(modelProject)
	}

	if p.Status == models.ProjectEnd {
		modelProject.Status = p.Status
		modelProject.EndAt = &timeNow
		// projectEnd verify
		if canEnd, err := pm.projectEndVerify(projectID); canEnd != true {
			return err
		}
		return pm.model.UpdateProject(modelProject)
	}

	UpdateConstraint := false
	originUser := ""
	if p.Owner != "" {
		if p.Owner != modelProject.Owner {
			UpdateConstraint = true
			originUser = modelProject.Owner
		}
		modelProject.Owner = p.Owner
	}
	if p.Name != "" {
		if item, err := pm.model.GetProjectByProjectName(p.Name); err == nil {
			if item.ID != projectID {
				return fmt.Errorf("项目名称不允许重复，请你确认后重试")
			}
		}
		modelProject.Name = p.Name
	}
	modelProject.Description = p.Description
	// if p.Owner changed, update project constraint
	if UpdateConstraint {
		// TODO: add project constraint for owner, group use taimei
		if err := dao.AddGroupUserConstraintValues("system", p.Owner, "projectID", []string{strconv.Itoa(int(projectID))}); err != nil {
			return err
		}
		if err := dao.DeleteGroupUserConstraintValues("system", originUser, "projectID", []string{strconv.Itoa(int(projectID))}); err != nil {
			return err
		}
	}
	return pm.model.UpdateProject(modelProject)
}

// ProjectList ...
func (pm *ProjectManager) ProjectList(projectID []int64, filter *models.ProejctFilterQuery) (*query.QueryResult, error) {
	log.Log.Debug("params: %+v", filter)
	projects, modelDatas, err := pm.model.ProjectListByIDs(projectID, filter)
	if projects == nil {
		return nil, nil
	}
	items := []*models.ProjectResponse{}
	for _, item := range modelDatas {
		rsp := pm.GetProjectResp(item.ID)
		items = append(items, rsp)
	}
	projects.Item = items
	return projects, err
}

// GetProjectInfo ...
func (pm *ProjectManager) GetProjectInfo(projectID int64) (*models.ProjectDetailResponse, error) {
	projectBaseInfo := pm.GetProjectResp(projectID)
	rsp := &models.ProjectDetailResponse{
		ProjectResponse: projectBaseInfo,
	}
	if codeRepos, err := pm.model.GetProjectAppCounts(projectID); err == nil {
		rsp.CodeRepos = codeRepos
	} else {
		logs.Warn("when GetProjectInfo, GetProjectAppCounts occur error: %s", err.Error())
	}
	// releases: return format
	// [{"count": 1, "env": "开发环境"}]
	if releases, err := pm.publishModel.GetPublishReleasesByProjectID(projectID); err == nil {
		rsp.Releases = releases
	} else {
		logs.Warn("when GetProjectInfo, GetPublishByProjectID occur error: %s", err.Error())
	}

	return rsp, nil
}

// CheckProjectUser ..
func (pm *ProjectManager) CheckProjectUser(projectID int64, user string, groupsAdminList []string) (bool, error) {
	project, err := pm.model.GetProjectByID(projectID)
	if err != nil {
		err = fmt.Errorf("when check project user, get project occur error: %s", err)
		log.Log.Error("%v", err)
		return false, err
	}
	if project.Owner == user {
		return true, nil
	}
	// TODO: deprecated delete project Department
	return utils.Contains(groupsAdminList, "taimei"), nil
}

// DeleteProject ...
func (pm *ProjectManager) DeleteProject(projectID int64) error {
	project, err := pm.model.GetProjectByID(projectID)
	if err != nil {
		return err
	}
	if project.Status != models.ProjectEnd {
		return fmt.Errorf("你需要先结束项目，然后再尝试执行删除操作")
	}
	if err := pm.model.DeleteProject(projectID); err != nil {
		if err != orm.ErrNoRows {
			return errors.NewInternalServerError().SetCause(err)
		}
		return errors.NewNotFound().SetCause(err)
	}
	return nil
}

// GetProjectMembers ..
func (pm *ProjectManager) GetProjectMembers(projectID int64) ([]*ProjectNumberRsp, error) {
	modelUsers, err := pm.model.GetProjectUsers(projectID)
	if err != nil {
		return nil, fmt.Errorf("when get project numbers, occur errro: %v", err.Error())
	}
	rsp := []*ProjectNumberRsp{}
	for _, u := range modelUsers {
		role, err := pm.userrolesModel.GetRoleByID(u.RoleID)
		if err != nil {
			log.Log.Error("when get project nubmers, GetRoleByID id: %v occur error: %v", u.RoleID, err.Error())
			continue
		}
		itemRsp := &ProjectNumberRsp{
			ProjectUser: u,
			Role:        role.Role,
		}
		rsp = append(rsp, itemRsp)
	}
	return rsp, nil
}

// AddProjectMembers ..
func (pm *ProjectManager) AddProjectMembers(projectID int64, request *ProjectNumberReq, groupName string) error {
	log.Log.Debug("request params: %+v", request)
	if _, err := pm.userrolesModel.GetRoleByID(request.RoleID); err != nil {
		log.Log.Error("when add project number, get role by id occur error: %v", err.Error())
		return fmt.Errorf("请选择有效的角色后重试")
	}
	pp := &models.ProjectUser{
		ProjectID: projectID,
		User:      request.User,
		RoleID:    request.RoleID,
	}
	if _, err := pm.model.CreateProjectUserIfNotExist(pp); err != nil {
		log.Log.Error("when add project number, create project user item occur error: %v", err.Error())
		return fmt.Errorf("添加项目成员失败，请重试")
	}
	//
	if err := dao.AddGroupUserConstraintValues(groupName, request.User, "project_id", []string{strconv.Itoa(int(projectID))}); err != nil {
		return err
	}
	return nil
}

// DeleteProjectMember ..
func (pm *ProjectManager) DeleteProjectMember(numberID int64, groupName string) error {
	item, err := pm.model.GetProjectUserByID(numberID)
	if err != nil {
		log.Log.Error("when delete project number, GetProjectUserByID occur error: %v", err.Error())
		return fmt.Errorf("删除失败: %s", err.Error())
	}

	item.MarkDeleted()
	// TODO: group use defaut: company
	if err := dao.DeleteGroupUserConstraintValues(groupName, item.User, "project_id", []string{strconv.Itoa(int(item.ProjectID))}); err != nil {
		return err
	}
	return pm.model.UpdateProjectUser(item)
}

/*  --- Project Pipeline --  start ------ */

// GetProjectPipelines ..
func (pm *ProjectManager) GetProjectPipelines(projectID int64) ([]*models.ProjectPipeline, error) {
	return pm.model.GetProjectPipelines(projectID)
}

// DeleteProjectPipeline ..
func (pm *ProjectManager) DeleteProjectPipeline(pipelineBindID int64) error {
	// TODO: verify pipeline 看是否有绑定此流程的项目且在【待执行】和【进行中】的，
	// 即非【已删除】和【已关闭】和【已结束】，则提示“该流程在试用中，无法删除”
	item, err := pm.model.GetProjectPipelineByID(pipelineBindID)
	if err != nil {
		log.Log.Error("when delete project pipeline, get Bindpipeline by ID occur error: %v", err.Error())
		return fmt.Errorf("删除失败，请重试")
	}

	item.MarkDeleted()
	return pm.model.UpdateProjectPipeline(item)
}
