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

package models

import (
	"github.com/go-atomci/atomci/utils/validate"
)

type ResourceType struct {
	Addons
	ResourceType        string                `orm:"column(resource_type);index; size(128);" json:"resource_type"`
	Description         string                `orm:"column(description);size(128)" json:"description"`
	ResourceOperations  []*ResourceOperation  `orm:"-" json:"resource_operations"`
	ResourceConstraints []*ResourceConstraint `orm:"-" json:"resource_constraints"`
}

func (t *ResourceType) TableName() string {
	return "sys_resource_type"
}

func (u *ResourceType) TableUnique() [][]string {
	return [][]string{
		{"ResourceType"},
		{"Description"},
	}
}

type ResourceTypeReq struct {
	ResourceType string `json:"resource_type"`
	Description  string `json:"description"`
}

func (v *ResourceTypeReq) Verify() error {
	v.ResourceType = validate.FormatString(v.ResourceType)
	v.Description = validate.FormatString(v.Description)
	if err := validate.ValidateName(v.ResourceType); err != nil {
		return err
	}
	if err := validate.ValidateDescription(v.Description); err != nil {
		return err
	}
	return nil
}

// ResourceOperation ..
type ResourceOperation struct {
	Addons
	ResourceType      string `orm:"column(resource_type);index" json:"resource_type"`
	ResourceOperation string `orm:"column(resource_operation)" json:"resource_operation"`
	Description       string `orm:"column(description)" json:"description"`
}

// TableName ..
func (t *ResourceOperation) TableName() string {
	return "sys_resource_operation"
}

// TableIndex ..
func (t *ResourceOperation) TableIndex() [][]string {
	return [][]string{
		{"ResourceType", "ResourceOperation"},
	}
}

func (u *ResourceOperation) TableUnique() [][]string {
	return [][]string{
		{"ResourceType", "ResourceOperation"},
		{"ResourceType", "Description"},
	}
}

type ResourceOperationReq struct {
	ResourceType      string `json:"resource_type"`
	ResourceOperation string `json:"resource_operation"`
	Description       string `json:"description"`
}

func (v *ResourceOperationReq) Verify() error {
	v.ResourceOperation = validate.FormatString(v.ResourceOperation)
	v.Description = validate.FormatString(v.Description)
	if err := validate.ValidateName(v.ResourceOperation); err != nil {
		return err
	}
	if err := validate.ValidateDescription(v.Description); err != nil {
		return err
	}
	return nil
}

type ResourceConstraint struct {
	Addons
	ResourceType       string `orm:"column(resource_type);index" json:"resource_type"`
	ResourceConstraint string `orm:"column(resource_constraint)" json:"resource_constraint"`
	Description        string `orm:"column(description)" json:"description"`
}

func (t *ResourceConstraint) TableName() string {
	return "sys_resource_constraint"
}

func (t *ResourceConstraint) TableIndex() [][]string {
	return [][]string{
		{"ResourceType", "ResourceConstraint"},
	}
}

func (u *ResourceConstraint) TableUnique() [][]string {
	return [][]string{
		{"ResourceType", "ResourceConstraint"},
		{"ResourceType", "Description"},
	}
}

type ResourceConstraintReq struct {
	ResourceType       string `json:"resource_type"`
	ResourceConstraint string `json:"resource_constraint"`
	Description        string `json:"description"`
}

func (v *ResourceConstraintReq) Verify() error {
	v.ResourceConstraint = validate.FormatString(v.ResourceConstraint)
	v.Description = validate.FormatString(v.Description)
	if err := validate.ValidateName(v.ResourceConstraint); err != nil {
		return err
	}
	if err := validate.ValidateDescription(v.Description); err != nil {
		return err
	}
	return nil
}

type ResourceReq struct {
	ResourceType        ResourceTypeReq
	ResourceOperations  []ResourceOperationReq
	ResourceConstraints []ResourceConstraintReq
}

type BatchResourceTypeReq struct {
	Resources []ResourceReq
}
