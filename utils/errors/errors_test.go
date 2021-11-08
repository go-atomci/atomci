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

package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		code := "oops"
		message := "don't panic"
		cause := fmt.Errorf("undefined reference")

		err := Error{}
		assert.Equal(t, 0, err.Status())

		assert.Equal(t, &err, err.SetCode(code))
		assert.Equal(t, code, err.Code())

		assert.Equal(t, &err, err.SetMessage(message))
		assert.Equal(t, message, err.Message())

		assert.Equal(t, &err, err.SetCause(cause))
		assert.Equal(t, cause, err.Cause())

		assert.NotEmpty(t, err.Error())
	})

	t.Run("4XX", func(t *testing.T) {
		errs := []*Error{
			NewBadRequest(),
			NewNotFound(),
			NewConflict(),
			NewUnauthorized(),
			NewForbidden(),
			NewMethodNotAllowed(),
		}
		for _, err := range errs {
			assert.NotNil(t, err)
			assert.True(t, err.Status() >= 400 && err.Status() < 500, fmt.Sprintf(`should be 4XX: %#v`, err))
		}
	})

	t.Run("5XX", func(t *testing.T) {
		errs := []*Error{
			NewInternalServerError(),
		}
		for _, err := range errs {
			assert.NotNil(t, err)
			assert.True(t, err.Status() >= 500 && err.Status() < 600, fmt.Sprintf(`should be 5XX: %#v`, err))
		}
	})
}
