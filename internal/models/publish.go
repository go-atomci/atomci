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

package models

import (
	"time"

	"github.com/go-atomci/atomci/utils/query"
)

// const publish-order status defined
const (
	Failed           = 0
	Success          = 1
	Running          = 2
	Pending          = 3
	END              = 4
	Closed           = 5
	UnKnown          = 6
	TerminateSuccess = 7
	TerminateFailed  = 8
	MergeFailed      = 9
	NotSupport       = 10
	Skipped          = -1
)

// const publish-order step defined based on pipeline steps
const (
	StepVerify = "verify"
	StepManual = "manual"
	StepBuild  = "build"
	StepDeploy = "deploy"
)

// ProejctReleaseFilterQuery ..
type ProejctReleaseFilterQuery struct {
	query.FilterQuery
	VersionNo     string `json:"versionNo"`
	Name          string `json:"name"`
	Step          string `json:"step"`
	Stage         int    `json:"stage"`
	Status        *int64 `json:"status"`
	Creator       string `json:"creator"`
	CreateAtStart string `json:"createAtStart"`
	CreateAtEnd   string `json:"createAtEnd"`
}

// PublishOperation ..
type PublishOperation struct {
	Manual      bool `json:"manual"`
	Build       bool `json:"build"`
	Terminate   bool `json:"terminate"`
	Deploy      bool `json:"deploy"`
	MergeBranch bool `json:"merge-branch"`

	NextStage bool `json:"next-stage"`
	BackTo    bool `json:"back-to"`
}

// Publish ..
type Publish struct {
	Addons
	StartAt                time.Time         `orm:"column(start_at);auto_now;type(datetime);null" json:"start_at"`
	EndAt                  *time.Time        `orm:"column(end_at);type(datetime);null" json:"end_at"`
	Name                   string            `orm:"column(name);size(65)" json:"name"`
	Creator                string            `orm:"column(creator);size(64)" json:"creator"`
	ProjectID              int64             `orm:"column(project_id)" json:"project_id"`
	StageID                int64             `orm:"column(stage_id)" json:"stage_id"`
	StageName              string            `orm:"column(stage_name);size(128)" json:"stage_name"`
	Step                   string            `orm:"column(step);size(128)" json:"step"`
	StepType               string            `orm:"column(step_type);size(64)" json:"step_type"`
	StepIndex              int               `orm:"column(step_index);size(64)" json:"step_index"`
	Status                 int64             `orm:"column(status)" json:"status"`
	PipelineID             int64             `orm:"column(pipeline_id)" json:"pipeline_id"`
	LastPipelineInstanceID int64             `orm:"column(last_pipeline_instance_id)" json:"last_pipeline_instance_id"`
	VersionNo              string            `orm:"column(version_no);size(64)" json:"version_no"`
	Operations             *PublishOperation `orm:"-" json:"operations"`
	NextStep               string            `orm:"-" json:"next_step"`
	Previous               string            `orm:"-" json:"previous"`
}

// TableName  ..
func (t *Publish) TableName() string {
	return "pub_publish"
}

// PublishApp ..
type PublishApp struct {
	Addons
	PublishID      int64  `orm:"column(publish_id)" json:"publish_id"`
	ProjectAppID   int64  `orm:"column(project_app_id)" json:"project_app_id"`
	BranchName     string `orm:"column(branch_name);size(64)" json:"branch_name"`
	CompileCommand string `orm:"column(compile_command);size(1024)" json:"compile_command"`
}

// TableName ...
func (t *PublishApp) TableName() string {
	return "pub_publish_app"
}

// PublishOperationLog ..
type PublishOperationLog struct {
	Addons
	PipelineInstanceID int64  `orm:"column(pipeline_instance_id)" json:"pipeline_instance_id"`
	StageInstanceID    int64  `orm:"column(stage_instance_id)" json:"stage_instance_id"`
	StepIndex          int    `orm:"column(step_index)" json:"step_index"`
	PublishID          int64  `orm:"column(publish_id)" json:"publish_id"`
	Creator            string `orm:"column(creator);size(64)" json:"creator"`
	Type               string `orm:"column(type);size(128);null" json:"type"`
	StageID            int64  `orm:"column(stage_id)" json:"stage_id"`
	Stage              string `orm:"column(stage);size(128)" json:"stage"`
	Step               string `orm:"column(step);size(128)" json:"step"`
	Status             int64  `orm:"column(status);" json:"status"`
	RunID              int64  `orm:"column(run_id);" json:"run_id"`
	JobName            string `orm:"column(job_name);size(128)" json:"job_name"`
	Code               string `orm:"column(code);size(128)" json:"code"`
	Message            string `orm:"column(message);size(256)" json:"message"`
}

// TableName ...
func (t *PublishOperationLog) TableName() string {
	return "pub_publish_operation"
}
