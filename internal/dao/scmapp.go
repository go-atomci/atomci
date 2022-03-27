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

	"github.com/go-atomci/atomci/internal/models"
	"github.com/go-atomci/atomci/utils/query"

	"github.com/astaxie/beego/orm"
)

// ScmAppModel ...
type ScmAppModel struct {
	ormer              orm.Ormer
	scmAppTableName    string
	AppBranchTableName string
}

// NewGitAppModel ...
func NewScmAppModel() (model *ScmAppModel) {
	return &ScmAppModel{
		ormer:              GetOrmer(),
		scmAppTableName:    (&models.ScmApp{}).TableName(),
		AppBranchTableName: (&models.AppBranch{}).TableName(),
	}
}

// CreateProjectAppIfNotExist ...
func (model *ScmAppModel) CreateScmAppIfNotExist(app *models.ScmApp) (int64, error) {
	created, id, err := model.ormer.ReadOrCreate(app, "name", "repo_id", "deleted")
	if err == nil {
		if !created {
			err = fmt.Errorf(fmt.Sprintf("app: %v existed in project", app.FullName))
		}
	}
	return id, err
}

func (model *ScmAppModel) GetScmApps() ([]*models.ScmApp, error) {
	app := []*models.ScmApp{}
	qs := model.ormer.QueryTable(model.scmAppTableName).Filter("deleted", false)
	// TODO: add scm app tags
	_, err := qs.All(&app)
	return app, err
}

func (model *ScmAppModel) GetScmAppByID(appID int64) (*models.ScmApp, error) {
	app := models.ScmApp{}
	err := model.ormer.QueryTable(model.scmAppTableName).
		Filter("deleted", false).
		Filter("id", appID).One(&app)
	return &app, err
}

// UpdateProjectApp ...
func (model *ScmAppModel) UpdateSCMApp(scmApp *models.ScmApp) error {
	_, err := model.ormer.Update(scmApp)
	return err
}

// DeleteProjectApp ...
func (model *ScmAppModel) DeleteSCMApp(scmAppID int64) error {
	app, err := model.GetScmAppByID(scmAppID)
	if err != nil {
		return err
	}
	app.MarkDeleted()
	_, err = model.ormer.Delete(app)
	return err
}

// GetCompileEnvsByPagination ..
func (model *ScmAppModel) GetScmAppsByPagination(filter *query.FilterQuery) (*query.QueryResult, error) {
	rst := &query.QueryResult{Item: []*models.ScmApp{}}
	queryCond := orm.NewCondition().AndCond(orm.NewCondition().And("deleted", false))

	if filterCond := query.FilterCondition(filter, filter.FilterKey); filterCond != nil {
		queryCond = queryCond.AndCond(filterCond)
	}
	qs := model.ormer.QueryTable(model.scmAppTableName).OrderBy("-create_at").SetCond(queryCond)
	count, err := qs.Count()
	if err != nil {
		return nil, err
	}
	if err = query.FillPageInfo(rst, filter.PageIndex, filter.PageSize, int(count)); err != nil {
		return nil, err
	}

	scmApplist := []*models.ScmApp{}
	_, err = qs.Limit(filter.PageSize, filter.PageSize*(filter.PageIndex-1)).All(&scmApplist)
	if err != nil {
		return nil, err
	}
	rst.Item = scmApplist

	return rst, nil
}

// CreateAppBranchIfNotExist ...
func (model *ScmAppModel) CreateAppBranchIfNotExist(branch *models.AppBranch) (int64, error) {
	created, id, err := model.ormer.ReadOrCreate(branch, "branch_name", "app_id", "deleted")
	if err == nil {
		if !created {
			err = fmt.Errorf(fmt.Sprintf("branch_name: %v existed in app branch table", branch.BranchName))
		}
	}
	return id, err
}

// UpdateAppBranch ...
func (model *ScmAppModel) UpdateAppBranch(branch *models.AppBranch) error {
	_, err := model.ormer.Update(branch)
	return err
}

// SoftDeleteAppBranch ...
func (model *ScmAppModel) SoftDeleteAppBranch(branch *models.AppBranch) error {
	branch.MarkDeleted()
	return model.UpdateAppBranch(branch)
}

// GetAppBranchesByPagination ...
func (model *ScmAppModel) GetAppBranchesByPagination(appID int64, filter *query.FilterQuery) (*query.QueryResult, error) {
	rst := &query.QueryResult{Item: []*models.AppBranch{}}
	queryCond := orm.NewCondition().AndCond(orm.NewCondition().And("deleted", false))

	queryCond = queryCond.AndCond(orm.NewCondition().And("app_id", appID))

	if filterCond := query.FilterCondition(filter, filter.FilterKey); filterCond != nil {
		queryCond = queryCond.AndCond(filterCond)
	}
	qs := model.ormer.QueryTable(model.AppBranchTableName).OrderBy("-create_at").SetCond(queryCond)
	count, err := qs.Count()

	if err != nil {
		return nil, err
	}
	if err = query.FillPageInfo(rst, filter.PageIndex, filter.PageSize, int(count)); err != nil {
		return nil, err
	}

	appList := []*models.AppBranch{}
	_, err = qs.Limit(filter.PageSize, filter.PageSize*(filter.PageIndex-1)).All(&appList)
	if err != nil {
		return nil, err
	}
	rst.Item = appList
	return rst, nil
}

// GetAppBranches ...
func (model *ScmAppModel) GetAppBranches(appID int64) ([]*models.AppBranch, error) {
	branches := []*models.AppBranch{}
	qs := model.ormer.QueryTable(model.AppBranchTableName).Filter("deleted", false)
	if appID != 0 {
		qs = qs.Filter("app_id", appID)
	}
	_, err := qs.All(&branches)
	return branches, err
}

// GetAppBranchByName ...
func (model *ScmAppModel) GetAppBranchByName(appID int64, branchName string) (*models.AppBranch, error) {
	branch := models.AppBranch{}
	err := model.ormer.QueryTable(model.AppBranchTableName).
		Filter("deleted", false).
		Filter("app_id", appID).
		Filter("branch_name", branchName).One(&branch)
	return &branch, err
}
