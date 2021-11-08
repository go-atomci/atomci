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
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"

	"github.com/go-atomci/atomci/common"
	"github.com/go-atomci/atomci/dao"
	"github.com/go-atomci/atomci/middleware"
	"github.com/go-atomci/atomci/middleware/log"
	"github.com/go-atomci/atomci/models"
	"github.com/go-atomci/atomci/utils"
	"github.com/go-atomci/atomci/utils/errors"
	"github.com/go-atomci/atomci/utils/query"
)

// BaseController wraps common methods for controllers to host API
type BaseController struct {
	beego.Controller
	User      string
	audit     models.Audit
	UserModel *models.User
}

// GetStringFromPath gets the param from path and returns it as string
func (b *BaseController) GetStringFromPath(key string) string {
	return b.Ctx.Input.Param(key)
}

// GetStringFromQuery gets the param from query and returns it as string
func (b *BaseController) GetStringFromQuery(key string) string {
	return b.Ctx.Input.Query(key)
}

// GetInt64FromPath gets the param from path and returns it as int64
func (b *BaseController) GetInt64FromPath(key string) (int64, error) {
	value := b.Ctx.Input.Param(key)
	return strconv.ParseInt(value, 10, 64)
}

// GetInt64FromQuery gets the param from query string and returns it as int64
func (b *BaseController) GetInt64FromQuery(key string) (int64, error) {
	value := b.Ctx.Input.Query(key)
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		convErr := err.(*strconv.NumError)
		if convErr.Num == "" {
			return -1, nil
		} else {
			return v, err
		}
	}
	return v, nil
}

// GetBoolFromQuery gets the param from query string and returns it as int64
func (b *BaseController) GetBoolFromQuery(key string) (bool, error) {
	value := b.Ctx.Input.Query(key)
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return boolValue, err
	}

	return boolValue, nil
}

// ServeResult serve result
func (b *BaseController) ServeResult(result Result) {
	b.Data["json"] = result
	b.ServeJSON()
}

// ServeError serve error
func (b *BaseController) ServeError(err error) {
	if err == nil {
		err = fmt.Errorf("nil")
	}
	var statusCode int
	var result Result
	switch srcErr := err.(type) {
	// error
	case *errors.Error:
		{
			errDetail := ""
			if srcErr.Cause() != nil {
				errDetail = srcErr.Cause().Error()
			}
			statusCode = srcErr.Status()
			result = NewErrorResult(srcErr.Code(), srcErr.Message(), errDetail)
		}
	// go error
	default:
		{
			statusCode = http.StatusInternalServerError
			result = NewErrorResult("InternalServerError", "internal server error", err.Error())
		}
	}
	b.Ctx.Output.SetStatus(statusCode)
	b.Data["json"] = result
	b.ServeJSON()
}

// HandleNotFound ...
func (b *BaseController) HandleNotFound(text string) {
	b.RenderError(http.StatusNotFound, text)
}

// HandleUnauthorized ...
func (b *BaseController) HandleUnauthorized(text string) {
	b.RenderError(http.StatusUnauthorized, fmt.Sprintf("Unauthorized: %v", text))
}

// HandleForbidden ...
func (b *BaseController) HandleForbidden(text string) {
	b.RenderError(http.StatusForbidden, fmt.Sprintf("Forbidden: %v", text))
}

// HandleBadRequest ...
func (b *BaseController) HandleBadRequest(text string) {
	b.RenderError(http.StatusBadRequest, text)
}

// HandleInternalServerError ...
func (b *BaseController) HandleInternalServerError(text string) {
	b.RenderError(http.StatusInternalServerError, text)
}

// HandleConflictError ...
func (b *BaseController) HandleConflictError(text string) {
	b.RenderError(http.StatusConflict, text)
}

// HandlePreconditionError StatusPreconditionFailed error
func (b *BaseController) HandlePreconditionError(text string) {
	b.RenderError(http.StatusPreconditionFailed, text)
}

// HandleNormalError ...
func (b *BaseController) HandleNormalError(code int, text string) {
	b.RenderError(code, text)
}

// Render returns nil as it won't render template
func (b *BaseController) Render() error {
	return nil
}

// RenderError provides shortcut to render http error
func (b *BaseController) RenderError(code int, text string) {
	http.Error(b.Ctx.ResponseWriter, text, code)
}

// DecodeJSONReq decodes a json request
func (b *BaseController) DecodeJSONReq(v interface{}) {
	err := json.Unmarshal(b.Ctx.Input.CopyBody(1<<32), v)
	if err != nil {
		log.Log.Error("Invalid json request: " + err.Error())
		b.CustomAbort(http.StatusBadRequest, "Invalid json request: "+err.Error())
	}
}

