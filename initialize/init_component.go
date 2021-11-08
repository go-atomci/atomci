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

package initialize

import (
	"github.com/go-atomci/atomci/dao"
	"github.com/go-atomci/atomci/middleware/log"
	"github.com/go-atomci/atomci/models"

	"github.com/astaxie/beego/orm"
)

type component struct {
	Name   string
	Type   string
	Params string
}

// Component ..
func Component() error {
	components := []component{
		{
			Name: "人工卡点",
			Type: "manual",
		},
		{
			Name: "构建",
			Type: "build",
		},
		{
			Name: "部署",
			Type: "deploy",
		},
	}
	for _, comp := range components {
		pipelineModel := dao.NewPipelineStageModel()
		_, err := pipelineModel.GetFlowComponentByType(comp.Type)
		if err != nil {
			if err == orm.ErrNoRows {
				component := &models.FlowComponent{
					Addons: models.NewAddons(),
					Name:   comp.Name,
					Type:   comp.Type,
					Params: comp.Params,
				}
				if err := pipelineModel.CreateFlowComponent(component); err != nil {
					log.Log.Error("when init component, occur error: %s", err.Error())
					return err
				}
			} else {
				return err
			}
		} else {
			log.Log.Debug("component type `%s` already exists, skip", comp.Type)
			continue
		}
	}
	return nil
}
