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

package apps

import "github.com/go-atomci/atomci/internal/models"

type ScmAppReq struct {
	// ProjectAppReq add app into project request body.
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

type ScmAppUpdateReq struct {
	BranchName   string `json:"branch_name"`
	Language     string `json:"language"`
	Name         string `json:"name"`
	Path         string `json:"path"`
	CompileEnvID int64  `json:"compile_env_id"`
	BuildPath    string `json:"build_path"`
	Dockerfile   string `json:"dockerfile"`
}

// SCMAppRsp ..
type SCMAppRsp struct {
	*models.ScmApp
	BranchHistoryList []string `json:"branch_history_list,omitempty"`
	CompileEnv        string   `json:"compile_env"`
}

// RepoProjectRsp ..
type RepoProjectRsp struct {
	RepoID   int64  `json:"repo_id"`
	Path     string `json:"path"`
	FullName string `json:"full_name"`
	Name     string `json:"name"`
}

type AppArrangConfig struct {
	Config string `json:"config,omitempty"`
}

// AppArrangeReq ..
type AppArrangeReq struct {
	ProjectAppID int64         `json:"project_app_id,omitempty"`
	CopyToEnvIDs []int64       `json:"copy_to_env_ids,omitempty"`
	Config       string        `json:"config,omitempty"`
	ImageMapings []ImageMaping `json:"image_mapings,omitempty"`
}

type ImageMaping struct {
	ID           int64  `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Image        string `json:"image,omitempty"`
	ProjectAppID int64  `json:"project_app_id,omitempty"`
	ImageTagType int64  `json:"image_tag_type,omitempty"`
	ArrangeID    int64  `json:"arrange_id,omitempty"`
}

type AppArrangeResp struct {
	ID           int64         `json:"id,omitempty"`
	Name         string        `json:"name,omitempty"`
	EnvID        int64         `json:"env_id,omitempty"`
	ProjectAppID int64         `json:"project_app_id,omitempty"`
	Config       string        `json:"config,omitempty"`
	ImageMapings []ImageMaping `json:"image_mapings,omitempty"`
}
