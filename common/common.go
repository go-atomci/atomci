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

package common

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/go-atomci/atomci/dao"
	"github.com/go-atomci/atomci/middleware/log"
)

// GetUserToken ..
func GetUserToken(user string) (string, error) {
	urModel := dao.NewUserRolesModel()
	userInfo, err := urModel.GetUserByName(user)
	if err != nil {
		log.Log.Error("when get %v token by name occur error: %s", user, err.Error())
		return "", err
	}
	return userInfo.Token, nil
}

var (
	HttpClient = &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).Dial,
			MaxIdleConns:        200,
			MaxIdleConnsPerHost: 200,
			IdleConnTimeout:     30 * time.Second,
			TLSHandshakeTimeout: 5 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 60 * time.Second,
	}
)

// SentHTTPRequest ...
func SentHTTPRequest(method, urlStr string, body io.Reader, token string) ([]byte, error) {
	rep, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	rep.Header.Set("Content-Type", "application/json")
	if token != "" {
		rep.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := HttpClient.Do(rep)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		return respBody, nil
	} else {
		return nil, fmt.Errorf("%s", respBody)
	}
}
