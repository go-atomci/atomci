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
	"github.com/go-atomci/atomci/internal/models"

	"github.com/astaxie/beego/orm"
)

// K8sClusterModel ...
type K8sClusterModel struct {
	ormer                orm.Ormer
	ApplicationTableName string
	NamespaceTableName   string
}

// NewK8sClusterModel ...
func NewK8sClusterModel() (model *K8sClusterModel) {
	return &K8sClusterModel{
		ormer:                GetOrmer(),
		ApplicationTableName: (&models.CaasApplication{}).TableName(),
	}
}

// GetApplicationsByProjectID ...
func (model *K8sClusterModel) GetApplicationsByProjectID(projectID int64) ([]*models.CaasApplication, error) {
	apps := []*models.CaasApplication{}
	qs := model.ormer.QueryTable(model.ApplicationTableName).
		Filter("project_id", projectID).
		Filter("deleted", false)
	_, err := qs.All(&apps)
	return apps, err
}

// GetApplication ...
func (model *K8sClusterModel) GetApplication(cluster, department, svcName string) ([]*models.CaasApplication, error) {
	apps := []*models.CaasApplication{}
	qs := model.ormer.QueryTable(model.ApplicationTableName).
		Filter("cluster", cluster).
		Filter("namespace", department).
		Filter("name", svcName).
		Filter("deleted", false)
	_, err := qs.All(&apps)
	if err != nil {
		return nil, err
	}
	return apps, err
}
