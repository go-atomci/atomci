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
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/go-atomci/atomci/internal/core/settings"
)

type Migration20220324 struct {
}

func (m Migration20220324) GetCreateAt() time.Time {
	return time.Date(2022, 3, 24, 0, 0, 0, 0, time.Local)
}

func (m Migration20220324) Upgrade(ormer orm.Ormer) error {
	pm := settings.NewSettingManager()
	k8sSettings, err := pm.GetIntegrateSettings([]string{"kubernetes"})
	if err != nil {
		return err
	}
	for _, setting := range k8sSettings {
		req := &setting.IntegrateSettingReq
		cfg := req.Config.(*settings.KubeConfig)
		if cfg.Type == "" {
			cfg.Type = settings.KubernetesConfig
			cfg.URL = ""
			err = pm.UpdateIntegrateSetting(req, setting.ID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
