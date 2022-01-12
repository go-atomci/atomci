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

package api

// Result ..
type Result interface{}

type SuccessResult struct {
	IsSuccess bool        `json:"IsSuccess"`
	Data      interface{} `json:"Data,omitempty"`
	Message   string      `json:"Message,omitempty"`
	ErrMsg    string      `json:"ErrMsg,omitempty"`
}

type ErrorResult struct {
	IsSuccess bool   `json:"IsSuccess"`
	ErrCode   string `json:"ErrCode,omitempty"`
	ErrMsg    string `json:"ErrMsg,omitempty"`
	ErrDetail string `json:"ErrDetail,omitempty"`
}

func NewResult(isSuccess bool, data interface{}, errMsg string) Result {
	return &SuccessResult{IsSuccess: isSuccess, Data: data, ErrMsg: errMsg}
}

func NewErrorResult(errCode, errMsg, errDetail string) Result {
	return &ErrorResult{
		IsSuccess: false,
		ErrCode:   errCode,
		ErrMsg:    errMsg,
		ErrDetail: errDetail,
	}
}

func NewSuccessResult(data ...interface{}) Result {
	var result interface{}
	if len(data) > 0 {
		result = data[0]
	}
	return &SuccessResult{
		IsSuccess: true,
		Data:      result,
	}
}
