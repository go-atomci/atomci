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

package pipelinemgr

import (
	"encoding/json"
)

// PipelineBaseReq ..
type PipelineBaseReq struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	ID          int64  `json:"id"`
	Enabled     bool   `json:"enabled"`
	IsDefault   bool   `json:"is_default"`
}

// PipelineSetupReq ..
type PipelineSetupReq struct {
	PipelineBaseReq
	Configs []*PipelineStageStruct `json:"configs"`
}

// PipelineConfig ..
type PipelineConfig []*PipelineStageStruct

// PipelineStageStruct ..
type PipelineStageStruct struct {
	StageID    int64         `json:"stage_id"`
	PipelineID int64         `json:"pipeline_id"`
	Index      int64         `json:"index"`
	ID         int64         `json:"id"`
	Steps      PipelineSteps `json:"steps"`
	Name       string        `json:"name,omitempty"`
	// instance part
	PipelineInstanceID int64 `json:"pipeline_instance_id,omitempty"`
}

// PipelineSteps ..
type PipelineSteps []*struct {
	Name        string     `json:"name"`
	Type        string     `json:"type"`
	TypeDisplay string     `json:"type_display,omitempty"`
	StepID      int64      `json:"step_id"`
	Index       int        `json:"index"`
	Driver      string     `json:"driver"`
	SubTask     []*subTask `json:"sub_task"`
}

type subTask struct {
	Index  int          `json:"index,omitempty"`
	Name   string       `json:"name,omitempty"`
	Type   string       `json:"type,omitempty"`
	Params []compileEnv `json:"params,omitempty"`
}

type compileEnv struct {
	Name            string `json:"name,omitempty"`
	Version         string `json:"version,omitempty"`
	Image           string `json:"image,omitempty"`
	WorkingDir      string `json:"working_dir,omitempty"`
	Command         string `json:"command,omitempty"`
	Args            string `json:"args,omitempty"`
	CompileCommpand string `json:"compile_commpand,omitempty"` // compile command, eg mvn, npm
}

// String ...
func (p *PipelineStageStruct) String() (string, error) {
	bytes, err := json.Marshal(p.Steps)
	return string(bytes), err
}

// Struct ..
func (config PipelineConfig) Struct(sc string) (PipelineConfig, error) {
	err := json.Unmarshal([]byte(sc), &config)
	return config, err
}

// Struct ...
func (p PipelineSteps) Struct(sc string) (PipelineSteps, error) {
	err := json.Unmarshal([]byte(sc), &p)
	return p, err
}

// StepTemplate ...
type StepTemplate []*struct {
	Index        int      `json:"index"`
	Enable       bool     `json:"enable"`
	Role         []string `json:"role"`
	Name         string   `json:"name"`
	SystemCode   string   `json:"system_code"`
	TargetBranch []string `json:"target_branch,omitempty"`
}

// Struct ...
func (p StepTemplate) Struct(sc string) (StepTemplate, error) {
	err := json.Unmarshal([]byte(sc), &p)
	return p, err
}
