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

package harbor

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/go-atomci/atomci/common"
)

func PingHarbor(addr, user, password string, https bool) error {
	method := "GET"
	protocol := "http"
	if https {
		protocol = "https"
	}

	urlStr := fmt.Sprintf("%s://%v/api/users", protocol, addr)
	if _, err := harborAuth(user, password, method, urlStr, nil); err != nil {
		return fmt.Errorf("harbor auth login failed: %v, please check url and account info", err.Error())
	}

	return nil
}

func harborAuth(user, password, method, urlStr string, body io.Reader) ([]byte, error) {
	rep, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	rep.Header.Set("Content-Type", "application/json")
	rep.SetBasicAuth(user, password)
	resp, err := common.HttpClient.Do(rep)
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
	} else if resp.StatusCode == http.StatusCreated {
		return respBody, nil
	} else if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("401 Unauthorized")
	}

	return nil, err
}
