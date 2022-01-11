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

package kuberes

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-atomci/atomci/constant"
	"github.com/go-atomci/atomci/internal/core/settings"
	"github.com/go-atomci/atomci/internal/dao"
	"github.com/go-atomci/atomci/internal/middleware/log"
	"github.com/go-atomci/atomci/internal/models"
	"github.com/go-atomci/atomci/pkg/kube"

	"github.com/go-atomci/atomci/utils/errors"

	"github.com/go-atomci/atomci/utils/query"
	"github.com/go-atomci/atomci/utils/validate"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"

	"github.com/astaxie/beego/orm"
)

type ContainerParam struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type RollingUpdateApp struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type AppParam struct {
	Name       string              `json:"name"` //appname
	Containers []ContainerParam    `json:"containers,omitempty"`
	Replicas   *intstr.IntOrString `json:"replicas,omitempty"`
}

type VersionWeight struct {
	Stage   string             `json:"stage"`
	Version string             `json:"version"`
	Weight  intstr.IntOrString `json:"weight"`
}

type AppItem struct {
	models.CaasApplication
	ReplicasConstrast string `json:"replicas_constrast,omitempty"`
	Status            string `json:"status,omitempty"`
	Pods              string `json:"pods,omitempty"`
	CreateAt          string `json:"create_at,omitempty"`
	UpdateAt          string `json:"update_at,omitempty"`
}

type AppPod struct {
	Pod    `json:",inline"`
	Weight int `json:"weight"`
}

type AppDetail struct {
	AppItem
	Services []*ServiceDetail `json:"services,omitempty"`
	Pods     []*AppPod        `json:"pods,omitempty"`
}

const (
	AppKindDaemonSet     = "daemonset"
	AppKindDeployment    = "deployment"
	AppKindStatefulSet   = "statefulset"
	LABLE_APPNAME_KEY    = "app"
	LABLE_APPVERSION_KEY = "version"

	ServiceKind              = "service"
	ConfigMapKind            = "configmap"
	SecretKind               = "secret"
	DescriptionAnnotationKey = "description"
	OwnerNameAnnotationKey   = "owner_name"
	DEFAULT_PROJECT_ID       = 0
	YamlSeparator            = "---\n"
	IngApiVersion            = "extensions/v1beta1"
)

type PatcherFunction func(app models.CaasApplication)
type NamespaceListFunction func() []string
type ResType string

const (
	ResTypeApp      ResType = "app"
	ResTypePod      ResType = "pod"
	ResTypeDeploy   ResType = "deploy"
	ResTypeTemplate ResType = "template"
	ResTypeImage    ResType = "image"
)

type AppRes struct {
	Cluster   string
	EnvID     int64
	ProjectID int64
	Client    kubernetes.Interface
	Appmodel  *dao.AppModel
}

type AppPodBasicParam struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	PodName   string `json:"pod_name"`
}

type AppPodDetail struct {
	models.CaasApplication `json:",inline"`
	Pod                    AppPod `json:"pod"`
}

func NewAppRes(cluster string, envID, projectID int64) (*AppRes, error) {
	if cluster == "" {
		return &AppRes{
			Cluster:   cluster,
			EnvID:     envID,
			Appmodel:  dao.NewAppModel(),
			ProjectID: projectID,
		}, nil
	}
	client, err := kube.GetClientset(cluster)
	if err != nil {
		if cluster != "" {
			return nil, errors.NewInternalServerError().SetCause(err)
		}
	}

	return &AppRes{
		Cluster:   cluster,
		EnvID:     envID,
		Appmodel:  dao.NewAppModel(),
		ProjectID: projectID,
		Client:    client,
	}, nil
}

