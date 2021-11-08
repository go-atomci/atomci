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
	"time"

	"github.com/go-atomci/atomci/models"
	"github.com/go-atomci/atomci/utils/query"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

var publishEnableFilterKeys = []string{
	"id",
	"name",
	"description",
	"creator",
}

var publishOpertaionEnableFilterKeys = []string{
	"creator",
	"stage",
	"step",
	"message",
}

// PublishModel ...
type PublishModel struct {
	ormer                     orm.Ormer
	publishTableName          string
	publishOpertaionTableName string
	publishAppTableName       string
	publishApplyTableName     string
}

// NewPublishModel ...
func NewPublishModel() (model *PublishModel) {
	return &PublishModel{
		ormer:                     GetOrmer(),
		publishTableName:          (&models.Publish{}).TableName(),
		publishOpertaionTableName: (&models.PublishOperationLog{}).TableName(),
		publishAppTableName:       (&models.PublishApp{}).TableName(),
	}
}

// GetPublishesList ...
func (model *PublishModel) GetPublishesList(projectID int64, filter *models.ProejctReleaseFilterQuery) (*query.QueryResult, []*models.Publish, error) {
	rst := &query.QueryResult{Item: []*models.Publish{}}

	ormCond := orm.NewCondition()
	ormCond = ormCond.And("project_id", projectID).And("deleted", false)

	if filter.Name != "" {
		ormCond = ormCond.And("name__icontains", filter.Name)
	}
	if filter.VersionNo != "" {
		ormCond = ormCond.And("version_no__icontains", filter.VersionNo)
	}
	if filter.Creator != "" {
		ormCond = ormCond.And("creator__icontains", filter.Creator)
	}

	if filter.Status != nil {
		ormCond = ormCond.And("status", filter.Status)
	}

	if filter.Stage != 0 {
		ormCond = ormCond.And("stage_id", filter.Stage)
	}

	if filter.Step != "" {
		ormCond = ormCond.And("step__icontains", filter.Step)
	}
	if filter.CreateAtStart != "" && filter.CreateAtEnd != "" {
		if createAtStart, err := time.Parse("2006-01-02", filter.CreateAtStart); err == nil {
			ormCond = ormCond.And("create_at__gte", createAtStart)
		} else {
			logs.Error("time parse error: %s", err.Error())
		}
		if createAtEnd, err := time.Parse("2006-01-02", filter.CreateAtEnd); err == nil {
			ormCond = ormCond.And("create_at__lte", createAtEnd)
		} else {
			logs.Error("time parse error: %s", err.Error())
		}
	}

	qs := model.ormer.QueryTable(model.publishTableName).OrderBy("-create_at").SetCond(ormCond)
	count, err := qs.Count()
	if err != nil {
		return nil, nil, err
	}
	if err = query.FillPageInfo(rst, filter.PageIndex, filter.PageSize, int(count)); err != nil {
		return nil, nil, err
	}

	publishList := []*models.Publish{}
	_, err = qs.Limit(filter.PageSize, filter.PageSize*(filter.PageIndex-1)).All(&publishList)
	if err != nil {
		return nil, nil, err
	}
	rst.Item = publishList

	return rst, publishList, nil
}

// GetRunninbPublishesByProjectID ..
func (model *PublishModel) GetRunninbPublishesByProjectID(projectID int64) ([]*models.Publish, error) {
	publishes := []*models.Publish{}
	queryCond := orm.NewCondition().And("project_id", projectID).And("deleted", false)

	// 获取运行中的版本
	queryCond = queryCond.AndCond(orm.NewCondition().AndNot("status__in", []int64{models.Closed, models.END}))
	qs := model.ormer.QueryTable(model.publishTableName).SetCond(queryCond)
	_, err := qs.All(&publishes)
	return publishes, err
}

// GetPublishReleasesByProjectID ..
func (model *PublishModel) GetPublishReleasesByProjectID(projectID int64) (interface{}, error) {
	var maps []orm.Params
	_, err := model.ormer.Raw("select stage_name as 'env', count(id) as 'count' from pub_publish where `project_id` = ? and deleted = 0 group by stage_name", projectID).Values(&maps)
	return maps, err
}

// GetPublishByID ...
func (model *PublishModel) GetPublishByID(publishID int64) (*models.Publish, error) {
	publish := models.Publish{}
	qs := model.ormer.QueryTable(model.publishTableName).Filter("deleted", false)
	if publishID != -1 {
		qs = qs.Filter("id", publishID)
	}
	err := qs.One(&publish)
	return &publish, err
}

// GetPublishByPipelineInstanceID ...
func (model *PublishModel) GetPublishByPipelineInstanceID(pipelineInstanceID int64) (*models.Publish, error) {
	publish := models.Publish{}
	qs := model.ormer.QueryTable(model.publishTableName).Filter("deleted", false).
		Filter("last_pipeline_instance_id", pipelineInstanceID)

	err := qs.One(&publish)
	return &publish, err
}

// GetPublishUnEndBystatus ...
func (model *PublishModel) GetPublishUnEndBystatus(status []int64, projectIDs []int64) ([]*models.Publish, error) {
	publishes := []*models.Publish{}
	_, err := model.ormer.QueryTable(model.publishTableName).
		Filter("deleted", false).
		Filter("project_id__in", projectIDs).
		Exclude("status__in", status).All(&publishes)
	return publishes, err
}

