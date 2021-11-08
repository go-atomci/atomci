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

package controllers

import (
	"fmt"

	"github.com/go-atomci/atomci/core/pipelinemgr"
	"github.com/go-atomci/atomci/core/settings"
	"github.com/go-atomci/atomci/dao"
	"github.com/go-atomci/atomci/initialize"
	"github.com/go-atomci/atomci/middleware/log"
	"github.com/go-atomci/atomci/models"
	"github.com/go-atomci/atomci/utils/errors"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

type InitController struct {
	BaseController
}

type InitResourceReq struct {
	Resources []initialize.BatchResourceTypeSpec `json:"resources"`
}

func (init *InitController) InitResource() {
	var req InitResourceReq
	init.DecodeJSONReq(&req)

	if err := dao.BatchCreateResourceType(initialize.ToBatchResourceTypeReq(req.Resources)); err != nil {
		init.HandleInternalServerError(err.Error())
		log.Log.Error("Init resource error: %s", err.Error())
		return
	}

	init.Data["json"] = NewResult(true, nil, "")
	init.ServeJSON()
}

type InitGatewayReq struct {
	Routers [][]string `json:"routers"`
}

func (init *InitController) InitGateway() {
	backend := init.GetStringFromPath(":backend")
	var req InitGatewayReq
	init.DecodeJSONReq(&req)

	if err := dao.DeleteGatewayRouteByBackend(backend); err != nil {
		init.HandleInternalServerError(err.Error())
		log.Log.Error("Init gateway error: %s", err.Error())
		return
	}
	for _, route := range req.Routers {
		if len(route) != 5 {
			init.HandleBadRequest(fmt.Sprintf("Invalid router parameter: %v", route))
			return
		}
		if err := dao.CreateGatewayRoute(route[0], route[1], route[2], route[3], route[4]); err != nil {
			if !errors.OrmError1062(err) {
				init.HandleInternalServerError(err.Error())
				log.Log.Error("Init gateway error: %s", err.Error())
				return
			}
		}
	}

	init.Data["json"] = NewResult(true, nil, "")
	init.ServeJSON()
}

func (init *InitController) InitUsers() {
	users, _ := dao.UserList()
	for _, user := range users {
		dao.InitSystemMember(user)
	}
	init.Data["json"] = NewResult(true, nil, "")
	init.ServeJSON()
}

func (init *InitController) InitGroups() {
	groups, err := dao.GroupList()
	if err != nil {
		init.HandleInternalServerError(err.Error())
		return
	}
	for _, group := range groups {
		log.Log.Debug("init group %v policies / roles", group.Group)
	}
	go func(groups []*models.Group) {
		for _, group := range groups {

			if group.Group != "admin" {
				// lib.InitGroupRoles(group.Group)
			}
		}
	}(groups)
	init.Data["json"] = NewResult(true, nil, "")
	init.ServeJSON()
}

func (init *InitController) InitCompileEnvs() {
	var req []settings.CompileEnvReq
	init.DecodeJSONReq(&req)

	settingModel := dao.NewSysSettingModel()
	for _, item := range req {
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
	init.Data["json"] = NewResult(true, nil, "")
	init.ServeJSON()
}

func (init *InitController) InitTaskTemplates() {
	req := []pipelinemgr.TaskTmplReq{}
	init.DecodeJSONReq(&req)
	pipeline := pipelinemgr.NewPipelineManager()

	for _, item := range req {
		_, err := pipeline.GetTaskTmplByName(item.Name)
		if err != nil {
			if err == orm.ErrNoRows {
				if err := pipeline.CreateTaskTmpl(&item, "admin"); err != nil {
					log.Log.Error("when init task template, occur error: %s", err.Error())
					continue
				}
			} else {
				logs.Warn("init task template occur error: %s", err.Error())
				continue
			}
		} else {
			log.Log.Debug("component type `%s` already exists, skip", item.Name)
			continue
		}
	}
	init.Data["json"] = NewResult(true, nil, "")
	init.ServeJSON()
}
