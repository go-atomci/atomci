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

// 资源类型
type ResourceController struct {
	BaseController
}

// ResourceTypeList ..
func (r *ResourceController) ResourceTypeList() {
	conValues, err := r.ResourceTypeConValues()
	if err != nil {
		r.HandleInternalServerError(err.Error())
		log.Log.Error("Get resource type list error: %s", err.Error())
		return
	}
	rsp, err := dao.ResourceTypeListFilter(conValues)
	if err != nil {
		r.HandleInternalServerError(err.Error())
		log.Log.Error("Get resource type list error: %s", err.Error())
		return
	}

	r.Data["json"] = NewResult(true, rsp, "")
	r.ServeJSON()
}

// ResourceOperationsList ..
func (r *ResourceController) ResourceOperationsList() {
	rsp, err := dao.GetResourceOperationList()
	if err != nil {
		r.HandleInternalServerError(err.Error())
		log.Log.Error("Get resource type list error: %s", err.Error())
		return
	}

	r.Data["json"] = NewResult(true, rsp, "")
	r.ServeJSON()
}

// CreateResourceType ..
func (r *ResourceController) CreateResourceType() {
	var req models.ResourceTypeReq
	r.DecodeJSONReq(&req)

	if err := req.Verify(); err != nil {
		r.HandleBadRequest(err.Error())
		log.Log.Error("Create resource type error: %s", err.Error())
		return
	}

	rsp, err := dao.CreateResourceType(req.ResourceType, req.Description)
	if err != nil {
		r.HandleInternalServerError(err.Error())
		log.Log.Error("Create resource type error: %s", err.Error())
		return
	}

	r.Data["json"] = NewResult(true, rsp, "")
	r.ServeJSON()
}

// GetResourceType ..
func (r *ResourceController) GetResourceType() {
	operationConValues, err := r.ResourceOperationConValues()
	if err != nil {
		r.HandleInternalServerError(err.Error())
		log.Log.Error("Get resource type list error: %s", err.Error())
		return
	}
	constraintConValues, err := r.ResourceConstraintConValues()
	if err != nil {
		r.HandleInternalServerError(err.Error())
		log.Log.Error("Get resource type list error: %s", err.Error())
		return
	}
	resourceType := r.GetStringFromPath(":resourceType")
	rsp, err := dao.GetResourceTypeDetail(resourceType, operationConValues, constraintConValues)
	if err != nil {
		r.HandleInternalServerError(err.Error())
		log.Log.Error("Get resource type error: %s", err.Error())
		return
	}

	r.Data["json"] = NewResult(true, rsp, "")
	r.ServeJSON()
}

// UpdateResourceType ..
func (r *ResourceController) UpdateResourceType() {
	resourceType := r.GetStringFromPath(":resourceType")
	var req models.ResourceTypeReq
	r.DecodeJSONReq(&req)
	req.ResourceType = resourceType

	if err := req.Verify(); err != nil {
		r.HandleBadRequest(err.Error())
		log.Log.Error("Update resource type error: %s", err.Error())
		return
	}

	if err := dao.UpdateResourceType(req.ResourceType, req.Description); err != nil {
		r.HandleInternalServerError(err.Error())
		log.Log.Error("Update resource type error: %s", err.Error())
		return
	}

	r.Data["json"] = NewResult(true, nil, "")
	r.ServeJSON()
}

// DeleteResourceType ..
func (r *ResourceController) DeleteResourceType() {
	resourceType := r.GetStringFromPath(":resourceType")

	if err := dao.DeleteResourceType(resourceType); err != nil {
		r.HandleInternalServerError(err.Error())
		log.Log.Error("Delete resource type error: %s", err.Error())
		return
	}
	r.Data["json"] = NewResult(true, nil, "")
	r.ServeJSON()
}

