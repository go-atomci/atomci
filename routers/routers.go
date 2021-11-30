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

	"github.com/go-atomci/atomci/controllers"
	"github.com/go-atomci/atomci/middleware/log"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
)

// Init init api router
func init() {
	publishAPI :=
		beego.NewNamespace("atomci/api",
			beego.NSNamespace("/v1",

				beego.NSRouter("/logout", &controllers.AuthController{}, "get:Logout"),
				beego.NSRouter("/login", &controllers.AuthController{}, "post:Authenticate"),
				beego.NSRouter("/getCurrentUser", &controllers.UserController{}, "get:GetCurrentUser"),

				beego.NSRouter("/audit", &controllers.AuditController{}, "get:AuditList"),

				beego.NSRouter("/init/users", &controllers.InitController{}, "post:InitUsers"),
				beego.NSRouter("/init/compileenvs", &controllers.InitController{}, "post:InitCompileEnvs"),
				beego.NSRouter("/init/tasktmpls", &controllers.InitController{}, "post:InitTaskTemplates"),
				beego.NSRouter("/init/groups", &controllers.InitController{}, "post:InitGroups"),
				beego.NSRouter("/init/resource", &controllers.InitController{}, "post:InitResource"),
				beego.NSRouter("/init/gateway/:backend", &controllers.InitController{}, "post:InitGateway"),

				beego.NSRouter("/resources", &controllers.ResourceController{}, "get:ResourceTypeList;post:CreateResourceType"),
				beego.NSRouter("/resources-operations", &controllers.ResourceController{}, "get:ResourceOperationsList"),
				beego.NSRouter("/resources/:resourceType", &controllers.ResourceController{}, "get:GetResourceType;put:UpdateResourceType;delete:DeleteResourceType"),
				beego.NSRouter("/resources/:resourceType/operations", &controllers.ResourceController{}, "post:AddResourceOperation"),
				beego.NSRouter("/resources/:resourceType/operations/:resourceOperation", &controllers.ResourceController{}, "put:UpdateResourceOperation;delete:DeleteResourceOperation"),
				beego.NSRouter("/resources/:resourceType/constraints", &controllers.ResourceController{}, "post:AddResourceConstraint"),
				beego.NSRouter("/resources/:resourceType/constraints/:resourceConstraint", &controllers.ResourceController{}, "put:UpdateResourceConstraint;delete:DeleteResourceConstraint"),

				beego.NSRouter("/users", &controllers.UserController{}, "get:UserList;post:CreateUser"),
				beego.NSRouter("/users/:user", &controllers.UserController{}, "get:GetUser;put:UpdateUser;delete:DeleteUser"),
				beego.NSRouter("/users/:user/resources/:resourceType/constraints/values", &controllers.UserController{}, "get:GetUserResourceConstraintValues"),

				beego.NSRouter("/groups", &controllers.GroupController{}, "get:GroupList"),
				beego.NSRouter("/groups/:group", &controllers.GroupController{}, "get:GetGroup;put:UpdateGroup;delete:DeleteGroup"),

				beego.NSRouter("/groups/:group/users", &controllers.GroupMemberController{}, "get:GroupUserList;post:AddGroupUsers"),
				beego.NSRouter("/groups/:group/users/:user", &controllers.GroupMemberController{}, "delete:RemoveGroupUser"),
				beego.NSRouter("/groups/:group/users/:user/roles", &controllers.GroupMemberController{}, "get:GroupUserRoleList;post:AddGroupUserRole"),
				beego.NSRouter("/groups/:group/users/:user/roles/:role", &controllers.GroupMemberController{}, "delete:RemoveGroupUserRole"),
				beego.NSRouter("/groups/:group/users/:user/constraints", &controllers.GroupMemberController{}, "get:GetGroupUserConstraint"),
				beego.NSRouter("/groups/:group/users/:user/constraints/:resourceConstraint", &controllers.GroupMemberController{}, "delete:DeleteGroupUserConstraint"),
				beego.NSRouter("/groups/:group/users/:user/constraints/:resourceConstraint/values", &controllers.GroupMemberController{}, "post:AddGroupUserConstraintValues;put:UpdateGroupUserConstraintValues;delete:DeleteGroupUserConstraintValues"),

				beego.NSRouter("/roles", &controllers.RoleController{}, "get:RoleList;post:CreateRole"),

				beego.NSRouter("/roles", &controllers.RoleController{}, "get:RoleList;post:CreateRole"),
				beego.NSRouter("/roles/:role", &controllers.RoleController{}, "get:GetRole;put:UpdateRole;delete:DeleteRole"),
				beego.NSRouter("/roles/:role/operations", &controllers.RoleController{}, "get:RoleOperationList;post:AddRoleOperation"),
				beego.NSRouter("/roles/:role/operations/:operationID", &controllers.RoleController{}, "delete:RemoveRoleOperation"),
				beego.NSRouter("/groups/:group/roles/:role/bundling", &controllers.RoleController{}, "get:RoleBundlingList;post:RoleBundling;delete:RoleUnbundling"),

				// PipelineStage
				beego.NSRouter("/pipelines/flow/components", &controllers.PipelineController{}, "get:GetFlowComponents"),
				beego.NSRouter("/pipelines/flow/steps", &controllers.PipelineController{}, "get:GetTaskTmpls;post:GetTaskTmplsByPagination"),
				beego.NSRouter("/pipelines/flow/steps/create", &controllers.PipelineController{}, "post:CreateTaskTmpl"),
				beego.NSRouter("/pipelines/flow/steps/:step_id", &controllers.PipelineController{}, "put:UpdateTaskTmpl;delete:DeleteTaskTmpl"),

				// Integrate Settings
				beego.NSRouter("/integrate/settings", &controllers.IntegrateController{}, "get:GetIntegrateSettings;post:GetIntegrateSettingsByPagination"),
				beego.NSRouter("/integrate/settings/create", &controllers.IntegrateController{}, "post:CreateIntegrateSetting"),
				beego.NSRouter("/integrate/settings/:id", &controllers.IntegrateController{}, "put:UpdateIntegrateSetting;delete:DeleteIntegrateSetting"),
				beego.NSRouter("/integrate/settings/verify", &controllers.IntegrateController{}, "post:VerifyIntegrateSetting"),
				beego.NSRouter("/clusters", &controllers.IntegrateController{}, "get:GetClusterIntegrateSettings"),
				// CompileEnv
				beego.NSRouter("/integrate/compile_envs", &controllers.IntegrateController{}, "get:GetCompileEnvs;post:GetCompileEnvsByPagination"),
				beego.NSRouter("/integrate/compile_envs/create", &controllers.IntegrateController{}, "post:CreateCompileEnv"),
				beego.NSRouter("/integrate/compile_envs/:id", &controllers.IntegrateController{}, "put:UpdateCompileEnv;delete:DeleteCompileEnv"),

				// Git Repository
				beego.NSRouter("/repos", &controllers.AppController{}, "get:GetRepos"),
				beego.NSRouter("/repos/:repo_id/projects", &controllers.AppController{}, "post:GetGitProjectsByRepoID"),

				// Project
				beego.NSRouter("/projects", &controllers.ProjectController{}, "post:ProjectList"),
				beego.NSRouter("/projects/create", &controllers.ProjectController{}, "post:Create"),
				beego.NSRouter("/projects/:project_id", &controllers.ProjectController{}, "put:Update;delete:Delete;get:GetProject"),
				beego.NSRouter("/projects/:project_id/checkProjectOwner", &controllers.ProjectController{}, "post:CheckProjetCreator"),
				beego.NSRouter("/projects/:project_id/clusters/:cluster/apps", &controllers.ProjectController{}, "post:GetAppserviceList"),
				beego.NSRouter("/clusters/:cluster/namespaces/:namespace/apps/:app", &controllers.ProjectController{}, "get:AppInspect;delete:AppDelete"),
				beego.NSRouter("/clusters/:cluster/namespaces/:namespace/apps/:app/log", &controllers.ProjectController{}, "get:PodLog"),
				beego.NSRouter("/clusters/:cluster/namespaces/:namespace/apps/:app/event", &controllers.ProjectController{}, "get:AppEvent"),
				beego.NSRouter("/clusters/:cluster/namespaces/:namespace/apps/:app/restart", &controllers.ProjectController{}, "post:AppRestart"),
				beego.NSRouter("/clusters/:cluster/namespaces/:namespace/apps/:app/scale", &controllers.ProjectController{}, "post:AppScale"),
				beego.NSRouter("/clusters/:cluster/namespaces/:namespace/pods/:podname/containernames/:containername", &controllers.TerminalController{}, "get:PodTerminal"),

				// Project Setup
				beego.NSRouter("/projects/:project_id/members", &controllers.ProjectController{}, "get:GetProjectMembers;put:AddProjectMember"),
				beego.NSRouter("/projects/:project_id/members/:id", &controllers.ProjectController{}, "delete:DeleteProjectMember"),

				// Project env
				beego.NSRouter("/projects/:project_id/envs", &controllers.ProjectController{}, "get:GetProjectEnvs;post:GetProjectEnvsByPagination"),
				beego.NSRouter("/projects/:project_id/envs/create", &controllers.ProjectController{}, "post:CreateProjectEnv"),
				beego.NSRouter("/projects/:project_id/envs/:env_id", &controllers.ProjectController{}, "put:UpdateProjectEnv;delete:DeleteProjectEnv"),

				// Project pipeline
				beego.NSRouter("/projects/:project_id/pipelines", &controllers.ProjectController{}, "get:GetProjectPipelines;post:GetPipelinesByPagination"),
				beego.NSRouter("/projects/:project_id/pipelines/create", &controllers.ProjectController{}, "post:CreatePipeline"),
				beego.NSRouter("/projects/:project_id/pipelines/:id", &controllers.ProjectController{}, "get:GetProjectPipeline;put:UpdatePipelineConfig;delete:DeleteProjectPipeline"),

				// Project App
				beego.NSRouter("/projects/:project_id/apps/create", &controllers.ProjectController{}, "post:CreateApp"),
				beego.NSRouter("/projects/:project_id/apps", &controllers.ProjectController{}, "get:GetApps;post:GetAppsByPagination"),
				beego.NSRouter("/projects/:project_id/apps/:app_id/:env_id/arrange", &controllers.AppController{}, "get:GetArrange;post:SetArrange"),
				beego.NSRouter("/arrange/yaml/parser", &controllers.AppController{}, "post:ParseArrangeYaml"),
				beego.NSRouter("/projects/:project_id/apps/:app_id/branches", &controllers.AppController{}, "post:GetAppBranches"),
				beego.NSRouter("/projects/:project_id/apps/:app_id/syncBranches", &controllers.AppController{}, "post:SyncAppBranches"),
				beego.NSRouter("/projects/:project_id/apps/:project_app_id", &controllers.ProjectController{}, "get:ProjectAppInfo;patch:SwitchProjectBranch;put:UpdateProjectApp;delete:DeleteProjectApp"),

				// Project stats
				beego.NSRouter("/projects/:project_id/publish/stats", &controllers.PipelineController{}, "post:GetPublishStats"),

				// Publish-Order / release
				beego.NSRouter("/projects/:project_id/publishes", &controllers.PublishController{}, "post:PublishList"),
				beego.NSRouter("/projects/:project_id/publishes/create", &controllers.PublishController{}, "post:Create"),
				beego.NSRouter("/projects/:project_id/publishes/:publish_id", &controllers.PublishController{}, "get:GetPublish;put:ClosePublish;delete:DeletePublish"),
				beego.NSRouter("/projects/:project_id/publishes/:publish_id/apps/can_added", &controllers.PublishController{}, "get:CanAddedApps"),
				beego.NSRouter("/projects/:project_id/publishes/:publish_id/apps/create", &controllers.PublishController{}, "post:AddPublishApp"),
				beego.NSRouter("/projects/:project_id/publishes/:publish_id/apps/:publish_app_id", &controllers.PublishController{}, "delete:DeletePublishApp"),
				beego.NSRouter("/projects/:project_id/publishes/:publish_id/audits", &controllers.PublishController{}, "post:GetOpertaionLogByPagination"),
				beego.NSRouter("/projects/:project_id/publishes/:publish_id/stages/:stage_id/back-to", &controllers.PublishController{}, "get:GetBackTo;post:TriggerBackTo"),
				beego.NSRouter("/projects/:project_id/publishes/:publish_id/stages/:stage_id/next-stage", &controllers.PublishController{}, "get:GetNextStage;post:TriggerNextStage"),

				// Publish pipeline
				beego.NSRouter("/pipelines/:project_id/publishes/:publish_id/stages/:stage_id/steps/:step_name", &controllers.PipelineController{}, "get:GetStepInfo;post:RunStep"),
				beego.NSRouter("/pipelines/:project_id/publishes/:publish_id/stages/:stage_id/steps/:step_name/callback", &controllers.PipelineController{}, "post:RunStepCallback"),
				beego.NSRouter("/pipelines/stages/:stage_id/jenkins-config", &controllers.PipelineController{}, "get:GetJenkinsConfig"),
			))

	beego.AddNamespace(publishAPI)

	beego.Get("/health", func(ctx *context.Context) {
		ctx.Output.Body([]byte("Ok"))
	})

	beego.ErrorController(&controllers.ErrorController{})

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
		result := controllers.NewErrorResult("Panic", fmt.Sprintf("%v", err), stack)
		ctx.Output.SetStatus(http.StatusInternalServerError)
		ctx.Output.JSON(result, hasIndent, false)
	}
}