// GetAppListByPagination ..
func (ar *AppRes) GetAppListByPagination(filterQuery *query.FilterQuery, projectID int64, cluster string) (*query.QueryResult, error) {
	appList := []AppItem{}

	res, err := ar.Appmodel.GetAppList(filterQuery, projectID, cluster, "")
	if err != nil {
		return nil, err
	}
	list, ok := res.Item.([]models.CaasApplication)
	if !ok {
		return nil, fmt.Errorf("data type is not right! ")
	}
	for _, item := range list {
		aitem := AppItem{}
		aitem.CaasApplication = item
		aitem.CreateAt = item.CreateAt.Format("2006-01-02 15:04:05")
		aitem.UpdateAt = item.UpdateAt.Format("2006-01-02 15:04:05")

		deploymentName := item.Name
		pods, status, _ := ar.GetDeployRuntime(item, deploymentName)
		aitem.Pods = pods
		aitem.Status = status
		appList = append(appList, aitem)
	}
	res.Item = appList

	return res, nil
}

func (ar *AppRes) GetDeployRuntime(app models.CaasApplication, deploymentName string) (string, string, error) {
	// TODO: current only support deployment
	v1Deployment, err := ar.Client.AppsV1().Deployments(app.Namespace).Get(deploymentName, metav1.GetOptions{})
	if err != nil {
		log.Log.Warn("get deployment error: %s", err.Error())
		return "", "", err
	}

	readyReplicas := v1Deployment.Status.ReadyReplicas
	replicas := v1Deployment.Status.Replicas
	pods := fmt.Sprintf("%v / %v", readyReplicas, replicas)
	status := "NotReady"
	if readyReplicas == replicas {
		status = "Running"
	} else if readyReplicas != 0 {
		status = "Warning"
	}
	return pods, status, nil
}

func (ar *AppRes) GetAppDetail(namespace, name string) (*AppDetail, error) {
	app, err := ar.Appmodel.GetAppByName(ar.Cluster, namespace, name)
	if err != nil {
		return nil, err
	}

	detail := AppDetail{}
	detail.CaasApplication = *app

	deploymentName := app.Name

	detail.ReplicasConstrast, detail.Status, err = ar.GetDeployRuntime(*app, deploymentName)
	if err != nil {
		return nil, err
	}
	detail.CreateAt = app.CreateAt.Format("2006-01-02 15:04:05")
	detail.UpdateAt = app.UpdateAt.Format("2006-01-02 15:04:05")

	// Pods
	detail.Pods, err = ar.getAppPodList(app, deploymentName)
	if err != nil {
		log.Log.Error("Get Pods information failed: %s", err.Error())
		return nil, err
	}

	nativeAppTemplate := NativeAppTemplate{}
	err = json.Unmarshal([]byte(app.Template), &nativeAppTemplate)
	if err != nil {
		log.Log.Error("app template json unmarshal error: %s", err.Error())
	}
	appDeploymentName := nativeAppTemplate.Deployment.Name
	if svc, err := ar.GetAppServiceDetail(namespace, appDeploymentName, getBestNodeIP(detail.Pods)); err != nil {
		log.Log.Warn("get service detail failed: %s", err.Error())
	} else {
		log.Log.Debug("app: %v's svc len: %v", app.Name, len(svc))
		detail.Services = svc
	}
	return &detail, nil
}
func (ar *AppRes) GetAppServiceDetail(namespace, appDeploymentName, nodeIP string) ([]*ServiceDetail, error) {
	return GetAppServices(ar.Client, ar.Cluster, namespace, appDeploymentName, nodeIP)
}

func (ar *AppRes) GetAppPodStatus(namespace, appName, podName string) (interface{}, error) {
	_, err := ar.Appmodel.GetAppByName(ar.Cluster, namespace, appName)
	if err != nil {
		return "", err
	}
	pod, err := ar.Client.CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return pod.Status, nil
}

