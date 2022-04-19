// Copyright 2017 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package beegoormadapter

import (
	"fmt"
	"runtime"

	"github.com/astaxie/beego/orm"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
)

type CasbinRule struct {
	Id    int
	Ptype string
	V0    string
	V1    string
	V2    string
	V3    string
	V4    string
	V5    string
}

func init() {
	orm.RegisterModel(new(CasbinRule))
}

const (
	defaultTableName = "casbin_rule"
)

// Adapter represents the Xorm adapter for policy storage.
type Adapter struct {
	driverName      string
	dataSourceName  string
	dataSourceAlias string
	tableName       string
	dbSpecified     bool
	o               orm.Ormer
}

// finalizer is the destructor for Adapter.
func finalizer(a *Adapter) {
}

// NewAdapter is the constructor for Adapter.
// dataSourceAlias: Database alias. ORM will use it to switch database.
// driverName: database driverName.
// dataSourceName: connection string
func NewAdapter(dataSourceAlias, driverName, dataSourceName string) (*Adapter, error) {
	a := &Adapter{}
	a.driverName = driverName
	a.dataSourceName = dataSourceName
	a.dataSourceAlias = dataSourceAlias
	a.tableName = defaultTableName

	err := a.open()

	if err != nil {
		return nil, err
	}

	// Call the destructor when the object is released.
	runtime.SetFinalizer(a, finalizer)

	return a, nil
}

func (a *Adapter) registerDataBase(aliasName, driverName, dataSource string, params ...int) error {
	err := orm.RegisterDataBase(aliasName, driverName, dataSource, params...)
	return err
}

func (a *Adapter) open() error {
	var err error

	err = a.registerDataBase(a.dataSourceAlias, a.driverName, a.dataSourceName)
	if err != nil {
		return err
	}

	a.o = orm.NewOrm()
	err = a.o.Using(a.dataSourceAlias)
	if err != nil {
		return err
	}

	err = a.createTable()
	if err != nil {
		return err
	}

	return nil
}

func (a *Adapter) close() {
	a.o = nil
}

func (a *Adapter) createTable() error {
	return orm.RunSyncdb(a.dataSourceAlias, false, false)
}

func (a *Adapter) dropTable() error {
	_, err := a.o.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %v;", a.tableName)).Exec()
	return err
}

func loadPolicyLine(line CasbinRule, model model.Model) {
	lineText := line.Ptype
	if line.V0 != "" {
		lineText += ", " + line.V0
	}
	if line.V1 != "" {
		lineText += ", " + line.V1
	}
	if line.V2 != "" {
		lineText += ", " + line.V2
	}
	if line.V3 != "" {
		lineText += ", " + line.V3
	}
	if line.V4 != "" {
		lineText += ", " + line.V4
	}
	if line.V5 != "" {
		lineText += ", " + line.V5
	}

	persist.LoadPolicyLine(lineText, model)
}

// LoadPolicy loads policy from database.
func (a *Adapter) LoadPolicy(model model.Model) error {
	var lines []CasbinRule
	_, err := a.o.QueryTable("casbin_rule").All(&lines)
	if err != nil {
		return err
	}

	for _, line := range lines {
		loadPolicyLine(line, model)
	}

	return nil
}

func savePolicyLine(ptype string, rule []string) CasbinRule {
	line := CasbinRule{}

	line.Ptype = ptype
	if len(rule) > 0 {
		line.V0 = rule[0]
	}
	if len(rule) > 1 {
		line.V1 = rule[1]
	}
	if len(rule) > 2 {
		line.V2 = rule[2]
	}
	if len(rule) > 3 {
		line.V3 = rule[3]
	}
	if len(rule) > 4 {
		line.V4 = rule[4]
	}
	if len(rule) > 5 {
		line.V5 = rule[5]
	}

	return line
}

// SavePolicy saves policy to database.
func (a *Adapter) SavePolicy(model model.Model) error {
	err := a.dropTable()
	if err != nil {
		return err
	}

	err = a.createTable()
	if err != nil {
		return err
	}

	var lines []CasbinRule

	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			lines = append(lines, line)
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			lines = append(lines, line)
		}
	}

	_, err = a.o.InsertMulti(len(lines), lines)
	return err
}

// AddPolicy adds a policy rule to the storage.
func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	line := savePolicyLine(ptype, rule)
	_, err := a.o.Insert(&line)
	return err
}

// RemovePolicy removes a policy rule from the storage.
func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	line := savePolicyLine(ptype, rule)
	_, err := a.o.Delete(&line, "ptype", "v0", "v1", "v2", "v3", "v4", "v5")
	return err
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	line := CasbinRule{}

	line.Ptype = ptype
	filter := []string{}
	filter = append(filter, "ptype")
	if fieldIndex <= 0 && 0 < fieldIndex+len(fieldValues) {
		line.V0 = fieldValues[0-fieldIndex]
		filter = append(filter, "v0")
	}
	if fieldIndex <= 1 && 1 < fieldIndex+len(fieldValues) {
		line.V1 = fieldValues[1-fieldIndex]
		filter = append(filter, "v1")
	}
	if fieldIndex <= 2 && 2 < fieldIndex+len(fieldValues) {
		line.V2 = fieldValues[2-fieldIndex]
		filter = append(filter, "v2")
	}
	if fieldIndex <= 3 && 3 < fieldIndex+len(fieldValues) {
		line.V3 = fieldValues[3-fieldIndex]
		filter = append(filter, "v3")
	}
	if fieldIndex <= 4 && 4 < fieldIndex+len(fieldValues) {
		line.V4 = fieldValues[4-fieldIndex]
		filter = append(filter, "v4")
	}
	if fieldIndex <= 5 && 5 < fieldIndex+len(fieldValues) {
		line.V5 = fieldValues[5-fieldIndex]
		filter = append(filter, "v5")
	}

	_, err := a.o.Delete(&line, filter...)
	return err
}
