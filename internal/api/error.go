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

import (
	"github.com/go-atomci/atomci/utils/errors"
)

type ErrorController struct {
	BaseController
}

func (this *ErrorController) Error404() {
	err := errors.NewNotFound()
	this.ServeError(err)
}

func (this *ErrorController) Error405() {
	err := errors.NewMethodNotAllowed()
	this.ServeError(err)
}