func (ar *AppRes) GetAppPodLog(namespace, appName, podName, containerName string) (string, error) {
	_, err := ar.Appmodel.GetAppByName(ar.Cluster, namespace, appName)
	if err != nil {
		return "", err
	}

	tailLines := int64(1000)
	body, err := ar.Client.CoreV1().Pods(namespace).GetLogs(podName, &apiv1.PodLogOptions{
		Container: containerName,
		TailLines: &tailLines,
	}).Do().Raw()

	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (ar *AppRes) InstallApp(
	namespace, tname string,
	template Template,
	eparam *ExtensionParam) error {
	CreateK8sNamespace(ar.Cluster, namespace)
	CreateRegistrySecret(ar.Cluster, namespace, ar.EnvID)
	if err := template.Validate(); err != nil {
		return errors.NewBadRequest().SetCause(err)
	}
	if err := template.Default(ar.EnvID).Deploy(ar.ProjectID, ar.EnvID, ar.Cluster, namespace, tname, eparam); err != nil {
		return errors.NewInternalServerError().SetCause(err)
	}
	return nil
}

func (ar *AppRes) UninstallApp(app models.CaasApplication) error {
	if app.Template == "" {
		return nil
	}
	template, err := CreateAppTemplateByApp(app)
	if err != nil {
		return err
	}
	kr := NewKubeAppRes(ar.Client, ar.Cluster, app.Namespace, app.Kind)

	return kr.DeleteAppResource(template)
}

func (ar *AppRes) DeleteApp(namespace, appname string) error {
	app, err := ar.Appmodel.GetAppByName(ar.Cluster, namespace, appname)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil
		} else {
			return err
		}
	}
	err = ar.UninstallApp(*app)
	if err != nil {
		return err
	}
	err = ar.Appmodel.DeleteApp(*app)
	return err
}

func (ar *AppRes) Restart(namespace, appname string) error {
	app, err := ar.Appmodel.GetAppByName(ar.Cluster, namespace, appname)
	if err != nil {
		if err == orm.ErrNoRows {
			log.Log.Warn("%s, %s, %s is not existed", ar.Cluster, namespace, appname)
			return nil
		}
		return err
	}
	// TODO: refactor
	// template, err := CreateAppTemplateByApp(*app)
	// if err != nil {
	// return err
	// }
	return NewKubeAppRes(ar.Client, ar.Cluster, namespace, app.Kind).Restart(appname)
}

func (ar *AppRes) ReconfigureApp(app models.CaasApplication, template AppTemplate) (*AppDetail, error) {
	kr := NewKubeAppRes(ar.Client, ar.Cluster, app.Namespace, app.Kind)
	exist, err := kr.CheckAppIsExisted(app.Name)
	if err != nil {
		return nil, errors.NewInternalServerError().SetCause(err)
	}
	if exist {
		// TODO: create app template by app bug
		oldTpl, err := CreateAppTemplateByApp(app)
		if err != nil {
			return nil, errors.NewInternalServerError().SetCause(err)
		}
		//update
		err = kr.UpdateAppResource(&app, template, oldTpl, true)
		if err != nil {
			return nil, errors.NewInternalServerError().SetCause(err)
		}
		log.Log.Warn("the app is reconfigured, cluster: %s, namespace: %s, appname: %s", ar.Cluster, app.Namespace, app.Name)
	} else {
		if err := template.UpdateAppObject(&app); err != nil {
			return nil, errors.NewBadRequest().SetCause(err)
		}
		//recreate
		err = kr.CreateAppResource(template)
		if err != nil {
			return nil, errors.NewInternalServerError().SetCause(err)
		}
		log.Log.Warn("the app is recreated, cluster: %s, namespace: %s, appname: %s", ar.Cluster, app.Namespace, app.Name)
	}
	// update app info
	err = ar.Appmodel.UpdateApp(&app, true)
	if err != nil {
		return nil, errors.NewInternalServerError().SetCause(err)
	}
	appDetail, err := ar.GetAppDetail(app.Namespace, app.Name)
	// TODO: bug!!!  GetAppDetail error
	if err != nil {
		return nil, errors.NewInternalServerError().SetCause(err)
	}

	return appDetail, nil
}

