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
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Ptype string `gorm:"size:512"`
	V0    string `gorm:"size:512"`
	V1    string `gorm:"size:512"`
	V2    string `gorm:"size:512"`
	V3    string `gorm:"size:512"`
	V4    string `gorm:"size:512"`
	V5    string `gorm:"size:512"`
}

// TableName ..
func (CasbinRule) TableName() string {
	return "casbin_rule"
}
