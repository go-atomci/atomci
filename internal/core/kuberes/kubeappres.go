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
	"fmt"

	"github.com/go-atomci/atomci/internal/middleware/log"
	"github.com/go-atomci/atomci/internal/models"

	// istio_types "istio/api/routing/v1alpha1"

	"github.com/astaxie/beego"
	apiv1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type KubeAppRes struct {
	DomainSuffix  string
	Namespace     string
	cluster       string
	client        kubernetes.Interface
	kubeAppHandle KubeAppInterface
}

type RollBackFunc func() error

func NewKubeAppRes(client kubernetes.Interface, cluster, namespace, kind string) *KubeAppRes {
	res := &KubeAppRes{
		cluster:   cluster,
		Namespace: namespace,
		client:    client,
	}
	if kind == AppKindStatefulSet {
		res.kubeAppHandle = NewStatefulRes(client, namespace)
	} else {
		//default is deployment
		res.kubeAppHandle = NewDeploymentRes(client, namespace)
	}
	return res
}

func (kr *KubeAppRes) CreateAppResource(template AppTemplate) error {
	// rollback if create resource failed
	rollbackFuncList := []RollBackFunc{}
	defer func() {
		for _, rollback := range rollbackFuncList {
			if err := rollback(); err != nil {
				log.Log.Error("%v", err)
			}
		}
	}()
	objMap, err := template.GenerateKubeObject(kr.cluster, kr.Namespace)
	if err != nil {
		return err
	}
	// create svc
	objs, existed := objMap[ServiceKind]
	if existed {
		svcList, ok := objs.([]*apiv1.Service)
		if !ok {
			return fmt.Errorf("service object list is not right")
		}
		// create or update
		if err := kr.CreateService(svcList); err != nil {
			return err
		}
		rollbackFuncList = append(rollbackFuncList, func() error { return kr.deleteService(svcList) })
	}
	// TODO: ingress resources disable
	// create app
	if obj, existed := objMap[template.GetAppKind()]; existed {
		if err := NewKubeAppValidator(kr.client, kr.Namespace).Validator(obj); err != nil {
			return err
		}
		// create or update
		if err := kr.kubeAppHandle.Create(obj); err != nil {
			return err
		}
	}
	rollbackFuncList = nil
	return nil
}

func (kr *KubeAppRes) UpdateAppResource(app *models.CaasApplication, new, old AppTemplate, all bool) error {
	err := new.UpdateAppObject(app)
	if err != nil {
		return err
	}
	objMap, err := new.GenerateKubeObject(app.Cluster, app.Namespace)
	if err != nil {
		return err
	}
	oldMap := make(map[string]interface{})
	if old != nil {
		oldMap, err = old.GenerateKubeObject(app.Cluster, app.Namespace)
		if err != nil {
			return err
		}
	}
	if obj, existed := objMap[new.GetAppKind()]; existed {
		if err := NewKubeAppValidator(kr.client, kr.Namespace).Validator(obj); err != nil {
			return err
		}
		// create or update
		if _, err = kr.kubeAppHandle.Update(*app, obj); err != nil {
			return err
		}
	}
	if !all {
		return nil
	}
	if err := kr.updateSvcResource(oldMap[ServiceKind], objMap[ServiceKind]); err != nil {
		beego.Warn("update service resource failed:", err)
	}

	return nil
}

func (kr *KubeAppRes) DeleteAppResource(template AppTemplate) error {
	objMap, err := template.GenerateKubeObject(kr.cluster, kr.Namespace)
	if err != nil && objMap == nil {
		return err
	}
	objs, existed := objMap[ServiceKind]
	if existed {
		svcList, ok := objs.([]*apiv1.Service)
		if !ok {
			return fmt.Errorf("service object list is not right")
		}
		// delete
		if err := kr.deleteService(svcList); err != nil {
			return fmt.Errorf("delete service error %v", err)
		}
	}
	// TODO: check app ingress resource
	return kr.kubeAppHandle.Delete(template.GetAppName())
}

