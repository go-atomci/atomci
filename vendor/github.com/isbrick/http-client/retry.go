// Copyright 2020 colynn.liu
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package httpclient

import "time"

// Retriable ..
type Retriable interface {
	NextInterval(retry int) time.Duration
}

type retrier struct {
	backoff Backoff
}

var _ Retriable = (*retrier)(nil)

// NewRetrier returns retrier with some backoff strategy
func NewRetrier(backoff Backoff) Retriable {
	return &retrier{
		backoff: backoff,
	}
}

// NextInterval returns next retriable time
func (r *retrier) NextInterval(retry int) time.Duration {
	return r.backoff.Next(retry)
}

type noRetrier struct {
}

var _ Retriable = (*noRetrier)(nil)

// NewNoRetrier returns a null object for retriable
func NewNoRetrier() Retriable {
	return &noRetrier{}
}

// NextInterval returns next retriable time, always 0
func (r *noRetrier) NextInterval(retry int) time.Duration {
	return 0 * time.Millisecond
}
