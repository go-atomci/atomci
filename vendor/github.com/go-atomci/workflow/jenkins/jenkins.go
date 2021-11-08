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

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"

	"github.com/go-atomci/workflow"
	"github.com/go-atomci/workflow/jenkins/templates"
)

// BaseContext Jenkins Pipeline base template
type BaseContext struct {
	Pipeline string
}

// jenkinsWorker Create/Get/Update jenkins job
type jenkinsWorker struct {
	url        string
	user       string
	token      string
	crumbKey   string
	crumbValue string
	jobName    string
}

type crumbIssuerResp struct {
	Crumb             string `json:"crumb,omitempty"`
	CrumbRequestField string `json:"crumb_request_field,omitempty"`
}

// Jenkins jenkins struct
type Jenkins struct {
	jenkinsWorker
	processor FlowProcessor
}

type Option func(*Jenkins)

func URL(url string) Option {
	return func(w *Jenkins) {
		w.url = url
	}
}

func JenkinsUser(user string) Option {
	return func(w *Jenkins) {
		w.user = user
	}
}

func JenkinsToken(token string) Option {
	return func(w *Jenkins) {
		w.token = token
	}
}

func JenkinsJob(jobName string) Option {
	return func(w *Jenkins) {
		w.jobName = jobName
	}
}

func Processor(processor FlowProcessor) Option {
	return func(w *Jenkins) {
		w.processor = processor
	}
}

// FlowProcessor ..
type FlowProcessor interface {
	// return run id and error
	Run(addr, user, token, crumbKey, crumbValue, jobName string, param []byte) (int64, error)
}

// NewJenkinsClient ...
func NewJenkinsClient(opts ...Option) (workflow.WorkFlow, error) {
	jenkins := Jenkins{}
	for _, opt := range opts {
		opt(&jenkins)
	}
	return &jenkins, nil
}

// Build create a jenkins build job & trigger it.
func (w *Jenkins) Build() (int64, error) {
	if err := w.crumbHeaderVerify(); err != nil {
		return 0, err
	}
	param, err := json.Marshal(w.processor)
	if err != nil {
		return 0, err
	}
	return w.processor.Run(w.url, w.user, w.token, w.crumbKey, w.crumbValue, w.jobName, param)
}

func (w *Jenkins) Ping() (string, error) {
	respHeader, _, err := w.getCrumbRequestHeader()
	if err != nil {
		return "", err
	}
	return respHeader.Get("X-Jenkins"), nil
}

func (w *Jenkins) getCrumbRequestHeader() (http.Header, *crumbIssuerResp, error) {
	url := fmt.Sprintf("%v/crumbIssuer/api/json", strings.TrimSuffix(w.url, "/"))
	respHeader, respBody, err := sentHTTPRequest("GET", w.user, w.token, "", "", url, nil)
	if err != nil {
		return respHeader, nil, err
	}
	respJSON := crumbIssuerResp{}
	if err := json.Unmarshal(respBody, &respJSON); err != nil {
		return respHeader, nil, err
	}
	return respHeader, &respJSON, nil
}

// Abort jenkins build job
func (w *Jenkins) Abort(runID int64) error {
	if err := w.crumbHeaderVerify(); err != nil {
		return err
	}
	url := fmt.Sprintf("%v/job/%v/%v/stop", strings.TrimSuffix(w.url, "/"), w.jobName, runID)
	_, _, err := sentHTTPRequest("POST", w.user, w.token, w.crumbKey, w.crumbValue, url, nil)
	if err != nil {
		return err
	}
	return nil
}

func (w *Jenkins) crumbHeaderVerify() error {
	if w.crumbKey == "" || w.crumbValue == "" {
		_, respBody, err := w.getCrumbRequestHeader()
		if err != nil {
			return err
		}
		if respBody.CrumbRequestField == "" {
			respBody.CrumbRequestField = "Jenkins-Crumb"
		}
		w.crumbKey = respBody.CrumbRequestField
		w.crumbValue = respBody.Crumb
	}
	return nil
}

// GeneratePipelineXMLStr ..
func GeneratePipelineXMLStr(templateStr string, context interface{}) (string, error) {
	pipelineTemplate := template.Must(template.New("pipline").Parse(templateStr))
	var pipelineBuf bytes.Buffer
	err := pipelineTemplate.Execute(&pipelineBuf, context)
	return pipelineBuf.String(), err
}

