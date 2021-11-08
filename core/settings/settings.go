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

package settings

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/go-atomci/atomci/dao"
	"github.com/go-atomci/atomci/middleware/log"
	"github.com/go-atomci/atomci/models"
	"github.com/go-atomci/atomci/pkg/harbor"
	"github.com/go-atomci/atomci/utils/query"
	"github.com/go-atomci/atomci/utils/validate"

	"github.com/astaxie/beego"
	"github.com/go-atomci/workflow/jenkins"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// SettingManager ...
type SettingManager struct {
	model *dao.SysSettingModel
}

// IntegrateSettingResponse create stage
type IntegrateSettingResponse struct {
	IntegrateSettingReq
	Creator  string     `json:"creator,omitempty"`
	CreateAt *time.Time `json:"create_at,omitempty"`
	UpdateAt *time.Time `json:"update_at,omitempty"`
	ID       int64      `json:"id,omitempty"`
}

// VerifyResponse   integrate verify
type VerifyResponse struct {
	Msg   string `json:"msg,omitempty"`
	Error error  `json:"error,omitempty"`
}

// IntegrateSettingReq ..
type IntegrateSettingReq struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Config      interface{} `json:"config,omitempty"`
	Type        string      `json:"type"`
}

// const variables
const (
	KubernetesType = "kubernetes"
	HarborType     = "harbor"
	JenkinsType    = "jenkins"
)

type Config struct{}

type BaseConfig struct {
	URL  string `json:"url,omitempty"`
	User string `json:"user,omitempty"`
}

type KubeConfig struct {
	URL  string `json:"url,omitempty"`
	Conf string `json:"conf,omitempty"`
}
type HarborConfig struct {
	BaseConfig
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}

type JenkinsConfig struct {
	BaseConfig
	Token     string `json:"token,omitempty"`
	WorkSpace string `json:"workspace,omitempty"`
}

func (intergrateItem *IntegrateSettingReq) String() (string, error) {
	bytes, err := json.Marshal(intergrateItem.Config)
	return string(bytes), err
}

// Struct ...
func (config *Config) Struct(sc string, settingType string) (interface{}, error) {
	switch strings.ToLower(settingType) {
	case "kubernetes":
		kube := &KubeConfig{}
		err := json.Unmarshal([]byte(sc), kube)
		return kube, err
	case "jenkins":
		jenkins := &JenkinsConfig{}
		err := json.Unmarshal([]byte(sc), jenkins)
		return jenkins, err
	case "harbor":
		harbor := &HarborConfig{}
		err := json.Unmarshal([]byte(sc), harbor)
		return harbor, err
	default:
		log.Log.Warn("this settings type %s is not support, return origin string", settingType)
		return sc, nil
	}
}

// NewSettingManager ...
func NewSettingManager() *SettingManager {
	return &SettingManager{
		model: dao.NewSysSettingModel(),
	}
}

// GetIntegrateSettings ..
func (pm *SettingManager) GetIntegrateSettings(integrateType string) ([]*IntegrateSettingResponse, error) {
	items, err := pm.model.GetIntegrateSettings(integrateType)
	if err != nil {
		log.Log.Error("get interate settings error: %s", err.Error())
		return nil, err
	}
	rsp := formatIntegrateSettingResponse(items)
	return rsp, err
}

// GetIntegrateSettingByID ..
func (pm *SettingManager) GetIntegrateSettingByID(id int64) (*IntegrateSettingResponse, error) {
	integrateSetting, err := pm.model.GetIntegrateSettingByID(id)
	if err != nil {
		log.Log.Error("when GetGetIntegrateSettingByID, get GetIntegrateSettingByID occur error: %s", err.Error())
		return nil, err
	}
	config := &Config{}
	return formatSignalIntegrateSetting(integrateSetting, config), err
}

