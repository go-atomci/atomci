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

package command

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/astaxie/beego"
	"github.com/spf13/cobra"
)

var (
	httpClient = &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).Dial,
			MaxIdleConns:        200,
			MaxIdleConnsPerHost: 200,
			IdleConnTimeout:     30 * time.Second,
			TLSHandshakeTimeout: 5 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 60 * time.Second,
	}
)

func sentHTTPRequest(method, urlStr, adminToken string, body io.Reader) ([]byte, error) {
	rep, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		beego.Error(fmt.Sprintf("Url: %v Method: %v Error: %v", urlStr, method, err.Error()))
		return nil, err
	}
	if adminToken == "" {
		beego.Error("you missed token, eg './cli init --token=[token-value]'")
		return nil, err
	}

	rep.Header.Set("Content-Type", "application/json")
	rep.Header.Set("Authorization", "Bearer "+adminToken)
	resp, err := httpClient.Do(rep)
	if err != nil {
		beego.Error(fmt.Sprintf("Url: %v Method: %v Error: %v", urlStr, method, err.Error()))
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		return respBody, nil
	} else if resp.StatusCode == http.StatusCreated {
		return respBody, nil
	}

	beego.Error(fmt.Errorf("Url: %v Method: %v StatusCode: %v respBody: %v", urlStr, method, resp.StatusCode, string(respBody)))
	return nil, fmt.Errorf(string(respBody))
}

type ResourceReq struct {
	Resources []BatchResourceTypeSpec `json:"resources"`
}

type BatchResourceTypeSpec struct {
	ResourceType       []string   `json:"resource_type"`
	ResourceOperation  [][]string `json:"resource_operation"`
	ResourceConstraint [][]string `json:"resource_constraint"`
}

var (
	HTTPPort, _ = beego.AppConfig.Int("default::httpport")
)

// init resource
func initResource(token string) error {
	resourceReq := ResourceReq{
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
					[]string{"*", "所有操作"},
					[]string{"UserLogin", "用户登录"},
					[]string{"UserLogout", "用户登出"},
					[]string{"GetCurrentUser", "获取当前用户信息"},
				},
				ResourceConstraint: [][]string{},
			},
			BatchResourceTypeSpec{
				ResourceType: []string{"audit", "操作审计"},
				ResourceOperation: [][]string{
					[]string{"*", "所有操作"},
					[]string{"AuditList", "获取操作审计列表"},
				},
				ResourceConstraint: [][]string{},
			},
			BatchResourceTypeSpec{
				ResourceType: []string{"user", "用户"},
				ResourceOperation: [][]string{
					[]string{"*", "所有操作"},
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

			// project pipeline
			BatchResourceTypeSpec{
				ResourceType: []string{"pipeline", "流程"},
				ResourceOperation: [][]string{
					[]string{"*", "所有操作"},
					[]string{"PipelineList", "获取流程列表"},
					[]string{"PipelineListByPagination", "获取流程分页列表"},
					[]string{"PipelineListByPagination", "恢复默认流程"},
					[]string{"PipelineCreate", "创建流程"},
					[]string{"PipelineUpdate", "更新流程基础信息"},
					[]string{"PipelineDelete", "删除流程"},
					[]string{"FlowComponentList", "获取基础组件列表"},
					[]string{"FlowStepList", "获取步骤列表"},
					[]string{"FlowStepListByPagination", "获取步骤分页列表"},
					[]string{"FlowStepCreate", "创建步骤"},
					[]string{"FlowStepUpdate", "更新步骤"},
					[]string{"FlowStepDelete", "删除步骤"},
				},
				ResourceConstraint: [][]string{},
			},
			BatchResourceTypeSpec{
				ResourceType: []string{"project", "项目"},
				ResourceOperation: [][]string{
					[]string{"*", "所有操作"},
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
					[]string{"ProjectPipelineInfo", "获取项目绑定流程信息"},
					[]string{"ProjectAppServiceStats", "获取项目应用统计"},
				},
				ResourceConstraint: [][]string{
					[]string{"project_id", "项目ID"},
				},
			},
			BatchResourceTypeSpec{
				ResourceType: []string{"publish", "流水线"},
				ResourceOperation: [][]string{
					[]string{"*", "所有操作"},
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
				ResourceType: []string{"system", "系统回调"},
				ResourceOperation: [][]string{
					[]string{"*", "所有操作"},
					[]string{"UpdateJobBuildResult", "更新构建结果"},
				},
				ResourceConstraint: [][]string{},
			},
		},
	}

	method := "POST"
	if HTTPPort == 0 {
		HTTPPort = 8080
	}
	urlStr := fmt.Sprintf("http://127.0.0.1:%v/atomci/api/v1/init/resource", HTTPPort)
	jsonData, err := json.Marshal(resourceReq)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(jsonData)
	if _, err := sentHTTPRequest(method, urlStr, token, body); err != nil {
		return err
	}
	return nil
}

