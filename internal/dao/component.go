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
	"github.com/go-atomci/atomci/internal/models"
	"github.com/go-atomci/atomci/utils/query"

	"github.com/astaxie/beego/orm"
)

// PipelineStageModel ...
type PipelineStageModel struct {
	ormer                     orm.Ormer
	PipelineTableName         string
	PipelineInstanceTableName string
	FlowComponentTableName    string
	TaskTmplTableName         string
}

// NewPipelineStageModel ...
func NewPipelineStageModel() (model *PipelineStageModel) {
	return &PipelineStageModel{
		ormer:                     GetOrmer(),
		PipelineTableName:         (&models.ProjectPipeline{}).TableName(),
		FlowComponentTableName:    (&models.FlowComponent{}).TableName(),
		TaskTmplTableName:         (&models.TaskTmpl{}).TableName(),
		PipelineInstanceTableName: (&models.PipelineInstance{}).TableName(),
	}
}

// GetPipelineInstanceConfigByID ..
func (model *PipelineStageModel) GetPipelineInstanceConfigByID(id int64) (*models.PipelineInstance, error) {
	pipeline := models.PipelineInstance{}
	qs := model.ormer.QueryTable(model.PipelineInstanceTableName).Filter("deleted", false)
	if err := qs.Filter("id", id).One(&pipeline); err != nil {
		return nil, err
	}
	return &pipeline, nil
}

// CreatePipelineInstance ...
func (model *PipelineStageModel) CreatePipelineInstance(pipeline *models.PipelineInstance) (int64, error) {
	return model.ormer.Insert(pipeline)
}

// GetFlowComponentByType ..
func (model *PipelineStageModel) GetFlowComponentByType(componentType string) (*models.FlowComponent, error) {
	component := models.FlowComponent{}
	err := model.ormer.QueryTable(model.FlowComponentTableName).Filter("deleted", false).
		Filter("type", componentType).One(&component)
	return &component, err
}

// GetFlowComponents ...
func (model *PipelineStageModel) GetFlowComponents() ([]*models.FlowComponent, error) {
	components := []*models.FlowComponent{}
	qs := model.ormer.QueryTable(model.FlowComponentTableName).Filter("deleted", false)
	_, err := qs.All(&components)
	if err != nil {
		return nil, err
	}
	return components, err
}

// GetTaskTmplByID ...
func (model *PipelineStageModel) GetTaskTmplByID(stepID int64) (*models.TaskTmpl, error) {
	step := models.TaskTmpl{}
	if err := model.ormer.QueryTable(model.TaskTmplTableName).
		Filter("deleted", false).
		Filter("id", stepID).One(&step); err != nil {
		return nil, err
	}
	return &step, nil
}

// GetTaskTmplByName ...
func (model *PipelineStageModel) GetTaskTmplByName(name string) (*models.TaskTmpl, error) {
	step := models.TaskTmpl{}
	if err := model.ormer.QueryTable(model.TaskTmplTableName).
		Filter("deleted", false).
		Filter("name", name).One(&step); err != nil {
		return nil, err
	}
	return &step, nil
}

// GetTaskTmpls ...
func (model *PipelineStageModel) GetTaskTmpls() ([]*models.TaskTmpl, error) {
	steps := []*models.TaskTmpl{}
	qs := model.ormer.QueryTable(model.TaskTmplTableName).Filter("deleted", false)
	_, err := qs.All(&steps)
	if err != nil {
		return nil, err
	}
	return steps, err
}

// GetTaskTmplsByPagination ..
func (model *PipelineStageModel) GetTaskTmplsByPagination(filter *query.FilterQuery) ([]*models.TaskTmpl, int64, error) {

	queryCond := orm.NewCondition().AndCond(orm.NewCondition().And("deleted", false))

	if filterCond := query.FilterCondition(filter, filter.FilterKey); filterCond != nil {
		queryCond = queryCond.AndCond(filterCond)
	}
	qs := model.ormer.QueryTable(model.TaskTmplTableName).OrderBy("-create_at").SetCond(queryCond)
	count, err := qs.Count()
	if err != nil {
		return nil, 0, err
	}

	stepList := []*models.TaskTmpl{}
	_, err = qs.Limit(filter.PageSize, filter.PageSize*(filter.PageIndex-1)).All(&stepList)
	if err != nil {
		return nil, 0, err
	}

	return stepList, count, nil
}

// CreateFlowComponent ...
func (model *PipelineStageModel) CreateFlowComponent(comp *models.FlowComponent) error {
	_, err := model.ormer.InsertOrUpdate(comp)
	return err
}

// CreateTaskTmpl ...
func (model *PipelineStageModel) CreateTaskTmpl(step *models.TaskTmpl) error {
	_, err := model.ormer.InsertOrUpdate(step)
	return err
}

// UpdateTaskTmpl ..
func (model *PipelineStageModel) UpdateTaskTmpl(step *models.TaskTmpl) error {
	_, err := model.ormer.Update(step)
	return err
}

// DeleteTaskTmpl ..
func (model *PipelineStageModel) DeleteTaskTmpl(stepID int64) error {
	step, err := model.GetTaskTmplByID(stepID)
	if err != nil {
		return err
	}
	step.MarkDeleted()
	_, err = model.ormer.Delete(step)
	return err
}
