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

	"github.com/go-atomci/atomci/internal/core/auth"
	"github.com/go-atomci/atomci/internal/dao"
	"github.com/go-atomci/atomci/internal/middleware/log"

	"golang.org/x/crypto/bcrypt"
)

// Provider a local authentication provider.
// TODO: support configuration later
type Provider struct{}

// NewProvider creates a new local authentication provider.
func NewProvider() auth.Provider {
	return &Provider{}
}

// Authenticate ..
func (p *Provider) Authenticate(loginUser, password string) (*auth.ExternalAccount, error) {
	userModel, err := dao.GetUser(loginUser)
	if err != nil {
		log.Log.Error("get user error: %v")
		return nil, fmt.Errorf("用户不存在或密码错误")
	}

	_, err = CompareHashAndPassword(userModel.Password, password)
	if err != nil {
		log.Log.Error("comparehas password, error: %v", err.Error())
		return nil, fmt.Errorf("用户不存在或密码错误")
	}
	return &auth.ExternalAccount{
		Name:  userModel.Name,
		Email: userModel.Email,
		User:  userModel.User,
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
