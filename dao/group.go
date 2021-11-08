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

	"github.com/go-atomci/atomci/models"

	"github.com/astaxie/beego/orm"
)

func GroupList() ([]*models.Group, error) {
	groupList := []*models.Group{}
	if _, err := GetOrmer().QueryTable("sys_group").OrderBy("-create_at").All(&groupList); err != nil {
		return nil, err
	}
	return groupList, nil
}

func GroupListByFilter(groups []string) ([]*models.Group, error) {
	groupList := []*models.Group{}
	querySeter := GetOrmer().QueryTable("sys_group").OrderBy("-create_at")
	if len(groups) != 0 {
		querySeter = querySeter.Filter("group__in", groups)
	}
	if _, err := querySeter.All(&groupList); err != nil {
		return nil, err
	}
	return groupList, nil
}

func GroupListByUserId(userId int64) ([]*models.Group, error) {
	groupList := []*models.Group{}
	sql := `select t.* from sys_group t where exists (select 1 from sys_group_user_rel t1 where t.id = t1.group_id and t1.user_id = ?)`
	ormer := GetOrmer()
	if _, err := ormer.Raw(sql, userId).QueryRows(&groupList); err != nil {
		return nil, err
	}
	return groupList, nil
}

func GetGroupByName(groupName string) (*models.Group, error) {
	var group models.Group
	if err := GetOrmer().QueryTable("sys_group").
		Filter("group", groupName).One(&group); err != nil {
		return nil, err
	}
	return &group, nil
}

func GetGroupById(id int64) (*models.Group, error) {
	var group models.Group
	if err := GetOrmer().QueryTable("sys_group").Filter("id", id).One(&group); err != nil {
		return nil, err
	}
	return &group, nil
}

func GetGroupDetailByName(groupName string) (*models.Group, error) {
	var group models.Group
	if err := GetOrmer().QueryTable("sys_group").
		Filter("group", groupName).One(&group); err != nil {
		return nil, err
	}

	// 获取用户组关联的用户列表
	users, err := GroupUserList(groupName)
	if err != nil {
		return nil, err
	}
	// 获取组角色列表
	roles, err := GroupRoleList(groupName)
	if err != nil {
		return nil, err
	}

	group.Users = users
	group.Roles = roles

	return &group, nil
}

func InsertGroup(group *models.Group) (int64, error) {
	return GetOrmer().Insert(group)
}

func UpdateGroup(groupName, description string) error {
	sql := `update sys_group set description=? where ` + "`group`" + `=? and description<>?`
	_, err := GetOrmer().Raw(sql, description, groupName, description).Exec()
	return err
}

func DeleteGroup(groupName string) error {
	// 删除用户组关联的用户
	sql := `delete from sys_group_role_user where ` + "`group`" + `=?`
	if _, err := GetOrmer().Raw(sql, groupName).Exec(); err != nil {
		return err
	}
	sql = `delete from sys_group_user_constraint where ` + "`group`" + `=?`
	if _, err := GetOrmer().Raw(sql, groupName).Exec(); err != nil {
		return err
	}
	// 删除组角色
	sql = `delete from sys_group_role where ` + "`group`" + `=?`
	if _, err := GetOrmer().Raw(sql, groupName).Exec(); err != nil {
		return err
	}
	sql = `delete from sys_group_role_operation where ` + "`group`" + `=?`
	if _, err := GetOrmer().Raw(sql, groupName).Exec(); err != nil {
		return err
	}

	// 删除用户组
	sql = `delete from sys_group where ` + "`group`" + `=?`
	if _, err := GetOrmer().Raw(sql, groupName).Exec(); err != nil {
		return err
	}

	return nil
}

// 用户组关联用户
func GroupUserList(group string) ([]*models.GroupRoleUserRsp, error) {
	// TODO: fix group by error colynn.
	// Error 1055: Expression #1 of SELECT list is not in GROUP BY clause
	// and contains nonaggregated column 'atomci.pgu.id' which is not functionally dependent
	// on columns in GROUP BY clause;
	// this is incompatible with sql_mode=only_full_group_by
	sql := `select pgu.id,pgu.create_at,pgu.update_at,pgu.user,pg.name,pg.email from sys_group_role_user as pgu join sys_user as pg on 
			pgu.user=pg.user where pgu.group=? group by pgu.user order by pgu.create_at desc;`
	users := []*models.GroupRoleUserRsp{}
	if _, err := GetOrmer().Raw(sql, group).QueryRows(&users); err != nil {
		return nil, err
	}
	for _, user := range users {
		sql = `select pgu.id,pgu.create_at,pgu.update_at,pgu.role,pgr.description from sys_group_role_user as pgu join 
				sys_group_role as pgr on pgu.role=pgr.role and pgu.group=pgr.group where pgu.group=? and pgu.user=? order by pgu.create_at desc`
		if _, err := GetOrmer().Raw(sql, group, user.User).QueryRows(&user.Roles); err != nil {
			return nil, err
		}
	}
	return users, nil
}

func GetGroupUserRoles(group, user string) ([]*models.UserGroupRole, error) {
	roles := []*models.UserGroupRole{}
	sql := `select * from sys_group_role_user as a 
		inner join sys_group_role as b on a.group = b.group and a.role = b.role 
		where a.group = ? and a.user = ?`
	if _, err := GetOrmer().Raw(sql, group, user).QueryRows(&roles); err != nil {
		if err != orm.ErrNoRows {
			return nil, err
		}
	}
	return roles, nil
}

