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

type Migration20230703 struct {
}

func (m Migration20230703) GetCreateAt() time.Time {
	return time.Date(2023, 7, 03, 0, 0, 0, 0, time.Local)
}

func (m Migration20230703) Upgrade(ormer orm.Ormer) error {

	addedCompileEnvs := []settings.CompileEnvReq{
		{
			Name:        "checkout",
			Image:       "colynn/checkout:latest",
			Description: "代码检出",
		},
	}

	// init compile envs
	_ = initCompileEnvs(addedCompileEnvs)

	// update origin jnlp image
	return updateCompileEnvs(ormer)
}

func updateCompileEnvs(ormer orm.Ormer) error {
	_, err := ormer.Raw("UPDATE `sys_compile_env` SET `image`='jenkins/inbound-agent:latest',`description`='默认jenkins jnlp agent' WHERE `name`='jnlp';").Exec()
	if err != nil {
		return err
	}
	return nil
}
