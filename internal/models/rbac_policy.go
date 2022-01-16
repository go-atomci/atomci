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
)

type RbacPolicy struct {
	Addons
	PType string `orm:"column(p_type);" json:"p_type"`
	V0    string `orm:"column(v0)" json:"v0"`
	V1    string `orm:"column(v1)" json:"v1"`
	V2    string `orm:"column(v2)" json:"v2"`
	V3    string `orm:"column(v3)" json:"v3"`
	V4    string `orm:"column(v4)" json:"v4"`
	V5    string `orm:"column(v5)" json:"v5"`
}

func (t *RbacPolicy) TableName() string {
	return "rbac_policy"
}

func Policy() []RbacPolicy {
	return []RbacPolicy{
		{Addons: Addons{ID: 1, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " admin", V1: " *", V2: " *", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 2, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/getCurrentUser", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 3, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/projects", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 4, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/projects/:id", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 5, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/projects/:id/apps", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 6, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/projects/:id/apps/:id/branches", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 7, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/projects/:id/apps/:id", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 8, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/projects/:id/arrange_env/:env/namespaces", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 9, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/projects/:id/apps/:id/:env/arrange", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 10, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/projects/:id/arrange_env/:env/bizclusters", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 11, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/projects/:id/arrange_env/:env/nodes", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 12, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/projects/:id/apps/:id/syncBranches", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 13, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/projects/:id/apps/:id", V2: " PATCH", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 14, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/projects/:id/apps/:id", V2: " PUT", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 15, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/pipelines/flow/stages", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 16, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/pipelines/:id/publishes/:pipeline_id/stages/:stage_id/steps/:step", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 17, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/pipelines/:id/publishes/:pipeline_id/stages/:stage_id/steps/:step", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 18, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/pipelines/flow/stages", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 19, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/projects/:id/publishes", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 20, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/projects/:id/publishes/:id/stages/:id/:step", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 21, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/projects/:id/publishes/:id/stages/:id/:step", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 22, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/projects/:id/stages/:id/publish-jobs/:id/deploy", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 23, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/projects/:id/publishes/:id/audits", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 24, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/projects/:id/publishes/:id", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 25, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/pipelines/stages/:id/jenkins-config", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 26, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/init/users", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 27, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/init/groups", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 28, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/init/resource", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 29, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/init/gateway/:backend", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 30, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/getCurrentUser", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 31, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/projects", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 32, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/projects/:id", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 33, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/projects/:id/apps", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 34, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/projects/:id/apps/:id/branches", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 35, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/projects/:id/apps/:id", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 36, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/projects/:id/arrange_env/:env/namespaces", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 37, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/projects/:id/apps/:id/:env/arrange", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 38, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/projects/:id/arrange_env/:env/bizclusters", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 39, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/projects/:id/arrange_env/:env/nodes", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 40, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/projects/:id/apps/:id/syncBranches", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 41, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/projects/:id/apps/:id", V2: " PATCH", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 42, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/projects/:id/apps/:id", V2: " PUT", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 43, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/pipelines/flow/stages", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 44, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/projects/create", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 45, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/projects/create", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 46, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/projects/:project_id/apps", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 47, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/projects/:project_id/envs", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 48, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/projects/:project_id/envs", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 49, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/projects/:project_id/pipelines", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 50, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/integrate/clusters", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 51, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/integrate/settings", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 52, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/repos", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 53, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/repos/:repo_id/projects", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 54, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/integrate/compile_envs", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 55, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " developer", V1: " /atomci/api/v1/projects/:project_id/apps/create", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 56, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/clusters/:cluster/namespaces/:namespace/apps/:app/restart", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 57, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/clusters/:cluster/namespaces/:namespace/apps/:app/scale", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 58, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/clusters/:cluster/namespaces/:namespace/pods/:podname/containernames/:containername", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 59, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/pipelines", V2: " PUT", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 60, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/apps/branches", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 61, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/create", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 62, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/apps/create", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 63, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/projects/:project_id/envs/create", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 64, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/clusters/:cluster/namespaces/:namespace/apps/:app", V2: " DELETE", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 65, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/pipelines/:id", V2: " DELETE", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 66, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id", V2: " DELETE", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 67, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/apps/:project_app_id", V2: " DELETE", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 68, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/members/:id", V2: " DELETE", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 69, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/pipelines/flow/steps", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 70, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/apps/:app_id/branches", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 71, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/apps", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 72, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/clusters/:cluster/namespaces/:namespace/apps/:app/event", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 73, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/clusters/:cluster/namespaces/:namespace/apps/:app", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 74, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/clusters/:cluster/namespaces/:namespace/apps/:app/log", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 75, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/apps/:app_id/:arrange_env/arrange", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 76, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/arrange_env/:arrange_env/bizclusters", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 77, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/pipelines/stages/:stage_id/jenkins-config", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 78, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/arrange_env/:arrange_env/namespaces", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 79, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/arrange_env/:arrange_env/nodes", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 80, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 81, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/apps/:project_app_id", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 82, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/apps", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 83, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/projects/:project_id/clusters/:cluster/apps", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 84, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/projects/:project_id/envs", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 85, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/projects/:project_id/envs", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 86, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/members", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 87, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/pipelines", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 88, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/projects/:project_id/pipelines", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 89, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/arrange/yaml/parser", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 90, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/projects/:project_id/pipelines/create", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 91, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/projects/:project_id/pipelines/:id", V2: " DELETE", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 92, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/projects/:project_id/pipelines/:id", V2: " PUT", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 93, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /caas/api/v1/projects/:project_id/apps/stats", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 94, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 95, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/pipelines/:id", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 96, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/publish/stats", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 97, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/apps/:app_id/:arrange_env/arrange", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 98, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/apps/:project_app_id", V2: " PATCH", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 99, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/apps/:app_id/syncBranches", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 100, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id", V2: " PUT", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 101, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/apps/:project_app_id", V2: " PUT", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 102, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /atomci/api/v1/projects/:project_id/envs/:env_id", V2: " PUT", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 103, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/members", V2: " PUT", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 104, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/publishes/:publish_id/apps/create", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 105, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/publishes/:publish_id", V2: " PUT", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 106, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/publishes/create", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 107, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/publishes/:publish_id", V2: " DELETE", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 108, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/publishes/:publish_id/apps/:publish_app_id", V2: " DELETE", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 109, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/publishes/:publish_id/stages/:stage_id/back-to", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 110, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/publishes/:publish_id/apps/can_added", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 111, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/pipelines/stages/:stage_id/jenkins-config", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 112, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/publishes/:publish_id/stages/:stage_id/next-stage", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 113, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/publishes/:publish_id/audits", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 114, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/publishes/:publish_id", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 115, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/publish/setup", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 116, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/pipelines/:project_id/publishes/:publish_id/stages/:stage_id/steps/:step_name", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 117, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/publishes", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 118, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/pipelines/:project_id/publishes/:publish_id/stages/:stage_id/steps/:step_name", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 119, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/pipelines/:project_id/publishes/:publish_id/stages/:stage_id/steps/:step_name/callback", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 120, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/publishes/:publish_id/stages/:stage_id/back-to", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 121, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/projects/:project_id/publishes/:publish_id/stages/:stage_id/next-stage", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 122, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/publish/setup", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 123, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/getCurrentUser", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 124, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/login", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 125, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/logout", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 126, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/pipelines/flow/components", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 127, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/pipelines/flow/stages/create", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 128, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/pipelines/flow/stages/:stage_id", V2: " DELETE", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 129, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/pipelines/flow/stages", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 130, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/pipelines/flow/stages", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 131, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/pipelines/flow/stages/:stage_id", V2: " PUT", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 132, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/pipelines/flow/steps/create", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 133, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/pipelines/flow/steps/:step_id", V2: " DELETE", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 134, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/pipelines/flow/steps", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 135, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/pipelines/flow/steps", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 136, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/pipelines/flow/steps/:step_id", V2: " PUT", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 137, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/pipelines/clusters", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 138, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/pipelines/:pipeline_id/setup", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 139, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/pipelines/create", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 140, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/pipelines/:pipeline_id", V2: " DELETE", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 141, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/pipelines", V2: " GET", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 142, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/pipelines", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 143, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/pipelines/:pipeline_id", V2: " PUT", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 144, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/pipelines/reset", V2: " POST", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 145, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "p", V0: " devManager", V1: " /publishctl/api/v1/pipelines/:pipeline_id/setup", V2: " PUT", V3: "", V4: "", V5: ""},
		{Addons: Addons{ID: 146, Deleted: false, CreateAt: time.Time{}, UpdateAt: time.Time{}, DeleteAt: nil}, PType: "g", V0: " admin", V1: " developer", V2: "", V3: "", V4: "", V5: ""},
	}
}
