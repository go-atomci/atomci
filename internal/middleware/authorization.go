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

package middleware

import (
	mycasbin "github.com/go-atomci/atomci/internal/middleware/casbin"
	"github.com/go-atomci/atomci/internal/middleware/log"

	"github.com/astaxie/beego/context"
)

// Authorization 鉴权
func Authorization(c *context.Context, username string) (bool, error) {
	e, err := mycasbin.NewCasbin()
	if err != nil {
		log.Log.Error("casbin new occur error: %v", err.Error())
		return false, err
	}
	urlPath := c.Request.URL.Path
	urlMethod := c.Request.Method
	res, err := e.Enforce(username, urlPath, urlMethod)
	// TODO: user constraint permission
	// based on  urlpath get resource type, then get resource constraint
	log.Log.Debug("role key: %s, path: %s, method: %s, res: %v", username, urlPath, urlMethod, res)
	return res, err
}
