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
	"encoding/json"
	"time"

	"github.com/go-atomci/atomci/internal/models"
)

// ResetPipelineReq ..
type ResetPipelineReq struct {
	CID int64 `json:"cID"`
}

// ManualStepReq ..
type ManualStepReq struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// RunBuildAppReq .
type RunBuildAppReq struct {
	Branch         string `json:"branch_name"`
	CompileCommand string `json:"compile_command"`
	ProjectAppID   int64  `json:"project_app_id"`
}

// BuildStepReq ..
type BuildStepReq struct {
	ActionName string            `json:"action_name,omitempty"`
	Apps       []*RunBuildAppReq `json:"apps,omitempty"`
	EnvVars    []EnvItem         `json:"env_vars,omitempty"`
}

// EnvItem env variable
type EnvItem struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

// BuildStepCallbackReq ..
type BuildStepCallbackReq struct {
	PublishJobID int64 `json:"publish_job_id"`
}

// RunDeployAppReq .
type RunDeployAppReq struct {
	ProjectAppID int64 `json:"project_app_id"`
	Gray         bool  `json:"gray"`
}

// DeployStepReq ..
type DeployStepReq struct {
	ActionName string             `json:"action_name"`
	Apps       []*RunDeployAppReq `json:"apps"`
}

// WeeklyDenyList ..
type WeeklyDenyList []*struct {
	StartTime string `json:"start_time"`
	Enable    bool   `json:"enable"`
	Name      string `json:"name"`
	EndTime   string `json:"end_time"`
}

// DenyEnvItem ..
type DenyEnvItem struct {
	Name   string `json:"name"`
	Env    string `json:"env"`
	Enable bool   `json:"deny"`
}

// DenyEnvList ..
type DenyEnvList []*DenyEnvItem

// ProjectWhiteIDs ..
type ProjectWhiteIDs []int64

// PublishSetupReq ..
type PublishSetupReq struct {
	Announcement    string          `json:"announcement"`
	DateEnabled     bool            `json:"date_enable"`
	DenyEndTime     string          `json:"deny_end_time"`
	DenyStartTime   string          `json:"deny_start_time"`
	DenyEnvList     DenyEnvList     `json:"deny_env_list"`
	WeeklyDenyList  WeeklyDenyList  `json:"pre_week_deny_list"`
	ProjectWhiteIDs ProjectWhiteIDs `json:"project_white_ids"`
}

// String ...
func (p *DenyEnvList) String() (string, error) {
	bytes, err := json.Marshal(p)
	return string(bytes), err
}

// Struct ...
func (p DenyEnvList) Struct(sc string) (DenyEnvList, error) {
	err := json.Unmarshal([]byte(sc), &p)
	return p, err
}

// String ...
func (w *WeeklyDenyList) String() (string, error) {
	bytes, err := json.Marshal(w)
	return string(bytes), err
}

// Struct ...
func (w WeeklyDenyList) Struct(ws string) (WeeklyDenyList, error) {
	err := json.Unmarshal([]byte(ws), &w)
	return w, err
}

// String ...
func (p *ProjectWhiteIDs) String() (string, error) {
	bytes, err := json.Marshal(p)
	return string(bytes), err
}

// Struct ...
func (p ProjectWhiteIDs) Struct(ws string) (ProjectWhiteIDs, error) {
	err := json.Unmarshal([]byte(ws), &p)
	return p, err
}

/* --- Sprint2 Flow-component/step/stage      ----      */

// TaskTmplReq create step
type TaskTmplReq struct {
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	SubTask     []SubTask `json:"sub_task"`
}

// TaskTmplResp ..
type TaskTmplResp struct {
	Name        string    `json:"name,omitempty"`
	Type        string    `json:"type,omitempty"`
	Description string    `json:"description,omitempty"`
	SubTask     []subTask `json:"sub_task,omitempty"`

	ID          int64     `json:"id,omitempty"`
	ComponentID int64     `json:"component_id,omitempty"`
	CreateAt    time.Time `json:"create_at,omitempty"`
	UpdateAt    time.Time `json:"update_at,omitempty"`
	Creator     string    `json:"creator,omitempty"`
	TypeDisplay string    `json:"type_display,omitempty"`
}

// String ...
func (f *TaskTmplReq) String() (string, error) {
	bytes, err := json.Marshal(f.SubTask)
	return string(bytes), err
}

// StepStruct ..
func StepStruct(sc string) ([]subTask, error) {
	subTask := []subTask{}
	err := json.Unmarshal([]byte(sc), &subTask)
	return subTask, err
}

// FlowStageReq create stage
type FlowStageReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Cluster     string `json:"cluster"`
	CIServer    string `json:"ci_server"`
	ArrangeEnv  string `json:"arrange_env"`
}

/* --- Response Part start ---  */

// StepRsp .. manual step defined
type StepRsp struct {
	Name    string `json:"name,omitempty"`
	Creator string `json:"creator,omitempty"`
	Message string `json:"message,omitempty"`
}