type GatewayReq struct {
	Routers [][]string `json:"routers"`
}

func initGateWayRoute(token string) error {
	gaetwayReq := GatewayReq{
		Routers: [][]string{
			[]string{"atomci/api/v1/login", "POST", "atomci", "auth", "UserLogin"},
			[]string{"atomci/api/v1/logout", "GET", "atomci", "auth", "UserLogout"},
			[]string{"atomci/api/v1/getCurrentUser", "GET", "atomci", "auth", "GetCurrentUser"},
			[]string{"atomci/api/v1/audit", "GET", "atomci", "audit", "AuditList"},
			[]string{"atomci/api/v1/init/users", "POST", "atomci", "init", "InitUsers"},
			[]string{"atomci/api/v1/init/groups", "POST", "atomci", "init", "InitGroups"},
			[]string{"atomci/api/v1/init/resource", "POST", "atomci", "init", "InitResource"},
			[]string{"atomci/api/v1/init/gateway/:backend", "POST", "atomci", "init", "InitGateway"},
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

			// pipelines
			[]string{"atomci/api/v1/pipelines/flow/components", "GET", "atomci", "pipeline", "FlowComponentList"},
			[]string{"atomci/api/v1/pipelines/flow/steps", "GET", "atomci", "pipeline", "FlowStepList"},
			[]string{"atomci/api/v1/pipelines/flow/steps", "POST", "atomci", "pipeline", "FlowStepListByPagination"},
			[]string{"atomci/api/v1/pipelines/flow/steps/create", "POST", "atomci", "pipeline", "FlowStepCreate"},
			[]string{"atomci/api/v1/pipelines/flow/steps/:step_id", "PUT", "atomci", "pipeline", "FlowStepUpdate"},
			[]string{"atomci/api/v1/pipelines/flow/steps/:step_id", "DELETE", "atomci", "pipeline", "FlowStepDelete"},
			[]string{"atomci/api/v1/pipelines", "GET", "atomci", "pipeline", "PipelineList"},
			[]string{"atomci/api/v1/pipelines", "POST", "atomci", "pipeline", "PipelineListByPagination"},
			[]string{"atomci/api/v1/pipelines/reset", "POST", "atomci", "pipeline", "ResetDefaultPipeline"},
			[]string{"atomci/api/v1/pipelines/create", "POST", "atomci", "pipeline", "PipelineCreate"},
			[]string{"atomci/api/v1/pipelines/:pipeline_id", "PUT", "atomci", "pipeline", "PipelineUpdate"},
			[]string{"atomci/api/v1/pipelines/:pipeline_id", "DELETE", "atomci", "pipeline", "PipelineDelete"},

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
			[]string{"atomci/api/v1/projects/:project_id/pipelines", "PUT", "atomci", "project", "BindProjectPipeline"},
			[]string{"atomci/api/v1/projects/:project_id/pipelines/:id", "DELETE", "atomci", "project", "DeleteBindPipeline"},
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
			[]string{"atomci/api/v1/projects/:project_id/publish/stats", "POST", "atomci", "project", "ProjectPublishStats"},
			[]string{"atomci/api/v1/projects/:project_id/pipelines/:id", "GET", "atomci", "project", "ProjectPipelineInfo"},
			// TODO: change to project env
			// []string{"atomci/api/v1/pipelines/flow/stages", "GET", "atomci", "pipeline", "FlowStageList"},
			// []string{"atomci/api/v1/pipelines/flow/stages", "POST", "atomci", "pipeline", "FlowStageListByPagination"},
			// []string{"atomci/api/v1/pipelines/flow/stages/create", "POST", "atomci", "pipeline", "FlowStageCreate"},
			// []string{"atomci/api/v1/pipelines/flow/stages/:stage_id", "PUT", "atomci", "pipeline", "FlowStageUpdate"},
			// []string{"atomci/api/v1/pipelines/flow/stages/:stage_id", "DELETE", "atomci", "pipeline", "FlowStageDelete"},

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
		},
	}

	method := "POST"
	if HTTPPort == 0 {
		HTTPPort = 8080
	}

	urlStr := fmt.Sprintf("http://127.0.0.1:%v/atomci/api/v1/init/gateway/atomci", HTTPPort)
	jsonData, err := json.Marshal(gaetwayReq)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(jsonData)
	if _, err := sentHTTPRequest(method, urlStr, token, body); err != nil {
		return err
	}
	return nil
}

