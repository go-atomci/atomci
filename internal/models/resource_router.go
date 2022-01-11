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

type GatewayRouter struct {
	Addons
	Router            string `orm:"column(router)" json:"router"`
	Method            string `orm:"column(method)" json:"method"`
	Backend           string `orm:"column(backend)" json:"backend"`
	ResourceType      string `orm:"column(resource_type)" json:"resource_type"`
	ResourceOperation string `orm:"column(resource_operation)" json:"resource_operation"`
}

func (t *GatewayRouter) TableName() string {
	return "sys_resource_router"
}

func (t *GatewayRouter) TableIndex() [][]string {
	return [][]string{
		{"Backend"},
		{"Router", "Method"},
	}
}

func (u *GatewayRouter) TableUnique() [][]string {
	return [][]string{
		{"Router", "Method"},
		{"ResourceType", "ResourceOperation"},
	}
}
