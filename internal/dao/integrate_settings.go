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

// SysSettingModel ...
type SysSettingModel struct {
	ormer                     orm.Ormer
	IntegrateSettingTableName string
	CompileEnvTableName       string
}

// NewSysSettingModel ...
func NewSysSettingModel() (model *SysSettingModel) {
	return &SysSettingModel{
		ormer:                     GetOrmer(),
		IntegrateSettingTableName: (&models.IntegrateSetting{}).TableName(),
		CompileEnvTableName:       (&models.CompileEnv{}).TableName(),
	}
}

// GetIntegrateSettingByID ...
func (model *SysSettingModel) GetIntegrateSettingByID(integrateSettingID int64) (*models.IntegrateSetting, error) {
	integrateSetting := models.IntegrateSetting{}
	qs := model.ormer.QueryTable(model.IntegrateSettingTableName).Filter("deleted", false)
	if err := qs.Filter("id", integrateSettingID).One(&integrateSetting); err != nil {
		return nil, err
	}
	return &integrateSetting, nil
}

func (model *SysSettingModel) GetIntegrateSettingByName(name, integrateType string) (*models.IntegrateSetting, error) {
	integrateSetting := models.IntegrateSetting{}
	qs := model.ormer.QueryTable(model.IntegrateSettingTableName).Filter("deleted", false)
	if err := qs.Filter("name", name).Filter("type", integrateType).One(&integrateSetting); err != nil {
		return nil, err
	}
	return &integrateSetting, nil
}

// GetIntegrateSettings ...
func (model *SysSettingModel) GetIntegrateSettings(integrateTypes []string) ([]*models.IntegrateSetting, error) {
	var integrateSettings []*models.IntegrateSetting
	qs := model.ormer.QueryTable(model.IntegrateSettingTableName).Filter("deleted", false)
	if len(integrateTypes) > 0 {
		qs = qs.Filter("type__in", integrateTypes)
	}
	_, err := qs.All(&integrateSettings)
	if err != nil {
		return nil, err
	}
	return integrateSettings, err
}

// GetIntegrateSettingsByPagination ..
func (model *SysSettingModel) GetIntegrateSettingsByPagination(filter *query.FilterQuery, intergrateTypes []string) (*query.QueryResult, []*models.IntegrateSetting, error) {
	rst := &query.QueryResult{Item: []*models.IntegrateSetting{}}
	queryCond := orm.NewCondition().AndCond(orm.NewCondition().And("deleted", false))

	if filterCond := query.FilterCondition(filter, filter.FilterKey); filterCond != nil {
		queryCond = queryCond.AndCond(filterCond)
	}
	qs := model.ormer.QueryTable(model.IntegrateSettingTableName).OrderBy("-create_at").SetCond(queryCond)
	if len(intergrateTypes) > 0 {
		qs = qs.Filter("type__in", intergrateTypes)
	}
	count, err := qs.Count()
	if err != nil {
		return nil, nil, err
	}
	if err = query.FillPageInfo(rst, filter.PageIndex, filter.PageSize, int(count)); err != nil {
		return nil, nil, err
	}

	settingList := []*models.IntegrateSetting{}
	_, err = qs.Limit(filter.PageSize, filter.PageSize*(filter.PageIndex-1)).All(&settingList)
	if err != nil {
		return nil, nil, err
	}
	rst.Item = settingList

	return rst, settingList, nil
}

// UpdateIntegrateSetting ..
func (model *SysSettingModel) UpdateIntegrateSetting(integrateSetting *models.IntegrateSetting) error {
	_, err := model.ormer.Update(integrateSetting)
	return err
}

// DeleteIntegrateSetting ..
func (model *SysSettingModel) DeleteIntegrateSetting(integrateSettingID int64) error {
	integrateSetting, err := model.GetIntegrateSettingByID(integrateSettingID)
	if err != nil {
		return err
	}
	integrateSetting.MarkDeleted()
	_, err = model.ormer.Update(integrateSetting)
	return err
}

// CreateIntegrateSetting ...
func (model *SysSettingModel) CreateIntegrateSetting(integrateSetting *models.IntegrateSetting) error {
	_, err := model.ormer.InsertOrUpdate(integrateSetting)
	return err
}

// GetCompileEnvByID ...
func (model *SysSettingModel) GetCompileEnvByID(integrateSettingID int64) (*models.CompileEnv, error) {
	integrateSetting := models.CompileEnv{}
	qs := model.ormer.QueryTable(model.CompileEnvTableName).Filter("deleted", false)
	if err := qs.Filter("id", integrateSettingID).One(&integrateSetting); err != nil {
		return nil, err
	}
	return &integrateSetting, nil
}

// GetCompileEnvByName ...
func (model *SysSettingModel) GetCompileEnvByName(compileEnvItem string) (*models.CompileEnv, error) {
	compileEnv := models.CompileEnv{}
	qs := model.ormer.QueryTable(model.CompileEnvTableName).Filter("deleted", false)
	if err := qs.Filter("name", compileEnvItem).One(&compileEnv); err != nil {
		return nil, err
	}
	return &compileEnv, nil
}

// GetCompileEnvs ...
func (model *SysSettingModel) GetCompileEnvs() ([]*models.CompileEnv, error) {
	integrateSettings := []*models.CompileEnv{}
	qs := model.ormer.QueryTable(model.CompileEnvTableName).Filter("deleted", false)

	_, err := qs.All(&integrateSettings)
	if err != nil {
		return nil, err
	}
	return integrateSettings, err
}

// GetCompileEnvsByPagination ..
func (model *SysSettingModel) GetCompileEnvsByPagination(filter *query.FilterQuery) (*query.QueryResult, []*models.CompileEnv, error) {
	rst := &query.QueryResult{Item: []*models.CompileEnv{}}
	queryCond := orm.NewCondition().AndCond(orm.NewCondition().And("deleted", false))

	if filterCond := query.FilterCondition(filter, filter.FilterKey); filterCond != nil {
		queryCond = queryCond.AndCond(filterCond)
	}
	qs := model.ormer.QueryTable(model.CompileEnvTableName).OrderBy("-create_at").SetCond(queryCond)
	count, err := qs.Count()
	if err != nil {
		return nil, nil, err
	}
	if err = query.FillPageInfo(rst, filter.PageIndex, filter.PageSize, int(count)); err != nil {
		return nil, nil, err
	}

	settingList := []*models.CompileEnv{}
	_, err = qs.Limit(filter.PageSize, filter.PageSize*(filter.PageIndex-1)).All(&settingList)
	if err != nil {
		return nil, nil, err
	}
	rst.Item = settingList

	return rst, settingList, nil
}

// UpdateCompileEnv ..
func (model *SysSettingModel) UpdateCompileEnv(integrateSetting *models.CompileEnv) error {
	_, err := model.ormer.Update(integrateSetting)
	return err
}

// DeleteCompileEnv ..
func (model *SysSettingModel) DeleteCompileEnv(integrateSettingID int64) error {
	integrateSetting, err := model.GetCompileEnvByID(integrateSettingID)
	if err != nil {
		return err
	}
	integrateSetting.MarkDeleted()
	_, err = model.ormer.Update(integrateSetting)
	return err
}

// CreateCompileEnv ...
func (model *SysSettingModel) CreateCompileEnv(integrateSetting *models.CompileEnv) error {
	_, err := model.ormer.InsertOrUpdate(integrateSetting)
	return err
}
