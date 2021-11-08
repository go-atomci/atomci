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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/go-atomci/atomci/middleware/log"
	"github.com/go-atomci/atomci/utils/validate"

	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	errors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
)

var metaAccessor = meta.NewAccessor()

const (
	default_replicas = 1
	INIT_APPNAME     = ""

	DefaultVersion = "latest"

	APP_PROTOCOL_HTTP = "http"
)

type ResObject struct {
	Namespace string
	Name      string

	Object  runtime.Object
	RawData []byte
}

type ContainerItem struct {
	Name  string `json:"name,omitempty"`
	Image string `json:"image,omitempty"`
}

type AppResourceItem struct {
	Kind string `json:"kind,omitempty"`
	Name string `json:"name,omitempty"`
}

type DeployConfig struct {
	DefaultPort     int32  `json:"default_port"` //默认服务端口
	DeployStrategy  string `json:"deploy_strategy"`
	ImagePullSecret string `json:"image_pull_secret"`
	HarborAddr      string `json:"harbor_addr"`
	Description     string `json:"description"`
}

// native template and api
type NativeTemplate struct {
	Template string       `json:"template"`
	Config   DeployConfig `json:"config"`
}

func NewNativeTemplate() *NativeTemplate {
	return &NativeTemplate{}
}

// set default value for template config
func (t *NativeTemplate) Default(envID int64) Template {
	if t.Config.ImagePullSecret == "" && envID != 0 {
		defSecret, harborURL, err := getDefaultPullSecretAndHarborAddr(envID)
		if err != nil {
			log.Log.Error("get default pull secret failed:", err)
		}
		t.Config.ImagePullSecret = defSecret
		t.Config.HarborAddr = harborURL
	}
	return t
}

func (t *NativeTemplate) Validate() error {
	resObjects, err := t.parser()
	if err != nil {
		return err
	}
	if len(t.Template) > TemplateMaxSize {
		return fmt.Errorf("the length of template can not be above %vKB", TemplateMaxSize>>10)
	}
	var svcList []*apiv1.Service
	podSpecList := make([]apiv1.PodSpec, 0)
	for _, obj := range resObjects {
		kind, err := metaAccessor.Kind(obj.Object)
		if err != nil {
			return err
		}
		switch strings.ToLower(kind) {
		case AppKindDeployment, AppKindStatefulSet:
			deploy := &v1.Deployment{}
			err := json.Unmarshal(obj.RawData, &deploy)
			if err != nil {
				return err
			}
			podSpecList = append(podSpecList, deploy.Spec.Template.Spec)
			continue
		case SecretKind, ConfigMapKind:
			continue
		case ServiceKind:
			svc := &apiv1.Service{}
			if err := json.Unmarshal(obj.RawData, svc); err != nil {
				return err
			}
			if err := validateService(svc); err != nil {
				return err
			}
			svcList = append(svcList, svc)
			continue
		default:
			return fmt.Errorf("the system does not support this resource kind:%s", strings.ToLower(kind))
		}
		// TODO: deploy schedule policy verify, replicaset
	}

	for _, podSpec := range podSpecList {
		for _, c := range podSpec.InitContainers {
			if err := t.validateContainer(c.Name, c, ""); err != nil {
				return err
			}
		}
		for _, c := range podSpec.Containers {
			if err := t.validateContainer(c.Name, c, ""); err != nil {
				return err
			}
		}
		//check vol
		for _, vol := range podSpec.Volumes {
			if strings.TrimSpace(vol.Name) == "" {
				return fmt.Errorf("volume name can not be empty")
			}
			if vol.PersistentVolumeClaim != nil {
				if strings.TrimSpace(vol.PersistentVolumeClaim.ClaimName) == "" {
					return fmt.Errorf("PVC name can not be empty")
				}
			}
		}
	}
	return nil
}

