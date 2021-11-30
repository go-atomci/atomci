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

package dao

import (
	"fmt"

	"github.com/astaxie/beego/logs"

	"github.com/go-atomci/atomci/constant"
	"github.com/go-atomci/atomci/models"
)

// 用户
func UserList() ([]*models.User, error) {
	userList := []*models.User{}
	if _, err := GetOrmer().QueryTable("sys_user").Exclude("user", constant.SystemAdminUser).
		OrderBy("-create_at").All(&userList); err != nil {
		return nil, err
	}
	return userList, nil
}

func GetUserList(userReq *models.UserReq) []*models.User {
	userList := []*models.User{}
	qs := GetOrmer().QueryTable("sys_user")
	if userReq != nil && userReq.User != "" {
		qs = qs.Filter("user", userReq.User)
	}
	if userReq != nil && userReq.Name != "" {
		qs = qs.Filter("name", userReq.Name)
	}
	if userReq != nil && userReq.Email != "" {
		qs = qs.Filter("email", userReq.Email)
	}
	if _, err := qs.OrderBy("-create_at").All(&userList); err != nil {
		logs.Error(err.Error())
	}
	return userList
}

func UserIsAdmin(userName string) bool {
	return GetOrmer().QueryTable("sys_group_role_user").
		Filter("group", constant.SystemGroup).
		Filter("user", userName).
		Filter("role", constant.SystemAdminRole).
		Exist()
}

func GetUserDetail(userName string) (*models.User, error) {
	if !UserExist(userName) {
		return nil, fmt.Errorf("user %v does not exist", userName)
	}
	var user models.User
	if err := GetOrmer().QueryTable("sys_user").Filter("user", userName).One(&user); err != nil {
		return nil, err
	}

	// 获取用户关联的组列表
	groups := []*models.Group{}
	var err error
	if userName == "admin" {
		groups, err = GroupList()
		if err != nil {
			logs.Error("group: %v error: %s", 1, err.Error())
		}
	} else {
		groups, err = GroupListByUserId(user.ID)
		if err != nil {
			logs.Error("group: %v error: %s", 1, err.Error())
		}
	}

	userGroups := []*models.Group{}
	userGroupRoles := []*models.UserGroupRole{}
	for _, group := range groups {
		userGroups = append(userGroups, group)
		roles, err := GetGroupUserRoles(group.Group, userName)
		if err != nil {
			logs.Error(err.Error())
		}
		for _, role := range roles {
			if role.Role == constant.SystemAdminRole {
				user.Admin = 1
				user.GroupAdmin = 1
				break
			}
			if role.Role == constant.DevAdminRole {
				user.GroupAdmin = 1
			}
		}
		userGroupRoles = append(userGroupRoles, roles...)

	}
	user.Groups = userGroups
	user.UserGroups = userGroupRoles

	return &user, nil
}

func GetUser(userName string) (*models.User, error) {
	if !UserExist(userName) {
		return nil, fmt.Errorf("user %v does not exist", userName)
	}
	var user models.User
	if err := GetOrmer().QueryTable("sys_user").
		Filter("user", userName).One(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(user *models.User) (int64, error) {
	return GetOrmer().Insert(user)
}

func UpdateUser(user *models.User) error {
	if !UserExist(user.User) {
		return fmt.Errorf("user %v does not exist", user.User)
	}
	_, err := GetOrmer().Update(user)
	return err
}

func DeleteUser(user *models.User) error {
	// 删除用户关联的组
	sql := `delete from sys_group_role_user where user=?`
	if _, err := GetOrmer().Raw(sql, user.User).Exec(); err != nil {
		return err
	}
	sql = `delete from sys_group_user_constraint where user=?`
	if _, err := GetOrmer().Raw(sql, user.User).Exec(); err != nil {
		return err
	}
	sql = `delete from sys_user where user=?`
	if _, err := GetOrmer().Raw(sql, user.User).Exec(); err != nil {
		return err
	}
	return nil
}

func GetUserByToken(token string) (*models.User, error) {
	var user models.User
	if err := GetOrmer().QueryTable("sys_user").
		Filter("token", token).One(&user); err != nil {
		return nil, err
	}
	if UserIsAdmin(user.User) {
		user.Admin = 1
	}
	return &user, nil
}

func UserExist(user string) bool {
	exist := GetOrmer().QueryTable("sys_user").Filter("user", user).Exist()
	return exist
}

// InitSystemMember create system user and role
func InitSystemMember(user *models.User) error {
	if !UserExist(user.User) {
		userId, err := CreateUser(user)
		if err != nil {
			logs.Error("create user error: %s", err.Error())
			return err
		}
		group, _ := GetGroupByName(constant.SystemGroup)
		groupId := group.ID
		if _, err := InsertGroupUserRel(groupId, userId); err != nil {
			logs.Error("insert group user rel group id: %v, user id: %v,  error: %s", groupId, userId, err.Error())
			return err
		}

		if err := GroupRoleBundling(&models.GroupRoleBundlingReq{
			Group: constant.SystemGroup,
			Role:  constant.SystemMemberRole,
			Users: []string{user.User},
		}); err != nil {
			logs.Error("group role bondling: %s", err.Error())
			return err
		}
	}
	return nil
}

//InitGroupUser init user and role
func InitGroupUser(user *models.User, group *models.Group, role string) error {
	if _, err := InsertGroupUserRel(group.ID, user.ID); err != nil {
		logs.Error(err.Error())
	}

	if err := GroupRoleBundling(&models.GroupRoleBundlingReq{
		Group: group.Group,
		Role:  role,
		Users: []string{user.User},
	}); err != nil {
		logs.Error(err.Error)
		return err
	}
	return nil
}

type UserResourceConstraintValues struct {
	Values []map[string][]string `json:"values"`
}

// GetUserResourceConstraintValues ..
func GetUserResourceConstraintValues(resourceType, user string) (UserResourceConstraintValues, error) {
	conValues := UserResourceConstraintValues{
		Values: []map[string][]string{},
	}
	// TODO: need rewrite
	resourceConstraintkey, err := GetResourceConstraintList(resourceType)
	if err != nil {
		logs.Error("get resource constraintlist error: %s", err.Error())
		return conValues, err
	}

	values, err := GetUserConstraintByKey(user, resourceConstraintkey)
	if err != nil {
		logs.Error("get user constraint by key error: %s", err.Error())
		return conValues, err
	}
	conValues.Values = append(conValues.Values, values)
	return conValues, nil
}
