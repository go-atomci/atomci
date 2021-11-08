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

package mycasbin

// CasbinRule ..
type CasbinRule struct {
	PType string `json:"p_type" gorm:"type:varchar(100);"`
	V0    string `json:"v0" orm:"size(100);"`
	V1    string `json:"v1" orm:"size(100);"`
	V2    string `json:"v2" orm:"size(100);"`
	V3    string `json:"v3" orm:"size(100);"`
	V4    string `json:"v4" orm:"size(100);"`
	V5    string `json:"v5" orm:"size(100);"`
}

// TableName ..
func (CasbinRule) TableName() string {
	return "casbin_rule"
}