// ManualStepResp ..
type ManualStepResp struct {
	PreviousStep *StepRsp `json:"previous_step"`
	CurrenStep   *StepRsp `json:"current_step"`
}

// PublishStepResp ...
type PublishStepResp struct {
	TargetBranch      []string `json:"target_branch,omitempty"`
	AppName           string   `json:"app_name,omitempty"`
	Type              string   `json:"type,omitempty"`
	ProjectAppID      int64    `json:"project_app_id,omitempty"`
	Language          string   `json:"language,omitempty"`
	BranchName        string   `json:"branch_name,omitempty"`
	BuildPath         string   `json:"build_path,omitempty"`
	CompileCommand    string   `json:"compile_command,omitempty"`
	BranchHistoryList []string `json:"branch_history_list,omitempty"`
}

// BuildStepResp ..
type BuildStepResp struct {
	Apps        []*PublishStepResp `json:"apps"`
	VersionNo   string             `json:"versionNo"`
	VersionName string             `json:"versionName"`
}

// DeployStepAppRsp ..
type DeployStepAppRsp struct {
	Type         string `json:"type"`
	Name         string `json:"name"`
	ProjectAppID int64  `json:"project_app_id"`
}

// AppMergeInfo ..
type AppMergeInfo struct {
	Name         string `json:"name"`
	SourceBranch string `json:"source_branch"`
	TargetBranch string `json:"target_branch"`
	Result       string `json:"result"`
	IsSuccess    bool   `json:"isSuccess"`
	Error        string `json:"error"`
}

// MergeBranchStepApp ..
type MergeBranchStepApp struct {
	*models.PublishApp
	TargetBranch []string `json:"target_branch"`
	AppName      string   `json:"app_name"`
}

// MergeBranchStepResp ..
type MergeBranchStepResp struct {
	*models.Publish
	Apps []*MergeBranchStepApp `json:"apps"`
}

// RunBuildAllParms there are all apps parms for jenkins pipeline job
type RunBuildAllParms struct {
	*RunBuildAppReq
	*models.ScmApp
	Release     string `json:"release"`
	MergeBranch bool   `json:"merge-branch"`
	ProjectID   int64
}

// RunDeployAllParms there are all apps parms for jenkins pipeline job
type RunDeployAllParms struct {
	*models.ScmApp
	*RunDeployAppReq
	ImageAddr string `json:"image_addr"`
	ProjectID int64
}

// AppParamsForCreatePublishJob ..
type AppParamsForCreatePublishJob struct {
	ProjectAppID int64  `json:"project_app_id"`
	Branch       string `json:"branch"`
	Path         string `json:"path"`
	ImageVersion string `json:"image_version"`
	Gray         bool   `json:"gray"`
	ImageAddr    string `json:"image_addr"`
}

// AppParamsForHealthCheck ..
type AppParamsForHealthCheck struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	FullName string `json:"full_name"`
}

// PublishJobBuildResult ..
type PublishJobBuildResult struct {
	AppID           string `json:"app_id"`
	ImageVersion    string `json:"image_version"`
	OriginRevision  string `json:"origin_revision"`
	CurrentRevision string `json:"current_revision"`
	Release         string `json:"release"`
}

// PublishJobRsp ..
type PublishJobRsp struct {
	*models.PublishJob
	Apps []*PublishJobAppRsp `json:"apps"`
}

// PublishJobAppRsp ..
type PublishJobAppRsp struct {
	*models.PublishJobApp
	AppName string `json:"app_name"`
}

// FakePublishJobRsp ..
type FakePublishJobRsp struct {
	Progress int    `json:"progress"`
	Status   string `json:"status"`
}

// JenkinsConfigRsp ..
type JenkinsConfigRsp struct {
	Jenkins string `json:"jenkins"`
}

// PublishStatsReq ..
type PublishStatsReq struct {
	AppIDs    []int64 `json:"app_ids"`
	StartTime string  `json:"start_time"`
	EndTime   string  `json:"end_time"`
	EnvIDs    []int64 `json:"env_ids"`
}

// PublishStatsRsp ..
type PublishStatsRsp struct {
	Time          string `json:"time"`
	BuildSuccess  int64  `json:"build_success"`
	BuildFailed   int64  `json:"build_failed"`
	DeploySuccess int64  `json:"deploy_success"`
	DeployFailed  int64  `json:"deploy_failed"`
	TotalSuccess  int64  `json:"total_success"`
	TotalFailed   int64  `json:"total_failed"`
}

// AutoTestReportRsp ..
type AutoTestReportRsp struct {
	SonarReportURL  string `json:"sonar_report_url"`
	AppName         string `json:"app_name"`
	CurrentRevision string `json:"current_revision"`
	BranchName      string `json:"branch_name"`
	JacocoReportURL string `json:"jacoco_report_url"`
}

type clusterItem struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

// ClusterListRsp ..
type ClusterListRsp struct {
	Jenkins []*clusterItem `json:"jenkins"`
	K8s     []*clusterItem `json:"k8s"`
}
