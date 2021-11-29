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
	"golang.org/x/crypto/bcrypt"

	"github.com/go-atomci/atomci/dao"
	"github.com/go-atomci/atomci/middleware/log"
	"github.com/go-atomci/atomci/models"
	"github.com/go-atomci/atomci/utils"
)

type UserController struct {
	BaseController
}

// GetCurrentUser ..
func (u *UserController) GetCurrentUser() {
	user := u.User
	res, err := dao.GetUserDetail(user)
	if err != nil {
		u.HandleInternalServerError(err.Error())
		log.Log.Error("Get user error: %s", err.Error())
		return
	}
	u.Data["json"] = NewResult(true, res, "")
	u.ServeJSON()
}

// UserList ..
func (u *UserController) UserList() {
	var req models.UserReq
	if u.Ctx.Input.RequestBody != nil && len(u.Ctx.Input.RequestBody) > 0 {
		u.DecodeJSONReq(&req)
	}
	res := dao.GetUserList(&req)
	for _, ires := range res {
		ires.Token = ""
	}
	u.Data["json"] = NewResult(true, res, "")
	u.ServeJSON()
}

// CreateUser ..
func (u *UserController) CreateUser() {
	var req models.UserReq
	u.DecodeJSONReq(&req)

	if err := req.Verify(); err != nil {
		u.HandleBadRequest(err.Error())
		log.Log.Error("Create user error error: %s", err.Error())
		return
	}

	// TODO: confirm req.password whether match
	passwordHash, err := generatePassword(req.Password)
	if err != nil {
		u.HandleBadRequest(err.Error())
		log.Log.Error("generate password hash error: %s", err.Error())
	}
	user := models.User{
		User:      req.User,
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(passwordHash),
		LoginType: models.LocalAuth,
		Token:     utils.MakeToken(),
	}
	if err := dao.InitSystemMember(&user); err != nil {
		u.HandleInternalServerError(err.Error())
		log.Log.Error("Create user error: %s", err.Error())
		return
	}

	u.Data["json"] = NewResult(true, nil, "")
	u.ServeJSON()
}

// GetUser ..
func (u *UserController) GetUser() {
	user := u.GetStringFromPath(":user")

	res, err := dao.GetUserDetail(user)
	if err != nil {
		u.HandleInternalServerError(err.Error())
		log.Log.Error("Get user error: %s", err.Error())
		return
	}
	u.Data["json"] = NewResult(true, res, "")
	u.ServeJSON()
}

// UpdateUser ..
func (u *UserController) UpdateUser() {
	user := u.GetStringFromPath(":user")
	var req models.UserReq
	u.DecodeJSONReq(&req)
	req.User = user

	if err := req.Verify(); err != nil {
		u.HandleBadRequest(err.Error())
		log.Log.Error("Update user error error: %s", err.Error())
		return
	}

	oldUser, err := dao.GetUser(user)
	if err != nil {
		u.HandleInternalServerError(err.Error())
		log.Log.Error("Update user error: %s", err.Error())
		return
	}
	if len(req.Password) > 0 {
		passwordHash, err := generatePassword(req.Password)
		if err != nil {
			u.HandleBadRequest(err.Error())
			log.Log.Error("generate password hash error: %s", err.Error())
		}
		oldUser.Password = string(passwordHash)
	}

	oldUser.Name = req.Name
	oldUser.Email = req.Email
	if err := dao.UpdateUser(oldUser); err != nil {
		u.HandleInternalServerError(err.Error())
		log.Log.Error("Update user error: %s", err.Error())
		return
	}

	u.Data["json"] = NewResult(true, nil, "")
	u.ServeJSON()
}

// DeleteUser ..
func (u *UserController) DeleteUser() {
	userName := u.GetStringFromPath(":user")

	user, err := dao.GetUser(userName)
	if err != nil {
		u.HandleInternalServerError(err.Error())
		log.Log.Error("Delete user error: %s", err.Error())
		return
	}

	if err := dao.DeleteUser(user); err != nil {
		u.HandleInternalServerError(err.Error())
		log.Log.Error("Delete user error: %s", err.Error())
		return
	}
	u.Data["json"] = NewResult(true, nil, "")
	u.ServeJSON()
}

// GetUserResourceConstraintValues ..
func (u *UserController) GetUserResourceConstraintValues() {
	userName := u.GetStringFromPath(":user")
	resourceType := u.GetStringFromPath(":resourceType")

	rsp, err := dao.GetUserResourceConstraintValues(resourceType, userName)
	if err != nil {
		u.HandleInternalServerError(err.Error())
		log.Log.Error("Get user resource constraint values error: %s", err.Error())
		return
	}
	u.Data["json"] = NewResult(true, rsp, "")
	u.ServeJSON()
}

func generatePassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
