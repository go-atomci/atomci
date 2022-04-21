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
	"github.com/drone/go-scm/scm/driver/gogs"
	"net/http"
	"strings"

	"github.com/drone/go-scm/scm/driver/gitea"

	"github.com/go-atomci/atomci/internal/middleware/log"
	"github.com/go-atomci/atomci/internal/models"
	"github.com/go-atomci/atomci/utils/query"

	"github.com/go-atomci/atomci/utils"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/gitee"
	"github.com/drone/go-scm/scm/driver/github"
	"github.com/drone/go-scm/scm/driver/gitlab"
	"github.com/drone/go-scm/scm/transport"
)

// NewScmProvider ..
func NewScmProvider(vcsType, vcsPath, token string) (*scm.Client, error) {
	var err error
	var client *scm.Client
	switch strings.ToLower(vcsType) {
	case "gitea", "gitlab", "gogs":
		if strings.HasSuffix(vcsPath, ".git") {
			vcsPath = strings.TrimSuffix(vcsPath, ".git")
		}

		vcsPathSplit := strings.Split(vcsPath, "://")
		// TODO: verify vcsPath, only support http, do not support git@gitlab.com:/dddd.git
		projectPathSplit := strings.Split(vcsPathSplit[1], "/")
		projectName := strings.Join(projectPathSplit[1:], "/")
		log.Log.Debug("git projectpathsplit: %s,\tprojectName: %s", projectPathSplit, projectName)

		schema := vcsPathSplit[0]
		gitRepo := strings.ToLower(vcsType)

		if "gitea" == gitRepo {
			client, err = gitea.New(schema + "://" + projectPathSplit[0])
			client.Client = &http.Client{
				Transport: &transport.BearerToken{
					Token: token,
				},
			}
		} else if "gitlab" == gitRepo {
			client, err = gitlab.New(schema + "://" + projectPathSplit[0])

			client.Client = &http.Client{
				Transport: &transport.PrivateToken{
					Token: token,
				},
			}
		} else {
			client, err = gogs.New(schema + "://" + projectPathSplit[0])
			client.Client = &http.Client{
				Transport: &transport.PrivateToken{
					Token: token,
				},
			}
		}
	case "github":
		client = github.NewDefault()

		client.Client = &http.Client{
			Transport: &transport.BearerToken{
				Token: token,
			},
		}

	case "gitee":
		client = gitee.NewDefault()

		client.Client = &http.Client{
			Transport: &transport.BearerToken{
				Token: token,
			},
		}

	default:
		err = fmt.Errorf("source code management system not configured")
	}
	return client, err
}

// SyncAppBranches ...
func (manager *AppManager) SyncAppBranches(appID int64) error {
	scmApp, _ := manager.scmAppModel.GetScmAppByID(appID)
	scmIntegrateResp, err := manager.settingsHandler.GetSCMIntegrateSettinByID(scmApp.RepoID)
	if err != nil {
		return err
	}
	if err != nil {
		log.Log.Error("getCompileEnvByID occur error: %v", err.Error())
		return fmt.Errorf("网络错误，请重试")
	}
	client, err := NewScmProvider(scmIntegrateResp.Type, scmApp.Path, scmIntegrateResp.Token)
	branchList := []*scm.Reference{}
	listOptions := scm.ListOptions{
		Page: 1,
		Size: 100,
	}
	got, res, err := client.Git.ListBranches(context.Background(), scmApp.FullName, listOptions)
	if err != nil {
		return fmt.Errorf("when get branches list from gitlab occur error: %s", err.Error())
	}
	branchList = append(branchList, got...)

	for i := 1; i < res.Page.Last; {
		listOptions.Page++
		got, _, err := client.Git.ListBranches(context.Background(), scmApp.FullName, listOptions)
		if err != nil {
			return fmt.Errorf("when get branches list from gitlab occur error: %s", err.Error())
		}
		branchList = append(branchList, got...)
		i++
	}
	for _, branch := range branchList {
		if strings.HasPrefix(branch.Name, "release_") {
			continue
		}
		originBranch, err := manager.scmAppModel.GetAppBranchByName(appID, branch.Name)
		if err != nil {
			if strings.Contains(err.Error(), "no row found") {
				err = nil
			} else {
				return fmt.Errorf("when get app branch occur error: %s", err.Error())
			}
		}
		if originBranch.BranchName == "" {
			appBranch := &models.AppBranch{
				BranchName: branch.Name,
				Path:       scmApp.Path,
				AppID:      appID,
			}
			if _, err := manager.scmAppModel.CreateAppBranchIfNotExist(appBranch); err != nil {
				return err
			}
		} else {
			originBranch.Path = scmApp.Path
			if err := manager.scmAppModel.UpdateAppBranch(originBranch); err != nil {
				return err
			}
		}
	}

	branchListInDB, err := manager.scmAppModel.GetAppBranches(appID)
	if err != nil {
		return err
	}
	branchNameList := []string{}
	for _, branch := range branchList {
		branchNameList = append(branchNameList, branch.Name)
	}
	for _, branchDBItem := range branchListInDB {
		if !utils.Contains(branchNameList, branchDBItem.BranchName) {
			manager.scmAppModel.SoftDeleteAppBranch(branchDBItem)
		}
	}
	return nil
}

