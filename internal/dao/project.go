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

package dao

import (
	"fmt"
	"time"

	"github.com/go-atomci/atomci/internal/models"
	"github.com/go-atomci/atomci/utils/query"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

var projectEnableFilterKeys = []string{
	"name",
	"description",
	"owner",
}

// ProjectModel ...
type ProjectModel struct {
	ormer                    orm.Ormer
	projectTableName         string
	projectEnvTableName      string
	projectPipelineTableName string
	projectUserTableName     string
	projectAppTableName      string
}

// NewProjectModel ...
func NewProjectModel() (model *ProjectModel) {
	return &ProjectModel{
		ormer:                    GetOrmer(),
		projectTableName:         (&models.Project{}).TableName(),
		projectEnvTableName:      (&models.ProjectEnv{}).TableName(),
		projectPipelineTableName: (&models.ProjectPipeline{}).TableName(),
		projectUserTableName:     (&models.ProjectUser{}).TableName(),
		projectAppTableName:      (&models.ProjectApp{}).TableName(),
	}
}

// ProjectListByIDs ...
func (model *ProjectModel) ProjectListByIDs(projectID []int64, filter *models.ProejctFilterQuery) (*query.QueryResult, []*models.Project, error) {
	rst := &query.QueryResult{Item: []*models.ProjectResponse{}}
	cond := orm.NewCondition()
	queryCond := cond.AndCond(cond.And("deleted", false))
	if projectID == nil {
		return nil, nil, nil
	}
	queryCond = queryCond.AndCond(cond.And("id__in", projectID))

	if filter.Name != "" {
		queryCond = queryCond.AndCond(cond.Or("name__icontains", filter.Name).Or("owner__icontains", filter.Name))
	}

	if filter.CreateAtStart != "" && filter.CreateAtEnd != "" {
		if createAtStart, err := time.Parse("2006-01-02", filter.CreateAtStart); err == nil {
			queryCond = queryCond.AndCond(orm.NewCondition().And("create_at__gte", createAtStart))
		} else {
			logs.Error("time parse error: %s", err.Error())
		}
		if createAtEnd, err := time.Parse("2006-01-02", filter.CreateAtEnd); err == nil {
			queryCond = queryCond.AndCond(orm.NewCondition().And("create_at__lte", createAtEnd))
		} else {
			logs.Error("time parse error: %s", err.Error())
		}

	}

	if filter.Status != nil {
		queryCond = queryCond.AndCond(orm.NewCondition().And("status", filter.Status))
	}
	qs := model.ormer.QueryTable(model.projectTableName).OrderBy("-create_at").SetCond(queryCond)
	count, err := qs.Count()
	if err != nil {
		return nil, nil, err
	}
	if err = query.FillPageInfo(rst, filter.PageIndex, filter.PageSize, int(count)); err != nil {
		return nil, nil, err
	}

	projectList := []*models.Project{}
	_, err = qs.Limit(filter.PageSize, filter.PageSize*(filter.PageIndex-1)).All(&projectList)
	if err != nil {
		return nil, nil, err
	}

	return rst, projectList, nil
}

// GetProjectByID ...
func (model *ProjectModel) GetProjectByID(projectID int64) (*models.Project, error) {
	project := models.Project{}
	qs := model.ormer.QueryTable(model.projectTableName).Filter("deleted", false)
	if projectID != -1 {
		qs = qs.Filter("id", projectID)
	}
	err := qs.One(&project)
	return &project, err
}

// GetProjects ...
func (model *ProjectModel) GetProjects() ([]*models.Project, error) {
	projects := []*models.Project{}
	qs := model.ormer.QueryTable(model.projectTableName).
		Filter("deleted", false)

	_, err := qs.All(&projects)
	return projects, err
}

// GetProjectByProjectName ...
func (model *ProjectModel) GetProjectByProjectName(name string) (*models.Project, error) {
	project := models.Project{}
	qs := model.ormer.QueryTable(model.projectTableName).Filter("deleted", false).
		Filter("name", name)

	err := qs.One(&project)
	return &project, err
}

// CreateProjectifNotExist ...
func (model *ProjectModel) CreateProjectifNotExist(project *models.Project) (int64, error) {
	created, id, err := model.ormer.ReadOrCreate(project, "name", "deleted")
	if err == nil {
		if !created {
			err = fmt.Errorf(fmt.Sprintf("project name existed:%s", project.Name))
		}
	}
	return id, err
}

func (model *ProjectModel) GetProjectEnvByID(stageID int64) (*models.ProjectEnv, error) {
	stage := models.ProjectEnv{}
	qs := model.ormer.QueryTable(model.projectEnvTableName).Filter("deleted", false)
	if err := qs.Filter("id", stageID).One(&stage); err != nil {
		return nil, err
	}
	return &stage, nil
}

// GetProjectEnvs ...
func (model *ProjectModel) GetProjectEnvs(projectID int64) ([]*models.ProjectEnv, error) {
	// TODO: verify proejctID is validate
	stages := []*models.ProjectEnv{}
	qs := model.ormer.QueryTable(model.projectEnvTableName).Filter("deleted", false).
		Filter("project_id", projectID)
	_, err := qs.All(&stages)
	if err != nil {
		return nil, err
	}
	return stages, err
}

// GetProjectEnvsByPagination ..
func (model *ProjectModel) GetProjectEnvsByPagination(filter *query.FilterQuery, projectID int64) (*query.QueryResult, error) {
	if projectID == 0 {
		return nil, fmt.Errorf("project_id %v is invalidate", projectID)
	}
	rst := &query.QueryResult{Item: []*models.ProjectEnv{}}
	queryCond := orm.NewCondition().AndCond(orm.NewCondition().And("deleted", false).And("project_id", projectID))

	if filterCond := query.FilterCondition(filter, filter.FilterKey); filterCond != nil {
		queryCond = queryCond.AndCond(filterCond)
	}
	qs := model.ormer.QueryTable(model.projectEnvTableName).OrderBy("-create_at").SetCond(queryCond)
	count, err := qs.Count()
	if err != nil {
		return nil, err
	}
	if err = query.FillPageInfo(rst, filter.PageIndex, filter.PageSize, int(count)); err != nil {
		return nil, err
	}

	stageList := []*models.ProjectEnv{}
	_, err = qs.Limit(filter.PageSize, filter.PageSize*(filter.PageIndex-1)).All(&stageList)
	if err != nil {
		return nil, err
	}
	rst.Item = stageList

	return rst, nil
}

// GetProjectEnvBycIDAndEnvTag ..
func (model *ProjectModel) GetProjectEnvBycIDAndEnvTag(env string, projectID int64) (*models.ProjectEnv, error) {
	stage := models.ProjectEnv{}
	err := model.ormer.QueryTable(model.projectEnvTableName).
		Filter("deleted", false).
		Filter("arrange_env", env).
		Filter("project_id", projectID).One(&stage)
	return &stage, err
}

// UpdateProjectEnv ..
func (model *ProjectModel) UpdateProjectEnv(stage *models.ProjectEnv) error {
	_, err := model.ormer.Update(stage)
	return err
}

// DeleteProjectEnv ..
func (model *ProjectModel) DeleteProjectEnv(stageID int64) error {
	stage, err := model.GetProjectEnvByID(stageID)
	if err != nil {
		return err
	}
	stage.MarkDeleted()
	_, err = model.ormer.Update(stage)
	return err
}

// CreateProjectEnv ...
func (model *ProjectModel) CreateProjectEnv(stage *models.ProjectEnv) error {
	_, err := model.ormer.InsertOrUpdate(stage)
	return err
}

// CreatePipeline ...
func (model *ProjectModel) CreatePipeline(pipeline *models.ProjectPipeline) (int64, error) {
	created, id, err := model.ormer.ReadOrCreate(pipeline, "project_id", "name", "deleted")
	if err == nil {
		if !created {
			err = fmt.Errorf(fmt.Sprintf("pipeline name existed:%s", pipeline.Name))
		}
	}
	return id, err
}

// GetProjectPipelines ...
func (model *ProjectModel) GetProjectPipelines(projectID int64) ([]*models.ProjectPipeline, error) {
	app := []*models.ProjectPipeline{}
	qs := model.ormer.QueryTable(model.projectPipelineTableName).Filter("deleted", false)
	if projectID != -1 {
		qs = qs.Filter("project_id", projectID)
	}
	_, err := qs.All(&app)
	return app, err
}

// UpdateProjectPipeline ...
func (model *ProjectModel) UpdateProjectPipeline(pp *models.ProjectPipeline) error {
	_, err := model.ormer.Update(pp)
	return err
}

// GetProjectPipelineByID ...
func (model *ProjectModel) GetProjectPipelineByID(pipelineID int64) (*models.ProjectPipeline, error) {
	app := models.ProjectPipeline{}
	err := model.ormer.QueryTable(model.projectPipelineTableName).
		Filter("deleted", false).
		Filter("id", pipelineID).One(&app)
	return &app, err
}

// GetDefaultPipeline ..
func (model *ProjectModel) GetDefaultPipeline(projectID int64) (*models.ProjectPipeline, error) {
	pipelineItem := models.ProjectPipeline{}
	err := model.ormer.QueryTable(model.projectPipelineTableName).
		Filter("deleted", false).
		Filter("project_id", projectID).
		Filter("is_default", true).One(&pipelineItem)
	return &pipelineItem, err
}

// GetPipelinesByPagination ..
func (model *ProjectModel) GetPipelinesByPagination(filter *query.FilterQuery, projectID int64) (*query.QueryResult, error) {
	if projectID == 0 {
		return nil, fmt.Errorf("project id: %v is invalidation", projectID)
	}
	rst := &query.QueryResult{Item: []*models.ProjectPipeline{}}
	queryCond := orm.NewCondition().AndCond(orm.NewCondition().And("deleted", false).And("project_id", projectID))

	if filterCond := query.FilterCondition(filter, filter.FilterKey); filterCond != nil {
		queryCond = queryCond.AndCond(filterCond)
	}
	qs := model.ormer.QueryTable(model.projectPipelineTableName).OrderBy("-create_at").SetCond(queryCond)
	count, err := qs.Count()
	if err != nil {
		return nil, err
	}
	if err = query.FillPageInfo(rst, filter.PageIndex, filter.PageSize, int(count)); err != nil {
		return nil, err
	}

	stepList := []*models.ProjectPipeline{}
	_, err = qs.Limit(filter.PageSize, filter.PageSize*(filter.PageIndex-1)).All(&stepList)
	if err != nil {
		return nil, err
	}
	rst.Item = stepList

	return rst, nil
}

// UpdateProject ...
func (model *ProjectModel) UpdateProject(project *models.Project) error {
	_, err := model.ormer.Update(project)
	return err
}

// DeleteProject ...
func (model *ProjectModel) DeleteProject(projectID int64) error {
	project, err := model.GetProjectByID(projectID)
	if err != nil {
		return err
	}
	project.MarkDeleted()
	// _, err = model.ormer.Delete(project)
	_, err = model.ormer.Update(project)
	return err
}

/* ------ Project User Part ------ */

// GetProjectUserByRoleType ..
func (model *ProjectModel) GetProjectUserByRoleType(projectID int64, roleType string) ([]*models.ProjectUser, error) {
	users := []*models.ProjectUser{}
	qs := model.ormer.QueryTable(model.projectUserTableName).Filter("deleted", false)
	if projectID != 0 {
		qs = qs.Filter("project_id", projectID)
	}
	if roleType != "" {
		qs = qs.Filter("role", roleType)
	}
	_, err := qs.All(&users)
	return users, err
}

// GetProjectUsersByRoles ..
func (model *ProjectModel) GetProjectUsersByRoles(projectID int64, roleTypes []string) ([]*models.ProjectUser, error) {
	users := []*models.ProjectUser{}
	qs := model.ormer.QueryTable(model.projectUserTableName).Filter("deleted", false)
	if projectID != 0 {
		qs = qs.Filter("project_id", projectID)
	}
	if roleTypes != nil {
		qs = qs.Filter("role__in", roleTypes)
	}
	_, err := qs.All(&users)
	return users, err
}

// GetProjectUserByID ..
func (model *ProjectModel) GetProjectUserByID(nubmerID int64) (*models.ProjectUser, error) {
	user := models.ProjectUser{}
	qs := model.ormer.QueryTable(model.projectUserTableName).Filter("deleted", false)
	if nubmerID != 0 {
		qs = qs.Filter("id", nubmerID)
	} else {
		return nil, fmt.Errorf("成员id: %d 是无效的", nubmerID)
	}
	err := qs.One(&user)
	return &user, err
}

// GetProjectUsers ..
func (model *ProjectModel) GetProjectUsers(projectID int64) ([]*models.ProjectUser, error) {
	users := []*models.ProjectUser{}
	qs := model.ormer.QueryTable(model.projectUserTableName).Filter("deleted", false)
	if projectID != 0 {
		qs = qs.Filter("project_id", projectID)
	}
	_, err := qs.All(&users)
	return users, err
}

// UpdateProjectUser ..
func (model *ProjectModel) UpdateProjectUser(user *models.ProjectUser) error {
	_, err := model.ormer.Update(user)
	return err
}

// GetProjectUsersByRoleTypeAndUserName ..
func (model *ProjectModel) GetProjectUsersByRoleTypeAndUserName(projectID int64, roleType string, userNames []string) ([]*models.ProjectUser, error) {
	users := []*models.ProjectUser{}
	qs := model.ormer.QueryTable(model.projectUserTableName).Filter("deleted", false)
	if projectID != 0 {
		qs = qs.Filter("project_id", projectID)
	}
	_, err := qs.Filter("role", roleType).
		Filter("username__in", userNames).All(&users)
	return users, err
}

// CreateProjectUserIfNotExist ...
func (model *ProjectModel) CreateProjectUserIfNotExist(user *models.ProjectUser) (int64, error) {
	created, id, err := model.ormer.ReadOrCreate(user, "project_id", "user", "role_id", "deleted")
	if err == nil {
		if !created {
			err = fmt.Errorf(fmt.Sprintf("user: %s existed in Project User", user.User))
		}
	}
	return id, err
}

// DeleteProjectUser ...
func (model *ProjectModel) DeleteProjectUser(user *models.ProjectUser) error {
	user.MarkDeleted()
	_, err := model.ormer.Delete(user)
	return err
}

/* --- Project APP Part --- */

// GetProjectAppsList ..
func (model *ProjectModel) GetProjectAppsList(projectID int64, filter *models.ProejctAppFilterQuery) (*query.QueryResult, []*models.ProjectApp, error) {
	rst := &query.QueryResult{Item: []*models.ProjectApp{}}

	ormCond := orm.NewCondition()
	ormCond = ormCond.And("project_id", projectID).And("deleted", false)

	if filter.Name != "" {
		ormCond = ormCond.And("name__icontains", filter.Name)
	}
	if filter.Creator != "" {
		ormCond = ormCond.And("creator__icontains", filter.Creator)
	}

	if filter.Path != "" {
		ormCond = ormCond.And("path__icontains", filter.Path)
	}

	if filter.Language != "" {
		ormCond = ormCond.And("language", filter.Language)
	}

	if filter.Type != "" {
		ormCond = ormCond.And("type", filter.Type)
	}

	if filter.CreateAtStart != "" && filter.CreateAtEnd != "" {
		if createAtStart, err := time.Parse("2006-01-02", filter.CreateAtStart); err == nil {
			ormCond = ormCond.And("create_at__gte", createAtStart)
		} else {
			logs.Error("time parse error: %s", err.Error())
		}
		if createAtEnd, err := time.Parse("2006-01-02", filter.CreateAtEnd); err == nil {
			ormCond = ormCond.And("create_at__lte", createAtEnd)
		} else {
			logs.Error("time parse error: %s", err.Error())
		}
	}
	qs := model.ormer.QueryTable(model.projectAppTableName).OrderBy("-create_at").SetCond(ormCond)
	count, err := qs.Count()

	if err != nil {
		return nil, nil, err
	}
	if err = query.FillPageInfo(rst, filter.PageIndex, filter.PageSize, int(count)); err != nil {
		return nil, nil, err
	}

	appList := []*models.ProjectApp{}
	_, err = qs.Limit(filter.PageSize, filter.PageSize*(filter.PageIndex-1)).All(&appList)
	if err != nil {
		return nil, nil, err
	}
	return rst, appList, nil
}

// CreateProjectAppIfNotExist ...
func (model *ProjectModel) CreateProjectAppIfNotExist(app *models.ProjectApp) (int64, error) {
	created, id, err := model.ormer.ReadOrCreate(app, "project_id", "name", "repo_id", "deleted")
	if err == nil {
		if !created {
			err = fmt.Errorf(fmt.Sprintf("app: %v existed in project", app.FullName))
		}
	}
	return id, err
}

// GetProjectAppsByIDs ...
func (model *ProjectModel) GetProjectAppsByIDs(projectID int64, projectAppIDs []int64) ([]*models.ProjectApp, error) {
	app := []*models.ProjectApp{}
	qs := model.ormer.QueryTable(model.projectAppTableName).Filter("deleted", false)
	_, err := qs.Filter("project_id", projectID).
		Filter("id__in", projectAppIDs).All(&app)
	return app, err
}

// GetProjectAppCounts ..
func (model *ProjectModel) GetProjectAppCounts(projectID int64) (int64, error) {
	qs := model.ormer.QueryTable(model.projectAppTableName).Filter("deleted", false)
	if projectID != -1 {
		qs = qs.Filter("project_id", projectID)
	}
	num, err := qs.Count()
	return num, err
}

// GetProjectApps ...
func (model *ProjectModel) GetProjectApps(projectID int64) ([]*models.ProjectApp, error) {
	app := []*models.ProjectApp{}
	qs := model.ormer.QueryTable(model.projectAppTableName).Filter("deleted", false)
	if projectID != -1 {
		qs = qs.Filter("project_id", projectID)
	}
	_, err := qs.All(&app)
	return app, err
}

// GetProjectApp ...
func (model *ProjectModel) GetProjectApp(projectAppID int64) (*models.ProjectApp, error) {
	app := models.ProjectApp{}
	qs := model.ormer.QueryTable(model.projectAppTableName).Filter("deleted", false)
	if projectAppID != 0 {
		qs = qs.Filter("id", projectAppID)
	}
	err := qs.One(&app)
	return &app, err
}

// UpdateProjectApp ...
func (model *ProjectModel) UpdateProjectApp(projectApp *models.ProjectApp) error {
	_, err := model.ormer.Update(projectApp)
	return err
}

// DeleteProjectApp ...
func (model *ProjectModel) DeleteProjectApp(projectAppID int64) error {
	app, err := model.GetProjectApp(projectAppID)
	if err != nil {
		return err
	}
	app.MarkDeleted()
	_, err = model.ormer.Delete(app)
	return err
}
