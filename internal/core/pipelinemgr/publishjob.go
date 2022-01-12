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

package pipelinemgr

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-atomci/atomci/internal/middleware/log"
	"github.com/go-atomci/atomci/internal/models"
	"github.com/go-atomci/atomci/utils"
)

// GetPublishStats ..
func (pm *PipelineManager) GetPublishStats(projectID int64, request *PublishStatsReq) ([]*PublishStatsRsp, error) {

	var startTime, endTime time.Time
	if request.StartTime > "" {
		startTime, _ = time.Parse("2006-01-2", request.StartTime)
	}
	if request.EndTime > "" {
		endTime, _ = time.Parse("2006-01-2", request.EndTime)
	}
	qs, err := pm.modelPublishJob.GetPublishJobByProjectIDFilters(projectID, request.AppIDs, request.EnvIDs)
	if err != nil {
		log.Log.Error("when GetPublishStats, GetPublishJobByProjectIDFilters occur error: %s", err.Error())
		return nil, err
	}

	emptyTime := time.Time{}
	days := 0
	if startTime != emptyTime && endTime != emptyTime {
		difference := endTime.Sub(startTime)
		days = int(difference.Hours() / 24)
	}

	log.Log.Debug("days: %v", days)
	rsp := []*PublishStatsRsp{}
	for i := 0; i <= days; i++ {
		oneDay := startTime.AddDate(0, 0, 1)
		oneQs := qs.Filter("create_at__gt", startTime).
			Filter("create_at__lt", oneDay)
		log.Log.Debug("filter startTime: %v ", startTime)
		buildQs := oneQs.Filter("job_type", "build")
		deployQs := oneQs.Filter("job_type", "deploy")

		buildSuccess, _ := buildQs.Filter("status", models.StatusSuccess).Count()
		buildFailed, _ := buildQs.Filter("status", models.StatusFailure).Count()

		deploySuccess, _ := deployQs.Filter("status", models.StatusSuccess).Count()
		deployFailed, _ := deployQs.Filter("status", models.StatusFailure).Count()

		totalSuccess, _ := oneQs.Filter("status", models.StatusSuccess).Count()
		totalFailed, _ := oneQs.Filter("status", models.StatusFailure).Count()

		item := &PublishStatsRsp{
			BuildSuccess:  buildSuccess,
			BuildFailed:   buildFailed,
			DeploySuccess: deploySuccess,
			DeployFailed:  deployFailed,
			TotalSuccess:  totalSuccess,
			TotalFailed:   totalFailed,
			Time:          startTime.Format("2006-01-02"),
		}
		rsp = append(rsp, item)
		startTime = oneDay
	}

	return rsp, nil
}

// CreatePublishJob ..
func (pm *PipelineManager) CreatePublishJob(projectID, publishID, stageID int64,
	operator string, jobType string,
	allAppsParms []*AppParamsForCreatePublishJob) (int64, error) {
	publishJob := &models.PublishJob{
		Operator:  operator,
		ProjectID: projectID,
		PublishID: publishID,
		EnvID:     stageID,
		Status:    models.StatusInit,
		JobType:   jobType,
	}
	id, err := pm.modelPublishJob.CreatePublishJobifNotExist(publishJob)
	if err != nil {
		return 0, err
	}
	for _, app := range allAppsParms {
		publishJobApp := &models.PublishJobApp{
			ProjectID:    projectID,
			PublishJobID: id,
			ProjectAPPID: app.ProjectAppID,
			BranchName:   app.Branch,
			BranchURL:    app.Path,
			ImageVersion: app.ImageVersion,
			Gray:         app.Gray,
			ImageAddr:    app.ImageAddr,
		}
		_, err := pm.modelPublishJob.CreateJobAppIfNotExist(publishJobApp)
		if err != nil {
			// TODO: add transaction processing
			log.Log.Error("crate publish id %d job app  occur error: %s", id, err)
			return 0, err
		}
	}
	return id, nil
}

// UpdatePublishJob update job runID into publishjob item
func (pm *PipelineManager) UpdatePublishJob(publishJobID, runID int64) error {
	modelPublishJob, err := pm.modelPublishJob.GetPublishJobByID(publishJobID)
	if err != nil {
		log.Log.Error("when update publishJob RunID, get publish job id occur error: %s", err)
		return err
	}
	if runID != 0 {
		modelPublishJob.RunID = runID
		modelPublishJob.Status = models.StatusRunning
		modelPublishJob.Progress = 10
	}

	err = pm.modelPublishJob.UpdatePublishJob(modelPublishJob)
	if err != nil {
		log.Log.Error("udpate publishjob runID occur error: %s", err)
		return err
	}
	return nil
}

// UpdatePublishJobStatus ..
func (pm *PipelineManager) UpdatePublishJobStatus(publishJobID int64, status string) error {
	publishJob, err := pm.modelPublishJob.GetPublishJobByID(publishJobID)
	if err != nil {
		log.Log.Error("when update publish job status, get publish job by id occur error: %s", err.Error())
		return fmt.Errorf("网络错误，请重试")
	}
	statusUpper := strings.ToUpper(status)
	jobEndStatus := []string{"SUCCESS", "INIT_FAILURE", "FAILURE"}
	if utils.Contains(jobEndStatus, publishJob.Status) {
		return fmt.Errorf("update publish job status: %v already was end status, skipped", publishJob.Status)
	}
	switch statusUpper {
	case "SUCCESS":
		publishJob.Status = statusUpper
	case "FAILURE":
		publishJob.Status = statusUpper
	default:
		log.Log.Error("status: %v, unexception, reset to FAILURE")
		publishJob.Status = "FAILURE"
	}
	return pm.modelPublishJob.UpdatePublishJob(publishJob)
}
