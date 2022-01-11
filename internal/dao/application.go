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

	"github.com/go-atomci/atomci/internal/models"
	util "github.com/go-atomci/atomci/utils"
	"github.com/go-atomci/atomci/utils/query"

	"github.com/astaxie/beego/orm"
)

var amLocker *util.SyncLocker

type AppModel struct {
	tOrmer    orm.Ormer
	TableName string
}

func init() {
	if amLocker == nil {
		amLocker = util.NewSyncLocker()
	}
}

func NewAppModel() *AppModel {
	return &AppModel{
		tOrmer:    GetOrmer(),
		TableName: (&models.CaasApplication{}).TableName(),
	}
}

func (am *AppModel) GetAppList(filterQuery *query.FilterQuery, projectID int64, cluster, namespace string) (*query.QueryResult, error) {
	rst := &query.QueryResult{Item: []models.CaasApplication{}}
	queryCond := orm.NewCondition().And("deleted", 0)

	if projectID == 0 {
		return nil, fmt.Errorf("project_id %v is invalid", projectID)
	}

	if cluster != "" {
		queryCond = queryCond.And("cluster", cluster)
	}
	if namespace != "" {
		queryCond = queryCond.And("namespace", namespace)
	}
	queryCond = queryCond.And("project_id", projectID)

	// TODO: queryCond need add.
	if filterQuery != nil {
		if filterQuery.FilterKey == "name" {
			queryCond = queryCond.And("name", filterQuery.FilterVal)
		}
	} else {
		// default filerQuery
		filterQuery.PageIndex = 1
		filterQuery.PageSize = 10
	}
	qs := am.tOrmer.QueryTable(am.TableName).OrderBy("-update_at").SetCond(queryCond)
	count, err := qs.Count()
	if err != nil {
		return nil, err
	}

	if err := query.FillPageInfo(rst, filterQuery.PageIndex, filterQuery.PageSize, int(count)); err != nil {
		return nil, err
	}

	appList := []models.CaasApplication{}
	_, err = qs.Limit(filterQuery.PageSize, filterQuery.PageSize*(filterQuery.PageIndex-1)).All(&appList)
	if err != nil {
		return nil, err
	}

	rst.Item = appList
	return rst, nil
}

func (am *AppModel) GetAppByName(cluster, namespace, name string) (*models.CaasApplication, error) {
	var app models.CaasApplication
	if err := am.tOrmer.QueryTable(am.TableName).
		Filter("name", name).
		Filter("cluster", cluster).
		Filter("namespace", namespace).
		Filter("deleted", 0).One(&app); err != nil {
		return nil, err
	}
	return &app, nil
}

func (am *AppModel) GetImage(cluster, namespace, name string) (string, error) {
	var app models.CaasApplication
	err := am.tOrmer.QueryTable(am.TableName).
		Filter("name", name).
		Filter("cluster", cluster).
		Filter("namespace", namespace).
		Filter("deleted", 0).One(&app, "image")
	return app.Image, err
}

func (am *AppModel) InsertApp(ins models.CaasApplication) error {
	_, err := am.tOrmer.Insert(&ins)
	return err
}

func (am *AppModel) CreateApp(ins models.CaasApplication) error {
	// TODO: verify cluster/namesapce is valid
	if am.AppExist(ins.Cluster, ins.Namespace, ins.Name) {
		return fmt.Errorf("app %s already exists in cluster %s", ins.Name, ins.Cluster)
	}
	ins.Addons = models.NewAddons()
	_, err := am.tOrmer.Insert(&ins)
	return err
}

func (am *AppModel) DeleteApp(app models.CaasApplication) error {
	_, err := am.tOrmer.Raw("UPDATE "+am.TableName+" SET deleted=1, delete_at=now() WHERE name=? AND cluster=? AND namespace=? AND deleted=0",
		app.Name, app.Cluster, app.Namespace).Exec()
	return err
}

func (am *AppModel) UpdateApp(ins *models.CaasApplication, updateTime bool) error {
	if ins == nil {
		return nil
	}
	amLocker.Lock(ins.Name)
	defer amLocker.Unlock(ins.Name)

	old, err := am.GetAppByName(ins.Cluster, ins.Namespace, ins.Name)
	if err != nil {
		return err
	}
	if old.UpdateAt != ins.UpdateAt {
		return fmt.Errorf("the application %s/%s/%s is updated by other routine", ins.Cluster, ins.Namespace, ins.Name)
	}

	// TODO: updatetime
	// if updateTime {
	// 	ins.Addons = ins.Addons.UpdateAddons()
	// } else {
	// 	ins.Addons = ins.Addons.FormatAddons()
	// }
	_, err = am.tOrmer.Update(ins)

	return err
}

func (am *AppModel) SetLabels(cluster, namespace, name, labels string) error {
	_, err := am.tOrmer.Raw("UPDATE "+am.TableName+" SET labels=? WHERE cluster=? AND namespace=? AND name=? AND deleted=0",
		labels, cluster, namespace, name).Exec()
	return err
}

func (am *AppModel) SetDeployStatus(cluster, namespace, name, status string) error {
	_, err := am.tOrmer.Raw("UPDATE "+am.TableName+" SET deploy_status=? WHERE name=? AND cluster=? AND namespace=? AND deleted=0",
		status, name, cluster, namespace).Exec()
	return err
}

func (am *AppModel) AppExist(cluster, namespace, name string) bool {
	return am.tOrmer.QueryTable(am.TableName).
		Filter("cluster", cluster).Filter("namespace", namespace).
		Filter("name", name).Filter("deleted", 0).Exist()
}
