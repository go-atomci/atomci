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

package initialize

var resourceReq = ResourceReq{
	Resources: []BatchResourceTypeSpec{
		BatchResourceTypeSpec{
			ResourceType: []string{"*", "所有资源"},
			ResourceOperation: [][]string{
				[]string{"*", "所有操作"},
			},
			ResourceConstraint: [][]string{
				[]string{"*", "所有约束"},
			},
		},
		BatchResourceTypeSpec{
			ResourceType: []string{"auth", "认证"},
			ResourceOperation: [][]string{
				[]string{"*", "认证所有操作"},
				[]string{"UserLogin", "用户登录"},
				[]string{"UserLogout", "用户登出"},
				[]string{"GetCurrentUser", "获取当前用户信息"},
			},
			ResourceConstraint: [][]string{},
		},
		BatchResourceTypeSpec{
			ResourceType: []string{"audit", "操作审计"},
			ResourceOperation: [][]string{
				[]string{"*", "操作审计所有操作"},
				[]string{"AuditList", "获取操作审计列表"},
			},
			ResourceConstraint: [][]string{},
		},
		BatchResourceTypeSpec{
			ResourceType: []string{"user", "用户"},
			ResourceOperation: [][]string{
				[]string{"*", "用户所有操作"},
				[]string{"UserList", "获取用户列表"},
				[]string{"CreateUser", "创建用户"},
				[]string{"GetUser", "获取用户详情"},
				[]string{"UpdateUser", "更新用户"},
				[]string{"DeleteUser", "删除用户"},
				[]string{"GetUserResourceConstraintValues", "获取用户资源约束的值"},
			},
			ResourceConstraint: [][]string{
				[]string{"user", "用户账号"},
			},
		},

		BatchResourceTypeSpec{
			ResourceType: []string{"project", "项目"},
			ResourceOperation: [][]string{
				[]string{"*", "项目所有操作"},
				[]string{"ProjectList", "获取项目列表"},
				[]string{"CreateProject", "创建项目"},
				[]string{"UpdateProject", "更新项目信息"},
				[]string{"DeleteProject", "删除项目"},
				[]string{"GetProject", "获取项目信息"},
				[]string{"CreateProjectApp", "项目添加应用"},
				[]string{"UpdateProjectApp", "更新项目应用"},
				[]string{"GetProjectApps", "获取项目应用列表"},
				[]string{"GetProjectApp", "获取项目应用详情"},
				[]string{"GetAppsByPagination", "获取项目应用分页列表"},
				[]string{"GetArrange", "获取应用编排"},
				[]string{"SetArrange", "设置应用编排"},
				[]string{"GetAppBranches", "获取应用分支"},
				[]string{"SyncAppBranches", "同步远程分支"},
				[]string{"SwitchProjectBranch", "切换项目应用的默认分支"},
				[]string{"DeleteProjectApp", "删除项目应用"},
				[]string{"GetRepos", "获取代码仓库列表"},
				[]string{"GetGitProjectsByRepoID", "获取代码仓库项目列表"},
				[]string{"GetProjectEnvs", "项目环境列表"},
				[]string{"GetProjectPipelinesByPagination", "项目流程分页列表"},

				// project app service
				[]string{"GetProjectAppServices", "应用服务列表"},

				// project pipeline
				[]string{"PipelineCreate", "创建项目流程"},
				[]string{"PipelineUpdate", "更新流程基础信息"},
				[]string{"ProjectPipelineInfo", "获取项目流程信息"},
				[]string{"PipelineDelete", "删除项目流程"},
				[]string{"FlowStepList", "获取任务模板列表"},

				[]string{"GetProjectEnvsByPagination", "项目环境分页列表"},
				[]string{"CreateProjectEnv", "新建项目环境"},
				[]string{"UpdateProjectEnv", "更新项目环境"},
				[]string{"ProjectAppServiceStats", "获取项目应用统计"},
			},
			ResourceConstraint: [][]string{
				[]string{"project_id", "项目ID"},
			},
		},
		BatchResourceTypeSpec{
			ResourceType: []string{"publish", "流水线"},
			ResourceOperation: [][]string{
				[]string{"*", "流水线所有操作"},
				[]string{"GetProjectPipelines", "项目流程列表"},
				[]string{"PublishList", "流水线列表"},
				[]string{"CreatePublishOrder", "创建流水线"},
				[]string{"GetPublish", "流水线详情"},
				[]string{"ClosePublish", "关闭流水线"},
				[]string{"DeletePublish", "删除流水线"},
				[]string{"GetCanAddedApps", "获取可添加应用列表"},
				[]string{"AddPublishApp", "版本添加应用"},
				[]string{"DeletePublishApp", "版本删除应用"},
				[]string{"GetOpertaionLogByPagination", "获取流水线操作日志"},
				[]string{"GetBackTo", "获取回退列表"},
				[]string{"TriggerBackTo", "触发流水线回退操作"},
				[]string{"GetNextStage", "获取流转列表"},
				[]string{"TriggerNextStage", "触发流水线流转操作"},
				[]string{"GetStepInfo", "获取步骤执行信息"},
				[]string{"RunStep", "触发步骤执行"},
				[]string{"RunStepCallback", "步骤执行回调"},
			},
			ResourceConstraint: [][]string{
				[]string{"project_id", "项目ID"},
				[]string{"publishID", "发布单ID"},
				[]string{"envID", "环境ID"},
			},
		},
		BatchResourceTypeSpec{
			ResourceType: []string{"system", "系统设置"},
			ResourceOperation: [][]string{
				[]string{"*", "系统设置所有操作"},
				[]string{"GetCompileEnvs", "编译环境列表"},
				[]string{"GetIntegrateClusters", "获取集成的集群列表"},
				[]string{"GetIntegrateSettings", "获取集成配置列表"},

				[]string{"FlowComponentList", "获取基础组件列表"},
				[]string{"FlowStepListByPagination", "获取任务模板分页列表"},
				[]string{"FlowStepCreate", "创建任务模板"},
				[]string{"FlowStepUpdate", "更新任务模板"},
				[]string{"FlowStepDelete", "删除任务模板"},
			},
			ResourceConstraint: [][]string{},
		},
	},
}

