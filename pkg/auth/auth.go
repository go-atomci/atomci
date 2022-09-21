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

package auth

// ExternalAccount contains queried information returned by an authenticate provider
// for an external account.
type ExternalAccount struct {
	// REQUIRED: The username of the account.
	User string
	// The nick name of the account.
	Name string
	// The email address of the account.
	Email string
	// Whether the user should be prompted as a site admin.
	Admin bool
}

// Provider defines an authenticate provider which provides ability to authentication against
// an external identity provider and query external account information.
type Provider interface {
	Authenticate(login, password string) (*ExternalAccount, error)
}
