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
	"os"

	"github.com/go-atomci/atomci/internal/middleware/log"
	"github.com/go-atomci/atomci/utils/errors"
)

func Init() {

	// 注册/更新资源
	initResource()

	// 初始化/更新路由
	initRouterItems()

	// 更新所有用户权限策略
	// TODO: confirm
	// func initUsers(){
	// users, _ := dao.UserList()
	// for _, user := range users {
	// 	dao.InitSystemMember(user)
	// }
	//}()

	// 初始化系统组
	if err := InitAdminUserAndGroup(); err != nil {
		if !errors.OrmError1062(err) {
			log.Log.Error(err.Error())
			os.Exit(2)
		}
	}

	/*
		TODO: Below resources just run once
	*/
	if err := Component(); err != nil {
		os.Exit(2)
	}

	// init compile envs
	initCompileEnvs()

	// init task tmpls
	if err := initTaskTemplates(); err != nil {
		os.Exit(2)
	}
}
