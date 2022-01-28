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

import (
	"github.com/go-atomci/atomci/utils"
)

// GitApp ...
type GitApp struct {
	Addons
	Name        string `orm:"column(name);size(64)" json:"name"`
	Description string `orm:"column(description);size(256)" json:"description"`
	Department  string `orm:"column(department);size(64)" json:"department"`
	Product     string `orm:"column(product);size(64)" json:"product"`
	Language    string `orm:"column(language);size(64)" json:"language"`
	Type        string `orm:"column(type);size(64)" json:"type"`
	VcsType     string `orm:"column(vcs_type);size(64)" json:"vcs_type"`
	Path        string `orm:"column(path);size(256)" json:"path"`
	BuildPath   string `orm:"column(build_path);size(256)" json:"build_path"`
	IsArranged  bool   `orm:"column(is_arranged);default(true)" json:"is_arranged"`
}

// TableName ...
func (t *GitApp) TableName() string {
	return "pub_gitapp"
}

// TableUnique ...
func (t *GitApp) TableUnique() [][]string {
	return [][]string{
		[]string{"Name", "Department"},
	}
}

// AppBranch ...
type AppBranch struct {
	Addons
	AppID      int64  `orm:"column(app_id);" json:"app_id"`
	BranchName string `orm:"column(branch_name);size(64)" json:"branch_name"`
	Path       string `orm:"column(path);size(256)" json:"path"`
}

// TableName ...
func (t *AppBranch) TableName() string {
	return "pub_app_branch"
}

// RepoServer ..
type RepoServer struct {
	Addons
	Type     string `orm:"column(type);" json:"type"`
	BaseURL  string `orm:"column(base_url);" json:"base_url"`
	User     string `orm:"column(user);" json:"user"`
	token    string `orm:"column(token);" json:"token"`
	password string `orm:"column(password);" json:"password"`
	CID      int64  `orm:"column(cid);" json:"cid"`
}

// TableName ...
func (t *RepoServer) TableName() string {
	return "pub_repo_server"
}

func (repo *RepoServer) SetToken(token string) {
	plainText := []byte(token)
	repo.token = string(utils.AesEny(plainText))
}

func (repo *RepoServer) GetToken() string {
	if len(repo.token) == 0 {
		return ""
	}
	return string(utils.AesEny([]byte(repo.token)))
}

func (repo *RepoServer) SetPassword(password string) {
	plainText := []byte(password)
	repo.password = string(utils.AesEny(plainText))
}

func (repo *RepoServer) GetPassword() string {
	if len(repo.password) == 0 {
		return ""
	}
	return string(utils.AesEny([]byte(repo.password)))
}