// getJobInfo
func (w *Jenkins) getJobInfo(runID int64) (*JobBaseInfo, error) {
	url := fmt.Sprintf("%v/job/%v/%v/api/json?pretty=true", strings.TrimSuffix(w.url, "/"), w.jobName, runID)
	_, respBody, err := sentHTTPRequest("GET", w.user, w.token, w.crumbKey, w.crumbValue, url, nil)
	if err != nil {
		return nil, err
	}
	respJSON := JobBaseInfo{}
	if err := json.Unmarshal(respBody, &respJSON); err != nil {
		return nil, err
	}
	return &respJSON, nil
}

// GetJobInfo ..
func (w *Jenkins) GetJobInfo(runID int64) (*workflow.JobInfo, error) {
	if err := w.crumbHeaderVerify(); err != nil {
		return nil, err
	}
	jobBaseInfo, err := w.getJobInfo(runID)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%v/job/%v/%v/wfapi/describe", strings.TrimSuffix(w.url, "/"), w.jobName, runID)
	_, respBody, err := sentHTTPRequest("GET", w.user, w.token, w.crumbKey, w.crumbValue, url, nil)
	if err != nil {
		return nil, err
	}
	detailInfo := JobDetailInfo{}
	if err := json.Unmarshal(respBody, &detailInfo); err != nil {
		return nil, err
	}

	return &workflow.JobInfo{
		Artifacts:         jobBaseInfo.Artifacts,
		Building:          jobBaseInfo.Building,
		Description:       jobBaseInfo.Description,
		DisplayName:       jobBaseInfo.DisplayName,
		Duration:          jobBaseInfo.Duration,
		EstimatedDuration: jobBaseInfo.EstimatedDuration,
		Executor:          jobBaseInfo.Executor,
		FullDisplayName:   jobBaseInfo.FullDisplayName,
		ID:                jobBaseInfo.ID,
		Number:            jobBaseInfo.Number,
		QueueID:           jobBaseInfo.QueueID,
		Result:            jobBaseInfo.Result,
		Status:            detailInfo.Status,
		DurationMillis:    detailInfo.DurationMillis,
		EndTimeMillis:     detailInfo.EndTimeMillis,
		StartTimeMillis:   detailInfo.StartTimeMillis,
		Stages:            detailInfo.Stages,
	}, nil
}

// Run CIPipeline
func (buildflow *CIContext) Run(addr, user, token, crumbKey, crumbValue, jobName string, param []byte) (int64, error) {
	// get configxml
	data := CIContext{}
	err := json.Unmarshal(param, &data)
	if err != nil {
		return 0, err
	}
	pipelineStrs, err := buildflow.GetCIPipelineXML(data)
	if err != nil {
		return 0, err
	}
	jw := jenkinsWorker{
		url:        addr,
		user:       user,
		token:      token,
		crumbKey:   crumbKey,
		crumbValue: crumbValue,
		jobName:    jobName,
	}
	err = jw.CreateOrUpdateJob(pipelineStrs)
	if err != nil {
		return 0, err
	}
	jobInfo, err := jw.GetJob()
	if err != nil {
		return 0, err
	}
	nextBuildNumber := jobInfo.NextBuildNumber
	return buildNow(addr, user, token, crumbKey, crumbValue, jobName, nextBuildNumber)
}

// Run Default Pipeline
func (deployflow *DeployContext) Run(addr, user, token, crumbKey, crumbValue, jobName string, param []byte) (int64, error) {
	// get configxml
	data := DeployContext{}
	err := json.Unmarshal(param, &data)
	if err != nil {
		return 0, err
	}
	pipelineStrs, err := deployflow.GetDeployPipelineXML(data)
	if err != nil {
		return 0, err
	}
	jw := jenkinsWorker{
		url:        addr,
		user:       user,
		token:      token,
		crumbKey:   crumbKey,
		crumbValue: crumbValue,
		jobName:    jobName,
	}
	err = jw.CreateOrUpdateJob(pipelineStrs)
	if err != nil {
		return 0, err
	}
	jobInfo, err := jw.GetJob()
	if err != nil {
		return 0, err
	}
	nextBuildNumber := jobInfo.NextBuildNumber
	return buildNow(addr, user, token, crumbKey, crumbValue, jobName, nextBuildNumber)
}

func buildNow(addr, user, token, crumbKey, crumbValue, jobName string, nextBuildNumber int64) (int64, error) {
	url := fmt.Sprintf("%v/job/%v/build?delay=0sec", strings.TrimSuffix(addr, "/"), jobName)
	// TODO: add debug log
	if _, _, err := sentHTTPRequest("POST", user, token, crumbKey, crumbValue, url, nil); err != nil {
		return 0, err
	}
	return nextBuildNumber, nil
}

