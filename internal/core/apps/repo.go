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

	"github.com/go-atomci/atomci/internal/core/settings"
	"github.com/go-atomci/atomci/internal/dao"
	"github.com/go-atomci/atomci/internal/middleware/log"
	"github.com/go-atomci/atomci/utils/query"

	"github.com/drone/go-scm/scm"
)

// AppManager ...
type AppManager struct {
	model           *dao.AppArrangeModel
	scmAppModel     *dao.ScmAppModel
	projectModel    *dao.ProjectModel
	settingsHandler *settings.SettingManager
}

// NewAppManager ...
func NewAppManager() *AppManager {
	return &AppManager{
		model:           dao.NewAppArrangeModel(),
		scmAppModel:     dao.NewScmAppModel(),
		projectModel:    dao.NewProjectModel(),
		settingsHandler: settings.NewSettingManager(),
	}
}

// AppBranches ...
func (manager *AppManager) AppBranches(appID int64, filter *query.FilterQuery) (*query.QueryResult, error) {
	return manager.scmAppModel.GetAppBranchesByPagination(appID, filter)
}

// GetScmProjectsByRepoID ..
func (manager *AppManager) GetScmProjectsByRepoID(repoID int64) (interface{}, error) {
	scmIntegrateResp, err := manager.settingsHandler.GetSCMIntegrateSettinByID(repoID)
	if err != nil {
		return nil, err
	}
	// 获取仓库项目需要授权，若配置的仓库是公共库，则直接返回
	if scmIntegrateResp.ScmAuthConf.Token == "" {
		return []*RepoProjectRsp{}, nil
	}
	scmClient, err := NewScmProvider(scmIntegrateResp.Type, scmIntegrateResp.ScmAuthConf.URL, scmIntegrateResp.ScmAuthConf.Token)
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

// 验证仓库源是否能正常连通，若无token，只要地址能通就行；若有token，则必须认证通过
func (manager *AppManager) VerifyRepoConnetion(scmType string, url string, token string) error {
	scmClient, err := NewScmProvider(scmType, url, token)
	if err != nil {
		return err
	}
	op := scm.ListOptions{
		Page: 0,
		Size: 1,
	}
	_, resp, err := scmClient.Organizations.List(context.Background(), op)
	if resp == nil && err != nil {
		return fmt.Errorf("连接源码仓库失败,错误信息：%s", err)
	}
	if token == "" {
		if (resp.Status >= 200 && resp.Status <= 299) || resp.Status == 401 {
			return nil
		}
		return fmt.Errorf("连接源码仓库失败,服务器返回：%d", resp.Status)
	} else {
		if resp.Status >= 200 && resp.Status <= 299 {
			return nil
		} else if resp.Status == 401 {
			return fmt.Errorf("连接源码仓库成功，但是认证失败,仓库返回：401")
		} else {
			return fmt.Errorf("连接源码仓库失败,服务器返回：%d", resp.Status)
		}
	}
}

// 验证仓库地址是否能正常连通,方式通过获取代码分支，若能正常获取，则表示通过
func (manager *AppManager) VerifyAppConnetion(repoID int64, url string, repo string) error {
	scmIntegrateResp, err := manager.settingsHandler.GetSCMIntegrateSettinByID(repoID)
	if err != nil {
		return err
	}
	token := scmIntegrateResp.ScmAuthConf.Token
	scmClient, err := NewScmProvider(scmIntegrateResp.Type, url, token)
	if err != nil {
		return err
	}
	listOptions := scm.ListOptions{
		Page: 0,
		Size: 1,
	}
	_, resp, err := scmClient.Git.ListBranches(context.Background(), repo, listOptions)
	if resp == nil && err != nil {
		return fmt.Errorf("连接源码仓库失败,错误信息：%s", err)
	}

	if resp.Status >= 200 && resp.Status <= 299 {
		return nil
	} else if resp.Status == 401 {
		return fmt.Errorf("连接源码仓库成功,但是认证失败,服务器返回：401")
	} else if resp.Status == 404 {
		return fmt.Errorf("连接源码仓库失败,仓库路径错误或当前Token没有权限")
	} else {
		return fmt.Errorf("连接源码仓库失败,服务器返回：%d", resp.Status)
	}
}
