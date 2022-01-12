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

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/go-atomci/atomci/internal/core/settings"
	"github.com/go-atomci/atomci/internal/dao"
	"github.com/go-atomci/atomci/internal/middleware/log"
	"github.com/go-atomci/atomci/internal/models"
)

var compileEnvs = []settings.CompileEnvReq{
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

func initCompileEnvs() error {
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
					log.Log.Error("when init compile env, occur error: %s", err.Error())
					continue
				}
			} else {
				logs.Warn("init compile env occur error: %s", err.Error())
				continue
			}
		} else {
			log.Log.Debug("component type `%s` already exists, skip", item.Name)
			continue
		}
	}
	return nil
}
