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

package workflow

import (
	"time"

	httpclient "github.com/isbrick/http-client"
)

// Driver ..
type Driver int

//
const (
	DriverJenkins Driver = iota + 1
)

func (d Driver) String() (s string) {
	switch d {
	case DriverJenkins:
		return "jenkins"
	default:
		return "unknown"
	}
}

// WorkFlow ..
type WorkFlow interface {
	Ping() (string, error)
	Build() (int64, error)
	Abort(RunID int64) error
	GetJobInfo(runID int64) (*JobInfo, error)
}

// HTTPClient defined http native client
var (
	timeout    = 1000 * time.Millisecond
	HTTPClient = httpclient.NewHClient(httpclient.WithHTTPTimeout(timeout))
)

// JobInfo ..
type JobInfo struct {
	Artifacts         []interface{} `json:"artifacts"`
	Building          bool          `json:"building"`
	Description       interface{}   `json:"description"`
	DisplayName       string        `json:"displayName"`
	Duration          int           `json:"duration"`
	EstimatedDuration int           `json:"estimatedDuration"`
	Executor          interface{}   `json:"executor"`
	FullDisplayName   string        `json:"fullDisplayName"`
	ID                string        `json:"id"`
	Number            int           `json:"number"`
	QueueID           int           `json:"queueId"`
	Result            string        `json:"result"`
	Status            string        `json:"status"`
	StartTimeMillis   int64         `json:"startTimeMillis"`
	EndTimeMillis     int64         `json:"endTimeMillis"`
	DurationMillis    int           `json:"durationMillis"`
	Stages            []Stage       `json:"stages"`
}

// Stage job's stage
type Stage struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	ExecNode            string `json:"execNode"`
	Status              string `json:"status"`
	StartTimeMillis     int64  `json:"startTimeMillis"`
	DurationMillis      int    `json:"durationMillis"`
	PauseDurationMillis int    `json:"pauseDurationMillis"`
}
