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
	"sync"

	"github.com/astaxie/beego"
	"github.com/go-atomci/atomci/internal/middleware/log"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	beegoormadapter "github.com/go-atomci/atomci/pkg/beego-orm-adapter"
	_ "github.com/go-sql-driver/mysql"
)

var casbinadapter *beegoormadapter.Adapter
var casbinadapterOnce sync.Once
var casbinErr error

// GetOrmer :set ormer singleton
func GetAdapter() (*beegoormadapter.Adapter, error) {
	casbinadapterOnce.Do(func() {
		casbinadapter, casbinErr = initAdapter()
	})
	return casbinadapter, casbinErr
}

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
	a, err := GetAdapter()
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer(rbacModel, a)
	if err != nil {
		log.Log.Error("casbin new enforcer error: %s", err.Error())
		return nil, err
	}
	if err := e.LoadPolicy(); err != nil {
		log.Log.Error("casbin rbac_model or policy init error, message: %v", err.Error())
		return e, err
	}
	return e, nil
}

func initAdapter() (*beegoormadapter.Adapter, error) {
	dsn := beego.AppConfig.String("DB::url")
	a, err := beegoormadapter.NewAdapter("casbin", "mysql", dsn)
	if err != nil {
		log.Log.Error("beego orm adapter error: %s", err.Error())
		return nil, err
	}
	return a, nil
}
