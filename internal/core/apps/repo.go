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

package apps

import (
	"context"
	"fmt"
	"github.com/go-atomci/atomci/internal/dao"
	"github.com/go-atomci/atomci/internal/middleware/log"
	"github.com/go-atomci/atomci/internal/models"
	"github.com/go-atomci/atomci/utils/query"

	"github.com/astaxie/beego/orm"
	"github.com/drone/go-scm/scm"
)

// AppManager ...
type AppManager struct {
	model        *dao.AppArrangeModel
	gitAppModel  *dao.GitAppModel
	projectModel *dao.ProjectModel
}

// NewAppManager ...
func NewAppManager() *AppManager {
	return &AppManager{
		model:        dao.NewAppArrangeModel(),
		gitAppModel:  dao.NewGitAppModel(),
		projectModel: dao.NewProjectModel(),
	}
}

// AppBranches ...
func (manager *AppManager) AppBranches(appID int64, filter *query.FilterQuery) (*query.QueryResult, error) {
	return manager.gitAppModel.GetAppBranchesByPagination(appID, filter)
}

// GetRepos ..
func (manager *AppManager) GetRepos(projectID int64) ([]*RepoServerRsp, error) {
	repos := []*models.RepoServer{}
	// TODO: support code repository defined,
	defaultRepos := []string{"gitlab", "github", "gitee", "gitea"}
	// defaultRepos := []string{"gitlab"}
	for _, item := range defaultRepos {
		_, err := manager.gitAppModel.GetRepoBycIDAndType(projectID, item)
		if err != nil {
			if err == orm.ErrNoRows {
				if err := manager.gitAppModel.CreateDefaultRepo(projectID, item); err != nil {
					log.Log.Error("create default repos failed: %v", err.Error())
					return nil, fmt.Errorf("网络异常，请重试")
				}
				_, err = manager.gitAppModel.GetRepoBycIDAndType(projectID, item)
				if err != nil {
					log.Log.Error("after create, get repos occur error: %v", err.Error())
					return nil, fmt.Errorf("网络异常，请重试")
				}
			} else {
				return nil, err
			}
		}
	}
	repos, err := manager.gitAppModel.GetReposByprojectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("网络异常，请重试")
	}
	rsp := []*RepoServerRsp{}
	for _, repoItem := range repos {
		itemRsp := &RepoServerRsp{
			SetupRepo: SetupRepo{
				User:    repoItem.User,
				BaseURL: repoItem.BaseURL,
			},
			Type:   repoItem.Type,
			RepoID: repoItem.ID,
		}
		rsp = append(rsp, itemRsp)
	}
	return rsp, nil
}

// SetRepoAndGetProjects ..
func (manager *AppManager) SetRepoAndGetProjects(cID, repoID int64, request *SetupRepo) (interface{}, error) {
	repoModel, err := manager.gitAppModel.GetRepoByID(repoID)
	if err != nil {
		return nil, err
	}
	if len(request.Token) > 0 {
		repoModel.Token = request.Token
		repoModel.User = request.User
		repoModel.BaseURL = request.BaseURL
		if err := manager.gitAppModel.UpdateRepo(repoModel); err != nil {
			log.Log.Error("when setRepoGetprojects, update repomodel failed: %v", err.Error())
		}
	} else {
		if len(repoModel.Token) == 0 && len(repoModel.BaseURL) == 0 {
			return nil, fmt.Errorf("首次同步，麻烦输入相关验证信息")
		}
	}

	scmClient, err := NewScmProvider(repoModel.Type, repoModel.BaseURL, repoModel.Token)
	if err != nil {
		log.Log.Error("init scm Client occur error: %v", err.Error())
		return nil, fmt.Errorf("网络错误，请重试")
	}
	listOptions := scm.ListOptions{
		Page: 1,
		Size: 100,
	}
	repoList := []*scm.Repository{}
	got, rsp, err := scmClient.Repositories.List(context.Background(), listOptions)
	if err != nil {
		return nil, fmt.Errorf("scmclient get repositories list error: %s", err.Error())
	}
	repoList = append(repoList, got...)
	for i := 1; i < rsp.Page.Last; {
		listOptions.Page++
		got, _, err := scmClient.Repositories.List(context.Background(), listOptions)
		if err != nil {
			return nil, fmt.Errorf("when get repositories list from gitlab occur error: %s", err.Error())
		}
		repoList = append(repoList, got...)
		i++
	}

	newRsp := []*RepoProjectRsp{}
	for _, item := range repoList {
		newItem := &RepoProjectRsp{
			Name:     item.Name,
			FullName: item.Namespace + "/" + item.Name,
			Path:     item.Clone,
			RepoID:   repoID,
		}
		newRsp = append(newRsp, newItem)
	}
	return newRsp, nil
}
