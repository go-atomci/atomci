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

type Audit struct {
	Addons
	User            string `orm:"column(user)" json:"user"`
	Method          string `orm:"column(method)" json:"method"`
	Operation       string `orm:"column(operation)" json:"operation"`
	OperationObject string `orm:"column(operation_object)" json:"operation_object"`
	OperationBody   string `orm:"column(operation_body);type(text)" json:"operation_body"`
	OperationStatus int    `orm:"column(operation_status)" json:"operation_status"`
}

func (t *Audit) TableName() string {
	return "sys_audit"
}