// GetIntegrateSettingsByPagination ..
func (pm *SettingManager) GetIntegrateSettingsByPagination(filter *query.FilterQuery) (*query.QueryResult, error) {
	queryResult, settingsList, err := pm.model.GetIntegrateSettingsByPagination(filter)
	if err != nil {
		return nil, err
	}
	rsp := formatIntegrateSettingResponse(settingsList)
	queryResult.Item = rsp
	return queryResult, err
}

// UpdateIntegrateSetting ..
func (pm *SettingManager) UpdateIntegrateSetting(request *IntegrateSettingReq, stepID int64) error {
	stageModel, err := pm.model.GetIntegrateSettingByID(stepID)
	if err != nil {
		return err
	}
	if request.Name != "" {
		stageModel.Name = request.Name
		if stageModel.Type == KubernetesType {
			if err := validate.ValidateKubernetesName(request.Name); err != nil {
				return err
			}
		}
	}

	if request.Description != "" {
		stageModel.Description = request.Description
	}

	if request.Type != "" {
		stageModel.Type = request.Type
	}

	config, err := request.String()
	if err != nil {
		log.Log.Error("json marshal error: %s", err.Error())
		return err
	}
	stageModel.Config = config
	if request.Type == KubernetesType {
		kube := &KubeConfig{}
		err := json.Unmarshal([]byte(config), kube)
		if err == nil {
			pm.createOrupateKubernetesConfig(request.Name, kube.Conf)
		} else {
			log.Log.Error("kuber conf format error:  %v", err.Error())
		}
	}
	return pm.model.UpdateIntegrateSetting(stageModel)
}

func (pm *SettingManager) createOrupateKubernetesConfig(clusterName, config string) error {
	configPath := beego.AppConfig.String("k8s::configPath")

	log.Log.Debug("configPath: %v", configPath)
	err := os.MkdirAll(configPath, 0766)
	if err != nil {
		log.Log.Error(fmt.Sprintf("Failed to make the k8sconfig dir: %v", err.Error()))
		return err
	}
	fileObj, err := os.OpenFile(configPath+"/"+clusterName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Log.Error(fmt.Sprintf("Failed to open the file: %v", err.Error()))
		return err
	}
	if _, err := io.WriteString(fileObj, config); err != nil {
		log.Log.Error(fmt.Sprintf("init K8S cluster %v configure failed: %v", clusterName, err.Error()))
		return err
	}
	log.Log.Debug(fmt.Sprintf("update K8S cluster %v configure successfully", clusterName))
	return nil
}

// VerifyIntegrateSetting ..
func (pm *SettingManager) VerifyIntegrateSetting(request *IntegrateSettingReq) VerifyResponse {
	resp := VerifyResponse{}
	config, err := request.String()
	if err != nil {
		log.Log.Error("json marshal error: %s", err.Error())
		resp.Error = err
		return resp
	}

	switch strings.ToLower(request.Type) {
	case KubernetesType:
		kube := &KubeConfig{}
		err := json.Unmarshal([]byte(config), kube)
		if err != nil {
			log.Log.Error("kuber conf format error:  %v", err.Error())
			resp.Error = err
			return resp
		}
		k8sconf, err := clientcmd.RESTConfigFromKubeConfig([]byte(kube.Conf))
		if err != nil {
			resp.Error = err
			return resp
		}
		clientset, err := kubernetes.NewForConfig(k8sconf)
		if err != nil {
			resp.Error = err
			return resp
		}
		k8sVersion, err := clientset.Discovery().ServerVersion()
		if err != nil {
			log.Log.Error("get kubernetes verison error: %s", err.Error())
			resp.Error = err
			return resp
		}
		msg := fmt.Sprintf("Connected to Kubernetes %s", k8sVersion.GitVersion)
		resp.Msg = msg
	case HarborType:
		harborConf := &HarborConfig{}
		err := json.Unmarshal([]byte(config), harborConf)
		if err != nil {
			log.Log.Error("harborConf conf format error:  %v", err.Error())
			resp.Error = err
		} else {
			log.Log.Debug("verify harbor conf: %v", harborConf)
			if err := harbor.PingHarbor(harborConf.URL, harborConf.User, harborConf.Password, false); err != nil {
				resp.Error = err
			} else {
				resp.Msg = "Connected to Harbor"
			}
		}
	case JenkinsType:
		jenkinsConf := &JenkinsConfig{}
		err := json.Unmarshal([]byte(config), jenkinsConf)
		if err != nil {
			log.Log.Error("jenkinsConf conf format error:  %v", err.Error())
			resp.Error = err
			return resp
		}
		log.Log.Debug("verify jenkins conf: %v", jenkinsConf)
		jClient, err := jenkins.NewJenkinsClient(
			jenkins.URL(jenkinsConf.URL),
			jenkins.JenkinsUser(jenkinsConf.User),
			jenkins.JenkinsToken(jenkinsConf.Token),
		)
		if err != nil {
			log.Log.Error("create jenkins client error: %s", err.Error())
			resp.Error = err
			return resp
		}

		pingInfo, err := jClient.Ping()
		if err != nil {
			resp.Error = err
		} else {
			resp.Msg = fmt.Sprintf("Connected to Jenkins %v", pingInfo)
		}
	default:
		resp.Error = fmt.Errorf("no support type: %s integrate setting", request.Type)
	}
	return resp
}

