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

	"github.com/astaxie/beego/orm"

	"github.com/go-atomci/atomci/models"
)

// AppArrangeModel ...
type AppArrangeModel struct {
	ormer                    orm.Ormer
	AppArrangeTableName      string
	AppImageMappingTableName string
}

// NewAppArrangeModel ...
func NewAppArrangeModel() (model *AppArrangeModel) {
	return &AppArrangeModel{
		ormer:                    GetOrmer(),
		AppArrangeTableName:      (&models.AppArrange{}).TableName(),
		AppImageMappingTableName: (&models.AppImageMapping{}).TableName(),
	}
}

// GetAppImageMappingItemByImage ..
func (model *AppArrangeModel) GetAppImageMappingItemByImage(arrangeID int64, image string) (*models.AppImageMapping, error) {
	imageMapping := &models.AppImageMapping{}
	qs := model.ormer.QueryTable(model.AppImageMappingTableName).Filter("deleted", false)
	if arrangeID == 0 {
		return nil, fmt.Errorf("args invalidate arrange id: %v", arrangeID)
	}
	qs = qs.Filter("arrange_id", arrangeID).Filter("image", image)
	err := qs.One(imageMapping)
	return imageMapping, err
}

// GetInvalidAppImageMappingItems ..
func (model *AppArrangeModel) GetInvalidAppImageMappingItems(arrangeID, projectAppID int64, validIDs []int64) ([]*models.AppImageMapping, error) {
	imageMapping := []*models.AppImageMapping{}
	qs := model.ormer.QueryTable(model.AppImageMappingTableName).Filter("deleted", false)
	if arrangeID == 0 {
		return nil, fmt.Errorf("args invalidate arrange id: %v", arrangeID)
	}
	qs = qs.Filter("arrange_id", arrangeID).Filter("project_app_id", projectAppID).Exclude("id__in", validIDs)
	_, err := qs.All(&imageMapping)
	return imageMapping, err
}

// GetAppImageMappingItemByID ...
func (model *AppArrangeModel) GetAppImageMappingItemByID(id int64) (*models.AppImageMapping, error) {
	imageMapping := &models.AppImageMapping{}
	qs := model.ormer.QueryTable(model.AppImageMappingTableName).Filter("deleted", false)
	if id == 0 {
		return nil, fmt.Errorf("args invalidate id: %v", id)
	}
	qs = qs.Filter("id", id)
	err := qs.One(imageMapping)
	return imageMapping, err
}

// GetAppImageMappingByArrangeID ...
func (model *AppArrangeModel) GetAppImageMappingByArrangeID(arrangeID int64) ([]*models.AppImageMapping, error) {
	imageMappings := []*models.AppImageMapping{}
	qs := model.ormer.QueryTable(model.AppImageMappingTableName).Filter("deleted", false)
	if arrangeID == 0 {
		return nil, fmt.Errorf("args invalidate arrange id: %v", arrangeID)
	}
	_, err := qs.Filter("arrange_id", arrangeID).All(&imageMappings)
	return imageMappings, err
}

// GetAppImageMappingByArrangeIDAndProjectAppID ...
func (model *AppArrangeModel) GetAppImageMappingByArrangeIDAndProjectAppID(arrangeID, projectAppID int64) (*models.AppImageMapping, error) {
	imageMapping := models.AppImageMapping{}
	qs := model.ormer.QueryTable(model.AppImageMappingTableName).Filter("deleted", false)
	if arrangeID == 0 || projectAppID == 0 {
		return nil, fmt.Errorf("args invalidate arrange id: %v, projectAppID: %v", arrangeID, projectAppID)
	}
	err := qs.Filter("arrange_id", arrangeID).Filter("project_app_id", projectAppID).One(&imageMapping)
	return &imageMapping, err
}

// InsertAppImageMapping ...
func (model *AppArrangeModel) InsertAppImageMapping(appImageMappingItem *models.AppImageMapping) (int64, error) {
	id, err := model.ormer.Insert(appImageMappingItem)
	return id, err
}

// UpdateAppImageMapping ...
func (model *AppArrangeModel) UpdateAppImageMapping(appImageMappingItem *models.AppImageMapping) error {
	_, err := model.ormer.Update(appImageMappingItem)
	return err
}

// DeleteAppImageMapping ...
func (model *AppArrangeModel) DeleteAppImageMapping(appImageMappingItem *models.AppImageMapping) error {
	_, err := model.ormer.Delete(appImageMappingItem)
	return err
}

// GetAppArrange ...
func (model *AppArrangeModel) GetAppArrange(appID, envID int64) (*models.AppArrange, error) {
	arrange := &models.AppArrange{}
	qs := model.ormer.QueryTable(model.AppArrangeTableName).Filter("deleted", false)
	if appID == 0 || envID == 0 {
		return nil, fmt.Errorf("args invalidate app id: %v, env id: %v", appID, envID)
	}
	qs = qs.Filter("project_app_id", appID).Filter("env_id", envID)
	err := qs.One(arrange)
	return arrange, err
}

// AppArrangeIsExisted check
func (model *AppArrangeModel) AppArrangeIsExisted(AppID int64, arrangeEnv string) bool {
	return model.ormer.QueryTable(model.AppArrangeTableName).Filter("deleted", false).Filter("publish_app_id", AppID).Filter("arrange_env", arrangeEnv).Exist()
}

// InsertOrUpdateAppArrange ...
func (model *AppArrangeModel) InsertAppArrange(arrange *models.AppArrange) error {
	_, err := model.ormer.Insert(arrange)
	return err
}

// DeleteAppArrange ...
func (model *AppArrangeModel) DeleteAppArrange(AppID, envID int64) error {
	arrange, err := model.GetAppArrange(AppID, envID)
	if err != nil {
		return err
	}
	arrange.MarkDeleted()
	_, err = model.ormer.Delete(arrange)
	return err
}

// UpdateAppArrange ...
func (model *AppArrangeModel) UpdateAppArrange(arrange *models.AppArrange) error {
	_, err := model.ormer.Update(arrange)
	return err
}
