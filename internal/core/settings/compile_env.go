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

package settings

import (
	"errors"

	"github.com/go-atomci/atomci/internal/middleware/log"
	"github.com/go-atomci/atomci/internal/models"
	"github.com/go-atomci/atomci/utils/query"
)

// CompileEnvReq ..
type CompileEnvReq struct {
	Name        string `json:"name,omitempty"`
	Image       string `json:"image,omitempty"`
	Command     string `json:"command,omitempty"`
	Args        string `json:"args,omitempty"`
	Description string `json:"description,omitempty"`
}

// GetCompileEnvs ..
func (pm *SettingManager) GetCompileEnvs() ([]*models.CompileEnv, error) {
	items, err := pm.model.GetCompileEnvs()
	if err != nil {
		log.Log.Error("get interate settings error: %s", err.Error())
		return nil, err
	}

	return items, err
}

// GetCompileEnvByID ..
func (pm *SettingManager) GetCompileEnvByID(id int64) (*models.CompileEnv, error) {
	compileEnv, err := pm.model.GetCompileEnvByID(id)
	if err != nil {
		log.Log.Error("when GetCompileEnvBy id: %v, occur error: %s", id, err.Error())
		return nil, err
	}
	return compileEnv, err
}

// GetCompileEnvByID ..
func (pm *SettingManager) GetCompileEnvByName(name string) (*models.CompileEnv, error) {
	compileEnv, err := pm.model.GetCompileEnvByName(name)
	if err != nil {
		log.Log.Error("when GetCompileEnvBy name: %v, occur error: %s", name, err.Error())
		return nil, err
	}
	return compileEnv, err
}

// GetCompileEnvsByPagination ..
func (pm *SettingManager) GetCompileEnvsByPagination(filter *query.FilterQuery) (*query.QueryResult, error) {
	queryResult, settingsList, err := pm.model.GetCompileEnvsByPagination(filter)
	if err != nil {
		return nil, err
	}
	queryResult.Item = settingsList
	return queryResult, err
}

// resetEnv clear env config
func resetEnv(env *string) {
	*env = ""
}

func compileEnvNameUnique(pm *SettingManager, name string, stepId int64) error {
	if len(name) == 0 {
		return errors.New("param `Name` is not allowed empty")
	}

	exists, _ := pm.model.GetCompileEnvByName(name)

	if exists != nil && (stepId == 0 || exists.ID != stepId) {
		return errors.New("环境名称 `" + name + "` 已经存在")
	}

	return nil
}

// UpdateCompileEnv ..
func (pm *SettingManager) UpdateCompileEnv(request *CompileEnvReq, stepID int64) error {
	compileEnv, err := pm.model.GetCompileEnvByID(stepID)
	if err != nil {
		return err
	}

	if err := compileEnvNameUnique(pm, request.Name, stepID); err != nil {
		return err
	}

	if request.Name != "" {
		compileEnv.Name = request.Name
	}

	if request.Args != "" {
		compileEnv.Args = request.Args
	} else {
		resetEnv(&compileEnv.Args)
	}

	if request.Command != "" {
		compileEnv.Command = request.Command
	} else {
		resetEnv(&compileEnv.Command)
	}

	if request.Description != "" {
		compileEnv.Description = request.Description
	} else {
		resetEnv(&compileEnv.Description)
	}

	if request.Image != "" {
		compileEnv.Image = request.Image
	}

	return pm.model.UpdateCompileEnv(compileEnv)
}

// CreateCompileEnv ..
func (pm *SettingManager) CreateCompileEnv(request *CompileEnvReq, creator string) error {

	if err := compileEnvNameUnique(pm, request.Name, 0); err != nil {
		return err
	}

	// TODO: verify req struct is valid
	newCompileEnv := &models.CompileEnv{
		Name:        request.Name,
		Description: request.Description,
		Creator:     creator,
		Image:       request.Image,
		Command:     request.Command,
		Args:        request.Args,
	}

	return pm.model.CreateCompileEnv(newCompileEnv)
}

// DeleteCompileEnv ..
func (pm *SettingManager) DeleteCompileEnv(stageID int64) error {
	// TODO: add compile env delete verify
	return pm.model.DeleteCompileEnv(stageID)
}