// CreateSCMApp ...
func (manager *AppManager) CreateSCMApp(item *ScmAppReq, creator string) (int64, error) {
	log.Log.Debug("request params: %+v", item)

	if item.BranchName == "" {
		// reset default value is master
		item.BranchName = "master"
	}

	if item.Dockerfile == "" {
		item.Dockerfile = "Dockerfile"
	}
	scmAppModel := models.ScmApp{
		Addons:       models.NewAddons(),
		Creator:      creator,
		CompileEnvID: item.CompileEnvID,
		Name:         item.Name,
		FullName:     item.FullName,
		Language:     item.Language,
		BranchName:   item.BranchName,
		Path:         item.Path,
		RepoID:       item.RepoID,
		BuildPath:    item.BuildPath,
		Dockerfile:   item.Dockerfile,
	}

	id, err := manager.scmAppModel.CreateScmAppIfNotExist(&scmAppModel)
	if err != nil {
		log.Log.Error("create scm app error: %s", err)
		return 0, err
	}

	return id, nil
}

// GetProjectAppsByPagination ..
func (manager *AppManager) GetScmApps() ([]*models.ScmApp, error) {
	return manager.scmAppModel.GetScmApps()
}

// GetProjectAppsByPagination ..
func (manager *AppManager) GetScmAppsByPagination(filter *query.FilterQuery) (*query.QueryResult, error) {
	return manager.scmAppModel.GetScmAppsByPagination(filter)
}

func (manager *AppManager) GetScmApp(appID int64) (*SCMAppRsp, error) {
	app, err := manager.scmAppModel.GetScmAppByID(appID)
	if err != nil {
		return nil, err
	}
	return manager.formatscmAppResp(app)
}

// UpdateProjectApp ..
func (manager *AppManager) UpdateProjectApp(scmAppID int64, req *ScmAppUpdateReq) error {
	log.Log.Debug("update app projectAppID: %v, params: %+v", scmAppID, req)
	if req.Name == "" {
		return fmt.Errorf("请输入有效的『仓库名』")
	}

	if req.Path == "" {
		return fmt.Errorf("请输入有效的『路径』")
	}
	scmApp, err := manager.scmAppModel.GetScmAppByID(scmAppID)
	if err != nil {
		return err
	}

	if req.BuildPath == "" {
		scmApp.BuildPath = "/"
	} else {
		scmApp.BuildPath = req.BuildPath
	}

	if req.Dockerfile == "" {
		scmApp.Dockerfile = "Dockerfile"
	} else {
		scmApp.Dockerfile = req.Dockerfile
	}

	scmApp.BranchName = req.BranchName
	scmApp.CompileEnvID = req.CompileEnvID
	scmApp.Language = req.Language
	scmApp.Name = req.Name
	scmApp.Path = req.Path
	return manager.scmAppModel.UpdateSCMApp(scmApp)
}

func (manager *AppManager) DeleteSCMApp(scmAppID int64) error {
	log.Log.Debug("delete project app, scmAppID: %v", scmAppID)

	_, err := manager.scmAppModel.GetScmAppByID(scmAppID)
	if err != nil {
		log.Log.Error("when delete scm app, get scm app occur error: %s", err.Error())
		return fmt.Errorf("当前代码库可能已经删除，请你刷新页面后重试")
	}

	// TODO: add publish order verify
	err = manager.scmAppModel.DeleteSCMApp(scmAppID)
	if err != nil {
		return err
	}
	// TODO: delete app service constraint
	return nil
}

func (manager *AppManager) formatscmAppResp(modelApp *models.ScmApp) (*SCMAppRsp, error) {
	compileEnvName := ""
	if modelApp.CompileEnvID != 0 {
		compileEnv, err := manager.settingsHandler.GetCompileEnvByID(modelApp.CompileEnvID)
		if err != nil {
			log.Log.Error("get compile env by id: %v error: %s", modelApp.CompileEnvID, err.Error())
		} else {
			compileEnvName = compileEnv.Name
		}
	}

	return &SCMAppRsp{
		ScmApp:     modelApp,
		CompileEnv: compileEnvName,
	}, nil

}
