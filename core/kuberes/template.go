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

package kuberes

import (
	"net/http"

	"github.com/astaxie/beego/orm"

	"github.com/go-atomci/atomci/dao"
	"github.com/go-atomci/atomci/models"
	"github.com/go-atomci/atomci/utils/errors"
)

type TemplateInfo struct {
	models.CaasTemplate
	CreateAt string `json:"create_at"`
	UpdateAt string `json:"update_at"`
	DeleteAt string `json:"delete_at"`
}

type TemplateRes struct {
	modelHandle *dao.TemplateModel
	listNSFunc  NamespaceListFunction
}

func NewTemplateRes(get NamespaceListFunction) *TemplateRes {
	return &TemplateRes{
		modelHandle: dao.NewTemplateModel(),
		listNSFunc:  get,
	}
}

// template interface, nativetemplate support this interface
type Template interface {
	Default(envID int64) Template
	Validate() error
	GetExample() []byte
	Deploy(projectid, envID int64, cluster, namespace, tname string, eparam *ExtensionParam) error
}

func NewTemplate() Template {
	return NewNativeTemplate()
}

func (tr *TemplateRes) CreateTemplate(template models.CaasTemplate) (*models.CaasTemplate, error) {
	texist, err := tr.modelHandle.GetTemplate(template.Namespace, template.Name)
	if texist != nil {
		return nil, errors.NewConflict().SetCode("TemplateAlreadyExists").SetMessage("template already exists")
	} else {
		if err != nil {
			if err != orm.ErrNoRows {
				return nil, errors.NewInternalServerError().SetCause(err)
			}
		}
	}
	temp, err := tr.modelHandle.CreateTemplate(template)
	if err != nil {
		return nil, errors.NewInternalServerError().SetCause(err)
	}

	return temp, nil
}

func (tr *TemplateRes) DeleteTemplate(namespace, name string) error {
	_, err := tr.modelHandle.GetTemplate(namespace, name)
	if err != nil {
		if err == orm.ErrNoRows {
			return errors.NewNotFound().SetCause(err)
		}
		return errors.NewInternalServerError().SetCause(err)
	}

	//TODO: if template is used by some apps, it can not be deleted
	err = tr.modelHandle.DeleteTemplate(namespace, name)
	if err != nil {
		return errors.NewInternalServerError().SetCause(err)
	}

	return nil
}

func (tr *TemplateRes) UpdateTemplate(template models.CaasTemplate) (*models.CaasTemplate, error) {
	told, err := tr.modelHandle.GetTemplate(template.Namespace, template.Name)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.NewNotFound().SetCause(err)
		}
		return nil, errors.NewInternalServerError().SetCause(err)
	}
	told.Kind = template.Kind
	told.Spec = template.Spec
	told.Description = template.Description

	err = tr.modelHandle.UpdateTemplate(*told)
	if err != nil {
		return nil, errors.NewInternalServerError().SetCause(err)
	}

	return told, nil
}

func (tr *TemplateRes) GetTemplateByName(namespace, name string) (*models.CaasTemplate, int, error) {
	template, err := tr.modelHandle.GetTemplate(namespace, name)
	if err != nil {
		if err == orm.ErrNoRows {
			return template, http.StatusNotFound, err
		}
		return template, http.StatusInternalServerError, err
	}

	return template, http.StatusOK, nil
}

func (tr *TemplateRes) GetTemplateByID(id int64) (*models.CaasTemplate, int, error) {
	template, err := tr.modelHandle.GetTemplateByID(id)
	if err != nil {
		if err == orm.ErrNoRows {
			return template, http.StatusNotFound, err
		}
		return template, http.StatusInternalServerError, err
	}

	return template, http.StatusOK, nil
}