var gaetwayReq = RouterReq{
	Routers: [][]string{
		[]string{"atomci/api/v1/login", "POST", "atomci", "auth", "UserLogin"},
		[]string{"atomci/api/v1/logout", "GET", "atomci", "auth", "UserLogout"},
		[]string{"atomci/api/v1/getCurrentUser", "GET", "atomci", "auth", "GetCurrentUser"},
		[]string{"atomci/api/v1/audit", "GET", "atomci", "audit", "AuditList"},
		[]string{"atomci/api/v1/users", "GET", "atomci", "user", "UserList"},
		[]string{"atomci/api/v1/users", "POST", "atomci", "user", "CreateUser"},
		[]string{"atomci/api/v1/users/:user", "GET", "atomci", "user", "GetUser"},
		[]string{"atomci/api/v1/users/:user", "PUT", "atomci", "user", "UpdateUser"},
		[]string{"atomci/api/v1/users/:user", "DELETE", "atomci", "user", "DeleteUser"},
		[]string{"atomci/api/v1/users/:user/resources/:resourceType/constraints/values", "GET", "atomci", "user", "GetUserResourceConstraintValues"},
		[]string{"atomci/api/v1/groups", "GET", "atomci", "group", "GroupList"},
		[]string{"atomci/api/v1/groups/:group", "GET", "atomci", "group", "GetGroup"},
		[]string{"atomci/api/v1/groups/:group", "PUT", "atomci", "group", "UpdateGroup"},
		[]string{"atomci/api/v1/groups/:group", "DELETE", "atomci", "group", "DeleteGroup"},
		[]string{"atomci/api/v1/groups/:group/users", "GET", "atomci", "group", "GroupUserList"},
		[]string{"atomci/api/v1/groups/:group/users", "POST", "atomci", "group", "AddGroupUsers"},
		[]string{"atomci/api/v1/groups/:group/users/:user", "PUT", "atomci", "group", "UpdateGroupUser"},
		[]string{"atomci/api/v1/groups/:group/users/:user", "DELETE", "atomci", "group", "RemoveGroupUser"},

		// app repo
		[]string{"atomci/api/v1/repos", "GET", "atomci", "repository", "GetRepos"},
		[]string{"atomci/api/v1/repos/:repo_id/projects", "POST", "atomci", "repository", "GetGitProjectsByRepoID"},

		// project
		[]string{"atomci/api/v1/projects", "POST", "atomci", "project", "ProjectList"},
		[]string{"atomci/api/v1/projects/create", "POST", "atomci", "project", "CreateProject"},
		[]string{"atomci/api/v1/projects/:project_id", "PUT", "atomci", "project", "UpdateProject"},
		[]string{"atomci/api/v1/projects/:project_id", "DELETE", "atomci", "project", "DeleteProject"},
		[]string{"atomci/api/v1/projects/:project_id", "GET", "atomci", "project", "GetProject"},
		[]string{"atomci/api/v1/projects/:project_id/pipelines", "GET", "atomci", "project", "GetProjectPipelines"},
		[]string{"atomci/api/v1/projects/:project_id/pipelines", "POST", "atomci", "project", "GetProjectPipelinesByPagination"},
		[]string{"atomci/api/v1/pipelines/flow/steps", "GET", "atomci", "project", "FlowStepList"},
		[]string{"atomci/api/v1/projects/:project_id/pipelines/create", "POST", "atomci", "project", "PipelineCreate"},
		[]string{"atomci/api/v1/projects/:project_id/pipelines/:id", "GET", "atomci", "project", "ProjectPipelineInfo"},
		[]string{"atomci/api/v1/projects/:project_id/pipelines/:id", "PUT", "atomci", "project", "PipelineUpdate"},
		[]string{"atomci/api/v1/projects/:project_id/pipelines/:id", "DELETE", "atomci", "project", "PipelineDelete"},
		[]string{"atomci/api/v1/projects/:project_id/apps/create", "POST", "atomci", "project", "CreateProjectApp"},
		[]string{"atomci/api/v1/projects/:project_id/apps", "GET", "atomci", "project", "GetProjectApps"},
		[]string{"atomci/api/v1/projects/:project_id/apps/:project_app_id", "GET", "atomci", "project", "GetProjectApp"},
		[]string{"atomci/api/v1/projects/:project_id/apps", "POST", "atomci", "project", "GetAppsByPagination"},
		[]string{"atomci/api/v1/projects/:project_id/apps/:app_id/:arrange_env/arrange", "GET", "atomci", "project", "GetArrange"},
		[]string{"atomci/api/v1/projects/:project_id/apps/:app_id/:arrange_env/arrange", "POST", "atomci", "project", "SetArrange"},
		[]string{"atomci/api/v1/projects/:project_id/apps/:app_id/branches", "POST", "atomci", "project", "GetAppBranches"},
		[]string{"atomci/api/v1/projects/:project_id/apps/:app_id/syncBranches", "POST", "atomci", "project", "SyncAppBranches"},
		[]string{"atomci/api/v1/projects/:project_id/apps/:project_app_id", "PUT", "atomci", "project", "UpdateProjectApp"},
		[]string{"atomci/api/v1/projects/:project_id/apps/:project_app_id", "PATCH", "atomci", "project", "SwitchProjectBranch"},
		[]string{"atomci/api/v1/projects/:project_id/apps/:project_app_id", "DELETE", "atomci", "project", "DeleteProjectApp"},
		[]string{"atomci/api/v1/projects/:project_id/clusters/:cluster/apps", "POST", "atomci", "project", "GetProjectAppServices"},
		[]string{"atomci/api/v1/projects/:project_id/publish/stats", "POST", "atomci", "project", "ProjectPublishStats"},
		[]string{"atomci/api/v1/projects/:project_id/envs", "GET", "atomci", "project", "GetProjectEnvs"},
		[]string{"atomci/api/v1/projects/:project_id/envs", "POST", "atomci", "project", "GetProjectEnvsByPagination"},
		[]string{"atomci/api/v1/projects/:project_id/envs/create", "POST", "atomci", "project", "CreateProjectEnv"},
		[]string{"atomci/api/v1/projects/:project_id/envs/:env_id", "PUT", "atomci", "project", "UpdateProjectEnv"},

		// publish
		[]string{"atomci/api/v1/projects/:project_id/publishes", "POST", "atomci", "publish", "PublishList"},
		[]string{"atomci/api/v1/projects/:project_id/publishes/create", "POST", "atomci", "publish", "CreatePublishOrder"},
		[]string{"atomci/api/v1/projects/:project_id/publishes/:publish_id", "GET", "atomci", "publish", "GetPublish"},
		[]string{"atomci/api/v1/projects/:project_id/publishes/:publish_id", "PUT", "atomci", "publish", "ClosePublish"},
		[]string{"atomci/api/v1/projects/:project_id/publishes/:publish_id", "DELETE", "atomci", "publish", "DeletePublish"},
		[]string{"atomci/api/v1/projects/:project_id/publishes/:publish_id/apps/can_added", "GET", "atomci", "publish", "GetCanAddedApps"},
		[]string{"atomci/api/v1/projects/:project_id/publishes/:publish_id/apps/create", "POST", "atomci", "publish", "AddPublishApp"},
		[]string{"atomci/api/v1/projects/:project_id/publishes/:publish_id/apps/:publish_app_id", "DELETE", "atomci", "publish", "DeletePublishApp"},
		[]string{"atomci/api/v1/projects/:project_id/publishes/:publish_id/audits", "POST", "atomci", "publish", "GetOpertaionLogByPagination"},
		[]string{"atomci/api/v1/projects/:project_id/publishes/:publish_id/stages/:stage_id/back-to", "GET", "atomci", "publish", "GetBackTo"},
		[]string{"atomci/api/v1/projects/:project_id/publishes/:publish_id/stages/:stage_id/back-to", "POST", "atomci", "publish", "TriggerBackTo"},
		[]string{"atomci/api/v1/projects/:project_id/publishes/:publish_id/stages/:stage_id/next-stage", "GET", "atomci", "publish", "GetNextStage"},
		[]string{"atomci/api/v1/projects/:project_id/publishes/:publish_id/stages/:stage_id/next-stage", "POST", "atomci", "publish", "TriggerNextStage"},
		[]string{"atomci/api/v1/pipelines/:project_id/publishes/:publish_id/stages/:stage_id/steps/:step_name", "GET", "atomci", "publish", "GetStepInfo"},
		[]string{"atomci/api/v1/pipelines/:project_id/publishes/:publish_id/stages/:stage_id/steps/:step_name", "POST", "atomci", "publish", "RunStep"},
		[]string{"atomci/api/v1/pipelines/:project_id/publishes/:publish_id/stages/:stage_id/steps/:step_name/callback", "POST", "atomci", "publish", "RunStepCallback"},

		// integrate
		[]string{"atomci/api/v1/integrate/compile_envs", "GET", "atomci", "system", "GetCompileEnvs"},
		[]string{"atomci/api/v1/integrate/clusters", "GET", "atomci", "system", "GetIntegrateClusters"},
		[]string{"atomci/api/v1/integrate/settings", "GET", "atomci", "system", "GetIntegrateSettings"},

		// task template
		[]string{"atomci/api/v1/pipelines/flow/components", "GET", "atomci", "system", "FlowComponentList"},
		[]string{"atomci/api/v1/pipelines/flow/steps", "POST", "atomci", "system", "FlowStepListByPagination"},
		[]string{"atomci/api/v1/pipelines/flow/steps/create", "POST", "atomci", "system", "FlowStepCreate"},
		[]string{"atomci/api/v1/pipelines/flow/steps/:step_id", "PUT", "atomci", "system", "FlowStepUpdate"},
		[]string{"atomci/api/v1/pipelines/flow/steps/:step_id", "DELETE", "atomci", "system", "FlowStepDelete"},
	},
}