// CreateIntegrateSetting ..
func (pm *SettingManager) CreateIntegrateSetting(request *IntegrateSettingReq, creator string) error {
	// TODO: verify config struct is valid
	config, err := request.String()
	if err != nil {
		log.Log.Error("json marshal error: %s", err.Error())
		return err
	}
	if request.Type == KubernetesType {
		if err := validate.ValidateKubernetesName(request.Name); err != nil {
			return err
		}
	}
	newIntegrateSetting := &models.IntegrateSetting{
		Name:        request.Name,
		Description: request.Description,
		Creator:     creator,
		Type:        request.Type,
		Config:      config,
	}

	if request.Type == KubernetesType {
		kube := &KubeConfig{}
		err := json.Unmarshal([]byte(config), kube)
		if err != nil {
			msg := fmt.Sprintf("kuber conf format error:  %v", err.Error())
			log.Log.Error(msg)
			return fmt.Errorf(msg)
		}

		if err := pm.createOrupateKubernetesConfig(request.Name, kube.Conf); err != nil {
			log.Log.Error("create or update k8s config file error: %s", err.Error())
		} else {
			log.Log.Debug("create or update k8s config file success.")
		}
	}
	return pm.model.CreateIntegrateSetting(newIntegrateSetting)
}

// DeleteIntegrateSetting ..
func (pm *SettingManager) DeleteIntegrateSetting(integrateID int64) error {
	// TODO: verify integrateID is referenced by project env or not.
	return pm.model.DeleteIntegrateSetting(integrateID)
}

func formatIntegrateSettingResponse(items []*models.IntegrateSetting) []*IntegrateSettingResponse {
	rsp := []*IntegrateSettingResponse{}
	config := &Config{}
	for _, item := range items {
		rsp = append(rsp, formatSignalIntegrateSetting(item, config))
	}
	return rsp
}

func formatSignalIntegrateSetting(item *models.IntegrateSetting, config *Config) *IntegrateSettingResponse {
	configJSON, err := config.Struct(item.Config, item.Type)
	if err != nil {
		log.Log.Error("parse config error: %s", err.Error())
	}
	return &IntegrateSettingResponse{
		Creator:  item.Creator,
		UpdateAt: &item.UpdateAt,
		CreateAt: &item.CreateAt,
		ID:       item.ID,
		IntegrateSettingReq: IntegrateSettingReq{
			Name:        item.Name,
			Description: item.Description,
			Type:        item.Type,
			Config:      configJSON,
		},
	}
}
