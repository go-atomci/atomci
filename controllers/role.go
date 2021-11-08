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
	"github.com/go-atomci/atomci/dao"
	"github.com/go-atomci/atomci/middleware/log"
	"github.com/go-atomci/atomci/models"
)

type RoleController struct {
	BaseController
}

// RoleList ..
func (r *RoleController) RoleList() {
	groupName := r.GetStringFromPath(":group")
	if groupName == "" {
		groupName = "system"
	}
	rsp, err := dao.GroupRoleList(groupName)
	if err != nil {
		r.HandleInternalServerError(err.Error())
		log.Log.Error("Get role list error: %s", err.Error())
		return
	}
	r.Data["json"] = NewResult(true, rsp, "")
	r.ServeJSON()
}

func (r *RoleController) GetRole() {
	groupName := r.GetStringFromPath(":group")
	roleName := r.GetStringFromPath(":role")

	rsp, err := dao.GetGroupRoleByName(groupName, roleName)
	if err != nil {
		r.HandleInternalServerError(err.Error())
		log.Log.Error("Get role error: %s", err.Error())
		return
	}
	r.Data["json"] = NewResult(true, rsp, "")
	r.ServeJSON()
}

func (r *RoleController) CreateRole() {
	var req models.GroupRoleReq
	r.DecodeJSONReq(&req)
	// group use system
	req.Group = "system"

	if err := req.Verify(); err != nil {
		r.HandleBadRequest(err.Error())
		log.Log.Error("Create role error: %s", err.Error())
		return
	}

	rsp, err := dao.CreateGroupRole(&req)
	if err != nil {
		r.HandleInternalServerError(err.Error())
		log.Log.Error("Create role error: %s", err.Error())
		return
	}

	r.Data["json"] = NewResult(true, rsp, "")
	r.ServeJSON()
}

func (r *RoleController) UpdateRole() {
	groupName := r.GetStringFromPath(":group")
	roleName := r.GetStringFromPath(":role")
	var req models.GroupRoleReq
	r.DecodeJSONReq(&req)
	req.Group = groupName
	req.Role = roleName

	if err := req.Verify(); err != nil {
		r.HandleBadRequest(err.Error())
		log.Log.Error("Update role error: %s", err.Error())
		return
	}

	if err := dao.UpdateGroupRole(&req); err != nil {
		r.HandleInternalServerError(err.Error())
		log.Log.Error("Update role error: %s", err.Error())
		return
	}

	r.Data["json"] = NewResult(true, nil, "")
	r.ServeJSON()
}

func (r *RoleController) DeleteRole() {
	groupName := r.GetStringFromPath(":group")
	roleName := r.GetStringFromPath(":role")

	if err := dao.DeleteGroupRole(groupName, roleName); err != nil {
		r.HandleInternalServerError(err.Error())
		log.Log.Error("Delete role error: %s", err.Error())
		return
	}

	r.Data["json"] = NewResult(true, nil, "")
	r.ServeJSON()
}

func (r *RoleController) RoleBundlingList() {
	groupName := r.GetStringFromPath(":group")
	roleName := r.GetStringFromPath(":role")
	rsp, err := dao.GroupRoleBundlingList(groupName, roleName)
	if err != nil {
		r.HandleInternalServerError(err.Error())
		log.Log.Error("Get role list error: %s", err.Error())
		return
	}
	r.Data["json"] = NewResult(true, rsp, "")
	r.ServeJSON()
}

func (r *RoleController) RoleBundling() {
	groupName := r.GetStringFromPath(":group")
	roleName := r.GetStringFromPath(":role")
	var req models.GroupRoleBundlingReq
	r.DecodeJSONReq(&req)
	req.Group = groupName
	req.Role = roleName

	if err := dao.GroupRoleBundling(&req); err != nil {
		r.HandleInternalServerError(err.Error())
		log.Log.Error("role bundling error: %s", err.Error())
		return
	}

	r.Data["json"] = NewResult(true, nil, "")
	r.ServeJSON()
}

func (r *RoleController) RoleUnbundling() {
	groupName := r.GetStringFromPath(":group")
	roleName := r.GetStringFromPath(":role")
	var req models.GroupRoleBundlingReq
	r.DecodeJSONReq(&req)
	req.Group = groupName
	req.Role = roleName

	if err := dao.GroupRoleUnbundling(&req); err != nil {
		r.HandleInternalServerError(err.Error())
		log.Log.Error("role unbundling error: %s", err.Error())
		return
	}

	r.Data["json"] = NewResult(true, nil, "")
	r.ServeJSON()
}

func (r *RoleController) RolePolicyList() {
	// groupName := r.GetStringFromPath(":group")
	// roleName := r.GetStringFromPath(":role")
	// TODO: need change role resources list
	rsp := []string{}
	r.Data["json"] = NewResult(true, rsp, "")
	r.ServeJSON()
}

func (r *RoleController) AddRolePolicy() {
	groupName := r.GetStringFromPath(":group")
	roleName := r.GetStringFromPath(":role")
	var req models.GroupRolePolicyReq
	r.DecodeJSONReq(&req)
	req.Group = groupName
	req.Role = roleName

	if err := dao.AddRoleOperation(&req); err != nil {
		r.HandleInternalServerError(err.Error())
		log.Log.Error("Add role policy error: %s", err.Error())
		return
	}

	r.Data["json"] = NewResult(true, nil, "")
	r.ServeJSON()
}

func (r *RoleController) RemoveRolePolicy() {
	groupName := r.GetStringFromPath(":group")
	roleName := r.GetStringFromPath(":role")
	var req models.GroupRolePolicyReq
	r.DecodeJSONReq(&req)
	req.Group = groupName
	req.Role = roleName

	if err := dao.DeleteGroupRolePolicy(&req); err != nil {
		r.HandleInternalServerError(err.Error())
		log.Log.Error("Remove role policy error: %s", err.Error())
		return
	}

	r.Data["json"] = NewResult(true, nil, "")
	r.ServeJSON()
}
