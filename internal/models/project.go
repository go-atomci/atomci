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

// const project status defined
const (
	_ = iota
	ProjectRuning
	ProjectEnd
)

// ProejctFilterQuery ..
type ProejctFilterQuery struct {
	query.FilterQuery
	Status        *int64 `json:"status,omitempty"`
	CreateAtStart string `json:"createAtStart"`
	CreateAtEnd   string `json:"createAtEnd"`
	Name          string `json:"name"`
}

// ProejctAppFilterQuery ..
type ProejctAppFilterQuery struct {
	query.FilterQuery
	Name          string `json:"name"`
	Type          string `json:"type"`
	Language      string `json:"language"`
	Path          string `json:"path"`
	Creator       string `json:"creator"`
	CreateAtStart string `json:"createAtStart"`
	CreateAtEnd   string `json:"createAtEnd"`
}

// Project ...
type Project struct {
	Addons
	Name        string     `orm:"column(name);size(64)" json:"name"`
	Description string     `orm:"column(description);size(256)" json:"description"`
	Status      int8       `orm:"column(status); default(1)" json:"status"`
	Owner       string     `orm:"column(owner);size(64)" json:"owner"`
	Creator     string     `orm:"column(creator);size(64)" json:"creator"`
	StartAt     time.Time  `orm:"column(start_at);auto_now;type(datetime);null" json:"start_at"`
	EndAt       *time.Time `orm:"column(end_at);type(datetime);null" json:"end_at"`
}

// TableName ...
func (t *Project) TableName() string {
	return "pub_project"
}

// ProjectResponse Project response info
type ProjectResponse struct {
	ID          int64      `json:"id"`
	CreateAt    time.Time  `json:"create_at"`
	StartAt     time.Time  `json:"start_at"`
	EndAt       *time.Time `json:"end_at"`
	UpdateAt    time.Time  `json:"update_at"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      int8       `json:"status"`
	Owner       string     `json:"owner"`
	Creator     string     `json:"creator"`
	Members     int        `json:"members"`
	MembersName []string   `json:"membersName"`
}

// ProjectDetailResponse ..
type ProjectDetailResponse struct {
	*ProjectResponse
	CodeRepos int64       `json:"code_repos"`
	Releases  interface{} `json:"releases"`
}

// ProjectUser ...
type ProjectUser struct {
	Addons
	ProjectID int64  `orm:"column(project_id)" json:"project_id"`
	User      string `orm:"column(user);size(64)" json:"user"`
	RoleID    int64  `orm:"column(role_id)" json:"role_id"`
}

// TableName ...
func (t *ProjectUser) TableName() string {
	return "pub_project_user"
}

// ProjectApp ...
type ProjectApp struct {
	Addons
	Creator           string   `orm:"column(creator);size(64);null" json:"creator"`
	ProjectID         int64    `orm:"column(project_id)" json:"project_id"`
	ScmID             int64    `orm:"column(scm_id)" json:"scm_id"`
	BranchHistoryList []string `orm:"-" json:"branch_history_list"`
}

// TableName ..
func (t *ProjectApp) TableName() string {
	return "pub_project_app"
}

// ProjectEnv the Basic Data of stages based on commpany
type ProjectEnv struct {
	Addons
	ProjectID   int64  `orm:"column(project_id)" json:"project_id"`
	Name        string `orm:"column(name);size(64)" json:"name"`
	Description string `orm:"column(description);size(256)" json:"description"`
	Cluster     int64  `orm:"column(cluster);" json:"cluster"`
	Namespace   string `orm:"column(namespace);size(256)" json:"namespace"`
	ArrangeEnv  string `orm:"column(arrange_env);size(64)" json:"arrange_env"`
	CIServer    int64  `orm:"column(ci_server);" json:"ci_server"`
	Registry    int64  `orm:"column(registry);" json:"registry"`
	Creator     string `orm:"column(creator);size(64)" json:"creator"`
}

// TableName ...
func (t *ProjectEnv) TableName() string {
	return "project_env"
}

// ProjectPipeline ...
type ProjectPipeline struct {
	Addons
	Name        string `orm:"column(name);size(64)" json:"name"`
	Description string `orm:"column(description);size(256)" json:"description"`
	Config      string `orm:"column(config);type(text)" json:"config"`
	Creator     string `orm:"column(creator);size(64)" json:"creator"`
	ProjectID   int64  `orm:"column(project_id)" json:"project_id"`
	IsDefault   bool   `orm:"column(is_default);default(false)" json:"is_default"`
}

// TableName ...
func (t *ProjectPipeline) TableName() string {
	return "project_pipeline"
}

// PipelineInstance ..
type PipelineInstance struct {
	Addons
	PipelineID  int64  `orm:"column(pipeline_id)" json:"pipeline_id"`
	PublishID   int64  `orm:"column(publish_id)" json:"publish_id"`
	Creator     string `orm:"column(creator);size(64)" json:"creator"`
	Name        string `orm:"column(name);size(64)" json:"name"`
	Description string `orm:"column(description);size(256)" json:"description"`
	Config      string `orm:"column(config);type(text)" json:"config"`
	ProjectID   int64  `orm:"column(project_id)" json:"project_id"`
	IsDefault   bool   `orm:"column(is_default);default(false)" json:"is_default"`
}

// TableName ...
func (t *PipelineInstance) TableName() string {
	return "pipeline_instance"
}
