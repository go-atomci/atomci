package migrations

import (
	"github.com/astaxie/beego/orm"
	"github.com/go-atomci/atomci/internal/core/settings"
	"time"
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
