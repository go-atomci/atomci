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

package api

import (
	"fmt"

	"github.com/go-atomci/atomci/internal/dao"
	mycasbin "github.com/go-atomci/atomci/internal/middleware/casbin"
	"github.com/go-atomci/atomci/internal/middleware/log"
	"github.com/go-atomci/atomci/internal/models"
)

type GroupController struct {
	BaseController
}

type GroupClusters struct {
	Clusters []string `json:"clusters"`
}

func (g *GroupController) GroupList() {
	conValues, err := g.GroupConValues()
	if err != nil {
		g.HandleInternalServerError(err.Error())
		log.Log.Error("Get group list error: %s", err.Error())
		return
	}
	rsp, err := dao.GroupListByFilter(conValues)
	if err != nil {
		g.HandleInternalServerError(err.Error())
		log.Log.Error("Get group list error: %s", err.Error())
		return
	}
	g.Data["json"] = NewResult(true, rsp, "")
	g.ServeJSON()
}

func (g *GroupController) GetGroup() {
	groupName := g.GetStringFromPath(":group")

	rsp, err := dao.GetGroupDetailByName(groupName)
	if err != nil {
		g.HandleInternalServerError(err.Error())
		log.Log.Error("Get group error: %s", err.Error())
		return
	}
	g.Data["json"] = NewResult(true, rsp, "")
	g.ServeJSON()
}

func (g *GroupController) UpdateGroup() {
	groupName := g.GetStringFromPath(":group")
	var req models.GroupReq
	g.DecodeJSONReq(&req)
	req.Group = groupName

	if err := req.Verify(); err != nil {
		g.HandleBadRequest(err.Error())
		log.Log.Error("Update group error: %s", err.Error())
		return
	}

	if err := dao.UpdateGroup(groupName, req.Description); err != nil {
		g.HandleInternalServerError(err.Error())
		log.Log.Error("Update group error: %s", err.Error())
		return
	}
	g.Data["json"] = NewResult(true, nil, "")
	g.ServeJSON()
}

func (g *GroupController) DeleteGroup() {
	groupName := g.GetStringFromPath(":group")

	if err := dao.DeleteGroup(groupName); err != nil {
		g.HandleInternalServerError(err.Error())
		log.Log.Error("Delete group error: %s", err.Error())
		return
	}
	g.Data["json"] = NewResult(true, nil, "")
	g.ServeJSON()
}

// 组用户
type GroupMemberController struct {
	BaseController
}

func (g *GroupMemberController) GroupUserList() {
	groupName := g.GetStringFromPath(":group")
	rsp, err := dao.GroupUserList(groupName)
	if err != nil {
		g.HandleInternalServerError(err.Error())
		log.Log.Error("Get user list error: %s", err.Error())
		return
	}

	g.Data["json"] = NewResult(true, rsp, "")
	g.ServeJSON()
}

func (g *GroupMemberController) AddGroupUsers() {
	groupName := g.GetStringFromPath(":group")
	var req models.GroupRoleUserReq
	g.DecodeJSONReq(&req)
	req.Group = groupName

	users := []*models.GroupRoleUser{}
	for _, user := range req.Users {
		req.Roles = append(req.Roles, "developer")
		for _, role := range req.Roles {
			users = append(users, &models.GroupRoleUser{
				Group: req.Group,
				User:  user,
				Role:  role,
			})
		}
		dao.AddGroupUserConstraintValues(groupName, user, "bizcluster", []string{groupName})
	}
	if err := dao.AddGroupUsers(users); err != nil {
		g.HandleInternalServerError(err.Error())
		log.Log.Error("Add group user error: %s", err.Error())
		return
	}

	g.Data["json"] = NewResult(true, nil, "")
	g.ServeJSON()
}

// RemoveGroupUser ..
func (g *GroupMemberController) RemoveGroupUser() {
	groupName := g.GetStringFromPath(":group")
	userName := g.GetStringFromPath(":user")

	if err := dao.RemoveGroupUsers(groupName, []string{userName}); err != nil {
		g.HandleInternalServerError(err.Error())
		log.Log.Error("Add group user error: %s", err.Error())
		return
	}

	g.Data["json"] = NewResult(true, nil, "")
	g.ServeJSON()
}

// GroupUserRoleList ..
func (g *GroupMemberController) GroupUserRoleList() {
	groupName := g.GetStringFromPath(":group")
	userName := g.GetStringFromPath(":user")

	rsp, err := dao.GetGroupUserRoles(groupName, userName)
	if err != nil {
		g.HandleInternalServerError(err.Error())
		log.Log.Error("Get user role list error: %s", err.Error())
		return
	}
	g.Data["json"] = NewResult(true, rsp, "")
	g.ServeJSON()
}

