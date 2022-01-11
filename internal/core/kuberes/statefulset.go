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
	"net/http"

	"github.com/go-atomci/atomci/internal/models"

	"github.com/astaxie/beego"
	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/util/labels"
)

type StatefulRes struct {
	Namespace string
	client    kubernetes.Interface
}

const emptyVersion = ""

func NewStatefulRes(client kubernetes.Interface, namespace string) KubeAppInterface {
	return &StatefulRes{
		Namespace: namespace,
		client:    client,
	}
}

func (kr *StatefulRes) Create(obj interface{}) error {
	state, ok := obj.(*v1.StatefulSet)
	if !ok {
		return fmt.Errorf("can not generate statefulset object")
	}
	beego.Info("creating statefulset application, " + state.Name)
	_, err := kr.client.AppsV1().StatefulSets(kr.Namespace).Create(state)
	return err
}

func (kr *StatefulRes) Update(app models.CaasApplication, obj interface{}) (int, error) {
	newSet, ok := obj.(*v1.StatefulSet)
	if !ok {
		return http.StatusBadRequest, fmt.Errorf("can not generate statefulset object")
	}
	_, err := kr.client.AppsV1().StatefulSets(app.Namespace).Get(GenerateDeployName(app.Name), metav1.GetOptions{})
	if err == nil {
		_, err = kr.client.AppsV1().StatefulSets(app.Namespace).Update(newSet)
	}
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("update statefulset error: %v", err.Error())
	}

	return http.StatusOK, nil
}

func (kr *StatefulRes) Status(appname string) (*AppStatus, error) {
	set, err := kr.client.AppsV1().StatefulSets(kr.Namespace).Get(GenerateDeployName(appname), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &AppStatus{
		ReadyReplicas:     set.Status.ReadyReplicas,
		AvailableReplicas: set.Status.CurrentReplicas + set.Status.UpdatedReplicas,
		UpdatedReplicas:   set.Status.UpdatedReplicas,
	}, nil
}

func (kr *StatefulRes) Delete(appname string) error {
	return kr.deleteStatefulSet(appname, emptyVersion)
}

func (kr *StatefulRes) DeletePods(options *metav1.ListOptions, appname string) error {
	listOptions := options
	err := error(nil)
	if listOptions == nil {
		if listOptions, err = kr.getPodOption(appname, emptyVersion); err != nil {
			if errors.IsNotFound(err) {
				return nil
			}
			return err
		}
	}
	if err := kr.client.CoreV1().Pods(kr.Namespace).DeleteCollection(&metav1.DeleteOptions{}, *listOptions); err != nil {
		return fmt.Errorf("delete pod list error %v", err)
	}

	return nil
}

// AppIsExisted ..
func (kr *StatefulRes) AppIsExisted(appname string) (bool, error) {
	_, err := kr.client.AppsV1beta1().StatefulSets(kr.Namespace).Get(GenerateDeployName(appname), metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return false, err
		} else {
			return false, nil
		}
	}
	return true, nil
}

func (kr *StatefulRes) Scale(appname string, replicas int) error {
	ds, err := kr.client.AppsV1beta1().StatefulSets(kr.Namespace).Get(GenerateDeployName(appname), metav1.GetOptions{})
	if err != nil {
		return err
	}
	*ds.Spec.Replicas = int32(replicas)
	if _, err := kr.client.AppsV1beta1().StatefulSets(kr.Namespace).Update(ds); err != nil {
		return err
	}

	return nil
}

func (kr *StatefulRes) Restart(appname string) error {
	dpname := GenerateDeployName(appname)
	dp, err := kr.client.AppsV1().Deployments(kr.Namespace).Get(dpname, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if _, exist := dp.Spec.Template.ObjectMeta.Annotations[RESTART_LABLE]; exist {
		delete(dp.Spec.Template.ObjectMeta.Annotations, RESTART_LABLE)
	} else {
		dp.Spec.Template.ObjectMeta.Annotations = labels.AddLabel(dp.Spec.Template.ObjectMeta.Annotations, RESTART_LABLE, RESTART_LABLE_VALUE)
	}
	if _, err := kr.client.AppsV1().Deployments(kr.Namespace).Update(dp); err != nil {
		return err
	}

	return nil
}

func (kr *StatefulRes) GetOwnerForPod(pod apiv1.Pod, ref *metav1.OwnerReference) interface{} {
	if ref == nil {
		return nil
	}
	set, err := kr.client.AppsV1().StatefulSets(pod.Namespace).Get(ref.Name, metav1.GetOptions{})
	if err != nil {
		return nil
	}
	if set.UID != ref.UID {
		return nil
	}
	return set
}

func (kr *StatefulRes) createService(svc *apiv1.Service, svcIsExist bool) error {
	if svcIsExist {
		//if existed, update
		if _, err := kr.client.CoreV1().Services(svc.Namespace).Update(svc); err != nil {
			return fmt.Errorf("update service error: %v", err)
		}
	} else {
		//if not exist, create
		if _, err := kr.client.CoreV1().Services(svc.Namespace).Create(svc); err != nil {
			return fmt.Errorf("create service error: %v", err)
		}
	}
	return nil
}

func (kr *StatefulRes) deleteStatefulSet(appname, version string) error {
	if err := kr.client.AppsV1().StatefulSets(kr.Namespace).Delete(GenerateDeployName(appname), &metav1.DeleteOptions{}); err != nil {
		if !errors.IsNotFound(err) {
			return fmt.Errorf("delete statefulset error %v", err)
		}
	}

	return kr.DeletePods(nil, appname)
}

func (kr *StatefulRes) getPodOption(appname, version string) (*metav1.ListOptions, error) {
	state, err := kr.client.AppsV1().StatefulSets(kr.Namespace).Get(GenerateDeployName(appname), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	selector, err := metav1.LabelSelectorAsSelector(state.Spec.Selector)
	if err != nil {
		return nil, err
	}
	return &metav1.ListOptions{LabelSelector: selector.String()}, nil
}
