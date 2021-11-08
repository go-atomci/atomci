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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-atomci/atomci/utils"

	"github.com/go-atomci/atomci/constant"
	"github.com/go-atomci/atomci/models"

	"github.com/astaxie/beego/logs"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// native app template and api
type NativeAppTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Deployment        *appsv1.Deployment    `json:"deployment,omitempty"`
	Services          []*apiv1.Service      `json:"services,omitempty"`
	Ingresses         []*extensions.Ingress `json:"ingresses,omitempty"`
	Config            DeployConfig          `json:"config"`
}

//context is file context of native template,
func CreateNativeAppTemplate(app models.CaasApplication) (AppTemplate, error) {
	native := &NativeAppTemplate{}
	if err := json.Unmarshal([]byte(app.Template), native); err != nil {
		return nil, err
	}
	return native, nil
}

// TODO: refactor
func (tp *NativeAppTemplate) GenerateAppObject(cluster, namespace, tplname string, projectID int64) (*models.CaasApplication, error) {
	app := &models.CaasApplication{
		Cluster:       cluster,
		Namespace:     namespace,
		Name:          tp.GetAppName(),
		Kind:          tp.GetAppKind(),
		TemplateName:  tplname,
		LabelSelector: tp.getAppLabelSelector(),
		ProjectID:     projectID,
		Addons: models.Addons{
			ID: 0,
		},
	}

	bItem, err := tp.String()
	if err != nil {
		return nil, err
	}
	mainImage := ""
	var containers []apiv1.Container
	replicas := int32(default_replicas)
	if tp.GetAppKind() == AppKindDeployment {
		containers = tp.Deployment.Spec.Template.Spec.Containers
		if tp.Deployment.Spec.Replicas != nil {
			replicas = *tp.Deployment.Spec.Replicas
		}
	}
	if len(containers) > 0 {
		mainImage = containers[0].Image
	}

	for _, container := range containers {
		if container.Name == app.Name {
			// main container
			mainImage = container.Image
			break
		}
	}
	app.Image = mainImage
	app.Replicas = int(replicas)
	app.Template = bItem
	return app, nil
}

func (tp *NativeAppTemplate) UpdateAppObject(app *models.CaasApplication) error {
	newapp, err := tp.GenerateAppObject(app.Cluster, app.Namespace, app.TemplateName, app.ProjectID)
	if err != nil {
		return err
	}
	inputspec, err := tp.String()
	if err != nil {
		return err
	}
	//check image
	if newapp.Image == "" {
		return fmt.Errorf("the application configure has no image")
	}
	//check
	if !(newapp.Replicas >= constant.ReplicasMin && newapp.Replicas <= constant.ReplicasMax) {
		return fmt.Errorf("replcas is not right, its valid range is [%v, %v]", constant.ReplicasMin, constant.ReplicasMax)
	}
	if app.Kind != tp.GetAppKind() {
		return fmt.Errorf("the kind of application can not be changed")
	}
	app.Template = inputspec
	app.Replicas = newapp.Replicas
	app.Image = newapp.Image
	return nil
}

func (tp *NativeAppTemplate) GenerateKubeObject(cluster, namespace string) (map[string]interface{}, error) {
	// translate template to kubernetes resource objects
	objs := make(map[string]interface{})
	switch tp.GetAppKind() {
	case AppKindDeployment:
		deploy := &appsv1.Deployment{
			TypeMeta:   tp.Deployment.TypeMeta,
			ObjectMeta: tp.Deployment.ObjectMeta,
			Spec:       tp.Deployment.Spec,
		}
		// TODO: fix....
		deploy.Name = genAppName(deploy.Name)
		deploy.Spec.Template = tp.newPodTemplateSpec(deploy.Spec.Template, "")
		deploy.ObjectMeta = tp.newAppObjectMeta(tp.Deployment.ObjectMeta,
			deploy.Spec.Template.Labels,
			namespace,
			deploy.Name,
			"")
		deploy.Spec.Selector = tp.newAppSelector(deploy.Spec.Selector, deploy.Spec.Template)
		objs[AppKindDeployment] = deploy
	default:
		logs.Warn("cant support this application kind:", tp.GetAppKind())
		return nil, fmt.Errorf("cant support this application kind: %s", tp.GetAppKind())
	}
	svcList := []*apiv1.Service{}
	for _, svc := range tp.Services {
		if err := NewKubeSvcValidator(cluster, namespace, tp.GetAppName()).Validator(svc); err != nil {
			return nil, err
		}
		svc.ObjectMeta = tp.newObjectMeta(svc.ObjectMeta, svc.Labels, namespace, svc.Name)
		svcList = append(svcList, svc)
	}
	if len(svcList) > 0 {
		objs[ServiceKind] = svcList
	}
	// TODO: clean ingress created
	return objs, nil
}

func (tp *NativeAppTemplate) GetAppName() string {
	return tp.Name
}