func (kr *KubeAppRes) updateSvcResource(oldObjList, newObjList interface{}) error {
	delSvcs := []*apiv1.Service{}
	newSvcs := []*apiv1.Service{}
	ok := false
	if newObjList != nil {
		newSvcs, ok = newObjList.([]*apiv1.Service)
		if !ok {
			return fmt.Errorf("service object list is not right")
		}
		// create or update
		if err := kr.CreateService(newSvcs); err != nil {
			return err
		}
	}
	if oldObjList != nil {
		oldSvcList, _ := oldObjList.([]*apiv1.Service)
		for _, old := range oldSvcList {
			found := false
			for _, new := range newSvcs {
				if old.Name == new.Name {
					found = true
					break
				}
			}
			if !found {
				delSvcs = append(delSvcs, old)
			}
		}
	}
	// delete svc if only the old app has it
	if err := kr.deleteService(delSvcs); err != nil {
		beego.Warn("delete old service error:%v", err)
	}
	return nil
}

func (kr *KubeAppRes) SetAppStatus(app *models.CaasApplication) {
	status, err := kr.kubeAppHandle.Status(app.Name)
	if err != nil {
		log.Log.Error("get application status failed:", err)
		return
	}
	app.ReadyReplicas = status.ReadyReplicas
	app.UpdatedReplicas = status.UpdatedReplicas
	app.AvailableReplicas = status.AvailableReplicas
	app.AvailableStatus = status.AvailableStatus
	app.Message = status.Message
}

func (kr *KubeAppRes) Restart(appname string) error {
	return kr.kubeAppHandle.Restart(appname)
}

func (kr *KubeAppRes) DeleteApplication(appname string) error {
	return kr.kubeAppHandle.Delete(appname)
}

func (kr *KubeAppRes) CheckAppIsExisted(appname string) (bool, error) {
	return kr.kubeAppHandle.AppIsExisted(appname)
}

func (kr *KubeAppRes) Scale(appname string, replicas int) error {
	return kr.kubeAppHandle.Scale(appname, replicas)
}

func (kr *KubeAppRes) CreateService(svcList []*apiv1.Service) error {
	for _, svc := range svcList {
		svcIsExisted := false
		old, err := kr.client.CoreV1().Services(svc.Namespace).Get(svc.Name, metav1.GetOptions{})
		if err == nil {
			svcIsExisted = true
			svc.ResourceVersion = old.ResourceVersion
			svc.Spec.ClusterIP = old.Spec.ClusterIP
			for i, port1 := range svc.Spec.Ports {
				for _, port2 := range old.Spec.Ports {
					if port1.TargetPort != port2.TargetPort {
						continue
					}
					svc.Spec.Ports[i].Name = port2.Name
					if svc.Spec.Type == apiv1.ServiceTypeNodePort && port1.NodePort == 0 {
						svc.Spec.Ports[i].NodePort = port2.NodePort
					}
				}
			}
			// change type to nodeport
			if svc.Spec.Type == apiv1.ServiceTypeNodePort &&
				svc.Spec.Type != old.Spec.Type &&
				len(svc.Spec.Ports) > 1 {
				for _, port := range svc.Spec.Ports {
					if port.NodePort == 0 {
						return fmt.Errorf("you cant change service type to nodeport if the service has two or more ports without nodeports!")
					}
				}
			}
			//check annotation: // if not existed in new svc, then add it
			updateMap(svc.Annotations, old.Annotations)
		} else {
			if !errors.IsNotFound(err) {
				return err
			}
		}
		if svcIsExisted {
			if _, err := kr.client.CoreV1().Services(svc.Namespace).Update(svc); err != nil {
				return fmt.Errorf("update service error: %v", err)
			}
		} else {
			if _, err := kr.client.CoreV1().Services(svc.Namespace).Create(svc); err != nil {
				return fmt.Errorf("create service error: %v", err)
			}
		}
	}
	return nil
}

func (kr *KubeAppRes) deleteService(svcList []*apiv1.Service) error {
	for _, svc := range svcList {
		if err := kr.client.CoreV1().Services(kr.Namespace).Delete(svc.Name, &metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				return fmt.Errorf("delete service error %v", err)
			}
		} else {
			beego.Warn(fmt.Sprintf("delete service %s successfully!", svc.Name))
		}
	}
	return nil
}

func deleteIstioEgressByIngress(cluster string, ingress *extensions.Ingress) error {
	if ingress == nil {
		return nil
	}
	// TODO:
	return nil
}

func updateMap(new, old map[string]string) {
	for k, v := range old {
		if _, existed := new[k]; !existed {
			new[k] = v
		}
	}
}
