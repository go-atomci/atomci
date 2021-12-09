package initialize

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/go-atomci/atomci/core/settings"
	"github.com/go-atomci/atomci/dao"
	"github.com/go-atomci/atomci/middleware/log"
	"github.com/go-atomci/atomci/models"
)

var compileEnvs = []settings.CompileEnvReq{
	{
		Name:        "jnlp",
		Image:       "colynn/jenkins-jnlp-agent:latest",
		Description: "",
	},
	{
		Name:    "kaniko",
		Image:   "colynn/kaniko-executor:debug",
		Command: "/bin/sh -c",
		Args:    "cat",
	},
	{
		Name:        "node",
		Image:       "node:12.12-alpine",
		Description: "nodejs编译环境",
	},
	{
		Name:    "maven",
		Image:   "maven:3.8.2-openjdk-8",
		Command: "/bin/sh -c",
		Args:    "cat",
	},
}

func initCompileEnvs() error {
	settingModel := dao.NewSysSettingModel()
	for _, item := range compileEnvs {
		_, err := settingModel.GetCompileEnvByName(item.Name)
		if err != nil {
			if err == orm.ErrNoRows {
				component := &models.CompileEnv{
					Addons:      models.NewAddons(),
					Name:        item.Name,
					Image:       item.Image,
					Command:     item.Command,
					Creator:     "admin", // create use 'admin'
					Args:        item.Args,
					Description: item.Description,
				}
				if err := settingModel.CreateCompileEnv(component); err != nil {
					log.Log.Error("when init compile env, occur error: %s", err.Error())
					continue
				}
			} else {
				logs.Warn("init compile env occur error: %s", err.Error())
				continue
			}
		} else {
			log.Log.Debug("component type `%s` already exists, skip", item.Name)
			continue
		}
	}
	return nil
}
