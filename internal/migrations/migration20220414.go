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

package migrations

import (
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/go-atomci/atomci/internal/core/apps"
	"github.com/go-atomci/atomci/internal/core/project"
	"github.com/go-atomci/atomci/internal/core/settings"
	"github.com/go-atomci/atomci/internal/middleware/log"
	"github.com/go-atomci/atomci/utils"
)

type Migration20220414 struct {
}

func (m Migration20220414) GetCreateAt() time.Time {
	return time.Date(2022, 4, 14, 0, 0, 0, 0, time.Local)
}

type repoServer struct {
	ID       int64  `orm:"column(id);"`
	Type     string `orm:"column(type);"`
	BaseURL  string `orm:"column(base_url);"`
	User     string `orm:"column(user);"`
	Token    string `orm:"column(token);"`
	Password string `orm:"column(password);"`
	CID      int64  `orm:"column(cid);"`
}

type projectApp struct {
	ID           int64  `orm:"column(id)"`
	Creator      string `orm:"column(creator)"`
	ProjectID    int64  `orm:"column(project_id)"`
	Name         string `orm:"column(name)"`
	FullName     string `orm:"column(full_name)"`
	Language     string `orm:"column(language)"`
	BranchName   string `orm:"column(branch_name)"`
	Path         string `orm:"column(path)"`
	RepoID       int64  `orm:"column(repo_id)"`
	CompileEnvID int64  `orm:"column(compile_env_id)"`
	BuildPath    string `orm:"column(build_path)"`
	Dockerfile   string `orm:"column(dockerfile)"`
}

func (m Migration20220414) Upgrade(ormer orm.Ormer) error {
	// 01. migrate pub_repo_server into  sys_integrate_setting table
	// 		* save repo_id mapping into integrate_setting id
	// 02. migrate pub_project_app into pub_scm_app
	//		* repo_id change to sys_integrate_setting id
	// 		* pub_project_app add scm_id
	// 03. clean pub_app_branch data
	// 04. pub_project_app delete unused column

	// 01.
	var repoServerItems []repoServer
	_, err := ormer.Raw("SELECT id,type,base_url,user,token,password,cid FROM pub_repo_server WHERE deleted=0;").QueryRows(&repoServerItems)
	if err != nil {
		if strings.Contains(err.Error(), " doesn't exist") {
			return nil
		}
		return fmt.Errorf("select pub_repo_server: %s", err.Error())
	}
	log.Log.Debug("reposerver len: %v", len(repoServerItems))
	repoItemsMapping, err := m.migrateRepoServerIntoIntegrateSetting(repoServerItems)
	if err != nil {
		return fmt.Errorf("repo items mappings, error: %s", err.Error())
	}
	log.Log.Debug("repoItems mapping item: %+v", repoItemsMapping)

	// 02.
	// get all project apps
	var appsRes []projectApp
	_, err = ormer.Raw("SELECT id,creator,project_id,name,full_name,language,branch_name,path,repo_id,compile_env_id,build_path,dockerfile FROM pub_project_app WHERE deleted=0;").QueryRows(&appsRes)
	if err != nil {
		return fmt.Errorf("select pub_project_app : %s", err.Error())
	}
	scmAppHander := apps.NewAppManager()
	projectAppHander := project.NewProjectManager()
	for _, projectApp := range appsRes {
		log.Log.Debug("project app repo_id: %v", projectApp.RepoID)
		item := apps.ScmAppReq{
			Name:         projectApp.Name,
			CompileEnvID: projectApp.CompileEnvID,
			Language:     projectApp.Language,
			Path:         projectApp.Path,
			RepoID:       repoItemsMapping[projectApp.RepoID],
			FullName:     projectApp.FullName,
			BranchName:   projectApp.BranchName,
			BuildPath:    projectApp.BuildPath,
			Dockerfile:   projectApp.Dockerfile,
		}
		scmID, err := scmAppHander.CreateSCMApp(&item, projectApp.Creator)
		if err != nil {
			log.Log.Warn("migrate pub_scm_app, name: %s, occur error: %s, skip item", projectApp.Name, err.Error())
			continue
		}
		log.Log.Debug("created scm id: %v", scmID)
		req := project.ProjectAppUpdateReq{
			ScmID: scmID,
		}
		if err := projectAppHander.UpdateProjectApp(projectApp.ProjectID, projectApp.ID, &req); err != nil {
			log.Log.Warn("migrate pub_project_app, update scm_id into project app(%v), occur error: %s, skip item", projectApp.ID, err.Error())
			continue
		}
		log.Log.Debug("project app id: %v, migrated succeed", projectApp.ID)
	}
	// 03
	if _, err := ormer.Raw("delete from pub_app_branch;").Exec(); err != nil {
		return err
	}

	// 04
	return nil
}

func (m Migration20220414) migrateRepoServerIntoIntegrateSetting(items []repoServer) (map[int64]int64, error) {
	resp := map[int64]int64{}
	sm := settings.NewSettingManager()
	for _, item := range items {
		integrateSettingName := fmt.Sprintf("%v-%v", item.Type, utils.GenerateRandomstring(5))
		req := settings.IntegrateSettingReq{
			Name:   integrateSettingName,
			Type:   item.Type,
			Config: m.generateRepoConf(item),
		}
		if err := sm.CreateIntegrateSetting(&req, "admin"); err != nil {
			log.Log.Warn("migrate repo_server: %v when create integrate error: %s", item.BaseURL, err.Error())
			continue
		}
		scmIntegrateItem, err := sm.GetIntegrateSettingByName(integrateSettingName, item.Type)
		if err != nil {
			log.Log.Warn("migrate repo_server: %v when get integratesetting item error: %s", item.BaseURL, err.Error())
			continue
		}
		resp[item.ID] = scmIntegrateItem.ID
	}
	return resp, nil
}

func (m Migration20220414) generateRepoConf(repoItem repoServer) (item interface{}) {
	switch strings.ToLower(repoItem.Type) {
	case "gitlab":
		item = settings.ScmAuthConf{
			User: repoItem.User,
			ScmBaseConfig: settings.ScmBaseConfig{
				URL:   repoItem.BaseURL,
				Token: repoItem.Token,
			},
		}
	case "gitee", "gitea", "github":
		item = settings.ScmBaseConfig{
			URL:   repoItem.BaseURL,
			Token: repoItem.Token,
		}
	}
	return item
}
