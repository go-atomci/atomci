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

package ldap

import (
	"fmt"

	"github.com/go-atomci/atomci/internal/core/auth"
	"github.com/go-atomci/atomci/internal/middleware/log"

	"github.com/astaxie/beego"
	ldap "github.com/colynn/go-ldap-client/v3"
)

// Provider a ldap authentication provider.
// TODO: support configuration later
type Provider struct{}

// NewProvider creates a new ldap authentication provider.
func NewProvider() auth.Provider {
	return &Provider{}
}

// Authenticate ..
func (p *Provider) Authenticate(user, password string) (*auth.ExternalAccount, error) {
	port, _ := beego.AppConfig.Int("ldap::port")
	client := &ldap.Client{
		Base:               beego.AppConfig.String("ldap::baseDN"),
		Host:               beego.AppConfig.String("ldap::host"),
		Port:               port,
		UseSSL:             false,
		BindDN:             beego.AppConfig.String("ldap::bindDN"),
		BindPassword:       beego.AppConfig.String("ldap::bindPassword"),
		UserFilter:         beego.AppConfig.String("ldap::userFilter"),
		GroupFilter:        "(memberUid=%s)",
		Attributes:         []string{"givenName", "sn", "mail", "sAMAccountName"},
		SkipTLS:            true,
		InsecureSkipVerify: true,
	}
	defer client.Close()

	resp := map[string]string{}
	authVerify, resp, err := client.Authenticate(user, password)
	if authVerify == false {
		if err != nil {
			log.Log.Error("authVerify error: %s", err.Error())
		}
		return nil, fmt.Errorf("域帐号或密码错误，或请联系管理员")
	}

	// TODO: resp add verification
	return &auth.ExternalAccount{
		Name:  resp["sn"] + resp["givenName"],
		User:  resp["sAMAccountName"],
		Email: resp["mail"],
	}, nil
}
