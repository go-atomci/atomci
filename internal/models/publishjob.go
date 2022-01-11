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

// PublishJob status const defined
const (
	StatusInit        = "INIT"
	StatusInitFailure = "INIT_FAILURE"
	StatusRunning     = "RUNNING"
	StatusFailure     = "FAILURE"
	StatusUnknown     = "UNKNOWN"
	StatusSuccess     = "SUCCESS"
	StatusAbort       = "ABORTED"
)

// publishjob job type
const (
	JobTypeBuild  = "build"
	JobTypeDeploy = "deploy"
)

// PublishJob ..
type PublishJob struct {
	Addons
	PublishID        int64  `orm:"column(publish_id)" json:"publish_id"`
	ProjectID        int64  `orm:"column(project_id)" json:"project_id"`
	Status           string `orm:"column(status);size(16)" json:"status"`
	RunID            int64  `orm:"column(run_id)" json:"run_id"`
	Progress         int    `orm:"column(progress)" json:"progress"`
	DurationInMillis int64  `orm:"column(duration_in_millis)" json:"duration_in_millis"`
	EnvID            int64  `orm:"column(stage_id)" json:"stage_id"`
	Operator         string `orm:"column(operator); size(64)" json:"operator"`
	JobType          string `orm:"column(job_type);size(64)" json:"job_type"`
}

// TableName ...
func (t *PublishJob) TableName() string {
	return "pub_publish_job"
}

// PublishJobApp ..
type PublishJobApp struct {
	Addons
	ProjectID    int64  `orm:"column(project_id)" json:"project_id"`
	PublishJobID int64  `orm:"column(publish_job_id)" json:"publish_job_id"`
	ProjectAPPID int64  `orm:"column(project_app_id)" json:"project_app_id"`
	BranchName   string `orm:"column(branch_name); size(64)" json:"branch_name"`
	BranchURL    string `orm:"column(branch_url); size(255)" json:"branch_url"`
	ImageAddr    string `orm:"column(image_addr);size(255)" json:"image_addr"`
	ImageVersion string `orm:"column(image_version);size(64)" json:"image_version"`
	Release      string `orm:"column(release);size(64)" json:"release"`
	Gray         bool   `orm:"column(gray)" json:"gray"`
}

// TableName ...
func (t *PublishJobApp) TableName() string {
	return "pub_publish_job_app"
}
