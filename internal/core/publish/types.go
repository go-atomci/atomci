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

package publish

import (
	"github.com/go-atomci/atomci/internal/models"
	"github.com/go-atomci/atomci/utils/query"
)

// PubllishReqApp ..
type PubllishReqApp struct {
	AppID          int64  `json:"app_id"`
	BranchName     string `json:"branch_name"`
	CompileCommand string `json:"compile_command"`
}

// PublishAddApps ..
type PublishAddApps struct {
	Apps []*PubllishReqApp `json:"apps"`
}

// PublishReq create publish-order request body
type PublishReq struct {
	Apps           []*PubllishReqApp `json:"apps"`
	Name           string            `json:"name"`
	BindPipelineID int64             `json:"bind_pipeline_id"`
	VersionNo      string            `json:"version_no"`
}

// PublishUpdate ..
type PublishUpdate struct {
	VersionNo string `json:"version_no"`
	Name      string `json:"name"`
}

// PublishReqFilterQuery ...
type PublishReqFilterQuery struct {
	query.FilterQuery
	Status  int64 `json:"status"`
	Deleted bool  `json:"deleted"`
}

// PublishInfoApp ..
type PublishInfoApp struct {
	*models.PublishApp
	Name     string `json:"name"`
	Language string `json:"language"`
	Type     string `json:"type"`
}

// PublishStep ..
type PublishStep struct {
	Type   string `json:"type"`
	Name   string `json:"name"`
	Status int64  `json:"status"`
	Index  int    `json:"index"`
}

// PublishInfoResp ...
type PublishInfoResp struct {
	*models.Publish
	PipelineName string                   `json:"pipeline_name"`
	Apps         []*PublishInfoApp        `json:"apps"`
	Operations   *models.PublishOperation `json:"operations"`
	Steps        []*PublishStep           `json:"steps"`
}

// TriggerBackToReq trigger publish order back to request body
type TriggerBackToReq struct {
	StageID int64  `json:"stage_id"`
	Message string `json:"message"`
}

// CreateOperationLogReq ..
type CreateOperationLogReq struct {
	Creator            string `json:"creator"`
	StageName          string `json:"stage_name"`
	StepName           string `json:"step_name"`
	Message            string `json:"message"`
	JobName            string `json:"job_name"`
	Type               string `json:"type"`
	PipelineInstanceID int64  `json:"pipeline_instance_id"`
	StageInstanceID    int64  `json:"stage_instance_id"`
	StepIndex          int    `json:"step_index"`
	PublishID          int64  `json:"publish_id"`
	StageID            int64  `json:"stage_id"`
	Status             int64  `json:"status"`
	RunID              int64  `json:"run_id"`
}

// CirculationRsp back-to/next-stage
type CirculationRsp struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
