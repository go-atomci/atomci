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
	"strings"

	"github.com/go-atomci/atomci/internal/models"
	"github.com/go-atomci/atomci/utils/query"

	"github.com/astaxie/beego/orm"
)

// ScmAppModel ...
type ScmAppModel struct {
	ormer               orm.Ormer
	scmAppTableName     string
	repoServerTableName string
	AppBranchTableName  string
}

// NewGitAppModel ...
func NewScmAppModel() (model *ScmAppModel) {
	return &ScmAppModel{
		ormer:               GetOrmer(),
		scmAppTableName:     (&models.ScmApp{}).TableName(),
		repoServerTableName: (&models.RepoServer{}).TableName(),
		AppBranchTableName:  (&models.AppBranch{}).TableName(),
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

// GetGitRepoByID ...
func (model *ScmAppModel) GetGitRepoByID(repoID int64) (*models.RepoServer, error) {
	app := models.RepoServer{}
	err := model.ormer.QueryTable(model.repoServerTableName).
		Filter("deleted", false).
		Filter("id", repoID).One(&app)
	return &app, err
}

// GetReposByprojectID ..
func (model *ScmAppModel) GetReposByprojectID(cID int64) ([]*models.RepoServer, error) {
	repos := []*models.RepoServer{}
	_, err := model.ormer.QueryTable(model.repoServerTableName).
		Filter("deleted", false).
		Filter("cid", cID).All(&repos)
	return repos, err
}

// GetRepoBycIDAndType ..
func (model *ScmAppModel) GetRepoBycIDAndType(cID int64, repoType string) (*models.RepoServer, error) {
	repo := models.RepoServer{}
	err := model.ormer.QueryTable(model.repoServerTableName).
		Filter("deleted", false).
		Filter("cid", cID).
		Filter("type", strings.ToLower(repoType)).One(&repo)
	return &repo, err
}

// GetRepoByID ..
func (model *ScmAppModel) GetRepoByID(repoID int64) (*models.RepoServer, error) {
	repo := models.RepoServer{}
	err := model.ormer.QueryTable(model.repoServerTableName).
		Filter("deleted", false).
		Filter("id", repoID).One(&repo)
	return &repo, err
}

// CreateDefaultRepo ..
func (model *ScmAppModel) CreateDefaultRepo(cID int64, repoType string) error {
	rs := &models.RepoServer{
		CID:  cID,
		Type: strings.ToLower(repoType),
	}
	if _, err := model.createRepo(rs); err != nil {
		return err
	}
	return nil
}

// UpdateRepo ...
func (model *ScmAppModel) UpdateRepo(repo *models.RepoServer) error {
	_, err := model.ormer.Update(repo)
	return err
}

func (model *ScmAppModel) createRepo(rs *models.RepoServer) (int64, error) {
	_, id, err := model.ormer.ReadOrCreate(rs, "type", "deleted", "cid")
	return id, err
}

// GetGitApps ...
func (model *ScmAppModel) GetGitApps(appIDs []int64) ([]*models.ScmApp, error) {
	apps := []*models.ScmApp{}
	qs := model.ormer.QueryTable(model.scmAppTableName).Filter("deleted", false)
	if appIDs != nil {
		qs = qs.Filter("id__in", appIDs)
	}

	_, err := qs.All(&apps)
	return apps, err
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