// AddResourceOperation ..
func (r *ResourceController) AddResourceOperation() {
	resourceType := r.GetStringFromPath(":resourceType")
	var req models.ResourceOperationReq
	r.DecodeJSONReq(&req)
	req.ResourceType = resourceType

	if err := req.Verify(); err != nil {
		r.HandleBadRequest(err.Error())
		log.Log.Error("Create resource operation error: %s", err.Error())
		return
	}

	err := dao.AddResourceOperation(req.ResourceType, req.ResourceOperation, req.Description)
	if err != nil {
		r.HandleInternalServerError(err.Error())
		log.Log.Error("Create resource operation error: %s", err.Error())
		return
	}

	r.Data["json"] = NewResult(true, nil, "")
	r.ServeJSON()
}

// UpdateResourceOperation ..
func (r *ResourceController) UpdateResourceOperation() {
	resourceType := r.GetStringFromPath(":resourceType")
	resourceOperation := r.GetStringFromPath(":resourceOperation")

	var req models.ResourceOperationReq
	r.DecodeJSONReq(&req)
	req.ResourceType = resourceType
	req.ResourceOperation = resourceOperation

	if err := req.Verify(); err != nil {
		r.HandleBadRequest(err.Error())
		log.Log.Error("Update resource operation error: %s", err.Error())
		return
	}

	if err := dao.UpdateResourceOperation(req.ResourceType, req.ResourceOperation, req.Description); err != nil {
		r.HandleInternalServerError(err.Error())
		log.Log.Error("Update resource operation error: %s", err.Error())
		return
	}

	r.Data["json"] = NewResult(true, nil, "")
	r.ServeJSON()
}

// DeleteResourceOperation ..
func (r *ResourceController) DeleteResourceOperation() {
	resourceType := r.GetStringFromPath(":resourceType")
	resourceOperation := r.GetStringFromPath(":resourceOperation")

	if err := dao.DeleteResourceOperation(resourceType, resourceOperation); err != nil {
		r.HandleInternalServerError(err.Error())
		log.Log.Error("Delete resource operation error: %s", err.Error())
		return
	}
	r.Data["json"] = NewResult(true, nil, "")
	r.ServeJSON()
}

// AddResourceConstraint ..
func (r *ResourceController) AddResourceConstraint() {
	resourceType := r.GetStringFromPath(":resourceType")
	var req models.ResourceConstraintReq
	r.DecodeJSONReq(&req)
	req.ResourceType = resourceType

	if err := req.Verify(); err != nil {
		r.HandleBadRequest(err.Error())
		log.Log.Error("Create resource constraint error: %s", err.Error())
		return
	}

	err := dao.AddResourceConstraint(req.ResourceType, req.ResourceConstraint, req.Description)
	if err != nil {
		r.HandleInternalServerError(err.Error())
		log.Log.Error("Create resource constraint error: %s", err.Error())
		return
	}

	r.Data["json"] = NewResult(true, nil, "")
	r.ServeJSON()
}

// UpdateResourceConstraint ..
func (r *ResourceController) UpdateResourceConstraint() {
	resourceType := r.GetStringFromPath(":resourceType")
	resourceConstraint := r.GetStringFromPath(":resourceConstraint")

	var req models.ResourceConstraintReq
	r.DecodeJSONReq(&req)
	req.ResourceType = resourceType
	req.ResourceConstraint = resourceConstraint

	if err := req.Verify(); err != nil {
		r.HandleBadRequest(err.Error())
		log.Log.Error("Update resource constraint error: %s", err.Error())
		return
	}

	if err := dao.UpdateResourceConstraint(req.ResourceType, req.ResourceConstraint, req.Description); err != nil {
		r.HandleInternalServerError(err.Error())
		log.Log.Error("Update resource constraint error: %s", err.Error())
		return
	}

	r.Data["json"] = NewResult(true, nil, "")
	r.ServeJSON()
}

// DeleteResourceConstraint ..
func (r *ResourceController) DeleteResourceConstraint() {
	resourceType := r.GetStringFromPath(":resourceType")
	resourceConstraint := r.GetStringFromPath(":resourceConstraint")

	if err := dao.DeleteResourceConstraint(resourceType, resourceConstraint); err != nil {
		r.HandleInternalServerError(err.Error())
		log.Log.Error("Delete resource constraint error: %s", err.Error())
		return
	}
	r.Data["json"] = NewResult(true, nil, "")
	r.ServeJSON()
}
