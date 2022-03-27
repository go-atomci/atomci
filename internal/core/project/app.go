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

	"github.com/go-atomci/atomci/internal/middleware/log"
	"github.com/go-atomci/atomci/internal/models"
	"github.com/go-atomci/atomci/utils/query"
)

// CreateProjectApp ...
func (pm *ProjectManager) CreateProjectApp(projectID int64, item *ProjectAppReq, creator string) error {
	log.Log.Debug("request params: %+v", item)

	projectAppModel := models.ProjectApp{
		Addons:    models.NewAddons(),
		Creator:   creator,
		ProjectID: projectID,
		ScmID:     item.SCMID,
	}

	_, err := pm.model.CreateProjectAppIfNotExist(&projectAppModel)
	if err != nil {
		log.Log.Error("create project app error: %s", err)
		return err
	}

	return nil
}

// GetProjectApps ..
func (pm *ProjectManager) GetProjectApps(projectID int64) ([]*ProjectAppRsp, error) {
	modelProjectApps, err := pm.model.GetProjectApps(projectID)
	if err != nil {
		return nil, err
	}
	return pm.formatProjectAppsResp(modelProjectApps)
}

// GetProjectAppsByPagination ..
func (pm *ProjectManager) GetProjectAppsByPagination(projectID int64, filter *models.ProejctAppFilterQuery) (*query.QueryResult, error) {
	apps, modelDatas, err := pm.model.GetProjectAppsList(projectID, filter)
	if err != nil {
		return nil, err
	}

	projectAppsRsp, err := pm.formatProjectAppsResp(modelDatas)
	if err != nil {
		return nil, err
	}
	apps.Item = projectAppsRsp

	return apps, nil
}

// DeleteProjectApp ...
func (pm *ProjectManager) DeleteProjectApp(projectAppID int64) error {
	log.Log.Debug("delete project app, projectAppID: %v", projectAppID)

	_, err := pm.model.GetProjectApp(projectAppID)
	if err != nil {
		log.Log.Error("when delete project app, get project app occur error: %s", err.Error())
		return fmt.Errorf("当前代码库可能已经删除，请你刷新页面后重试")
	}

	// TODO: add publish order verify
	err = pm.model.DeleteProjectApp(projectAppID)
	if err != nil {
		return err
	}
	// TODO: delete app service constraint
	return nil
}

// UpdateProjectApp ..
func (pm *ProjectManager) UpdateProjectApp(projectID, projectAppID int64, req *ProjectAppUpdateReq) error {
	_, err := pm.model.GetProjectAppByScmID(projectID, req.ScmID)
	if err == nil {
		return fmt.Errorf("already exist scmid: %v register", req.ScmID)
	}
	projectApp, err := pm.model.GetProjectApp(projectAppID)
	if err != nil {
		return err
	}
	projectApp.ScmID = req.ScmID
	return pm.model.UpdateProjectApp(projectApp)
}