func (ar *AppRes) RollingUpdateApp(namespace, appname string, param []ContainerParam) error {
	app, err := ar.Appmodel.GetAppByName(ar.Cluster, namespace, appname)
	if err != nil {
		if err == orm.ErrNoRows {
			return errors.NewNotFound().SetCause(err)
		} else if err == orm.ErrMultiRows {
			return errors.NewConflict().SetCause(err)
		} else {
			return errors.NewInternalServerError().SetCause(err)
		}
	}
	template, err := CreateAppTemplateByApp(*app)
	if err != nil {
		return errors.NewInternalServerError().SetCause(err)
	}
	kr := NewKubeAppRes(ar.Client, ar.Cluster, namespace, app.Kind)
	if err = kr.UpdateAppResource(app, template.Image(param), nil, false); err != nil {
		return errors.NewInternalServerError().SetCause(err)
	}
	log.Log.Debug(fmt.Sprintf("new image for %s/%s/%s is %s!", ar.Cluster, namespace, appname, app.Image))
	if err = ar.Appmodel.UpdateApp(app, true); err != nil {
		return errors.NewInternalServerError().SetCause(err)
	}

	return nil
}

func (ar *AppRes) ScaleApp(namespace, appname string, replicas int) error {
	item, err := ar.Appmodel.GetAppByName(ar.Cluster, namespace, appname)
	if err != nil {
		return err
	}
	template, err := CreateAppTemplateByApp(*item)
	if err != nil {
		return err
	}
	// TODO: need refactor
	kr := NewKubeAppRes(ar.Client, ar.Cluster, namespace, item.Kind)
	if err := kr.Scale(item.Name, replicas); err != nil {
		return err
	}
	tplStr, err := template.Replicas(replicas).String()
	if err != nil {
		return err
	}
	item.Replicas = replicas
	item.Template = tplStr
	return ar.Appmodel.UpdateApp(item, true)
}

func (ar *AppRes) SetDeployStatus(namespace, appname, status string) error {
	return ar.Appmodel.SetDeployStatus(ar.Cluster, namespace, appname, status)
}

func (ar *AppRes) getAppPodList(app *models.CaasApplication, deploymentName string) ([]*AppPod, error) {
	podList, err := GetPods(ar.Client, ar.Cluster, app.Namespace, deploymentName, app.Replicas)
	if err != nil {
		log.Log.Error("Get Pods information failed: " + err.Error())
		return nil, err
	}
	appPodList := []*AppPod{}
	for _, item := range podList {
		pod := AppPod{
			Weight: 0,
		}
		pod.Pod = *item
		averWeight := models.DEFAULT_WEIGHT
		if pod.Status == string(apiv1.PodRunning) {
			pod.Weight = averWeight
		}
		appPodList = append(appPodList, &pod)
	}
	return appPodList, nil
}

func (ar *AppRes) SetLabels(namespace, name string, labels map[string]string) error {
	app, err := ar.Appmodel.GetAppByName(ar.Cluster, namespace, name)
	if err != nil {
		if err == orm.ErrNoRows {
			log.Log.Warn(fmt.Sprintf("application(%s/%s/%s) is not existed!", ar.Cluster, namespace, name))
			return nil
		}
		return err
	}
	if err := validate.ValidateLabels(constant.K8S_RESOURCE_TYPE_APP, labels); err != nil {
		return err
	}
	labelStr, err := json.Marshal(labels)
	if err != nil {
		return err
	}
	if string(labelStr) != app.Labels {
		return ar.Appmodel.SetLabels(ar.Cluster, namespace, name, string(labelStr))
	}
	return nil
}

