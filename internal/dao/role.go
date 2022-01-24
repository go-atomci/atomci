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
	mycasbin "github.com/go-atomci/atomci/internal/middleware/casbin"
	"github.com/go-atomci/atomci/internal/middleware/log"
	"github.com/go-atomci/atomci/internal/models"
)

// 组角色
func GroupRoleList(group string) ([]*models.GroupRole, error) {
	roles := []*models.GroupRole{}
	sql := `select * from sys_group_role where ` + "`group`" + `=? order by create_at desc`
	if _, err := GetOrmer().Raw(sql, group).QueryRows(&roles); err != nil {
		return nil, err
	}
	return roles, nil
}

func GetGroupRoleByName(group, role string) (*models.GroupRole, error) {
	var groupRole models.GroupRole
	if err := GetOrmer().QueryTable("sys_group_role").
		Filter("group", group).Filter("role", role).Filter("role", role).
		One(&groupRole); err != nil {
		return nil, err
	}
	users, err := GroupRoleBundlingList(group, role)
	if err != nil {
		return nil, err
	}
	groupRole.Users = users

	return &groupRole, nil
}

func GetRoleOperationsByRoleName(role string) ([]*models.GroupRoleOperation, error) {
	roleOperations := []*models.GroupRoleOperation{}
	if _, err := GetOrmer().QueryTable("sys_group_role_operation").
		Filter("group", "system").
		Filter("role", role).
		All(&roleOperations); err != nil {
		return nil, err
	}
	return roleOperations, nil
}

