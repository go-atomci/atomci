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

package routers

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/go-atomci/atomci/internal/api"
	"github.com/go-atomci/atomci/internal/middleware/log"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
)

// RegisterRoutes init api router
func RegisterRoutes() {
	publishAPI :=
		beego.NewNamespace("atomci/api",
			beego.NSNamespace("/v1",

				beego.NSRouter("/logout", &api.AuthController{}, "get:Logout"),
				beego.NSRouter("/login", &api.AuthController{}, "post:Authenticate"),
				beego.NSRouter("/getCurrentUser", &api.UserController{}, "get:GetCurrentUser"),

				beego.NSRouter("/audit", &api.AuditController{}, "get:AuditList"),

				beego.NSRouter("/resources", &api.ResourceController{}, "get:ResourceTypeList;post:CreateResourceType"),
				beego.NSRouter("/resources-operations", &api.ResourceController{}, "get:ResourceOperationsList"),
				beego.NSRouter("/resources/:resourceType", &api.ResourceController{}, "get:GetResourceType;put:UpdateResourceType;delete:DeleteResourceType"),
				beego.NSRouter("/resources/:resourceType/operations", &api.ResourceController{}, "post:AddResourceOperation"),
				beego.NSRouter("/resources/:resourceType/operations/:resourceOperation", &api.ResourceController{}, "put:UpdateResourceOperation;delete:DeleteResourceOperation"),
				beego.NSRouter("/resources/:resourceType/constraints", &api.ResourceController{}, "post:AddResourceConstraint"),
				beego.NSRouter("/resources/:resourceType/constraints/:resourceConstraint", &api.ResourceController{}, "put:UpdateResourceConstraint;delete:DeleteResourceConstraint"),

				beego.NSRouter("/users", &api.UserController{}, "get:UserList;post:CreateUser"),
				beego.NSRouter("/users/:user", &api.UserController{}, "get:GetUser;put:UpdateUser;delete:DeleteUser"),
				beego.NSRouter("/users/:user/resources/:resourceType/constraints/values", &api.UserController{}, "get:GetUserResourceConstraintValues"),

				beego.NSRouter("/groups", &api.GroupController{}, "get:GroupList"),
				beego.NSRouter("/groups/:group", &api.GroupController{}, "get:GetGroup;put:UpdateGroup;delete:DeleteGroup"),

				beego.NSRouter("/groups/:group/users", &api.GroupMemberController{}, "get:GroupUserList;post:AddGroupUsers"),
				beego.NSRouter("/groups/:group/users/:user", &api.GroupMemberController{}, "delete:RemoveGroupUser"),
				beego.NSRouter("/groups/:group/users/:user/roles", &api.GroupMemberController{}, "get:GroupUserRoleList;post:AddGroupUserRole"),
				beego.NSRouter("/groups/:group/users/:user/roles/:role", &api.GroupMemberController{}, "delete:RemoveGroupUserRole"),
				beego.NSRouter("/groups/:group/users/:user/constraints", &api.GroupMemberController{}, "get:GetGroupUserConstraint"),
				beego.NSRouter("/groups/:group/users/:user/constraints/:resourceConstraint", &api.GroupMemberController{}, "delete:DeleteGroupUserConstraint"),
				beego.NSRouter("/groups/:group/users/:user/constraints/:resourceConstraint/values", &api.GroupMemberController{}, "post:AddGroupUserConstraintValues;put:UpdateGroupUserConstraintValues;delete:DeleteGroupUserConstraintValues"),

				beego.NSRouter("/roles", &api.RoleController{}, "get:RoleList;post:CreateRole"),

				beego.NSRouter("/roles", &api.RoleController{}, "get:RoleList;post:CreateRole"),
				beego.NSRouter("/roles/:role", &api.RoleController{}, "get:GetRole;put:UpdateRole;delete:DeleteRole"),
				beego.NSRouter("/roles/:role/operations", &api.RoleController{}, "get:RoleOperationList;post:AddRoleOperation"),
				beego.NSRouter("/roles/:role/operations/:operationID", &api.RoleController{}, "delete:RemoveRoleOperation"),
				beego.NSRouter("/groups/:group/roles/:role/bundling", &api.RoleController{}, "get:RoleBundlingList;post:RoleBundling;delete:RoleUnbundling"),

				// PipelineStage
				beego.NSRouter("/pipelines/flow/components", &api.PipelineController{}, "get:GetFlowComponents"),
				beego.NSRouter("/pipelines/flow/steps", &api.PipelineController{}, "get:GetTaskTmpls;post:GetTaskTmplsByPagination"),
				beego.NSRouter("/pipelines/flow/steps/create", &api.PipelineController{}, "post:CreateTaskTmpl"),
				beego.NSRouter("/pipelines/flow/steps/:step_id", &api.PipelineController{}, "put:UpdateTaskTmpl;delete:DeleteTaskTmpl"),

				// Integrate Settings
				beego.NSRouter("/integrate/settings", &api.IntegrateController{}, "get:GetIntegrateSettings;post:GetIntegrateSettingsByPagination"),
				beego.NSRouter("/integrate/settings/create", &api.IntegrateController{}, "post:CreateIntegrateSetting"),
				beego.NSRouter("/integrate/settings/:id", &api.IntegrateController{}, "put:UpdateIntegrateSetting;delete:DeleteIntegrateSetting"),
				beego.NSRouter("/integrate/settings/verify", &api.IntegrateController{}, "post:VerifyIntegrateSetting"),
				beego.NSRouter("/integrate/clusters", &api.IntegrateController{}, "get:GetClusterIntegrateSettings"),
				// CompileEnv
				beego.NSRouter("/integrate/compile_envs", &api.IntegrateController{}, "get:GetCompileEnvs;post:GetCompileEnvsByPagination"),
				beego.NSRouter("/integrate/compile_envs/create", &api.IntegrateController{}, "post:CreateCompileEnv"),
				beego.NSRouter("/integrate/compile_envs/:id", &api.IntegrateController{}, "put:UpdateCompileEnv;delete:DeleteCompileEnv"),

				// Git Repository
				beego.NSRouter("/repos", &api.AppController{}, "get:GetRepos"),
				beego.NSRouter("/repos/:repo_id/projects", &api.AppController{}, "post:GetGitProjectsByRepoID"),

				// Project
				beego.NSRouter("/projects", &api.ProjectController{}, "post:ProjectList"),
				beego.NSRouter("/projects/create", &api.ProjectController{}, "post:Create"),
				beego.NSRouter("/projects/:project_id", &api.ProjectController{}, "put:Update;delete:Delete;get:GetProject"),
				beego.NSRouter("/projects/:project_id/checkProjectOwner", &api.ProjectController{}, "post:CheckProjetCreator"),
				beego.NSRouter("/projects/:project_id/clusters/:cluster/apps", &api.ProjectController{}, "post:GetAppserviceList"),
				beego.NSRouter("/clusters/:cluster/namespaces/:namespace/apps/:app", &api.ProjectController{}, "get:AppInspect;delete:AppDelete"),
				beego.NSRouter("/clusters/:cluster/namespaces/:namespace/apps/:app/log", &api.ProjectController{}, "get:PodLog"),
				beego.NSRouter("/clusters/:cluster/namespaces/:namespace/apps/:app/event", &api.ProjectController{}, "get:AppEvent"),
				beego.NSRouter("/clusters/:cluster/namespaces/:namespace/apps/:app/restart", &api.ProjectController{}, "post:AppRestart"),
				beego.NSRouter("/clusters/:cluster/namespaces/:namespace/apps/:app/scale", &api.ProjectController{}, "post:AppScale"),
				beego.NSRouter("/clusters/:cluster/namespaces/:namespace/pods/:podname/containernames/:containername", &api.TerminalController{}, "get:PodTerminal"),

				// Project Setup
				beego.NSRouter("/projects/:project_id/members", &api.ProjectController{}, "get:GetProjectMembers;put:AddProjectMember"),
				beego.NSRouter("/projects/:project_id/members/:id", &api.ProjectController{}, "delete:DeleteProjectMember"),

				// Project env
				beego.NSRouter("/projects/:project_id/envs", &api.ProjectController{}, "get:GetProjectEnvs;post:GetProjectEnvsByPagination"),
				beego.NSRouter("/projects/:project_id/envs/create", &api.ProjectController{}, "post:CreateProjectEnv"),
				beego.NSRouter("/projects/:project_id/envs/:env_id", &api.ProjectController{}, "put:UpdateProjectEnv;delete:DeleteProjectEnv"),

				// Project pipeline
				beego.NSRouter("/projects/:project_id/pipelines", &api.ProjectController{}, "get:GetProjectPipelines;post:GetPipelinesByPagination"),
				beego.NSRouter("/projects/:project_id/pipelines/create", &api.ProjectController{}, "post:CreatePipeline"),
				beego.NSRouter("/projects/:project_id/pipelines/:id", &api.ProjectController{}, "get:GetProjectPipeline;put:UpdatePipelineConfig;delete:DeleteProjectPipeline"),

				// Project App
				beego.NSRouter("/projects/:project_id/apps/create", &api.ProjectController{}, "post:CreateApp"),
				beego.NSRouter("/projects/:project_id/apps", &api.ProjectController{}, "get:GetApps;post:GetAppsByPagination"),
				beego.NSRouter("/projects/:project_id/apps/:app_id/:env_id/arrange", &api.AppController{}, "get:GetArrange;post:SetArrange"),
				beego.NSRouter("/arrange/yaml/parser", &api.AppController{}, "post:ParseArrangeYaml"),
				beego.NSRouter("/projects/:project_id/apps/:app_id/branches", &api.AppController{}, "post:GetAppBranches"),
				beego.NSRouter("/projects/:project_id/apps/:app_id/syncBranches", &api.AppController{}, "post:SyncAppBranches"),
				beego.NSRouter("/projects/:project_id/apps/:project_app_id", &api.ProjectController{}, "get:ProjectAppInfo;patch:SwitchProjectBranch;put:UpdateProjectApp;delete:DeleteProjectApp"),

				// Project stats
				beego.NSRouter("/projects/:project_id/publish/stats", &api.PipelineController{}, "post:GetPublishStats"),

				// Publish-Order / release
				beego.NSRouter("/projects/:project_id/publishes", &api.PublishController{}, "post:PublishList"),
				beego.NSRouter("/projects/:project_id/publishes/create", &api.PublishController{}, "post:Create"),
				beego.NSRouter("/projects/:project_id/publishes/:publish_id", &api.PublishController{}, "get:GetPublish;put:ClosePublish;delete:DeletePublish"),
				beego.NSRouter("/projects/:project_id/publishes/:publish_id/apps/can_added", &api.PublishController{}, "get:CanAddedApps"),
				beego.NSRouter("/projects/:project_id/publishes/:publish_id/apps/create", &api.PublishController{}, "post:AddPublishApp"),
				beego.NSRouter("/projects/:project_id/publishes/:publish_id/apps/:publish_app_id", &api.PublishController{}, "delete:DeletePublishApp"),
				beego.NSRouter("/projects/:project_id/publishes/:publish_id/audits", &api.PublishController{}, "post:GetOpertaionLogByPagination"),
				beego.NSRouter("/projects/:project_id/publishes/:publish_id/stages/:stage_id/back-to", &api.PublishController{}, "get:GetBackTo;post:TriggerBackTo"),
				beego.NSRouter("/projects/:project_id/publishes/:publish_id/stages/:stage_id/next-stage", &api.PublishController{}, "get:GetNextStage;post:TriggerNextStage"),

				// Publish pipeline
				beego.NSRouter("/pipelines/:project_id/publishes/:publish_id/stages/:stage_id/steps/:step_name", &api.PipelineController{}, "get:GetStepInfo;post:RunStep"),
				beego.NSRouter("/pipelines/:project_id/publishes/:publish_id/stages/:stage_id/steps/:step_name/callback", &api.PipelineController{}, "post:RunStepCallback"),
				beego.NSRouter("/pipelines/stages/:stage_id/jenkins-config", &api.PipelineController{}, "get:GetJenkinsConfig"),
			))

	beego.AddNamespace(publishAPI)

	beego.Get("/health", func(ctx *context.Context) {
		ctx.Output.Body([]byte("Ok"))
	})

	beego.ErrorController(&api.ErrorController{})

	// setup panic recover
	beego.BConfig.RecoverFunc = func(ctx *context.Context) {
		err := recover()
		if err == nil {
			return
		}
		if err == beego.ErrAbort {
			return
		}
		log.Log.Critical("The request is:", ctx.Input.Method(), ctx.Input.URL())
		log.Log.Critical("Handler crashed error:", err)
		var stack string
		for i := 1; ; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			logs.Critical(fmt.Sprintf("%s:%d", file, line))
			stack = stack + fmt.Sprintf("%s:%d\n", file, line)
		}
		hasIndent := beego.BConfig.RunMode != beego.PROD
		result := api.NewErrorResult("Panic", fmt.Sprintf("%v", err), stack)
		ctx.Output.SetStatus(http.StatusInternalServerError)
		ctx.Output.JSON(result, hasIndent, false)
	}
}
