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
	v1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/util/labels"
)

const (
	RESTART_LABLE       = "caas_restart"
	RESTART_LABLE_VALUE = "true"
)

type KubeAppInterface interface {
	Create(obj interface{}) error
	Update(app models.CaasApplication, obj interface{}) (int, error)
	Status(appname string) (*AppStatus, error)
	Delete(appname string) error
	DeletePods(selector *metav1.ListOptions, appname string) error
	AppIsExisted(appname string) (bool, error)
	Scale(appname string, replicas int) error
	Restart(appname string) error
	GetOwnerForPod(pod apiv1.Pod, ref *metav1.OwnerReference) interface{}
}

type AppStatus struct {
	ReadyReplicas     int32
	UpdatedReplicas   int32
	AvailableReplicas int32
	AvailableStatus   string
	Message           string
}

type DeploymentRes struct {
	Namespace string
	client    kubernetes.Interface
}

func NewDeploymentRes(client kubernetes.Interface, namespace string) KubeAppInterface {
	return &DeploymentRes{
		Namespace: namespace,
		client:    client,
	}
}

func (kr *DeploymentRes) Create(obj interface{}) error {
	dp, ok := obj.(*v1.Deployment)
	if !ok {
		return fmt.Errorf("can not generate deployment object")
	}
	beego.Info("creating deployment, " + dp.Name)
	_, err := kr.client.AppsV1().Deployments(kr.Namespace).Create(dp)
	return err
}

func (kr *DeploymentRes) Update(app models.CaasApplication, obj interface{}) (int, error) {
	newDp, ok := obj.(*v1.Deployment)
	if !ok {
		return http.StatusBadRequest, fmt.Errorf("can not generate deployment object")
	}
	_, err := kr.client.AppsV1().Deployments(app.Namespace).Get(GenerateDeployName(app.Name), metav1.GetOptions{})
	if err == nil {
		_, err = kr.client.AppsV1().Deployments(app.Namespace).Update(newDp)
	}
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("update deployment error: %v", err.Error())
	}
	return http.StatusOK, nil
}

// Status ..
func (kr *DeploymentRes) Status(appname string) (*AppStatus, error) {
	deployment, err := kr.client.AppsV1().Deployments(kr.Namespace).Get(GenerateDeployName(appname), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	status := &AppStatus{
		ReadyReplicas:     deployment.Status.ReadyReplicas,
		AvailableReplicas: deployment.Status.AvailableReplicas,
		UpdatedReplicas:   deployment.Status.UpdatedReplicas,
	}
	for _, condition := range deployment.Status.Conditions {
		if condition.Type == v1.DeploymentAvailable {
			status.AvailableStatus = string(condition.Status)
			status.Message = condition.Message
			break
		}
	}
	return status, nil
}

// DeletePods ..
func (kr *DeploymentRes) DeletePods(options *metav1.ListOptions, appname string) error {
	err := error(nil)
	listOptions := options
	if listOptions == nil {
		if listOptions, err = kr.getPodOption(appname); err != nil {
			return err
		}
	}

	// delete old application's pod which has no version label
	podList, err := kr.client.CoreV1().Pods(kr.Namespace).List(*listOptions)
	if err != nil {
		return fmt.Errorf("delete pod list error %v", err)
	}
	for _, pod := range podList.Items {
		if _, existed := pod.Labels[LABLE_APPVERSION_KEY]; !existed {
			if err := kr.client.CoreV1().Pods(kr.Namespace).Delete(pod.Name, &metav1.DeleteOptions{}); err != nil {
				return fmt.Errorf("delete pod list error %v", err)
			}
		}
	}
	return nil
}

func (kr *DeploymentRes) Delete(appname string) error {
	listOptions, err := kr.getPodOption(appname)
	if err != nil {
		if !errors.IsNotFound(err) {
			return fmt.Errorf("get pod option error %v", err)
		} else {
			return nil
		}
	}
	// TODO:  delete app, 未包含version, 导致删除 relicaset时找不到 资源，
	// the server could not find the requested resource cluster: dev, namespace: default, name: nginx-deployment01
	if err := kr.client.AppsV1().Deployments(kr.Namespace).Delete(GenerateDeployName(appname), &metav1.DeleteOptions{}); err != nil {
		if !errors.IsNotFound(err) {
			return fmt.Errorf("delete deployment error %v", err)
		}
	}
	return kr.DeletePods(listOptions, appname)
}

// AppIsExisted ..
func (kr *DeploymentRes) AppIsExisted(appname string) (bool, error) {
	_, err := kr.client.AppsV1().Deployments(kr.Namespace).Get(GenerateDeployName(appname), metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return false, err
		}
		return false, nil

	}
	return true, nil
}

func (kr *DeploymentRes) Scale(appname string, replicas int) error {
	dpname := GenerateDeployName(appname)
	ds, err := kr.client.AppsV1().Deployments(kr.Namespace).Get(dpname, metav1.GetOptions{})
	if err != nil {
		return err
	}
	*ds.Spec.Replicas = int32(replicas)
	if _, err := kr.client.AppsV1().Deployments(kr.Namespace).Update(ds); err != nil {
		return err
	}

	return nil
}

func (kr *DeploymentRes) Restart(appname string) error {
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

func (kr *DeploymentRes) GetOwnerForPod(pod apiv1.Pod, ref *metav1.OwnerReference) interface{} {
	if ref == nil {
		return nil
	}
	rs, err := kr.client.ExtensionsV1beta1().ReplicaSets(pod.Namespace).Get(ref.Name, metav1.GetOptions{})
	if err != nil || rs.UID != ref.UID {
		beego.Warn(fmt.Sprintf("Cannot get replicaset %s for pod %s: %v", ref.Name, pod.Name, err))
		return nil
	}
	// Now find the Deployment that owns that ReplicaSet.
	depRef := metav1.GetControllerOf(rs)
	if depRef == nil {
		return nil
	}
	// We can't look up by UID, so look up by Name and then verify UID.
	// Don't even try to look up by Name if it's the wrong Kind.
	if depRef.Kind != v1beta1.SchemeGroupVersion.WithKind("Deployment").Kind {
		return nil
	}
	d, err := kr.client.AppsV1().Deployments(pod.Namespace).Get(depRef.Name, metav1.GetOptions{})
	if err != nil {
		return nil
	}
	if d.UID != depRef.UID {
		return nil
	}
	return d
}

func (kr *DeploymentRes) getPodOption(appname string) (*metav1.ListOptions, error) {
	dp, err := kr.client.AppsV1().Deployments(kr.Namespace).Get(GenerateDeployName(appname), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	selector, err := metav1.LabelSelectorAsSelector(dp.Spec.Selector)
	if err != nil {
		return nil, err
	}
	return &metav1.ListOptions{LabelSelector: selector.String()}, nil
}
