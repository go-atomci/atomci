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

// Option represents the client options
type Option func(*HClient)

// WithHTTPTimeout sets hystrix timeout
func WithHTTPTimeout(timeout time.Duration) Option {
	return func(c *HClient) {
		c.timeout = timeout
	}
}

// WithRetryCount sets the retry count for the hystrixHTTPClient
func WithRetryCount(retryCount int) Option {
	return func(c *HClient) {
		c.retryCount = retryCount
	}
}

// WithRetrier sets the strategy for retrying
func WithRetrier(retrier Retriable) Option {
	return func(c *HClient) {
		c.retrier = retrier
	}
}
