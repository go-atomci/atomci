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
	glog "log"

	"github.com/go-atomci/atomci/middleware/log"
	tools "github.com/go-atomci/atomci/utils"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

// NewCasbin ..
func NewCasbin() (*casbin.Enforcer, error) {
	// TODO: changet to csv tmp, later add mysql apter
	// databaseURL := beego.AppConfig.String("DB::url")
	// Apter, err := gormadapter.NewAdapter("mysql", databaseURL, true)
	// if err != nil {
	// 	return nil, err
	// }

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

	rbacPolicyPath := tools.EnsureAbs("conf/rbac_policy.csv")
	rbacPolicy := fileadapter.NewAdapter(rbacPolicyPath)

	// TODO: change to csv tmp, enable mysql apter later
	// e, err := casbin.NewEnforcer(rbacConf, Apter)
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