func CreateK8sNamespace(cluster, namespace string) error {
	client, err := kube.GetClientset(cluster)
	if err != nil {
		return err
	}
	_, err = client.CoreV1().Namespaces().Create(&apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
			Labels: map[string]string{
				"name": namespace,
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func CreateRegistrySecret(cluster, namespace string, envID int64) error {
	client, err := kube.GetClientset(cluster)
	if err != nil {
		log.Log.Warning(fmt.Sprintf("create registry secret failed: %v", err.Error()))
		return err
	}
	// TODO: refactor code combine
	projectEnv, err := dao.NewProjectModel().GetProjectEnvByID(envID)
	if err != nil {
		log.Log.Error("when create registry secret get project env by id: %v, error: %s", envID, err.Error())
		return err
	}
	integrateSettingRegistry, err := settings.NewSettingManager().GetIntegrateSettingByID(projectEnv.Registry)
	if err != nil {
		log.Log.Error("when create registry secret get integrate setting by id: %v, error: %s", projectEnv.Registry, err.Error())
		return err
	}

	var registryAddr, registryUser, registryPassword, registryAuth string
	registryName := strings.ToLower(integrateSettingRegistry.Name)
	if registryConf, ok := integrateSettingRegistry.Config.(*settings.RegistryConfig); ok {
		registryAddr = registryConf.URL
		registryPassword = registryConf.Password
		registryUser = registryConf.User
		registryAuth = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", registryConf.User, registryConf.Password)))
	} else {
		log.Log.Error("parse integrate setting registry config error")
		return fmt.Errorf("parse integrate setting registry config error")
	}

	registrySecretName := fmt.Sprintf("registry-%v", registryName)
	registryInfo := make(map[string]interface{})
	registryInfo[registryAddr] = map[string]string{
		"username": registryUser,
		"password": registryPassword,
		"auth":     registryAuth,
	}
	auth, _ := json.Marshal(registryInfo)
	registrySec, err := client.CoreV1().Secrets(namespace).Get(registrySecretName, metav1.GetOptions{})
	if err != nil {
		if !k8serrors.IsNotFound(err) {
			return err
		}
		_, err = client.CoreV1().Secrets(namespace).Create(&apiv1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
				Name:      registrySecretName,
			},
			Type: apiv1.SecretTypeDockercfg,
			Data: map[string][]byte{
				".dockercfg": auth,
			},
		})
	} else {
		if string(registrySec.Data[".dockercfg"]) == string(auth) {
			return nil
		}
		registrySec.Data = map[string][]byte{".dockercfg": auth}
		_, err = client.CoreV1().Secrets(namespace).Update(registrySec)
	}
	if err != nil {
		log.Log.Warning(fmt.Sprintf("set registry secret failed: %v", err.Error()))
	}
	return err
}

func getKubeResNumber(res string) (int64, error) {
	bind := map[string]int64{"ki": 1 / (2 ^ 10), "mi": 1, "gi": (2 ^ 10), "ti": (2 ^ 20), "pi": (2 ^ 30), "ei": (2 ^ 40)}
	ints := map[string]int64{"k": 1 / (10 ^ 3), "m": 1, "g": (10 ^ 3), "t": (10 ^ 6), "p": (10 ^ 9), "e": (10 ^ 12)}
	//default g
	dest := strings.TrimSpace(strings.ToLower(res))
	for key, value := range bind {
		if strings.HasSuffix(dest, key) {
			nb, err := strconv.Atoi((strings.TrimRight(dest, key)))
			if err != nil {
				return 0, err
			}
			return int64(nb) * value, nil
		}
	}
	for key, value := range ints {
		if strings.HasSuffix(dest, key) {
			nb, err := strconv.Atoi((strings.TrimRight(dest, key)))
			if err != nil {
				return 0, err
			}
			return int64(nb) * value, nil
		}
	}
	nb, err := strconv.Atoi(dest)
	if err != nil {
		return 0, err
	}
	return int64(nb) * (10 ^ 3), nil
}

func getBestNodeIP(pods []*AppPod) string {
	if len(pods) > 0 {
		return pods[0].NodeIP
	}
	return ""
}

type AppEvent struct {
	EventLevel   string `json:"event_level"`
	EventObject  string `json:"event_object"`
	EventType    string `json:"event_type"`
	EventMessage string `json:"event_message"`
	EventTime    string `json:"event_time"`
}

func (ar *AppRes) GetAppEvent(namespace, appName string) ([]AppEvent, error) {
	appEvents := []AppEvent{}
	appResourceName := appName
	// TODO: current only support deployment
	eventList, err := GetEventList(ar.Client, ar.Cluster, namespace, appResourceName)
	if err != nil {
		return nil, err
	}

	for _, ievent := range eventList {
		appEvents = append(appEvents, AppEvent{
			EventLevel:   ievent.EventType,
			EventObject:  ievent.ObjectName,
			EventType:    ievent.Reason,
			EventMessage: ievent.Message,
			EventTime:    ievent.LastTimestamp.Format("2006-01-02 15:04:05"),
		})
	}

	return appEvents, nil
}
