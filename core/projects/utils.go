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

package projects

import (
	"fmt"

	"github.com/go-atomci/atomci/middleware/log"
	"github.com/go-atomci/atomci/models"
	"github.com/go-atomci/atomci/utils"
	"github.com/go-atomci/atomci/utils/errors"
)

func (pm *ProjectManager) projectEndVerify(projectID int64) (bool, error) {
	publishes, err := pm.publishModel.GetRunninbPublishesByProjectID(projectID)
	if err != nil {
		log.Log.Error("when end project, get publish release occur error:%s", err.Error())
		return false, fmt.Errorf("校验是否存在运行中的流水线时失败，请稍后重试")
	}
	if len(publishes) > 0 {
		return false, fmt.Errorf("存在运行中的流水线，请删除或归档流水线后重试")
	}
	// project app service
	applications, err := pm.k8sModel.GetApplicationsByProjectID(projectID)
	if err != nil {
		return false, fmt.Errorf("校验是否存在运行中的服务时失败，请稍后重试")
	}
	if len(applications) > 0 {
		return false, fmt.Errorf("存在运行中的服务，请至『我的项目』-『服务列表』删除服务后重试")
	}
	return true, nil
}

func (pm *ProjectManager) getProjectUser(projectID int64, roleType string) []*models.ProjectUser {
	modelUsers, err := pm.model.GetProjectUserByRoleType(projectID, roleType)
	if err != nil {
		log.Log.Error("get project User failed: roleType: %s, msg: %s", roleType, err)
	}
	return modelUsers
}

func (pm *ProjectManager) getProjectUsers(projectID int64) []*models.ProjectUser {
	modelUsers, err := pm.model.GetProjectUsers(projectID)
	if err != nil {
		log.Log.Error("get project User failed,msg: %s", err)
	}
	return modelUsers
}

func (pm *ProjectManager) createProjectUser(users []string, role string, projectID int64) error {
	for _, user := range users {
		userModel := models.ProjectUser{
			Addons:    models.NewAddons(),
			User:      user,
			ProjectID: projectID,
		}
		_, err := pm.model.CreateProjectUserIfNotExist(&userModel)
		if err != nil {
			return errors.NewBadRequest().SetMessage("创建项目用户失败: %s, msg: %s", user, err)
		}
	}
	return nil
}

func getNewAndNeedDeleteUsers(originUsers, requestUsers []string) ([]string, []string) {
	newUsers := []string{}
	deleteUsers := []string{}

	for _, user := range requestUsers {
		if !utils.Contains(originUsers, user) {
			newUsers = append(newUsers, user)
		}
	}
	for _, user := range originUsers {
		if !utils.Contains(requestUsers, user) {
			deleteUsers = append(deleteUsers, user)
		}
	}
	return newUsers, deleteUsers
}

// Base on roles create/delete project users
func (pm *ProjectManager) updateProjectUser(addedUsers, deleteUsers []string, role string, projectID int64) error {
	if len(addedUsers) > 0 {
		if err := pm.createProjectUser(addedUsers, role, projectID); err != nil {
			return err
		}
	}
	if len(deleteUsers) > 0 {
		if err := pm.deleteProjectUser(deleteUsers, role, projectID); err != nil {
			return err
		}
	}
	return nil
}

func (pm *ProjectManager) deleteProjectUser(users []string, role string, projectID int64) error {
	modelProjectUsers, err := pm.model.GetProjectUsersByRoleTypeAndUserName(projectID, role, users)
	if err != nil {
		return err
	}
	for _, modelUser := range modelProjectUsers {
		if err := pm.model.DeleteProjectUser(modelUser); err != nil {
			return err
		}
	}
	return nil
}

func (pm *ProjectManager) formatProjectAppResp(modelApp *models.ProjectApp) (*ProjectAppRsp, error) {
	// Get App Branches
	branches, err := pm.gitAppModel.GetAppBranches(modelApp.ID)
	if err != nil {
		return nil, err
	}
	// TODO: branchList get need commbined
	branchList := []string{}
	for _, branch := range branches {
		branchList = append(branchList, branch.BranchName)
	}
	if len(branchList) == 0 {
		branchList = []string{"master"}
	}
	compileEnvName := ""
	if modelApp.CompileEnvID != 0 {
		compileEnv, err := pm.settingModel.GetCompileEnvByID(modelApp.CompileEnvID)
		if err != nil {
			log.Log.Error("get compile env by id: %v error: %s", modelApp.CompileEnvID, err.Error())
		} else {
			compileEnvName = compileEnv.Name
		}
	}

	return &ProjectAppRsp{
		ProjectApp:        modelApp,
		CompileEnv:        compileEnvName,
		BranchHistoryList: branchList,
	}, nil

}

func (pm *ProjectManager) formatProjectAppsResp(modelApps []*models.ProjectApp) ([]*ProjectAppRsp, error) {
	projectAppsRsp := []*ProjectAppRsp{}
	for _, app := range modelApps {
		// Get App Branches
		branches, err := pm.gitAppModel.GetAppBranches(app.ID)
		if err != nil {
			return nil, err
		}
		branchList := []string{}
		for _, branch := range branches {
			branchList = append(branchList, branch.BranchName)
		}
		if len(branchList) == 0 {
			branchList = []string{"master"}
		}
		// TODO: Switch outside of the For loop， reduce mysql query
		// compileEnv, err := pm.settingModel.GetCompileEnvByID(app.CompileEnvID)
		// var compileEnvName string
		// if err != nil {
		// 	log.Log.Error("get compile env by id: %v error: %s", app.CompileEnvID, err.Error())
		// 	compileEnvName = ""
		// } else {
		// 	compileEnvName = compileEnv.Name
		// }
		projectAppRsp := &ProjectAppRsp{
			ProjectApp: app,
			// CompileEnv:        compileEnvName,
			BranchHistoryList: branchList,
		}
		projectAppsRsp = append(projectAppsRsp, projectAppRsp)
	}

	return projectAppsRsp, nil
}
