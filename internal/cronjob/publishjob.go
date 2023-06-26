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

package cronjob

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-atomci/atomci/internal/core/pipelinemgr"
	"github.com/go-atomci/atomci/internal/dao"
	"github.com/go-atomci/atomci/internal/middleware/log"
	"github.com/go-atomci/atomci/internal/models"
	"github.com/go-atomci/atomci/utils"
	"github.com/go-atomci/workflow"

	"github.com/go-atomci/workflow/jenkins"
)

// RunPublishJobServer ..
func RunPublishJobServer() {
	go func() {
		for {
			syncAllPublishJobStatus()
			time.Sleep(time.Minute * 2)
		}
	}()
}

func syncAllPublishJobStatus() {
	log.Log.Debug("start sync publish job status..")
	newPublishJob := dao.NewPublishJobModel()
	newPublish := dao.NewPublishModel()
	publishJobs := getRunningPublishJob(newPublishJob)
	pipeline := pipelinemgr.NewPipelineManager()

	for _, job := range publishJobs {
		if job.RunID == 0 {
			if err := justUpdateModelPublishJobStatus(job, models.StatusInitFailure, newPublishJob); err != nil {
				log.Log.Error("sync publish job id: %d, run id: %d occur error: %s", job.ID, job.RunID, err.Error())
			}
			continue
		}
		if err := updatePublishJobStatus(job, newPublishJob, newPublish, pipeline); err != nil {
			log.Log.Error("sync publish job id: %d, run id: %d, occur error: %s", job.ID, job.RunID, err.Error())
			continue
		}
	}
}

func getRunningPublishJob(newPublishJob *dao.PublishJobModel) []*models.PublishJob {
	publishJobs, err := newPublishJob.GetPublishJobsByFilter([]string{models.StatusRunning, models.StatusInit}, []string{models.JobTypeBuild, models.JobTypeDeploy})
	if err != nil {
		log.Log.Error("when sync publish job, get publish jobs occur error: %s", err.Error())
		return nil
	}
	return publishJobs
}
func justUpdateModelPublishJobStatus(job *models.PublishJob, status string, newPublishJob *dao.PublishJobModel) error {
	job.Status = status
	return newPublishJob.UpdatePublishJob(job)
}

func updatePublishJobStatus(job *models.PublishJob, newPublishJob *dao.PublishJobModel,
	newPublish *dao.PublishModel, pipeline *pipelinemgr.PipelineManager) error {
	log.Log.Info("sync publish job: %d, runID: %d", job.ID, job.RunID)
	switch job.JobType {
	case models.JobTypeBuild:
		jobName := fmt.Sprintf("atomci_%v_%v_%v", job.ProjectID, job.PublishID, job.EnvID)
		var publishStatus int
		var err error
		job, publishStatus, err = getPipelineJobStatus(jobName, job, pipeline)
		if err != nil {
			return err
		}
		// publish Order update
		updatePublishOrderStatus(job.PublishID, publishStatus, newPublish)
		return newPublishJob.UpdatePublishJob(job)
	case models.JobTypeDeploy:
		jobName := fmt.Sprintf("atomci_%v_%v", job.ProjectID, job.EnvID)
		var publishStatus int
		var err error
		job, publishStatus, err = getPipelineJobStatus(jobName, job, pipeline)
		if err != nil {
			return err
		}
		// publish Order update
		updatePublishOrderStatus(job.PublishID, publishStatus, newPublish)
		return newPublishJob.UpdatePublishJob(job)

	default:
		log.Log.Error("publish job type: %v is not support currently", job.JobType)
	}
	return nil
}

