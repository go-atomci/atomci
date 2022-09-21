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
	"encoding/json"
	"net/http"
	"time"

	"github.com/astaxie/beego"

	"github.com/go-atomci/atomci/constant"
	"github.com/go-atomci/atomci/internal/dao"
	"github.com/go-atomci/atomci/internal/middleware"
	"github.com/go-atomci/atomci/internal/models"
	"github.com/go-atomci/atomci/utils"

	mycasbin "github.com/go-atomci/atomci/internal/middleware/casbin"
	"github.com/go-atomci/atomci/internal/middleware/log"

	"github.com/go-atomci/atomci/pkg/auth"
	"github.com/go-atomci/atomci/pkg/auth/ldap"
	"github.com/go-atomci/atomci/pkg/auth/local"
)

// AuthController .operations about login/logout
type AuthController struct {
	beego.Controller
}

// LoginReq ..
type LoginReq struct {
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	LoginType int    `json:"login_type,omitempty"`
}

// LdapUserInfo ..
type LdapUserInfo struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	User  string `json:"user,omitempty"`
}

// Authenticate ..
func (a *AuthController) Authenticate() {
	req := LoginReq{}
	err := json.Unmarshal(a.Ctx.Input.CopyBody(1<<32), &req)
	if err != nil {
		log.Log.Error("Invalid json request: " + err.Error())
		a.CustomAbort(http.StatusBadRequest, "Invalid json request: "+err.Error())
	}

	var loginProvider auth.Provider
	switch req.LoginType {
	case models.LocalAuth:
		userModel, err := dao.GetUser(req.Username)
		if err != nil {
			log.Log.Error("get user error: " + err.Error())
			a.CustomAbort(http.StatusBadRequest, "用户不存在或密码错误")
		}
		loginProvider = local.NewProvider(
			local.Name(userModel.Name),
			local.Email(userModel.Email),
			local.User(userModel.User),
			local.Password(userModel.Password),
		)
	case models.LDAPAuth:
		port, _ := beego.AppConfig.Int("ldap::port")
		loginProvider = ldap.NewProvider(
			ldap.BaseDN(beego.AppConfig.String("ldap::baseDN")),
			ldap.Host(beego.AppConfig.String("ldap::host")),
			ldap.Port(port),
			ldap.BindDN(beego.AppConfig.String("ldap::bindDN")),
			ldap.BindPassword(beego.AppConfig.String("ldap::bindPassword")),
			ldap.UserFilter(beego.AppConfig.String("ldap::userFilter")),
		)
	default:
		log.Log.Error("login_type is %v, not support", req.LoginType)
		http.Error(a.Ctx.ResponseWriter, "不支持此类型的登录，请联系管理员", http.StatusInternalServerError)
		return
	}

	externalAccountInfo, authErr := loginProvider.Authenticate(req.Username, req.Password)
	if authErr == nil {
		log.Log.Debug("externalAccountInfo user: %s", externalAccountInfo.User)
	} else {
		log.Log.Error("login authenticate error: %v", authErr.Error())
		http.Error(a.Ctx.ResponseWriter, "用户不存在或密码错误", http.StatusInternalServerError)
		return
	}

	// init default user and role constant.SystemMemberRole
	_, err = createOrUpdateUser(externalAccountInfo, req.LoginType)
	if err != nil {
		log.Log.Error(err.Error())
		http.Error(a.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	e, err := mycasbin.NewCasbin()
	if err != nil {
		log.Log.Error("add user role, new casbin instance error: %s", err.Error())
		http.Error(a.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := e.AddRoleForUser(externalAccountInfo.User, constant.SystemMemberRole); err != nil {
		log.Log.Error("add %v user %v error: %s", constant.SystemMemberRole, externalAccountInfo.User, err.Error())
	}

	if err := e.SavePolicy(); err != nil {
		log.Log.Error("save casbin policy error: %s", err.Error())
	}

	// TODO: change role name `admin` to real
	token, err := middleware.JwtAuth(externalAccountInfo.User, "admin")
	if err != nil {
		log.Log.Error(err.Error())
		http.Error(a.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	a.Data["json"] = NewSuccessResult(map[string]interface{}{"user": externalAccountInfo.User, "token": token})
	a.ServeJSON()
}

func createOrUpdateUser(userAttributes *auth.ExternalAccount, loginType int) (*models.User, error) {
	userName := userAttributes.User
	email := userAttributes.Email
	realName := userAttributes.Name
	user, _ := dao.GetUser(userName)
	token := utils.MakeToken()

	if user == nil {
		user = &models.User{
			User:      userName,
			Email:     email,
			Name:      realName,
			LoginType: loginType,
			Token:     token,
		}
		dao.InitSystemMember(user)
	} else {
		user.Email = email
		user.Name = realName
		user.LoginType = loginType
		user.LastLoginTime, _ = time.Parse("2006-01-02 15:04:05", time.Now().Local().Format("2006-01-02 15:04:05"))
		dao.UpdateUser(user)
	}
	return user, nil
}

// Logout ..
func (a *AuthController) Logout() {
	a.DestroySession()
	a.Data["json"] = NewSuccessResult()
	a.ServeJSON()
}