// GetCIPipelineXML ..
func (buildflow *CIContext) GetCIPipelineXML(context CIContext) (string, error) {
	pipelineTemplate := template.Must(template.New("pipline").Parse(templates.CIPipeline))
	var pipelineBuf bytes.Buffer
	err := pipelineTemplate.Execute(&pipelineBuf, context)
	if err != nil {
		return "", err
	}
	baseContext := BaseContext{
		Pipeline: pipelineBuf.String(),
	}
	baseTemplate := template.Must(template.New("base").Parse(templates.BaseXML))
	var baseBuf bytes.Buffer
	err = baseTemplate.Execute(&baseBuf, baseContext)

	return baseBuf.String(), err
}

// GetDeployPipelineXML ..
func (deployflow *DeployContext) GetDeployPipelineXML(context DeployContext) (string, error) {
	pipelineTemplate := template.Must(template.New("pipline").Parse(templates.DeployPipeline))
	var pipelineBuf bytes.Buffer
	err := pipelineTemplate.Execute(&pipelineBuf, context)
	if err != nil {
		return "", err
	}
	baseContext := BaseContext{
		Pipeline: pipelineBuf.String(),
	}
	baseTemplate := template.Must(template.New("base").Parse(templates.DeployBaseXML))
	var baseBuf bytes.Buffer
	err = baseTemplate.Execute(&baseBuf, baseContext)

	return baseBuf.String(), err
}

// GetJob ..
func (jw *jenkinsWorker) GetJob() (*Job, error) {
	var job Job
	url := fmt.Sprintf("%v/job/%v/api/json", strings.TrimSuffix(jw.url, "/"), jw.jobName)
	_, rspBody, err := sentHTTPRequest("GET", jw.user, jw.token, jw.crumbKey, jw.crumbValue, url, nil)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(rspBody, &job); err != nil {
		return nil, err
	}
	return &job, err
}

// CreateOrUpdateJob ..
func (jw *jenkinsWorker) CreateOrUpdateJob(configXML string) error {
	if _, err := jw.GetJob(); err != nil {
		if strings.Contains(err.Error(), "404") {
			if err = jw.CreateJob(configXML); err != nil {
				return err
			}
		}
		return err
	}
	if err := jw.UpdateJob(configXML); err != nil {
		// TODO: add log
		return err
	}
	return nil
}

// CreateJob ..
func (jw *jenkinsWorker) CreateJob(configXML string) error {
	url := fmt.Sprintf("%v/createItem?name=%v", strings.TrimSuffix(jw.url, "/"), jw.jobName)
	payload := bytes.NewBuffer([]byte(configXML))
	if _, _, err := sentHTTPRequest("POST", jw.user, jw.token, jw.crumbKey, jw.crumbValue, url, payload); err != nil {
		return err
	}
	return nil
}

// UpdateJob ..
func (jw *jenkinsWorker) UpdateJob(configXML string) error {
	url := fmt.Sprintf("%v/job/%v/config.xml", strings.TrimSuffix(jw.url, "/"), jw.jobName)
	payload := bytes.NewBuffer([]byte(configXML))
	if _, _, err := sentHTTPRequest("POST", jw.user, jw.token, jw.crumbKey, jw.crumbValue, url, payload); err != nil {
		return err
	}
	return nil
}

// sentHTTPRequest ..
func sentHTTPRequest(method, user, token, crumbKey, crumbValue, url string, body io.Reader) (http.Header, []byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Accept", "application/xml")
	if crumbKey != "" && crumbValue != "" {
		req.Header.Set(crumbKey, crumbValue)
	}
	req.SetBasicAuth(user, token)
	rsp, err := workflow.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer rsp.Body.Close()

	respBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return rsp.Header, nil, err
	}

	if rsp.StatusCode == http.StatusOK || rsp.StatusCode == http.StatusCreated {
		return rsp.Header, respBody, nil
	} else if rsp.StatusCode == http.StatusNotFound {
		return rsp.Header, nil, fmt.Errorf("404 not found")
	} else if rsp.StatusCode == http.StatusUnauthorized {
		return rsp.Header, nil, fmt.Errorf("401 unauthorized")
	} else if rsp.StatusCode == http.StatusBadRequest {
		return rsp.Header, nil, fmt.Errorf("400 badRequest")
	}
	return rsp.Header, nil, fmt.Errorf(string(respBody))
}
