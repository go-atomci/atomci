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

import "time"

const (
	DEFAULT_WEIGHT = 1
)

type CaasApplication struct {
	Name              string `orm:"column(name)" json:"name"`
	ProjectID         int64  `orm:"column(project_id)" json:"project_id"`
	Env               string `orm:"column(env)" json:"env"`
	Kind              string `orm:"column(kind)" json:"kind"`
	TemplateName      string `orm:"column(template_name)" json:"template_name"`
	Cluster           string `orm:"column(cluster)" json:"cluster"`
	Namespace         string `orm:"column(namespace)" json:"namespace"`
	Replicas          int    `orm:"column(replicas)" json:"replicas"`
	Image             string `orm:"column(image)" json:"image"`
	Template          string `orm:"column(template);type(text)" json:"template,omitempty"`
	LabelSelector     string `orm:"column(label_selector)" json:"label_selector,omitempty"`
	StatusReplicas    int32  `orm:"column(status_replicas)" json:"status_replicas"`
	ReadyReplicas     int32  `orm:"column(ready_replicas)" json:"ready_replicas"`
	UpdatedReplicas   int32  `orm:"column(updated_replicas)" json:"updated_replicas"`
	AvailableReplicas int32  `orm:"column(available_replicas)" json:"available_replicas"`
	AvailableStatus   string `orm:"column(available_status)" json:"available_status"`
	Message           string `orm:"column(message)" json:"message"`
	DeployStatus      string `orm:"column(deploy_status)" json:"deploy_status"`
	Labels            string `orm:"column(labels);type(text)" json:"labels"`
	Addons
}

func (t *CaasApplication) TableName() string {
	return "caas_application"
}

type CaasTemplate struct {
	Name        string `orm:"column(name)" json:"name"`
	Namespace   string `orm:"column(namespace)" json:"namespace"`
	Description string `orm:"column(description)" json:"description,omitempty"`
	Spec        string `orm:"column(spec);type(text)" json:"spec"` //TemplateSpec
	Kind        string `orm:"column(kind)" json:"kind"`
	Addons
}

func (t *CaasTemplate) TableName() string {
	return "caas_template"
}

func (ons Addons) UpdateAddons() Addons {
	ons.CreateAt, _ = time.Parse("2006-01-02 15:04:05", ons.CreateAt.Format("2006-01-02 15:04:05"))
	ons.UpdateAt, _ = time.Parse("2006-01-02 15:04:05", time.Now().Local().Format("2006-01-02 15:04:05"))
	return ons
}
