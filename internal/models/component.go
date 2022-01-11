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

// FlowComponent ...
type FlowComponent struct {
	Addons
	Name   string `orm:"column(name);size(64)" json:"name"`
	Type   string `orm:"column(type)" json:"type"`
	Params string `orm:"column(params);size(1024)" json:"params"`
}

// TableName ...
func (t *FlowComponent) TableName() string {
	return "pub_flow_component"
}

// TaskTmpl ...
type TaskTmpl struct {
	Addons
	ComponentID int64  `orm:"column(component_id)" json:"component_id"`
	Creator     string `orm:"column(creator);size(64)" json:"creator"`
	Name        string `orm:"column(name);size(64)" json:"name"`
	Description string `orm:"column(description);size(256)" json:"description"`
	Type        string `orm:"column(type);size(64)" json:"type"`
	TypeDisplay string `orm:"column(type_display);size(128)" json:"type_display"`
	Params      string `orm:"column(params);size(1024)" json:"params"`
	SubTask     string `orm:"column(sub_task);size(4096)" json:"sub_task"`
}

// TableName ...
func (t *TaskTmpl) TableName() string {
	return "pub_flow_step"
}

// CompileEnv ...
type CompileEnv struct {
	Addons
	Name        string `orm:"column(name);size(64)" json:"name"`
	Image       string `orm:"column(image);size(256)" json:"image"`
	Command     string `orm:"column(command);size(128)" json:"command"`
	Args        string `orm:"column(args);size(128)" json:"args"`
	Creator     string `orm:"column(creator);size(64)" json:"creator"`
	Description string `orm:"column(description);size(256)" json:"description"`
}

// TableName ...
func (t *CompileEnv) TableName() string {
	return "sys_compile_env"
}