func (t *NativeTemplate) GetContainerImages() ([]*ContainerItem, error) {
	resObjects, err := t.parser()
	if err != nil {
		return nil, err
	}
	if len(t.Template) > TemplateMaxSize {
		return nil, fmt.Errorf("the length of template can not be above %vKB", TemplateMaxSize>>10)
	}
	podSpecList := []apiv1.PodSpec{}
	for _, obj := range resObjects {
		kind, err := metaAccessor.Kind(obj.Object)
		if err != nil {
			return nil, err
		}
		switch strings.ToLower(kind) {
		case AppKindDeployment, AppKindStatefulSet:
			deploy := &v1.Deployment{}
			err := json.Unmarshal(obj.RawData, &deploy)
			if err != nil {
				return nil, err
			}
			podSpecList = append(podSpecList, deploy.Spec.Template.Spec)
			continue
		default:
			log.Log.Debug("ignore the kind of %v 's parse", strings.ToLower(kind))
			continue
		}
	}

	rsp := []*ContainerItem{}
	for _, podSpec := range podSpecList {
		for _, c := range podSpec.InitContainers {
			err := t.validateContainer(c.Name, c, "")
			if err == nil {
				rsp = append(rsp, &ContainerItem{Name: c.Name, Image: c.Image})
			}
		}
		for _, c := range podSpec.Containers {
			err := t.validateContainer(c.Name, c, "")
			if err == nil {
				rsp = append(rsp, &ContainerItem{Name: c.Name, Image: c.Image})
			}
		}
	}

	return rsp, nil
}

// GetAppResourceNames ..
// TODO: current support deployment/statefulset daemonset parse
func (t *NativeTemplate) GetAppResourceNames() ([]AppResourceItem, error) {
	resObjects, err := t.parser()
	if err != nil {
		return nil, err
	}
	if len(t.Template) > TemplateMaxSize {
		return nil, fmt.Errorf("the length of template can not be above %vKB", TemplateMaxSize>>10)
	}
	appResourceItems := []AppResourceItem{}
	for _, obj := range resObjects {
		kind, err := metaAccessor.Kind(obj.Object)
		if err != nil {
			return nil, err
		}
		switch strings.ToLower(kind) {
		case AppKindDeployment, AppKindStatefulSet, AppKindDaemonSet:
			deploy := &v1.Deployment{}
			err := json.Unmarshal(obj.RawData, &deploy)
			if err != nil {
				return nil, err
			}
			appResourceItems = append(appResourceItems, AppResourceItem{
				Name: deploy.Name,
				Kind: deploy.Kind,
			})
		default:
			log.Log.Info("obj name: %v kind: %v was skipped parse", obj.Name, kind)
		}
	}
	return appResourceItems, nil
}

func (t *NativeTemplate) parser() ([]*ResObject, error) {
	infoList := []*ResObject{}
	buffer := bytes.NewBuffer([]byte(t.Template))
	d := yaml.NewYAMLOrJSONDecoder(buffer, 4096)
	for {
		ext := runtime.RawExtension{}
		if err := d.Decode(&ext); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error parsing: %v", err)
		}
		// TODO: This needs to be able to handle object in other encodings and schemas.
		ext.Raw = bytes.TrimSpace(ext.Raw)
		if len(ext.Raw) == 0 || bytes.Equal(ext.Raw, []byte("null")) {
			continue
		}
		//if err := ValidateSchema(ext.Raw, v.Schema); err != nil {
		//	return fmt.Errorf("error validating %q: %v", v.Source, err)
		//}
		info, err := infoForData(ext.Raw)
		if err != nil {
			continue
		}
		infoList = append(infoList, info)
	}
	if len(infoList) == 0 {
		return nil, fmt.Errorf("the template has no resolvable resource objects!")
	}
	return infoList, nil
}

// InfoForData creates an Info object for the given data. An error is returned
// if any of the decoding or client lookup steps fail. Name and namespace will be
// set into Info if the mapping's MetadataAccessor can retrieve them.
func infoForData(data []byte) (*ResObject, error) {
	obj, _, err := unstructured.UnstructuredJSONScheme.Decode(data, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to decode: %v", err)
	}
	name, _ := metaAccessor.Name(obj)
	namespace, _ := metaAccessor.Namespace(obj)
	ret := &ResObject{
		Namespace: namespace,
		Name:      name,
		Object:    obj,
		RawData:   data,
	}

	return ret, nil
}

