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

package migrations

import (
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/go-atomci/atomci/internal/core/pipelinemgr"
	"github.com/go-atomci/atomci/internal/core/settings"
	"github.com/go-atomci/atomci/internal/dao"
	"github.com/go-atomci/atomci/internal/middleware/log"
	"github.com/go-atomci/atomci/internal/models"
)

type Migration20220415 struct {
}

func (m Migration20220415) GetCreateAt() time.Time {
	return time.Date(2022, 4, 15, 0, 0, 0, 0, time.Local)
}

func (m Migration20220415) Upgrade(ormer orm.Ormer) error {
	// init component
	_ = initComponent()

	compileEnvs := []settings.CompileEnvReq{
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

	// init compile envs
	_ = initCompileEnvs(compileEnvs)

	// init task tmpls
	_ = initTaskTemplates()
	return nil
}

type component struct {
	Name   string
	Type   string
	Params string
}

// initComponent ..
func initComponent() error {
	components := []component{
		{
			Name: "人工卡点",
			Type: "manual",
		},
		{
			Name: "构建",
			Type: "build",
		},
		{
			Name: "部署",
			Type: "deploy",
		},
	}
	for _, comp := range components {
		pipelineModel := dao.NewPipelineStageModel()
		_, err := pipelineModel.GetFlowComponentByType(comp.Type)
		if err != nil {
			if err == orm.ErrNoRows {
				component := &models.FlowComponent{
					Addons: models.NewAddons(),
					Name:   comp.Name,
					Type:   comp.Type,
					Params: comp.Params,
				}
				if err := pipelineModel.CreateFlowComponent(component); err != nil {
					log.Log.Warn("when init component, occur error: %s", err.Error())
				}
			} else {
				log.Log.Warn("when init component, occur error: %s", err.Error())
			}
		} else {
			log.Log.Debug("component type `%s` already exists, skip", comp.Type)
		}
	}
	return nil
}

func initCompileEnvs(compileEnvs []settings.CompileEnvReq) error {
	settingModel := dao.NewSysSettingModel()
	for _, item := range compileEnvs {
		_, err := settingModel.GetCompileEnvByName(item.Name)
		if err != nil {
			if err == orm.ErrNoRows {
				component := &models.CompileEnv{
					Addons:      models.NewAddons(),
					Name:        item.Name,
					Image:       item.Image,
					Command:     item.Command,
					Creator:     "admin", // create use 'admin'
					Args:        item.Args,
					Description: item.Description,
				}
				if err := settingModel.CreateCompileEnv(component); err != nil {
					log.Log.Warn("when init compile env, occur error: %s", err.Error())
				}
			} else {
				log.Log.Warn("init compile env occur error: %s", err.Error())
			}
		} else {
			log.Log.Debug("init compile env `%s` already exists, skip", item.Name)
		}
	}
	return nil
}

func initTaskTemplates() error {

	taskTmpls := []pipelinemgr.TaskTmplReq{
		{
			Name:        "应用构建",
			Type:        "build",
			Description: "用于应用构建",
			SubTask: []pipelinemgr.SubTask{
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

	pipeline := pipelinemgr.NewPipelineManager()

	for _, item := range taskTmpls {
		_, err := pipeline.GetTaskTmplByName(item.Name)
		if err != nil {
			if err == orm.ErrNoRows {
				if err := pipeline.CreateTaskTmpl(&item, "admin"); err != nil {
					log.Log.Error("when init task template, occur error: %s", err.Error())
				}
			} else {
				logs.Warn("init task template occur error: %s", err.Error())
			}
		} else {
			log.Log.Debug("init task template `%s` already exists, skip", item.Name)
		}
	}
	return nil
}