func initUsers(token string) error {
	method := "POST"
	if HTTPPort == 0 {
		HTTPPort = 8080
	}
	urlStr := fmt.Sprintf("http://127.0.0.1:%v/atomci/api/v1/init/users", HTTPPort)
	if _, err := sentHTTPRequest(method, urlStr, token, nil); err != nil {
		return err
	}
	return nil
}

type CompileEnvReq struct {
	Name        string `json:"name,omitempty"`
	Image       string `json:"image,omitempty"`
	Command     string `json:"command,omitempty"`
	Args        string `json:"args,omitempty"`
	Description string `json:"description,omitempty"`
}

func initCompileEnvs(token string) error {
	method := "POST"
	if HTTPPort == 0 {
		HTTPPort = 8080
	}

	compileEnvs := []CompileEnvReq{
		{
			Name:        "jnlp",
			Image:       "colynn/jenkins-jnlp-agent:latest",
			Description: "",
		},
		{
			Name:    "kaniko",
			Image:   "colynn/kaniko-executor:debug",
			Command: "/bin/sh -c",
			Args:    "cat",
		},
		{
			Name:        "node",
			Image:       "node:12.12-alpine",
			Description: "nodejs编译环境",
		},
		{
			Name:    "maven",
			Image:   "maven:3.8.2-openjdk-8",
			Command: "/bin/sh -c",
			Args:    "cat",
		},
	}

	jsonData, err := json.Marshal(compileEnvs)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(jsonData)

	urlStr := fmt.Sprintf("http://127.0.0.1:%v/atomci/api/v1/init/compileenvs", HTTPPort)
	if _, err := sentHTTPRequest(method, urlStr, token, body); err != nil {
		return err
	}
	return nil
}

type tasktemplate struct {
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	SubTask     []subTask `json:"sub_task"`
}
type subTask struct {
	Index int    `json:"index,omitempty"`
	Name  string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"`
}

func initTaskTemplates(token string) error {
	method := "POST"
	if HTTPPort == 0 {
		HTTPPort = 8080
	}

	compileEnvs := []tasktemplate{
		{
			Name:        "应用构建",
			Type:        "build",
			Description: "用于应用构建",
			SubTask: []subTask{
				{
					Index: 1,
					Type:  "checkout",
					Name:  "检出代码",
				},
				{
					Index: 2,
					Type:  "compile",
					Name:  "编译",
				},
				{
					Index: 3,
					Type:  "build-image",
					Name:  "制作镜像",
				},
			},
		},
		{
			Name:        "应用部署",
			Type:        "deploy",
			Description: "用于应用部署健康检查",
		},
		{
			Name:        "人工卡点",
			Type:        "manual",
			Description: "人工卡点",
		},
	}

	jsonData, err := json.Marshal(compileEnvs)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(jsonData)

	urlStr := fmt.Sprintf("http://127.0.0.1:%v/atomci/api/v1/init/tasktmpls", HTTPPort)
	if _, err := sentHTTPRequest(method, urlStr, token, body); err != nil {
		return err
	}
	return nil
}

func init() {
	var token string
	var initCmd = &cobra.Command{
		Use:   "init passport",
		Short: "--token=[token] init resource operation/router",
		Long:  `init resource operation/router`,
		Run: func(cmd *cobra.Command, args []string) {
			// 注册应用资源
			initResource(token)
			// 初始化网关
			initGateWayRoute(token)
			// 更新所有用户权限策略
			initUsers(token)

			// initCompileEnv
			initCompileEnvs(token)

			// init task tmpls
			initTaskTemplates(token)
		},
	}
	initCmd.Flags().StringVarP(&token, "token", "t", "", "User token")
	rootCmd.AddCommand(initCmd)

}
