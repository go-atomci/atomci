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

package projects

import (
	"github.com/go-atomci/atomci/models"
)

// ProjectReq create project request body
type ProjectReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      int8   `json:"status"`
}

// ProjectUpdateReq ..
type ProjectUpdateReq struct {
	ProjectReq
	Owner string `json:"owner"`
}

// ProjectAppUpdateReq ..
type ProjectAppUpdateReq struct {
	BranchName   string `json:"branch_name"`
	Language     string `json:"language"`
	Name         string `json:"name"`
	Path         string `json:"path"`
	CompileEnvID int64  `json:"compile_env_id"`
	BuildPath    string `json:"build_path"`
	Dockerfile   string `json:"dockerfile"`
}

// ProjectAppBranchUpdateReq ..
type ProjectAppBranchUpdateReq struct {
	BranchName string `json:"branch_name"`
	AppID      int64  `json:"app_id"`
}

// ProjectAppBranchCreateReq .
type ProjectAppBranchCreateReq struct {
	BranchName  string `json:"branch_name"`
	ProjectApps []struct {
		ProjectAppID int64 `json:"project_app_id"`
		AppID        int64 `json:"app_id"`
	} `json:"project_apps"`
	TargetBranch string `json:"target_branch"`
}

// ProjectPipelineReq ..
type ProjectPipelineReq struct {
	PipelineID int64 `json:"pipeline_id"`
}

// ProjectNumberReq ..
type ProjectNumberReq struct {
	RoleID int64  `json:"role_id"`
	User   string `json:"user"`
}

/* ------ response start ------  */

// ProjectAppReq add app into project request body.
type ProjectAppReq struct {
	Name         string `json:"name"`
	CompileEnvID int64  `json:"compile_env_id"`
	Language     string `json:"language"`
	Path         string `json:"path"`
	RepoID       int64  `json:"repo_id"`
	FullName     string `json:"full_name"`
	BranchName   string `json:"branch_name"`
	BuildPath    string `json:"build_path"`
	Dockerfile   string `json:"dockerfile"`
}

// ProjectAppRsp ..
type ProjectAppRsp struct {
	*models.ProjectApp
	BranchHistoryList []string `json:"branch_history_list,omitempty"`
	CompileEnv        string   `json:"compile_env"`
}

// ProjectPipelineRsp ..
type ProjectPipelineRsp struct {
	*models.ProjectPipeline
	Name string `json:"name"`
}

// ProjectNumberRsp ..
type ProjectNumberRsp struct {
	*models.ProjectUser
	Role string `json:"role"`
}
