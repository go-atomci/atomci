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
	"net/http"
	"strings"

	"github.com/go-atomci/atomci/internal/middleware/log"
	"github.com/go-atomci/atomci/internal/models"

	"github.com/go-atomci/atomci/utils"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/github"
	"github.com/drone/go-scm/scm/driver/gitlab"
	"github.com/drone/go-scm/scm/transport"
)

// NewScmProvider ..
func NewScmProvider(vcsType, vcsPath, user, token string) (*scm.Client, error) {
	var err error
	var client *scm.Client
	switch strings.ToLower(vcsType) {
	case "gitlab":
		if strings.HasSuffix(vcsPath, ".git") {
			vcsPath = strings.Replace(vcsPath, ".git", "", -1)
		}
		// TODO: verify vcsPath, only support http, do not support git@gitlab.com:/dddd.git
		projectPathSplit := strings.Split(strings.Split(vcsPath, "://")[1], "/")
		projectName := strings.Join(projectPathSplit[1:], "/")
		log.Log.Debug("git projectpathsplit: %s,\tprojectName: %s", projectPathSplit, projectName)

		schema := strings.Split(vcsPath, "://")[0]
		client, err = gitlab.New(schema + "://" + projectPathSplit[0])
		client.Client = &http.Client{
			Transport: &transport.PrivateToken{
				Token: token,
			},
		}
	case "github":
		if strings.HasSuffix(vcsPath, ".git") {
			vcsPath = strings.Replace(vcsPath, ".git", "", -1)
		}
		// TODO: verify vcsPath, only support http, do not support git@github.com:/dddd.git
		projectPathSplit := strings.Split(strings.Split(vcsPath, "://")[1], "/")
		projectName := strings.Join(projectPathSplit[1:], "/")
		log.Log.Debug("git projectpathsplit: %s,\tprojectName: %s", projectPathSplit, projectName)

		// TODO: github does not work
		schema := strings.Split(vcsPath, "://")[0]
		client, err = github.New(schema + "://" + projectPathSplit[0])
		client.Client = &http.Client{
			Transport: &transport.PrivateToken{
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
	projectApp, err := manager.projectModel.GetProjectApp(appID)
	repoModel, err := manager.gitAppModel.GetRepoByID(projectApp.RepoID)
	if err != nil {
		log.Log.Error("GetRepoByID occur error: %v", err.Error())
		return fmt.Errorf("网络错误，请重试")
	}
	client, err := NewScmProvider(repoModel.Type, projectApp.Path, repoModel.User, repoModel.Token)
	branchList := []*scm.Reference{}
	listOptions := scm.ListOptions{
		Page: 1,
		Size: 100,
	}
	got, res, err := client.Git.ListBranches(context.Background(), projectApp.FullName, listOptions)
	if err != nil {
		return fmt.Errorf("when get branches list from gitlab occur error: %s", err.Error())
	}
	branchList = append(branchList, got...)

	for i := 1; i < res.Page.Last; {
		listOptions.Page++
		got, _, err := client.Git.ListBranches(context.Background(), projectApp.FullName, listOptions)
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
		originBranch, err := manager.gitAppModel.GetAppBranchByName(appID, branch.Name)
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
				Path:       projectApp.Path,
				AppID:      appID,
			}
			if _, err := manager.gitAppModel.CreateAppBranchIfNotExist(appBranch); err != nil {
				return err
			}
		} else {
			originBranch.Path = projectApp.Path
			if err := manager.gitAppModel.UpdateAppBranch(originBranch); err != nil {
				return err
			}
		}
	}

	branchListInDB, err := manager.gitAppModel.GetAppBranches(appID)
	if err != nil {
		return err
	}
	branchNameList := []string{}
	for _, branch := range branchList {
		branchNameList = append(branchNameList, branch.Name)
	}
	for _, branchDBItem := range branchListInDB {
		if !utils.Contains(branchNameList, branchDBItem.BranchName) {
			manager.gitAppModel.SoftDeleteAppBranch(branchDBItem)
		}
	}
	return nil
}
