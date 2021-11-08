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

	"github.com/go-atomci/atomci/models"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// PublishJobModel ...
type PublishJobModel struct {
	ormer                  orm.Ormer
	publishJobTableName    string
	publishJobAppTableName string
}

// NewPublishJobModel ...
func NewPublishJobModel() (model *PublishJobModel) {
	return &PublishJobModel{
		ormer:                  GetOrmer(),
		publishJobTableName:    (&models.PublishJob{}).TableName(),
		publishJobAppTableName: (&models.PublishJobApp{}).TableName(),
	}
}

// GetPublishJobByStageID ..
func (model *PublishJobModel) GetPublishJobByStageID(projectID, stageID int64) ([]*models.PublishJob, error) {
	publishJobsModel := []*models.PublishJob{}
	_, err := model.ormer.QueryTable(model.publishJobTableName).
		Filter("project_id", projectID).
		Filter("stage_id", stageID).
		OrderBy("-create_at").
		Limit(10).
		Filter("Deleted", false).All(&publishJobsModel)
	return publishJobsModel, err
}

// CreatePublishJobifNotExist ...
func (model *PublishJobModel) CreatePublishJobifNotExist(publishJob *models.PublishJob) (int64, error) {
	id, err := model.ormer.Insert(publishJob)
	return id, err
}

// UpdatePublishJob ...
func (model *PublishJobModel) UpdatePublishJob(publishjob *models.PublishJob) error {
	_, err := model.ormer.Update(publishjob)
	return err
}

// GetPublishJobByProjectIDFilters ..
func (model *PublishJobModel) GetPublishJobByProjectIDFilters(projectID int64, appIDs, envIDs []int64) (orm.QuerySeter, error) {
	publishJobIDs := []int64{}
	qs := model.ormer.QueryTable(model.publishJobTableName).Filter("project_id", projectID).Filter("Deleted", false)
	if len(appIDs) > 0 {
		publishJobIDs = []int64{}
		jobApps := []*models.PublishJobApp{}
		_, err := model.ormer.QueryTable(model.publishJobAppTableName).Filter("project_app_id__in", appIDs).
			Filter("Deleted", false).All(&jobApps)
		if err != nil {
			logs.Error("when GetPublishJobByProjectIDFilters, get publishJobApp occur error: %s", err.Error())
			return nil, err
		}
		for _, app := range jobApps {
			publishJobIDs = append(publishJobIDs, app.PublishJobID)
		}
		logs.Debug("publishJobIDs: %v", publishJobIDs)
		if len(publishJobIDs) == 0 {
			publishJobIDs = []int64{-1}
		}
		qs = qs.Filter("id__in", publishJobIDs)
	}

	if len(envIDs) > 0 {
		qs = qs.Filter("stage_id__in", envIDs)
	}

	return qs, nil
}

// GetPublishJobByID ..
func (model *PublishJobModel) GetPublishJobByID(ID int64) (*models.PublishJob, error) {
	publishJobModel := &models.PublishJob{}
	err := model.ormer.QueryTable(model.publishJobTableName).Filter("id", ID).Filter("Deleted", false).One(publishJobModel)
	return publishJobModel, err
}

// GetLastPublishJobByPublishID ..
func (model *PublishJobModel) GetLastPublishJobByPublishID(publishID int64) (*models.PublishJob, error) {
	publishJobModel := &models.PublishJob{}
	err := model.ormer.QueryTable(model.publishJobTableName).
		Filter("publish_id", publishID).
		OrderBy("-create_at").
		Filter("deleted", false).One(publishJobModel)
	return publishJobModel, err
}

// GetPublishJobsByFilter For PublishJob Serer sync publish/publish job status
func (model *PublishJobModel) GetPublishJobsByFilter(status []string, jobType []string) ([]*models.PublishJob, error) {
	publishJobsModel := []*models.PublishJob{}
	_, err := model.ormer.QueryTable(model.publishJobTableName).
		Filter("status__in", status).
		Filter("job_type__in", jobType).
		Filter("Deleted", false).All(&publishJobsModel)
	return publishJobsModel, err
}

// GetCurrentRunningBuildJob For Trigger publishOrder build verify, running job include init and running
func (model *PublishJobModel) GetCurrentRunningBuildJob(projectID, stageID, publishID int64, status []string, jobType string) ([]*models.PublishJob, error) {
	publishJobsModel := []*models.PublishJob{}
	_, err := model.ormer.QueryTable(model.publishJobTableName).
		Filter("status__in", status).
		Filter("stage_id", stageID).
		Filter("job_type", jobType).
		Filter("project_id", projectID).
		Filter("publish_id", publishID).
		Filter("Deleted", false).All(&publishJobsModel)
	return publishJobsModel, err
}

// GetCurrentRunningJob For Trigger publishOrder publish verify, running job include init and running
func (model *PublishJobModel) GetCurrentRunningJob(projectID, stageID int64, status []string, jobType string) ([]*models.PublishJob, error) {
	publishJobsModel := []*models.PublishJob{}
	_, err := model.ormer.QueryTable(model.publishJobTableName).
		Filter("status__in", status).
		Filter("stage_id", stageID).
		Filter("job_type", jobType).
		Filter("project_id", projectID).
		Filter("Deleted", false).All(&publishJobsModel)
	return publishJobsModel, err
}

/* --- PublishJob APP Part --- */

// CreateJobAppIfNotExist ...
func (model *PublishJobModel) CreateJobAppIfNotExist(app *models.PublishJobApp) (int64, error) {
	created, id, err := model.ormer.ReadOrCreate(app, "publish_job_id", "project_app_id", "deleted")
	if err == nil {
		if !created {
			err = fmt.Errorf(fmt.Sprintf("app id: %v existed in publishJob", app.ProjectAPPID))
		}
	}
	return id, err
}

// GetPublishJobApps ..
func (model *PublishJobModel) GetPublishJobApps(publishJobID int64) ([]*models.PublishJobApp, error) {
	jobAppsModel := []*models.PublishJobApp{}
	_, err := model.ormer.QueryTable(model.publishJobAppTableName).Filter("publish_job_id", publishJobID).Filter("Deleted", false).All(&jobAppsModel)
	return jobAppsModel, err
}

// GetPublishJobApp ..
func (model *PublishJobModel) GetPublishJobApp(publishJobID, AppID int64) (*models.PublishJobApp, error) {
	jobAppModel := &models.PublishJobApp{}
	err := model.ormer.QueryTable(model.publishJobAppTableName).Filter("Deleted", false).
		Filter("publish_job_id", publishJobID).
		Filter("project_app_id", AppID).One(jobAppModel)
	return jobAppModel, err
}

// UpdatePublishJobApp ...
func (model *PublishJobModel) UpdatePublishJobApp(publishJobApp *models.PublishJobApp) error {
	_, err := model.ormer.Update(publishJobApp)
	return err
}

// GetPublishJobAppByID ..
func (model *PublishJobModel) GetPublishJobAppByID(ID int64) (*models.PublishJobApp, error) {
	JobAppModel := &models.PublishJobApp{}
	err := model.ormer.QueryTable(model.publishJobAppTableName).Filter("id", ID).Filter("Deleted", false).One(JobAppModel)
	return JobAppModel, err
}
