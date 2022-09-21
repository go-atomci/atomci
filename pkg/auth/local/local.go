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

package local

import (
	"fmt"

	"github.com/go-atomci/atomci/pkg/auth"

	"golang.org/x/crypto/bcrypt"
)

// Provider a local authentication provider.
// TODO: support configuration later
type Provider struct {
	name     string
	email    string
	user     string
	password string
}

// NewProvider creates a new local authentication provider.
func NewProvider(opts ...Option) auth.Provider {
	provider := &Provider{}
	for _, opt := range opts {
		opt(provider)
	}
	return provider
}

// Authenticate ..
func (p *Provider) Authenticate(loginUser, password string) (*auth.ExternalAccount, error) {
	_, err := CompareHashAndPassword(p.password, password)
	if err != nil {
		return nil, fmt.Errorf("comparehas password, error: %v", err.Error())
	}
	return &auth.ExternalAccount{
		Name:  p.name,
		Email: p.email,
		User:  p.user,
	}, nil

}

// CompareHashAndPassword ..
func CompareHashAndPassword(e string, p string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(e), []byte(p))
	if err != nil {
		return false, err
	}
	return true, nil
}