func (t *NativeTemplate) GenNativeAppTemplate(namespace, appname string) ([]*NativeAppTemplate, []ResObject, error) {
	resObjects, err := t.parser()
	if err != nil {
		return nil, nil, err
	}
	var tplList []*NativeAppTemplate
	var noAppObjList []ResObject
	otherObjList := []ResObject{}
	for _, obj := range resObjects {
		kind, err := metaAccessor.Kind(obj.Object)
		if err != nil {
			return nil, nil, err
		}
		obj.Namespace = namespace
		switch strings.ToLower(kind) {
		case AppKindDeployment:
			extend := t.Config
			deploy := &v1.Deployment{}
			err := json.Unmarshal(obj.RawData, &deploy)
			if err != nil {
				return nil, nil, err
			}
			appmeta := deploy.ObjectMeta
			if appname != INIT_APPNAME {
				appmeta.Name = appname
			}
			deploy.Namespace = namespace
			tplList = append(tplList, &NativeAppTemplate{
				TypeMeta:   deploy.TypeMeta,
				ObjectMeta: appmeta,
				Deployment: deploy,
				Config:     extend,
			})
		default:
			noAppObjList = append(noAppObjList, *obj)
		}
	}
	// match service
	for _, obj := range noAppObjList {
		hasOwner := false
		kind, _ := metaAccessor.Kind(obj.Object)
		if strings.ToLower(kind) == ServiceKind {
			svc := &apiv1.Service{}
			err := json.Unmarshal(obj.RawData, svc)
			if err != nil {
				return nil, nil, err
			}
			svc.Namespace = namespace
			app := getServiceOwner(svc.Spec.Selector, tplList)
			if app != nil {
				hasOwner = true
				app.Services = append(app.Services, svc)
			}
		}
		if !hasOwner {
			otherObjList = append(otherObjList, obj)
		}
	}
	// TODO: add ingress type
	// match ingress
	// noAppObjList = otherObjList
	// otherObjList = []ResObject{}
	// for _, obj := range noAppObjList {
	// 	hasOwner := false
	// 	kind, _ := metaAccessor.Kind(obj.Object)
	// 	if strings.ToLower(kind) == IngressKind {
	// 		ing := &extensions.Ingress{}
	// 		err := json.Unmarshal(obj.RawData, ing)
	// 		if err != nil {
	// 			return nil, nil, err
	// 		}
	// 		ing.Namespace = namespace
	// 		app := getIngressOwner(ing, tplList)
	// 		if app != nil {
	// 			hasOwner = true
	// 			app.Ingresses = append(app.Ingresses, ing)
	// 		}
	// 	}
	// 	if !hasOwner {
	// 		otherObjList = append(otherObjList, obj)
	// 	}
	// }
	return tplList, otherObjList, nil
}

func (t *NativeTemplate) Deploy(projectID, envID int64, cluster, namespace, tname string, eparam *ExtensionParam) error {
	ar, err := NewAppRes(cluster, envID, projectID)
	if err != nil {
		return err
	}
	tplList, otherObjList, err := t.GenNativeAppTemplate(namespace, INIT_APPNAME)
	if err != nil {
		return err
	}
	// create otherObjList
	t.CreateNoAppResource(ar.Client, cluster, namespace, otherObjList)
	var appTplList []AppTemplate
	for _, tpl := range tplList {
		appTplList = append(appTplList, tpl)
	}
	return DeployAppTemplates(appTplList, projectID, envID, cluster, namespace, tname, eparam)
}

func (t *NativeTemplate) GetExample() []byte {
	spec := "apiVersion: apps/v1\n" +
		"kind: Deployment\n" +
		"metadata:\n" +
		"  name: helloworld\n" +
		"  namespace: default\n" +
		"spec:\n" +
		"  template:\n" +
		"    metadata:\n" +
		"      labels:\n" +
		"        app: helloworld\n" +
		"    spec:\n" +
		"      containers:\n" +
		"      - name: helloworld\n" +
		"        image: harbor.example.com/default/helloworld:latest\n" +
		"        imagePullPolicy: Always\n" +
		"      serviceAccountName: default\n"
	return []byte(spec)
}