// CreatePublishifNotExist ...
func (model *PublishModel) CreatePublishifNotExist(publish *models.Publish) (int64, error) {
	created, id, err := model.ormer.ReadOrCreate(publish, "version_no", "name", "deleted", "project_id")
	if err == nil {
		if !created {
			err = fmt.Errorf(fmt.Sprintf("publish name existed:%s", publish.Name))
		}
	}
	return id, err
}

// UpdatePublish ...
func (model *PublishModel) UpdatePublish(publish *models.Publish) error {
	_, err := model.ormer.Update(publish)
	return err
}

// DeletePublish ...
func (model *PublishModel) DeletePublish(publishID int64) error {
	publish, err := model.GetPublishByID(publishID)
	if err != nil {
		return err
	}
	publish.MarkDeleted()
	_, err = model.ormer.Delete(publish)
	return err
}

// GetOperationLogByInstanceIDAndStageIDStepType ...
func (model *PublishModel) GetOperationLogByInstanceIDAndStageIDStepType(instanceID, stageID int64, stepType int) ([]*models.PublishOperationLog, error) {
	operationLogs := []*models.PublishOperationLog{}
	_, err := model.ormer.QueryTable(model.publishOpertaionTableName).
		Filter("deleted", false).
		OrderBy("-create_at").
		Filter("pipeline_instance_id", instanceID).
		Filter("stage_id", stageID).
		Filter("step_index", stepType).All(&operationLogs)
	return operationLogs, err
}

// GetOperationLogsByPublishID ...
func (model *PublishModel) GetOperationLogsByPublishID(publishID int64, filter *query.FilterQuery) (*query.QueryResult, error) {
	rst := &query.QueryResult{Item: []*models.PublishOperationLog{}}
	queryCond := orm.NewCondition().AndCond(orm.NewCondition().And("deleted", false)).
		AndCond(orm.NewCondition().And("publish_id", publishID))

	if filterCond := query.FilterCondition(filter, filter.FilterKey); filterCond != nil {
		queryCond = queryCond.AndCond(filterCond)
	}
	qs := model.ormer.QueryTable(model.publishOpertaionTableName).OrderBy("-id").SetCond(queryCond)
	count, err := qs.Count()
	if err != nil {
		return nil, err
	}
	if err = query.FillPageInfo(rst, filter.PageIndex, filter.PageSize, int(count)); err != nil {
		return nil, err
	}

	logList := []*models.PublishOperationLog{}
	_, err = qs.Limit(filter.PageSize, filter.PageSize*(filter.PageIndex-1)).All(&logList)
	if err != nil {
		return nil, err
	}
	rst.Item = logList

	return rst, nil
}

// CreatePublishOperation ...
func (model *PublishModel) CreatePublishOperation(item *models.PublishOperationLog) error {
	_, err := model.ormer.InsertOrUpdate(item)
	return err
}

/* --- Publishes APP Part --- */

// CreatePublishAppIfNotExist ...
func (model *PublishModel) CreatePublishAppIfNotExist(app *models.PublishApp) (int64, error) {
	created, id, err := model.ormer.ReadOrCreate(app, "publish_id", "project_app_id", "deleted")
	if err == nil {
		if !created {
			err = fmt.Errorf(fmt.Sprintf("project app id: %v existed in publish", app.ProjectAppID))
		}
	}
	return id, err
}

// GetPublishAppByPublishIDAndAppID ..
func (model *PublishModel) GetPublishAppByPublishIDAndAppID(publishID, appID int64) (*models.PublishApp, error) {
	app := models.PublishApp{}
	qs := model.ormer.QueryTable(model.publishAppTableName).Filter("deleted", false)
	if publishID == 0 || appID == 0 {
		return nil, fmt.Errorf("publish_id or app_id must be a valid num, but current publish_id: %v app_id: %v", publishID, appID)
	}
	err := qs.Filter("publish_id", publishID).
		Filter("project_app_id", appID).One(&app)
	return &app, err
}

// UpdatePublishApp ...
func (model *PublishModel) UpdatePublishApp(publishApp *models.PublishApp) error {
	_, err := model.ormer.Update(publishApp)
	return err
}

// GetPublishAppsByID ..
func (model *PublishModel) GetPublishAppsByID(publishID int64) ([]*models.PublishApp, error) {
	apps := []*models.PublishApp{}
	qs := model.ormer.QueryTable(model.publishAppTableName).Filter("deleted", false)
	if publishID != -1 {
		qs = qs.Filter("publish_id", publishID)
	}
	_, err := qs.All(&apps)
	return apps, err
}

// GetPublishApp ...
func (model *PublishModel) GetPublishApp(publishAppID int64) (*models.PublishApp, error) {
	app := models.PublishApp{}
	qs := model.ormer.QueryTable(model.publishAppTableName).Filter("deleted", false)
	if publishAppID != -1 {
		qs = qs.Filter("id", publishAppID)
	}
	err := qs.One(&app)
	return &app, err
}

// DeletePublishApp ...
func (model *PublishModel) DeletePublishApp(publishAppID int64) error {
	app, err := model.GetPublishApp(publishAppID)
	if err != nil {
		return err
	}
	app.MarkDeleted()
	_, err = model.ormer.Delete(app)
	return err
}