// Validate validates v if it implements interface validation.ValidFormer
func (b *BaseController) Validate(v interface{}) {
	validator := validation.Validation{}
	isValid, err := validator.Valid(v)
	if err != nil {
		b.CustomAbort(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if !isValid {
		message := ""
		for _, e := range validator.Errors {
			message += fmt.Sprintf("%s %s \n", e.Field, e.Message)
		}
		b.CustomAbort(http.StatusBadRequest, message)
	}
}

// DecodeJSONReqAndValidate does both decoding and validation
func (b *BaseController) DecodeJSONReqAndValidate(v interface{}) {
	b.DecodeJSONReq(v)
	b.Validate(v)
}

// GetFilterQuery get page filter query
func (b *BaseController) GetFilterQuery() *query.FilterQuery {
	filter := query.FilterQuery{
		IsLike: true,
	}
	b.DecodeJSONReq(&filter)

	return &filter
}

func (b *BaseController) getAuthHeader() string {
	authHeader := b.Ctx.Input.Header("Authorization")
	log.Log.Debug("auth header: %v", authHeader)
	strList := strings.Split(authHeader, " ")
	token := ""
	if len(strList) == 2 && strList[0] == "Bearer" {
		token = strList[1]
	}
	if token == "" {
		urlPath := b.Controller.Ctx.Request.URL.Path
		if strings.Contains(urlPath, "/containernames/") {
			token, err := common.GetUserToken("admin")
			if err == nil {
				return token
			} else {
				log.Log.Error("get user by namd admin error: %s", err.Error())
			}
		}
		return ""
	}
	return token
}

func (b *BaseController) getUserName() string {
	token := b.getAuthHeader()

	log.Log.Debug("get token from header: %s", token)
	user := ""
	var err error
	if strings.Contains(token, ".") {
		user, err = middleware.JwtParse(b.Controller.Ctx, token)
		if err != nil {
			return ""
		}
	} else if len(token) == 16 {
		userModel, err := dao.GetUserByToken(token)
		if err != nil {
			log.Log.Error("get user by token error: %s", err.Error())
		}
		user = userModel.User
	}
	return user
}

// 校验用户是否已登录
func (b *BaseController) checkUserLogin() string {
	username := b.getUserName()
	log.Log.Debug("username: %s", username)
	if username == "" {
		b.HandleUnauthorized("用户登录已失效，请退出后重新登录")
		return ""
	}

	user, err := dao.GetUserDetail(username)
	if err != nil {
		log.Log.Error("check user login: %v", err.Error())
		// b.HandleInternalServerError(err.Error())
		return ""
	}
	b.Data["username"] = username
	b.UserModel = user
	b.User = user.User
	return username
}

// Prepare inits security context and project manager from request
// context
func (b *BaseController) Prepare() {
	user := b.checkUserLogin()
	if user == "" {
		b.ServeError(errors.NewUnauthorized().SetCause(fmt.Errorf("user is empty, missing header maybe")))
		return
	}

	constraint := map[string]string{}
	for key, value := range b.Ctx.Input.Params() {
		if key == ":splat" || strings.Split(key, "")[0] != ":" {
			continue
		}
		constraint[strings.Split(key, ":")[1]] = value
	}

	auth, err := middleware.Authorization(b.Controller.Ctx, user)
	if err != nil {
		log.Log.Error(err.Error())
		b.HandleInternalServerError(err.Error())
		return
	} else if !auth {
		beego.Warn(fmt.Sprintf("user %v permission denied, the request path is: %v", user, b.Ctx.Request.URL.Path))
		b.HandleForbidden("permission denied")
		return
	}

	operationObject, _ := json.Marshal(constraint)
	b.audit = models.Audit{
		User:            user,
		Method:          b.Ctx.Input.Method(),
		Operation:       b.Ctx.Input.URL(),
		OperationObject: string(operationObject),
		OperationBody:   string(b.Ctx.Input.CopyBody(1 << 32)),
	}
}

func (b *BaseController) AuditBlackList(operation string) bool {
	blackList := []string{}
	for _, op := range blackList {
		if operation == op {
			return true
		}
	}
	return false
}

func (b *BaseController) Finish() {
	if b.audit.Method == "GET" || b.AuditBlackList(b.audit.Operation) || b.audit.User == "" || b.audit.Operation == "" {
		return
	}
	var status int
	if b.Ctx.ResponseWriter.Status == 0 {
		status = http.StatusOK
	} else if b.Ctx.ResponseWriter.Status == 401 || b.Ctx.ResponseWriter.Status == 403 {
		return
	} else {
		status = b.Ctx.ResponseWriter.Status
	}
	b.audit.OperationStatus = status
	b.audit.Addons = models.NewAddons()
	if err := dao.AuditInsert(&b.audit); err != nil {
		log.Log.Error(fmt.Sprintf("audit insert error: %v", err.Error()))
	}
}

func (b *BaseController) ResourceTypeConValues() ([]string, error) {
	var res []string
	conValues, err := dao.GetUserResourceConstraintValues("resource", b.User)
	if err != nil {
		log.Log.Error(err.Error())
		return nil, err
	}
	for _, val := range conValues.Values {
		if _, ok := val["*"]; ok {
			res = []string{}
			return res, nil
		} else {
			if _, ok := val["resourceType"]; ok {
				for _, value := range val["resourceType"] {
					if value == "*" {
						res = []string{}
						return res, nil
					} else {
						res = append(res, value)
					}
				}
			}
		}
	}
	return res, nil
}

func (b *BaseController) ResourceOperationConValues() ([]string, error) {
	var res []string
	conValues, err := dao.GetUserResourceConstraintValues("resource", b.User)
	if err != nil {
		log.Log.Error(err.Error())
		return nil, err
	}
	for _, val := range conValues.Values {
		if _, ok := val["*"]; ok {
			res = []string{}
			return res, nil
		} else {
			if _, ok := val["resourceOperation"]; ok {
				for _, value := range val["resourceOperation"] {
					if value == "*" {
						res = []string{}
						return res, nil
					} else {
						res = append(res, value)
					}
				}
			}
		}
	}
	return res, nil
}

func (b *BaseController) ResourceConstraintConValues() ([]string, error) {
	var res []string
	conValues, err := dao.GetUserResourceConstraintValues("resource", b.User)
	if err != nil {
		log.Log.Error(err.Error())
		return nil, err
	}
	for _, val := range conValues.Values {
		if _, ok := val["*"]; ok {
			res = []string{}
			return res, nil
		} else {
			if _, ok := val["resourceConstraint"]; ok {
				for _, value := range val["resourceConstraint"] {
					if value == "*" {
						res = []string{}
						return res, nil
					} else {
						res = append(res, value)
					}
				}
			}
		}
	}
	return res, nil
}

func (b *BaseController) GroupConValues() ([]string, error) {
	var res []string
	conValues, err := dao.GetUserResourceConstraintValues("group", b.User)
	if err != nil {
		log.Log.Error(err.Error())
		return nil, err
	}
	for _, val := range conValues.Values {
		if _, ok := val["*"]; ok {
			res = []string{}
			return res, nil
		} else {
			if _, ok := val["group"]; ok {
				for _, value := range val["group"] {
					if value == "*" {
						res = []string{}
						return res, nil
					} else {
						res = append(res, value)
					}
				}
			}
		}
	}
	return res, nil
}

func (b *BaseController) IsSysAdmin() bool {
	return dao.UserIsAdmin(b.User)
}

func (b *BaseController) IsGroupAdmin() int {
	return b.UserModel.GroupAdmin
}

func (b *BaseController) UserGroup() string {
	return "system"
}

// Projects ..
func (b *BaseController) Projects() ([]int64, error) {
	var projectIDStrs []string
	conValues, err := dao.GetUserResourceConstraintValues("project", b.User)
	if err != nil {
		log.Log.Error("when get projects, GetResourceConstraintValues occur error: %s ", err.Error())
		return nil, err
	}
	log.Log.Debug("conValues: %+v", conValues)
	var projectIDs []int64
	for _, val := range conValues.Values {
		if _, ok := val["*"]; ok {
			return b.getProjectIDs()
		}
		if _, ok := val["project_id"]; ok {
			if utils.Contains(val["project_id"], "*") {
				return b.getProjectIDs()
			}
			projectIDStrs = val["project_id"]
		} else {
			continue
		}
	}
	for _, s := range projectIDStrs {
		projectID, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			log.Log.Warn("when get project Constraint， str parse to int occur error: %s", err.Error())
			continue
		}
		projectIDs = append(projectIDs, projectID)
	}
	log.Log.Debug("project IDs: %+v", projectIDs)
	return projectIDs, nil
}

func (b *BaseController) getProjectIDs() ([]int64, error) {
	var projectIDs []int64
	projects, err := dao.NewProjectModel().GetProjects()
	if err != nil {
		log.Log.Error("when get project Constraint, get project by cid occur error: %s", err.Error())
		return nil, nil
	}
	for _, project := range projects {
		projectIDs = append(projectIDs, project.ID)
	}
	return projectIDs, nil
}
