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
	"github.com/go-atomci/atomci/constant"
	"github.com/go-atomci/atomci/dao"
	"github.com/go-atomci/atomci/middleware/log"
	"github.com/go-atomci/atomci/models"
	"github.com/go-atomci/atomci/utils"

	"golang.org/x/crypto/bcrypt"
)

// 初始化系统用户组
func InitAdminUserAndGroup() error {
	// 初始化系统用户组
	groupId, err := initSystemGroup()
	if err != nil {
		log.Log.Warning(err.Error())
	}

	// 初始化系统管理员
	userId, err := initAdminUser()
	if err != nil {
		log.Log.Warning(err.Error())
	}

	if _, err := dao.InsertGroupUserRel(groupId, userId); err != nil {
		log.Log.Warning(err.Error())
	}

	// 初始化系统角色，创建系统管理员
	if err := initSystemRole(); err != nil {
		log.Log.Warning(err.Error())
	}

	return nil
}

func initSystemGroup() (int64, error) {
	group, _ := dao.GetGroupByName(constant.SystemGroup)
	if group == nil {
		return dao.InsertGroup(&models.Group{
			Group:       constant.SystemGroup,
			Level:       "system",
			ParentId:    0,
			Description: "系统用户组",
		})
	}
	return group.ID, nil
}

func generateDefaultPassword() (string, error) {
	var hash []byte
	var err error
	if hash, err = bcrypt.GenerateFromPassword([]byte(constant.AdminDefaultPassword), bcrypt.DefaultCost); err != nil {
		return "", err
	}
	return string(hash), nil
}

func initAdminUser() (int64, error) {
	user, _ := dao.GetUser(constant.SystemAdminUser)
	password, err := generateDefaultPassword()
	if err != nil {
		return 0, err
	}
	if user == nil {
		return dao.CreateUser(&models.User{
			User:     constant.SystemAdminUser,
			Name:     constant.SystemAdminUser,
			Token:    utils.MakeToken(),
			Password: password,
		})
	}
	return user.ID, nil
}

// 初始化系统角色和管理员用户
func initSystemRole() error {
	roles := []models.GroupRoleReq{
		{
			Group:       constant.SystemGroup,
			Role:        constant.SystemAdminRole,
			Description: "超级管理员",
			Operations:  []string{"*"},
		},
		{
			Group:       constant.SystemGroup,
			Role:        constant.SystemMemberRole,
			Description: "普通成员",
			// TODO: change to real resouce operation
			Operations: []string{"CreateProject"},
		},
		{
			Group:       constant.SystemGroup,
			Role:        constant.CompanyAdminRole,
			Description: "项目管理员",
			// TODO: change to real resouce operation
			Operations: []string{"CreateProject"},
		},
	}
	for _, role := range roles {
		if _, err := dao.CreateGroupRole(&role); err != nil {
			log.Log.Warning(err.Error())
		}
	}

	if err := dao.GroupRoleBundling(&models.GroupRoleBundlingReq{
		Group: constant.SystemGroup,
		Role:  constant.SystemAdminRole,
		Users: []string{constant.SystemAdminUser},
	}); err != nil {
		log.Log.Warning(err.Error())
		return err
	}
	return nil
}