// AddGroupUserRole ..
func (g *GroupMemberController) AddGroupUserRole() {
	// use default group system
	groupName := "system"
	userName := g.GetStringFromPath(":user")
	var req models.GroupRoleUserReq
	g.DecodeJSONReq(&req)

	users := []*models.GroupRoleUser{}
	e, err := mycasbin.NewCasbin()
	if err != nil {
		log.Log.Error("add user role, new casbin instance error: %s", err.Error())
		g.HandleInternalServerError(err.Error())
		return
	}
	for _, role := range req.Roles {
		users = append(users, &models.GroupRoleUser{
			Group: groupName,
			User:  userName,
			Role:  role,
		})
		if _, err := e.AddRoleForUser(userName, role); err != nil {
			log.Log.Error("add role user error: %s", err.Error())
		}
	}
	if err := dao.AddGroupUsers(users); err != nil {
		log.Log.Error("Add user role error: %s", err.Error())
		g.HandleInternalServerError(err.Error())
		return
	}
	if err := e.SavePolicy(); err != nil {
		log.Log.Error("save casbin policy error: %s", err.Error())
		g.HandleInternalServerError(err.Error())
		return
	}
	g.Data["json"] = NewResult(true, nil, "")
	g.ServeJSON()
}

// RemoveGroupUserRole ..
func (g *GroupMemberController) RemoveGroupUserRole() {
	groupName := "system"
	userName := g.GetStringFromPath(":user")
	roleName := g.GetStringFromPath(":role")

	if err := dao.GroupRoleUnbundling(&models.GroupRoleBundlingReq{
		Group: groupName,
		Users: []string{userName},
		Role:  roleName,
	}); err != nil {
		g.HandleInternalServerError(err.Error())
		log.Log.Error("Remove user role error: %s", err.Error())
		return
	}

	e, err := mycasbin.NewCasbin()
	if err != nil {
		log.Log.Error("add user role, new casbin instance error: %s", err.Error())
		g.HandleInternalServerError(err.Error())
		return
	}
	if _, err := e.DeleteRoleForUser(userName, roleName); err != nil {
		log.Log.Error("delete role user error: %s", err.Error())
	}

	if err := e.SavePolicy(); err != nil {
		log.Log.Error("save casbin policy error: %s", err.Error())
		g.HandleInternalServerError(err.Error())
		return
	}
	g.Data["json"] = NewResult(true, nil, "")
	g.ServeJSON()
}

// GetGroupUserConstraint ..
func (g *GroupMemberController) GetGroupUserConstraint() {
	groupName := g.GetStringFromPath(":group")
	userName := g.GetStringFromPath(":user")

	constraints, err := dao.GetGroupUserConstraint(groupName, userName)
	if err != nil {
		g.HandleInternalServerError(err.Error())
		log.Log.Error(fmt.Errorf("get user constraint error: %s", err.Error()).Error())
		return
	}

	g.Data["json"] = NewResult(true, constraints, "")
	g.ServeJSON()
}

// DeleteGroupUserConstraint ..
func (g *GroupMemberController) DeleteGroupUserConstraint() {
	groupName := g.GetStringFromPath(":group")
	userName := g.GetStringFromPath(":user")
	resourceConstraint := g.GetStringFromPath(":resourceConstraint")

	if err := dao.DeleteGroupUserConstraint(groupName, userName, resourceConstraint); err != nil {
		g.HandleInternalServerError(err.Error())
		log.Log.Error(fmt.Errorf("Delete user constraint error: %s", err.Error()).Error())
		return
	}

	g.Data["json"] = NewResult(true, nil, "")
	g.ServeJSON()
}

func (g *GroupMemberController) AddGroupUserConstraintValues() {
	groupName := g.GetStringFromPath(":group")
	userName := g.GetStringFromPath(":user")
	resourceConstraint := g.GetStringFromPath(":resourceConstraint")
	var req []string
	g.DecodeJSONReq(&req)

	if err := dao.AddGroupUserConstraintValues(groupName, userName, resourceConstraint, req); err != nil {
		g.HandleInternalServerError(err.Error())
		log.Log.Error(fmt.Errorf("Add user constraint values error: %s", err.Error()).Error())
		return
	}

	g.Data["json"] = NewResult(true, nil, "")
	g.ServeJSON()
}

func (g *GroupMemberController) UpdateGroupUserConstraintValues() {
	groupName := g.GetStringFromPath(":group")
	userName := g.GetStringFromPath(":user")
	resourceConstraint := g.GetStringFromPath(":resourceConstraint")
	var req []string
	g.DecodeJSONReq(&req)

	if err := dao.UpdateGroupUserConstraintValues(groupName, userName, resourceConstraint, req); err != nil {
		g.HandleInternalServerError(err.Error())
		log.Log.Error(fmt.Errorf("Add user constraint values error: %s", err.Error()).Error())
		return
	}

	g.Data["json"] = NewResult(true, nil, "")
	g.ServeJSON()
}

func (g *GroupMemberController) DeleteGroupUserConstraintValues() {
	groupName := g.GetStringFromPath(":group")
	userName := g.GetStringFromPath(":user")
	resourceConstraint := g.GetStringFromPath(":resourceConstraint")
	var req []string
	g.DecodeJSONReq(&req)

	if err := dao.DeleteGroupUserConstraintValues(groupName, userName, resourceConstraint, req); err != nil {
		g.HandleInternalServerError(err.Error())
		log.Log.Error(fmt.Errorf("Delete user constraint values error: %s", err.Error()).Error())
		return
	}

	g.Data["json"] = NewResult(true, nil, "")
	g.ServeJSON()
}
