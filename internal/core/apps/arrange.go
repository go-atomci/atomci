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
	"fmt"

	"github.com/go-atomci/atomci/internal/middleware/log"
	"github.com/go-atomci/atomci/internal/models"
	"github.com/go-atomci/atomci/utils/errors"

	"github.com/astaxie/beego/orm"
)

// GetArrange ...
func (manager *AppManager) GetArrange(ProjectAppID, envID int64) (*AppArrangeResp, error) {
	arrange, err := manager.model.GetAppArrange(ProjectAppID, envID)
	if err != nil || arrange.Config == "" {
		if err != orm.ErrNoRows && err != nil {
			return nil, errors.NewInternalServerError().SetCause(err)
		}
		return nil, nil
	}
	imageMapings := make([]ImageMaping, 0)

	appImageMappings, err := manager.model.GetAppImageMappingByArrangeID(arrange.ID)
	if err != nil {
		log.Log.Error("when get arrange, get app image mapping error: %s", err.Error())
		return nil, err
	}
	for _, item := range appImageMappings {
		item := ImageMaping{
			ID:           item.ID,
			Name:         item.Name,
			Image:        item.Image,
			ProjectAppID: item.ProjectAppID,
			ImageTagType: item.ImageTagType,
			ArrangeID:    item.ArrangeID,
		}
		imageMapings = append(imageMapings, item)
	}
	return &AppArrangeResp{
		ID:           arrange.ID,
		EnvID:        arrange.EnvID,
		ProjectAppID: arrange.ProjectAppID,
		Config:       arrange.Config,
		ImageMapings: imageMapings,
	}, nil
}

// GetRealArrange do not return template
func (manager *AppManager) GetRealArrange(appID, envID int64) (*models.AppArrange, error) {
	arrange, err := manager.model.GetAppArrange(appID, envID)
	if err != nil {
		log.Log.Error("get app arrange app id: %v, env id: %s occur error: %s", appID, envID, err.Error())
		return nil, err
	}

	if arrange.Config == "" {
		return nil, fmt.Errorf("app id: %v  env id: %v arrange did not setup", appID, envID)
	}
	return arrange, nil
}

// SetArrange ...
func (manager *AppManager) SetArrange(
	projectAppID int64,
	arrangeEnvID int64,
	request *AppArrangeReq,
) error {
	_, err := manager.projectModel.GetProjectApp(projectAppID)
	if err != nil {
		if err == orm.ErrNoRows {
			return errors.NewNotFound().SetCause(err)
		}
		return errors.NewInternalServerError().SetCause(err)
	}
	request.CopyToEnvIDs = append(request.CopyToEnvIDs, arrangeEnvID)
	if len(request.CopyToEnvIDs) > 0 {
		for _, item := range request.CopyToEnvIDs {
			apparrangeModel := genrateAppArrangeModel(projectAppID, item, request.Config)
			// create or update arrange with the config
			id, err := manager.createOrUpdateAppConfig(apparrangeModel)
			if err != nil {
				return err
			}
			log.Log.Debug("update app arrnage item id: %v", id)
			validItemIDs := []int64{}
			// TODO: add verify, one arrange only support add one image mapping item.
			for _, imageMappingitem := range request.ImageMapings {
				if imageMappingitem.ProjectAppID == 0 {
					log.Log.Debug("item: %v did not join with project app: 0, skip", imageMappingitem.Name)
					if err := manager.deleteAppImageMapping(imageMappingitem.ArrangeID, imageMappingitem.Image); err != nil {
						log.Log.Error("delete origin app image mapping error: %s", err.Error())
					}
					continue
				}
				appImageMappingModel := generateAppMappingModel(id, imageMappingitem)
				if imageMappingitem.ID == 0 {
					mappingID, err := manager.createAppMapping(appImageMappingModel)
					if err != nil {
						log.Log.Error("create app mapping occur error: %v", err.Error())
						return err
					}
					validItemIDs = append(validItemIDs, mappingID)
				} else {
					validItemIDs = append(validItemIDs, imageMappingitem.ID)
					item, err := manager.model.GetAppImageMappingItemByID(imageMappingitem.ID)
					if err != nil {
						log.Log.Error("get app mapping by %v occur error: %v", id, err.Error())
						return err
					}
					appImageMappingModel.Addons = item.Addons
					err = manager.model.UpdateAppImageMapping(&appImageMappingModel)
					if err != nil {
						log.Log.Error("update app mapping by %v occur error: %v", id, err.Error())
						return err
					}
				}
			}

			invalidItems, err := manager.model.GetInvalidAppImageMappingItems(id, projectAppID, validItemIDs)
			if err != nil {
				return err
			}
			for _, invalidItem := range invalidItems {
				if err := manager.model.DeleteAppImageMapping(invalidItem); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (manager *AppManager) createOrUpdateAppConfig(newArrange models.AppArrange) (int64, error) {
	// create or update arrange with the config
	oldArrange, err := manager.model.GetAppArrange(newArrange.ProjectAppID, newArrange.EnvID)
	if err == nil {
		newArrange.Addons = oldArrange.Addons
		err = manager.model.UpdateAppArrange(&newArrange)
	} else if err == orm.ErrNoRows {
		newArrange.Addons = models.NewAddons()
		err = manager.model.InsertAppArrange(&newArrange)
		if err != nil {
			log.Log.Error("when insert occur error: %s", err.Error())
		}
	} else {
		log.Log.Error("get app arrange error: %s", err.Error())
	}
	return newArrange.ID, err
}

func (manager *AppManager) createAppMapping(imageMappingItem models.AppImageMapping) (int64, error) {
	// create or update arrange with the config
	imageMappingItem.Addons = models.NewAddons()
	id, err := manager.model.InsertAppImageMapping(&imageMappingItem)
	if err != nil {
		log.Log.Error("when insert occur error: %s", err.Error())
	}

	return id, err
}

func (manager *AppManager) deleteAppImageMapping(arrangeID int64, image string) error {
	imageMapping, err := manager.model.GetAppImageMappingItemByImage(arrangeID, image)
	if err != nil {
		return err
	}
	return manager.model.DeleteAppImageMapping(imageMapping)
}

func genrateAppArrangeModel(appID, envID int64, config string) models.AppArrange {
	return models.AppArrange{
		EnvID:        envID,
		ProjectAppID: appID,
		Config:       config,
	}
}

func generateAppMappingModel(arrangeID int64, imageMapping ImageMaping) models.AppImageMapping {
	return models.AppImageMapping{
		ArrangeID:    arrangeID,
		Image:        imageMapping.Image,
		ProjectAppID: imageMapping.ProjectAppID,
		ImageTagType: imageMapping.ImageTagType,
		Name:         imageMapping.Name,
	}
}