func getPipelineJobStatus(jobName string, job *models.PublishJob, pipeline *pipelinemgr.PipelineManager) (*models.PublishJob, int, error) {
	jenkinsInfo, err := pipeline.GetCIConfig(job.EnvID)
	if err != nil {
		log.Log.Error("get Jenkins Config occur error: %s", err.Error())
		return nil, 0, err
	}
	addr, user, token := jenkinsInfo[0], jenkinsInfo[1], jenkinsInfo[2]

	workFlowProvider, err := jenkins.NewJenkinsClient(
		jenkins.URL(addr),
		jenkins.JenkinsUser(user),
		jenkins.JenkinsToken(token),
		jenkins.JenkinsJob(jobName),
	)

	if err != nil {
		log.Log.Error("create workflow Client occur error: %s", err.Error())
		return nil, 0, err
	}

	jobDetail, err := workFlowProvider.GetJobInfo(job.RunID)
	if err != nil {
		if strings.Contains(err.Error(), "404 not found") {
			// when job 404/ set jobDetail default value.
			jobDetail = &workflow.JobInfo{
				Result:         models.StatusUnknown,
				DurationMillis: 0,
				Stages:         []workflow.Stage{},
			}
		} else {
			return nil, 0, err
		}
	}

	// get job result and execute time
	status, durationInMillis := jobDetail.Result, jobDetail.DurationMillis
	if status == "" {
		status = jobDetail.Status
	}

	// get processes
	jobStages := jobDetail.Stages
	log.Log.Debug("jenkins job pipeline stage length: %v", len(jobStages))
	finishedStages := 0
	stageLength := 5
	for _, stage := range jobStages {
		if stage.Status == "SUCCESS" {
			finishedStages++
		}
	}
	if len(jobStages) < 2 {
		finishedStages = 1
	} else {
		stageLength = len(jobStages) + 4
		finishedStages = finishedStages + 4
	}
	log.Log.Debug("current finishedStages: %v, stageLength: %v", finishedStages, stageLength)
	progress := finishedStages * 100 / stageLength

	var jobStatus string
	publishStatus := models.Running
	switch status {
	case "SUCCESS":
		jobStatus = models.StatusSuccess
		publishStatus = models.Success
	case "FAILURE", "FAILED":
		jobStatus = models.StatusFailure
		publishStatus = models.Failed
	case "ABORTED":
		jobStatus = models.StatusAbort
		publishStatus = models.TerminateSuccess
	case "IN_PROGRESS":
		jobStatus = models.StatusRunning
	default:
		log.Log.Warn("job's status is undesired value: %v, reset job status to: 'UNKNOWN'", status)
		jobStatus = models.StatusUnknown
		publishStatus = models.UnKnown
	}
	if job.Progress < progress {
		job.Progress = progress
	}
	job.DurationInMillis = int64(durationInMillis)
	job.Status = jobStatus
	return job, publishStatus, nil
}

func updatePublishOrderStatus(publishID int64, publishStatus int, newPublish *dao.PublishModel) {
	// publish Order update
	modelPublishItem, err := newPublish.GetPublishByID(publishID)
	if err != nil {
		log.Log.Error("when get publishOrder model, occur error: %s", err.Error())
	}
	if utils.Contains([]string{"build", "deploy"}, modelPublishItem.StepType) && publishStatus != models.Running {
		modelPublishItem.Status = int64(publishStatus)
		err := newPublish.UpdatePublish(modelPublishItem)
		if err != nil {
			log.Log.Error("when update publishOrder status, occur error: %s", err.Error())
		}
		// create publish operation log
		if err := createPublishOperationLog(modelPublishItem, int64(publishStatus), newPublish); err != nil {
			log.Log.Error("after update publish order status, create publish operation log occur error: %s", err.Error())
		}
	} else {
		log.Log.Warn("current publishOrder id: %v 's stepType %v is not %v, Or publishStaus is not running, skip status update", publishID, modelPublishItem.StepType, fmt.Sprintf("%v, %v", models.StepBuild, models.StepDeploy))
	}
}

func createPublishOperationLog(publish *models.Publish, status int64, newPublish *dao.PublishModel) error {
	operationLog := &models.PublishOperationLog{
		Creator:   "system",
		Stage:     publish.StageName,
		StageID:   publish.StageID,
		Step:      publish.Step,
		Message:   "",
		Status:    status,
		PublishID: publish.ID,
	}
	if err := newPublish.CreatePublishOperation(operationLog); err != nil {
		return err
	}
	return nil
}