func CreateGroupRole(req *models.GroupRoleReq) (*models.GroupRole, error) {
	role, _ := GetGroupRoleByName(req.Group, req.Role)
	if role == nil {
		sql := `insert into sys_group_role(` + "`group`" + `,role,description) values(?,?,?)`
		if _, err := GetOrmer().Raw(sql, req.Group, req.Role, req.Description).Exec(); err != nil {
			return nil, err
		}
	}

	if err := AddRoleOperation(&models.GroupRoleOperationReq{
		Group:      req.Group,
		Role:       req.Role,
		Operations: req.Operations,
	}); err != nil {
		log.Log.Error("when create group role, add role operation error: %s", err.Error())
		return nil, err
	}

	log.Log.Debug("req operations length: %v", len(req.Operations))
	err := GenerateCasbinrules(req.Role, req.Operations)
	if err != nil {
		log.Log.Error("generate casbin rules error: %s", err.Error())
	}

	role, err = GetGroupRoleByName(req.Group, req.Role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func UpdateGroupRole(req *models.GroupRoleReq) error {
	sql := `update sys_group_role set description=? where ` + "`group`" + `=? and role=? and description<>?`
	_, err := GetOrmer().Raw(sql, req.Description, req.Group, req.Role, req.Description).Exec()
	return err
}

func DeleteGroupRole(group, role string) error {
	sql := `delete from sys_group_role_operation where ` + "`group`" + `=? and role=?`
	if _, err := GetOrmer().Raw(sql, group, role).Exec(); err != nil {
		return err
	}
	sql = `delete from sys_group_role where ` + "`group`" + `=? and role=?`
	if _, err := GetOrmer().Raw(sql, group, role).Exec(); err != nil {
		return err
	}
	sql = `delete from sys_group_role_user where ` + "`group`" + `=? and role=?`
	if _, err := GetOrmer().Raw(sql, group, role).Exec(); err != nil {
		return err
	}
	return nil
}

// 组角色绑定用户
func GroupRoleBundlingList(group, role string) ([]*models.GroupRoleBundlingUser, error) {
	userList := []*models.GroupRoleBundlingUser{}
	sql := `select pgu.user,pu.name,pgu.id,pgu.create_at,pgu.update_at from sys_group_role_user as pgu 
			join sys_user as pu on pgu.user=pu.user where pgu.group=? and pgu.role=? order by pgu.create_at desc`
	if _, err := GetOrmer().Raw(sql, group, role).QueryRows(&userList); err != nil {
		return userList, err
	}
	return userList, nil
}

func GroupRoleBundling(req *models.GroupRoleBundlingReq) error {
	for _, user := range req.Users {
		sql := `insert ignore into sys_group_role_user(` + "`group`" + `,user,role) values(?,?,?)`
		if _, err := GetOrmer().Raw(sql, req.Group, user, req.Role).Exec(); err != nil {
			return err
		}
	}
	return nil
}

func GroupRoleUnbundling(req *models.GroupRoleBundlingReq) error {
	for _, user := range req.Users {
		sql := `delete from sys_group_role_user where ` + "`group`" + `=? and user=? and role=?`
		if _, err := GetOrmer().Raw(sql, req.Group, user, req.Role).Exec(); err != nil {
			return err
		}
	}
	return nil
}

func AddRoleOperation(req *models.GroupRoleOperationReq) error {
	if len(req.Operations) > 0 {
		values := ""
		for index, operationID := range req.Operations {
			if index == 0 {
				values = fmt.Sprintf("('%v','%v',%v)", req.Group, req.Role, operationID)
			} else {
				values = values + "," + fmt.Sprintf("('%v','%v',%v)", req.Group, req.Role, operationID)
			}
		}
		// TODO: add casbin items;
		sql := `insert ignore into sys_group_role_operation(` + "`group`" + `,role, operation_id) values` + values
		if _, err := GetOrmer().Raw(sql).Exec(); err != nil {
			return err
		}

		if err := GenerateCasbinrules(req.Role, req.Operations); err != nil {
			return err
		}
	}
	return nil
}

func DeleteGroupRolePolicy(req *models.GroupRoleOperationReq) error {
	if len(req.Operations) > 0 {
		values := ""
		for index, police := range req.Operations {
			if index == 0 {
				values = fmt.Sprintf("'%v'", police)
			} else {
				values = values + "," + fmt.Sprintf("'%v'", police)
			}
		}
		sql := `delete from sys_group_role_operation where ` + `role=? and operation_id in (` + values + `)`
		if _, err := GetOrmer().Raw(sql, req.Role).Exec(); err != nil {
			return err
		}
		if err := DeleteCasbinrules(req.Role, req.Operations); err != nil {
			return err
		}
	}
	return nil
}

func DeleteCasbinrules(role string, operations []int64) error {
	resourceRouterItems, err := getRouterItembyOperationsID(operations)
	if err != nil {
		return err
	}
	if len(resourceRouterItems) > 0 {
		casbinRules := generateCasbinRules(resourceRouterItems, role)
		e, err := mycasbin.NewCasbin()
		if err != nil {
			log.Log.Error("new casbin instance error: %s", err.Error())
			return err
		}
		log.Log.Debug("role: %s, casbin rules length: %v", role, len(casbinRules))
		removeFlag, err := e.RemovePolicies(casbinRules)
		if err != nil {
			log.Log.Error("remove policys error: %s", err.Error())
		}
		log.Log.Info("remove policy to casbin rule, flag: %v", removeFlag)
		if err := e.SavePolicy(); err != nil {
			log.Log.Error("save casbin policy error: %s", err.Error())
			return err
		}
	}
	return nil
}

func GenerateCasbinrules(role string, operations []int64) error {
	resourceRouterItems, err := getRouterItembyOperationsID(operations)
	if err != nil {
		return err
	}
	if len(resourceRouterItems) > 0 {
		casbinRules := generateCasbinRules(resourceRouterItems, role)
		e, err := mycasbin.NewCasbin()
		if err != nil {
			log.Log.Error("new casbin instance error: %s", err.Error())
			return err
		}
		log.Log.Debug("role: %s, casbin rules length: %v", role, len(casbinRules))

		// it seems beegormadapter doesn't implement batch adapter, replaced AddPolicies by AddPolicy
		for _, value := range casbinRules {
			_, err = e.AddPolicy(value)
			if err != nil {
				log.Log.Error("add policys error: %s", err.Error())
			}
		}
		//addFlag, err := e.AddPolicies(casbinRules)
		if err != nil {
			log.Log.Error("add policys error: %s", err.Error())
		}
		//log.Log.Info("add policy to casbin rule, flag: %v", addFlag)
		if err := e.SavePolicy(); err != nil {
			log.Log.Error("save casbin policy error: %s", err.Error())
			return err
		}
	}
	return nil
}

func getRouterItembyOperationsID(operationsID []int64) ([]*models.GatewayRouter, error) {
	resOperationItems, err := GetResourceOperationByIDs(operationsID)
	if err != nil {
		log.Log.Error("when get resource operation by ids occur error: %s", err.Error())
		return nil, err
	}

	resTypeOperationsMapping := orderByResourceType(resOperationItems)
	resourceRouterItems := []*models.GatewayRouter{}
	for key, item := range resTypeOperationsMapping {
		resRouterItems, err := GetResourceRouterItems(key, item)
		if err != nil {
			log.Log.Error("when create group role, get resource router items error: %s", err.Error())
			continue
		}
		resourceRouterItems = append(resourceRouterItems, resRouterItems...)
	}
	return resourceRouterItems, nil

}
func orderByResourceType(res []*models.ResourceOperation) map[string][]string {
	resMap := map[string][]string{}
	for _, item := range res {
		resMap[item.ResourceType] = append(resMap[item.ResourceType], item.ResourceOperation)
	}
	return resMap
}
