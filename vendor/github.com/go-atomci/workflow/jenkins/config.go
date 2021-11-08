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

package jenkins

import "github.com/go-atomci/workflow"

// CommonContext ..
type CommonContext struct {
	JenkinsSlaveWorkspace string
	AccessToken           string
	RepoCnf               string
	AtomCIServer          string
}

// CallbackRequest ..
type CallbackRequest struct {
	Token string `json:"token"`
	URL   string `json:"url"`
	Body  string `json:"body"`
}

// StepItem Jenkins Step item struct defined
type StepItem struct {
	Name          string
	Command       string
	ContainerName string
}

// EnvItem build env
type EnvItem struct {
	Key   string
	Value interface{}
}

// ContainerEnv compile runtime env
type ContainerEnv struct {
	Name       string
	Image      string
	WorkingDir string
	CommandArr []string
	ArgsArr    []string
}

// CIContext ..
type CIContext struct {
	CommonContext
	RegistryAddr       string
	DockerConfig       string
	DockerAuth         string
	Stages             string
	EnvVars            []EnvItem
	ContainerTemplates []ContainerEnv
	CallBack           CallbackRequest
}

// DeployContext ..
type DeployContext struct {
	CommonContext
	EnvVars            []EnvItem
	ContainerTemplates []ContainerEnv
	HealthCheckItems   []*StepItem
	CallBack           CallbackRequest
}

// CustomScriptItem ...
type CustomScriptItem struct {
	StepItem
	Type string `json:"type"`
}

// Job for GetJob
type Job struct {
	Class   string `json:"_class"`
	Actions []struct {
		Class string `json:"_class,omitempty"`
	} `json:"actions"`
	Description           string        `json:"description"`
	DisplayName           string        `json:"displayName"`
	DisplayNameOrNull     interface{}   `json:"displayNameOrNull"`
	FullDisplayName       string        `json:"fullDisplayName"`
	FullName              string        `json:"fullName"`
	Name                  string        `json:"name"`
	URL                   string        `json:"url"`
	Buildable             bool          `json:"buildable"`
	Builds                []interface{} `json:"builds"`
	Color                 string        `json:"color"`
	FirstBuild            interface{}   `json:"firstBuild"`
	HealthReport          []interface{} `json:"healthReport"`
	InQueue               bool          `json:"inQueue"`
	KeepDependencies      bool          `json:"keepDependencies"`
	LastBuild             interface{}   `json:"lastBuild"`
	LastCompletedBuild    interface{}   `json:"lastCompletedBuild"`
	LastFailedBuild       interface{}   `json:"lastFailedBuild"`
	LastStableBuild       interface{}   `json:"lastStableBuild"`
	LastSuccessfulBuild   interface{}   `json:"lastSuccessfulBuild"`
	LastUnstableBuild     interface{}   `json:"lastUnstableBuild"`
	LastUnsuccessfulBuild interface{}   `json:"lastUnsuccessfulBuild"`
	NextBuildNumber       int64         `json:"nextBuildNumber"`
	Property              []struct {
		Class string `json:"_class"`
	} `json:"property"`
	QueueItem       interface{} `json:"queueItem"`
	ConcurrentBuild bool        `json:"concurrentBuild"`
	ResumeBlocked   bool        `json:"resumeBlocked"`
}

// JobBaseInfo ..
type JobBaseInfo struct {
	Class   string `json:"_class"`
	Actions []struct {
		Class  string `json:"_class,omitempty"`
		Causes []struct {
			Class            string `json:"_class"`
			ShortDescription string `json:"shortDescription"`
			UserID           string `json:"userId"`
			UserName         string `json:"userName"`
		} `json:"causes,omitempty"`
	} `json:"actions"`
	Artifacts         []interface{} `json:"artifacts"`
	Building          bool          `json:"building"`
	Description       interface{}   `json:"description"`
	DisplayName       string        `json:"displayName"`
	Duration          int           `json:"duration"`
	EstimatedDuration int           `json:"estimatedDuration"`
	Executor          interface{}   `json:"executor"`
	FullDisplayName   string        `json:"fullDisplayName"`
	ID                string        `json:"id"`
	KeepLog           bool          `json:"keepLog"`
	Number            int           `json:"number"`
	QueueID           int           `json:"queueId"`
	Result            string        `json:"result"`
	Timestamp         int64         `json:"timestamp"`
	URL               string        `json:"url"`
	ChangeSets        []interface{} `json:"changeSets"`
	Culprits          []interface{} `json:"culprits"`
	NextBuild         interface{}   `json:"nextBuild"`
	PreviousBuild     struct {
		Number int    `json:"number"`
		URL    string `json:"url"`
	} `json:"previousBuild"`
}

// JobDetailInfo ..
type JobDetailInfo struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
	ID                  string           `json:"id"`
	Name                string           `json:"name"`
	Status              string           `json:"status"`
	StartTimeMillis     int64            `json:"startTimeMillis"`
	EndTimeMillis       int64            `json:"endTimeMillis"`
	DurationMillis      int              `json:"durationMillis"`
	QueueDurationMillis int              `json:"queueDurationMillis"`
	PauseDurationMillis int              `json:"pauseDurationMillis"`
	Stages              []workflow.Stage `json:"stages"`
}
