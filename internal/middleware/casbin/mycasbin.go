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

import (
	"github.com/astaxie/beego"
	glog "log"

	"github.com/go-atomci/atomci/internal/middleware/log"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v2"
)

// NewCasbin ..
func NewCasbin() (*casbin.Enforcer, error) {

	rbacModel, err := model.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _
# g2 = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
# m = g(r.sub, p.sub) && r.obj == p.obj && (r.act == p.act || p.act == "*") || r.sub == "admin"
m = g(r.sub, p.sub) && keyMatch2(r.obj,p.obj) && (r.act == p.act || p.act == "*") || r.sub == "admin"
`)
	if err != nil {
		glog.Fatalf("error: model: %s", err)
	}

	dbUrl := beego.AppConfig.String("DB::url")
	rbacPolicy, _ := gormadapter.NewAdapter("mysql", dbUrl, true)

	e, err := casbin.NewEnforcer(rbacModel, rbacPolicy)
	if err != nil {
		log.Log.Error("casbin new enforcer error: %s", err.Error())
		return nil, err
	}
	if err := e.LoadPolicy(); err == nil {
		return e, err
	}
	log.Log.Error("casbin rbac_model or policy init error, message: %v", err)
	return nil, err
}