func AddGroupUsers(groupUsers []*models.GroupRoleUser) error {
	sql := `insert ignore into sys_group_role_user(` + "`group`" + `,user,role) values(?,?,?)`
	for _, user := range groupUsers {
		if _, err := GetOrmer().Raw(sql, user.Group, user.User, user.Role).Exec(); err != nil {
			return err
		}
	}
	return nil
}

func RemoveGroupUsers(group string, users []string) error {
	if _, err := GetOrmer().QueryTable("sys_group_role_user").
		Filter("group", group).Filter("user__in", users).Delete(); err != nil {
		return err
	}
	return nil
}

func GetGroupUserConstraint(group, user string) (map[string][]string, error) {
	constraints := []*models.GroupUserConstraint{}
	if _, err := GetOrmer().QueryTable("sys_group_user_constraint").
		Filter("group", group).Filter("user", user).
		All(&constraints); err != nil {
		return nil, err
	}
	res := map[string][]string{}
	for _, con := range constraints {
		if _, ok := res[con.Constraint]; ok {
			res[con.Constraint] = append(res[con.Constraint], con.Value)
		} else {
			res[con.Constraint] = []string{con.Value}
		}
	}
	return res, nil
}

func DeleteGroupUserConstraint(group, user, constraint string) error {
	sql := `delete from sys_group_user_constraint where ` + "`group`" + `=? and user=? and ` + "`constraint`" + `=?`
	_, err := GetOrmer().Raw(sql, group, user, constraint).Exec()
	return err
}

func AddGroupUserConstraintValues(group, user, constraint string, conValues []string) error {
	for _, val := range conValues {
		sql := `insert ignore into sys_group_user_constraint(` + "`group`" + `,user,` + "`constraint`" + `,value) values(?,?,?,?)`
		if _, err := GetOrmer().Raw(sql, group, user, constraint, val).Exec(); err != nil {
			return err
		}
	}
	return nil
}

func UpdateGroupUserConstraintValues(group, user, constraint string, conValues []string) error {
	values := ""
	for index, val := range conValues {
		if index == 0 {
			values = fmt.Sprintf("'%v'", val)
		} else {
			values = values + "," + fmt.Sprintf("'%v'", val)
		}
	}
	delValues := []string{}
	sql := `select value from sys_group_user_constraint where ` + "`group`" + `=? and user=? and ` + "`constraint`" + `=? and value not in (` + values + `)`
	if _, err := GetOrmer().Raw(sql, group, user, constraint).QueryRows(&delValues); err != nil {
		return err
	}
	if len(delValues) > 0 {
		values = ""
		for index, val := range delValues {
			if index == 0 {
				values = fmt.Sprintf("'%v'", val)
			} else {
				values = values + "," + fmt.Sprintf("'%v'", val)
			}
		}
		sql = `delete from sys_group_user_constraint where ` + "`group`" + `=? and user=? and ` + "`constraint`" + `=? and value in (` + values + `)`
		if _, err := GetOrmer().Raw(sql, group, user, constraint).Exec(); err != nil {
			return err
		}
	}
	values = ""
	for index, val := range conValues {
		if index == 0 {
			values = fmt.Sprintf("('%v','%v','%v','%v')", group, user, constraint, val)
		} else {
			values = values + "," + fmt.Sprintf("('%v','%v','%v','%v')", group, user, constraint, val)
		}
	}
	sql = `insert ignore into sys_group_user_constraint(` + "`group`" + `,user,` + "`constraint`" + `,value) values` + values
	if _, err := GetOrmer().Raw(sql).Exec(); err != nil {
		return err
	}

	return nil
}

func DeleteGroupUserConstraintValues(group, user, constraint string, conValues []string) error {
	for _, val := range conValues {
		sql := `delete from sys_group_user_constraint where ` + "`group`" + `=? and user=? and ` + "`constraint`" + `=? and value=?`
		if _, err := GetOrmer().Raw(sql, group, user, constraint, val).Exec(); err != nil {
			return err
		}
	}

	return nil
}

func GetGroupUserConstraintWithFilter(group, user string, constraintList []string) (map[string][]string, error) {
	res := map[string][]string{}
	if len(constraintList) == 0 {
		return res, nil
	}
	constraints := []*models.GroupUserConstraint{}
	if _, err := GetOrmer().QueryTable("sys_group_user_constraint").
		Filter("group", group).
		Filter("user", user).
		Filter("constraint__in", constraintList).
		GroupBy("constraint", "value").
		All(&constraints); err != nil {
		return nil, err
	}
	for _, con := range constraints {
		if _, ok := res[con.Constraint]; ok {
			res[con.Constraint] = append(res[con.Constraint], con.Value)
		} else {
			res[con.Constraint] = []string{con.Value}
		}
	}
	return res, nil
}

func InsertGroupUserRel(groupId int64, userId int64) (int64, error) {
	return GetOrmer().InsertOrUpdate(&models.GroupUserRel{
		GroupId: groupId,
		UserId:  userId,
		Addons:  models.NewAddons(),
	})
}