func (tp *NativeAppTemplate) GetAppKind() string {
	return strings.ToLower(tp.Kind)
}

func (tp *NativeAppTemplate) String() (string, error) {
	ctx, err := json.Marshal(tp)
	if err != nil {
		return "", err
	}

	return string(ctx), nil
}

func (tp *NativeAppTemplate) Image(param []ContainerParam) AppTemplate {
	for _, item := range param {
		podSpec := &apiv1.PodTemplateSpec{}
		switch tp.GetAppKind() {
		case AppKindDeployment:
			podSpec = &tp.Deployment.Spec.Template
		}
		for index, ctn := range podSpec.Spec.Containers {
			if item.Name == ctn.Name {
				podSpec.Spec.Containers[index].Image = item.Image
				break
			}
		}
	}
	return tp
}

func (tp *NativeAppTemplate) Replicas(replicas int) AppTemplate {
	num := int32(replicas)
	switch tp.GetAppKind() {
	case AppKindDeployment:
		tp.Deployment.Spec.Replicas = &num
	}
	return tp
}

func (tp *NativeAppTemplate) DefaultLabel() AppTemplate {
	return tp
}

func (tp *NativeAppTemplate) newPodTemplateSpec(spec apiv1.PodTemplateSpec, podVersion string) apiv1.PodTemplateSpec {
	spec.ObjectMeta = tp.newAppObjectMeta(spec.ObjectMeta,
		spec.Labels,
		spec.ObjectMeta.Namespace,
		spec.ObjectMeta.Name,
		podVersion)

	if tp.Config.ImagePullSecret != "" {
		podContainers := []ContainerItem{}
		for _, item := range spec.Spec.Containers {
			podContainers = append(podContainers, ContainerItem{
				Image: item.Image,
				Name:  item.Name,
			})
		}

		for _, item := range spec.Spec.InitContainers {
			podContainers = append(podContainers, ContainerItem{
				Name:  item.Name,
				Image: item.Image,
			})
		}
		IncludeEnvImage := false
		for _, item := range podContainers {
			if strings.Contains(item.Image, tp.Config.HarborAddr) {
				IncludeEnvImage = true
				break
			}
		}
		if IncludeEnvImage {
			spec.Spec.ImagePullSecrets = []apiv1.LocalObjectReference{{Name: tp.Config.ImagePullSecret}}
		}
	}
	return spec
}

func (tp *NativeAppTemplate) newAppSelector(old *metav1.LabelSelector, podTemplate apiv1.PodTemplateSpec) *metav1.LabelSelector {
	selector := old
	if selector == nil {
		selector = &metav1.LabelSelector{}
	}
	if selector.MatchLabels == nil {
		selector.MatchLabels = podTemplate.Labels
	}
	return selector
}

func (tp *NativeAppTemplate) newObjectMeta(old metav1.ObjectMeta, defLabel map[string]string, namespace, name string) metav1.ObjectMeta {
	meta := old
	meta.Name = name
	meta.Namespace = namespace
	if meta.Annotations == nil {
		meta.Annotations = make(map[string]string)
	}
	// set ann
	meta.Annotations[OwnerNameAnnotationKey] = tp.GetAppName()

	if meta.Labels == nil {
		meta.Labels = make(map[string]string)
	}

	//set labelif
	for k, v := range defLabel {
		meta.Labels[k] = v
	}
	return meta
}

func (tp *NativeAppTemplate) newAppObjectMeta(old metav1.ObjectMeta, defLabel map[string]string, namespace, name, podVersion string) metav1.ObjectMeta {
	meta := tp.newObjectMeta(old, defLabel, namespace, name)
	return meta
}

func (tp *NativeAppTemplate) getAppLabelSelector() string {
	labels := make(map[string]string)
	for k, v := range tp.getAppPodLabel() {
		labels[k] = v
	}
	selector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: labels,
	})
	if err != nil {
		logs.Warn("get application label failed:", err)
		return ""
	}
	return selector.String()
}

func (tp *NativeAppTemplate) getAppPodLabel() map[string]string {
	return tp.Deployment.Spec.Template.Labels
}

//generate ingress object
func (tp *NativeAppTemplate) genDefaultIngressObjects(namespace, appname, domainSuffix string) []*extensions.Ingress {
	// TODO: generate default ingress objects
	return tp.Ingresses
}

func ingressRuleIsExisted(defIng *extensions.Ingress, ings []*extensions.Ingress) bool {
	if defIng == nil {
		return false
	}
	// check host
	getSameHostRules := func() (*extensions.IngressRule, *extensions.IngressRule) {
		for _, rule := range defIng.Spec.Rules {
			for _, ing := range ings {
				for _, r := range ing.Spec.Rules {
					if r.Host == rule.Host {
						return &rule, &r
					}
				}
			}
		}
		return nil, nil
	}
	destRule, sameRule := getSameHostRules()
	if destRule == nil || sameRule == nil {
		return false
	}
	if destRule.HTTP == nil || sameRule.HTTP == nil {
		return false
	}
	for _, path1 := range destRule.HTTP.Paths {
		for _, path2 := range destRule.HTTP.Paths {
			if utils.PathsIsEqual(path1.Path, path2.Path) {
				return true
			}
		}
	}
	return false
}