func (t *NativeTemplate) validateContainer(appname string, c apiv1.Container, appkind string) error {
	if c.Image == "" {
		return fmt.Errorf("the container(%s) has no image", c.Name)
	}
	// TODO: refactor
	// if appkind == AppKindDeployment {
	// 	if c.Resources.Limits.Memory().Value() == 0 ||
	// 		c.Resources.Limits.Cpu().Value() == 0 {
	// 		return fmt.Errorf("the container has no CPU or memory limits!")
	// 	}
	// }
	for _, arg := range c.Command {
		if arg == "" {
			return fmt.Errorf("the commond of container is not right, it can not be empty")
		}
	}
	for _, env := range c.Env {
		if errs := validation.IsEnvVarName(env.Name); len(errs) != 0 {
			log.Log.Error("application env name of container is not right, !", appname, errs)
			return fmt.Errorf("the env name if not right, it must match [-._a-zA-Z][-._a-zA-Z0-9]*")
		}
	}
	validateHealth := func(probe *apiv1.Probe) error {
		if probe == nil {
			return nil
		}
		if probe.HTTPGet == nil && probe.TCPSocket == nil {
			return fmt.Errorf("incomplete configuration of health probe")
		}
		if probe.HTTPGet != nil {
			if err := validate.ValidatePortNum(int32(probe.HTTPGet.Port.IntValue())); err != nil {
				return err
			}
		}
		if probe.TCPSocket != nil {
			if err := validate.ValidatePortNum(int32(probe.TCPSocket.Port.IntValue())); err != nil {
				return err
			}
		}
		if probe.PeriodSeconds < 0 || probe.TimeoutSeconds < 0 ||
			probe.InitialDelaySeconds < 0 ||
			probe.SuccessThreshold < 0 ||
			probe.FailureThreshold < 0 {
			return fmt.Errorf("params of health probe must be equal or above 0")
		}
		return nil
	}
	// check health
	if c.LivenessProbe != nil {
		if err := validateHealth(c.LivenessProbe); err != nil {
			return err
		}
	}
	if c.ReadinessProbe != nil {
		if err := validateHealth(c.ReadinessProbe); err != nil {
			return err
		}
	}
	return nil
}

func (t *NativeTemplate) CreateNoAppResource(client kubernetes.Interface, cluster, namespace string, objs []ResObject) {
	resMap := map[string]kubeResInterface{}
	kr := NewKubeAppRes(client, cluster, namespace, "")
	svcList := &kubeServices{
		kubeAppHandler: kr,
	}
	// TODO: refactor
	// ingList := &kubeIngesses{
	// 	kubeAppHandler: kr,
	// }
	configs := configList{}
	secrets := secretList{}
	for _, obj := range objs {
		kind, err := metaAccessor.Kind(obj.Object)
		if err != nil {
			log.Log.Error("%v:", err)
			continue
		}
		switch strings.ToLower(kind) {
		case ServiceKind:
			svc := &apiv1.Service{}
			err := json.Unmarshal(obj.RawData, svc)
			if err != nil {
				log.Log.Error("%v", err)
				continue
			}
			svc.Namespace = namespace
			svcList.serviceList = append(svcList.serviceList, svc)
			resMap[ServiceKind] = svcList
		// case IngressKind:
		// 	ing := &extensions.Ingress{}
		// 	err := json.Unmarshal(obj.RawData, ing)
		// 	if err != nil {
		// 		log.Log.Error("%v", err)
		// 		continue
		// 	}
		// 	ing.Namespace = namespace
		// 	ingList.ingressList = append(ingList.ingressList, ing)
		// 	resMap[IngressKind] = ingList
		case ConfigMapKind:
			conf := &apiv1.ConfigMap{}
			err := json.Unmarshal(obj.RawData, conf)
			if err != nil {
				log.Log.Error("%v", err)
				continue
			}
			conf.Namespace = namespace
			configs = append(configs, conf)
			resMap[ConfigMapKind] = configs
		case SecretKind:
			sec := &apiv1.Secret{}
			err := json.Unmarshal(obj.RawData, sec)
			if err != nil {
				log.Log.Error("unmarshal virgin data to secret type failed:", err)
				continue
			}
			sec.Namespace = namespace
			secrets = append(secrets, sec)
			resMap[SecretKind] = secrets
		default:
			log.Log.Warn("dont support this resource kind", obj.Object.GetObjectKind())
		}
	}
	for kind, handler := range resMap {
		if err := handler.create(client); err != nil {
			log.Log.Warn("create "+kind+"failed:", err)
		}
	}
}

type kubeResInterface interface {
	create(client kubernetes.Interface) error
}

type kubeServices struct {
	kubeAppHandler *KubeAppRes
	serviceList    []*apiv1.Service
}

