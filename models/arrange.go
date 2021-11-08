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

// AppArrange ...
type AppArrange struct {
	Addons
	Name         string `orm:"column(name);size(64)" json:"name"`
	ArrangeEnv   string `orm:"column(arrange_env);size(64)" json:"arrange_env"`
	EnvID        int64  `orm:"column(env_id);" json:"env_id"`
	ProjectAppID int64  `orm:"column(project_app_id)" json:"project_app_id"`
	Config       string `orm:"column(config);type(text)" json:"config"`
}

// TableName ...
func (t *AppArrange) TableName() string {
	return "pub_app_env_arrange"
}

// TableUnique ...
func (t *AppArrange) TableUnique() [][]string {
	return [][]string{
		[]string{"ProjectAppID", "EnvID"},
	}
}

// TableIndex ...
func (t *AppArrange) TableIndex() [][]string {
	return [][]string{
		[]string{"ProjectAppID", "EnvID"},
	}
}

const (
	OriginTag = iota + 1
	SystemDefaultTag
	LatestTag
)

type AppImageMapping struct {
	Addons
	ArrangeID    int64  `orm:"column(arrange_id)" json:"arrange_id"`
	Name         string `orm:"column(name);size(64)" json:"name"`
	Image        string `orm:"column(image);size(128)" json:"image"`
	ProjectAppID int64  `orm:"column(project_app_id)" json:"project_app_id"`
	ImageTagType int64  `orm:"column(iamge_tag_type)" json:"image_tag_type"`
}

// TableName ...
func (t *AppImageMapping) TableName() string {
	return "app_image_mapping"
}