func genAppName(virginName string) string {
	return GenerateDeployName(virginName)
}

func GenDefaultIngressObject(svc *apiv1.Service, ingObjMeta metav1.ObjectMeta, defPort int32, domainName, domainSuffix string) *extensions.Ingress {
	// TODO: refactor
	// if svc == nil {
	// 	return nil
	// }
	// if domainSuffix == "" {
	// 	logs.Warn("domain suffix is empty, so cant generate ingress object!")
	// 	return nil
	// }
	// if len(svc.Spec.Ports) == 0 {
	// 	return nil
	// }
	// ingress := &extensions.Ingress{
	// 	TypeMeta: metav1.TypeMeta{
	// 		Kind:       "Ingress",
	// 		APIVersion: IngApiVersion,
	// 	},
	// }
	// var rules []extensions.IngressRule
	// for _, port := range svc.Spec.Ports {
	// 	dPort := port.Port
	// 	// set ingress rule
	// 	ruleValue := extensions.IngressRuleValue{
	// 		HTTP: &extensions.HTTPIngressRuleValue{
	// 			Paths: []extensions.HTTPIngressPath{
	// 				extensions.HTTPIngressPath{
	// 					Backend: extensions.IngressBackend{
	// 						ServiceName: svc.Name,
	// 						ServicePort: intstr.IntOrString{
	// 							Type:   intstr.Int,
	// 							IntVal: dPort,
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	}
	// 	hostSuffix := ""
	// 	if defPort != port.Port {
	// 		hostSuffix = strconv.Itoa(int(port.Port))
	// 	}
	// 	rules = append(rules, extensions.IngressRule{
	// 		Host:             GenerateIngressHost(domainName, domainSuffix, hostSuffix),
	// 		IngressRuleValue: ruleValue,
	// 	})
	// }
	// ingress.ObjectMeta = ingObjMeta
	// //ingress.ObjectMeta.ResourceVersion = ingrv
	// var specRules []extensions.IngressRule
	// // update rule
	// for _, rule := range rules {
	// 	ir := -1
	// 	for i, item := range ingress.Spec.Rules {
	// 		if item.Host == rule.Host {
	// 			ir = i
	// 			break
	// 		}
	// 	}
	// 	if ir == -1 {
	// 		// just append a rule
	// 		ingress.Spec.Rules = append(ingress.Spec.Rules, rule)
	// 	} else {
	// 		// update rule: add a path or update path's backend
	// 		// the method may cause some garbage paths
	// 		drule := &ingress.Spec.Rules[ir]
	// 		for _, spath := range rule.HTTP.Paths {
	// 			ipath := -1
	// 			for i, dpath := range drule.HTTP.Paths {
	// 				if dpath.Path == spath.Path {
	// 					ipath = i
	// 				}
	// 			}
	// 			if ipath == -1 {
	// 				// add a path
	// 				drule.HTTP.Paths = append(drule.HTTP.Paths, spath)
	// 			} else {
	// 				// update path's backend
	// 				drule.HTTP.Paths[ipath].Backend = spath.Backend
	// 			}
	// 		}
	// 	}
	// }
	// // delete ingress rule which service port is not existed in service
	// for _, rule := range ingress.Spec.Rules {
	// 	if rule.HTTP != nil {
	// 		for _, path := range rule.HTTP.Paths {
	// 			found := false
	// 			for _, port := range svc.Spec.Ports {
	// 				if path.Backend.ServicePort.IntValue() == int(port.Port) {
	// 					found = true
	// 					break
	// 				}
	// 			}
	// 			if found {
	// 				specRules = append(specRules, rule)
	// 				break
	// 			}
	// 		}
	// 	} else {
	// 		specRules = append(specRules, rule)
	// 	}
	// }
	// // delete tls host if host is not existed in rule
	// var specTLS []extensions.IngressTLS
	// for _, item := range ingress.Spec.TLS {
	// 	var hosts []string
	// 	for _, host := range item.Hosts {
	// 		existed := false
	// 		for _, rule := range specRules {
	// 			if rule.Host == host {
	// 				existed = true
	// 				break
	// 			}
	// 		}
	// 		if existed {
	// 			hosts = append(hosts, host)
	// 		}
	// 	}
	// 	if len(hosts) != 0 {
	// 		specTLS = append(specTLS, extensions.IngressTLS{
	// 			Hosts:      hosts,
	// 			SecretName: item.SecretName,
	// 		})
	// 	}
	// }

	// // set spec rules
	// ingress.Spec.Rules = specRules
	// ingress.Spec.TLS = specTLS
	// // set default annotation
	// kubeutil.SetCreatedDefaultAnno(ingress)
	// return ingress
	return nil
}