type kubeIngesses struct {
	kubeAppHandler *KubeAppRes
	ingressList    []*extensions.Ingress
}

type configList []*apiv1.ConfigMap

type secretList []*apiv1.Secret

func (svcs kubeServices) create(client kubernetes.Interface) error {
	if len(svcs.serviceList) > 0 {
		return svcs.kubeAppHandler.CreateService(svcs.serviceList)
	}
	return nil
}

func (confs configList) create(client kubernetes.Interface) error {
	for _, config := range confs {
		isExist := false
		_, err := client.CoreV1().ConfigMaps(config.Namespace).Get(config.Name, metav1.GetOptions{})
		if err == nil {
			isExist = true
		} else {
			if !errors.IsNotFound(err) {
				log.Log.Warn("config map: %s get occur error: %s", config.Name, err.Error())
			}
		}
		if isExist {
			if _, err := client.CoreV1().ConfigMaps(config.Namespace).Update(config); err != nil {
				log.Log.Warn("update configmap error: %v", err)
			}
		} else {
			if _, err := client.CoreV1().ConfigMaps(config.Namespace).Create(config); err != nil {
				log.Log.Warn("create configmap error: %v", err)
			}
		}
	}
	return nil
}

func (secrets secretList) create(client kubernetes.Interface) error {
	for _, secret := range secrets {
		isExist := false
		oldsecret, err := client.CoreV1().Secrets(secret.Namespace).Get(secret.Name, metav1.GetOptions{})
		if err == nil {
			isExist = true
			secret.ResourceVersion = oldsecret.ResourceVersion
		} else {
			if !errors.IsNotFound(err) {
				log.Log.Warn("get secret name: %v occur error: %s", secret.Name, err.Error())
			}
		}
		if isExist {
			if _, err := client.CoreV1().Secrets(secret.Namespace).Update(secret); err != nil {
				log.Log.Warn("update configmap error: %v", err)
			}
		} else {
			if _, err := client.CoreV1().Secrets(secret.Namespace).Create(secret); err != nil {
				log.Log.Warn("create configmap error: %v", err)
			}
		}
	}
	return nil
}

func (ings kubeIngesses) create(client kubernetes.Interface) error {
	// TODO: refactor
	// if len(ings.ingressList) > 0 {
	// 	return ings.kubeAppHandler.CreateIngress(ings.ingressList, nil, true)
	// }
	return nil
}

func getServiceOwner(selector map[string]string, appTpls []*NativeAppTemplate) *NativeAppTemplate {
	for _, item := range appTpls {
		found := true
		podLabel := item.getAppPodLabel()
		for key, value := range selector {
			if podLabel[key] != value {
				found = false
				break
			}
		}
		if found {
			return item
		}
	}
	return nil
}

func getIngressOwner(ing *extensions.Ingress, apptpls []*NativeAppTemplate) *NativeAppTemplate {
	svcNames := getIngressService(ing)
	if len(svcNames) != 1 {
		return nil
	}
	for _, item := range apptpls {
		found := false
		for _, svc := range item.Services {
			if svcNames[0] == svc.Name {
				found = true
				break
			}
		}
		if found {
			return item
		}
	}
	return nil
}

func getIngressService(ing *extensions.Ingress) []string {
	names := []string{}
	if ing == nil {
		return names
	}
	label := make(map[string]interface{})
	if ing.Spec.Backend != nil {
		label[ing.Spec.Backend.ServiceName] = nil
	}
	for _, rule := range ing.Spec.Rules {
		if rule.HTTP == nil {
			continue
		}
		for _, path := range rule.HTTP.Paths {
			label[path.Backend.ServiceName] = nil
		}
	}
	for name, _ := range label {
		names = append(names, name)
	}
	return names
}

func validateService(svc *apiv1.Service) error {
	if svc == nil {
		return nil
	}
	for _, port := range svc.Spec.Ports {
		err := validate.ValidatePortNum(port.Port)
		if err == nil {
			err = validate.ValidatePortNum(int32(port.TargetPort.IntValue()))
		}
		if err != nil {
			return err
		}
		if apiv1.ServiceType(svc.Spec.Type) == apiv1.ServiceTypeNodePort {
			if err = validate.ValidateNodePortNum(port.NodePort); err != nil {
				return err
			}
		}
	}
	return nil
}
