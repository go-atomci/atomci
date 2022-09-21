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

	"github.com/go-atomci/atomci/pkg/auth"

	"github.com/astaxie/beego"
	ldap "github.com/colynn/go-ldap-client/v3"
)

// Provider a ldap authentication provider.
// TODO: support configuration later
type Provider struct {
	baseDN       string
	host         string
	port         int
	bindDN       string
	bindPassword string
	userFilter   string
}

// NewProvider creates a new ldap authentication provider.
func NewProvider(opts ...Option) auth.Provider {
	provider := &Provider{}
	for _, opt := range opts {
		opt(provider)
	}
	return provider
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

	authVerify, resp, err := client.Authenticate(user, password)
	if !authVerify {
		return nil, fmt.Errorf("authVerify error: %v", err)
	}

	// TODO: resp add verification
	return &auth.ExternalAccount{
		Name:  resp["sn"] + resp["givenName"],
		User:  resp["sAMAccountName"],
		Email: resp["mail"],
	}, nil
}
