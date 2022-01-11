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

	"github.com/astaxie/beego/logs"
	apps "k8s.io/api/apps/v1beta1"
	v1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	PodStatusRunning  = "Running"
	PodStatusNotReady = "NotReady"
)

type Pod struct {
	Name         string            `json:"name"`
	Namespace    string            `json:"namespace"`
	Version      string            `json:"version"`
	NodeIP       string            `json:"node_ip"`
	PodIP        string            `json:"pod_ip"`
	Status       string            `json:"status"`
	Message      string            `json:"message"`
	RestartCount int32             `json:"restart_count"`
	StartTime    string            `json:"start_time"`
	Labels       map[string]string `json:"labels"`
	Containers   []*PodContainer   `json:"containers"`
}

type PodContainer struct {
	Name           string `json:"name"`
	Image          string `json:"image"`
	CpuRequests    int64  `json:"cpu_requests"`
	CpuLimits      int64  `json:"cpu_limits"`
	MemoryRequests int64  `json:"memory_requests"`
	MemoryLimits   int64  `json:"memory_limits"`
}

func GetPods(client kubernetes.Interface, cluster, namespace, appName string, replicas int) ([]*Pod, error) {
	labelSelector := ""
	deployment, error := client.AppsV1().Deployments(namespace).Get(appName, metav1.GetOptions{})
	if error != nil {
		log.Log.Error(fmt.Sprintf("Get %v deploy on namespace %v in cluster %v Error: %v", appName, namespace, cluster, error.Error()))
	} else if deployment != nil {
		selector := deployment.Spec.Selector.MatchLabels

		for k, v := range selector {
			labelSelector += k + "=" + v + ","
		}
		len := len(labelSelector) - 1
		labelSelector = labelSelector[:len]
	}

	if labelSelector == "" {
		labelSelector = LABLE_APPNAME_KEY + "=" + appName
	}
	k8sPods, err := client.CoreV1().Pods(namespace).List(metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Log.Error(fmt.Sprintf("Get %v pods on namespace %v in cluster %v Error: %v", labelSelector, namespace, cluster, err.Error()))
		return nil, err
	}
	pods := []*Pod{}
	podNum := len(k8sPods.Items)
	noPodIPNum := 0
	for _, k8spod := range k8sPods.Items {
		if k8spod.Status.PodIP == "" {
			if k8spod.Status.Reason == "Evicted" {
				logs.Warn(k8spod.Name, "is Evicted, and will be deleted,", "for", k8spod.Status.Message)
				client.CoreV1().Pods(namespace).Delete(k8spod.Name, &metav1.DeleteOptions{})
				continue
			}
			// if pod number is larger than replicas*2, it is not normal, dont need to check
			added := false
			if podNum > (replicas << 1) {
				noPodIPNum++
				added = true
				if noPodIPNum > (replicas<<1) && replicas != 0 {
					continue
				}
			}
			if isTrashPod(client, k8spod) {
				// trash pod, filtered
				if added {
					noPodIPNum--
				}
				gracePeriod := int64(0)
				client.CoreV1().Pods(namespace).Delete(k8spod.Name, &metav1.DeleteOptions{
					GracePeriodSeconds: &gracePeriod,
				})
				continue
			}
		}
		pod := podConv(k8spod)
		pods = append(pods, pod)
	}
	return pods, nil
}

func podConv(k8sPod v1.Pod) *Pod {
	pod := &Pod{
		Name:      k8sPod.ObjectMeta.Name,
		Namespace: k8sPod.ObjectMeta.Namespace,
		Version:   GetResourceVersion(&k8sPod, ResTypePod, ""),
		NodeIP:    k8sPod.Status.HostIP,
		PodIP:     k8sPod.Status.PodIP,
	}
	pod.Status, pod.Message, pod.RestartCount = getPodStatus(k8sPod)
	if k8sPod.Status.StartTime != nil {
		pod.StartTime = k8sPod.Status.StartTime.Format("2006-01-02 15:04:05")
	}

	for _, k8sContainer := range k8sPod.Spec.Containers {
		container := podContainerConv(k8sContainer)
		pod.Containers = append(pod.Containers, container)
	}

	return pod
}

func podContainerConv(k8scontainer v1.Container) *PodContainer {
	container := &PodContainer{
		Name:           k8scontainer.Name,
		Image:          k8scontainer.Image,
		CpuRequests:    k8scontainer.Resources.Requests.Cpu().MilliValue(),
		CpuLimits:      k8scontainer.Resources.Limits.Cpu().MilliValue(),
		MemoryRequests: k8scontainer.Resources.Requests.Memory().Value(),
		MemoryLimits:   k8scontainer.Resources.Limits.Memory().Value(),
	}
	return container
}

func getPodStatus(pod v1.Pod) (status string, message string, restartCount int32) {
	status = PodStatusRunning
	if pod.Status.Phase == v1.PodRunning {
		for _, c := range pod.Status.Conditions {
			if c.Type == v1.PodReady {
				if c.Status == v1.ConditionFalse {
					status = PodStatusNotReady
					message = fmt.Sprintf("%s: %s", c.Reason, c.Message)
					break
				}
			}
		}
	} else {
		status = PodStatusNotReady
		for _, c := range pod.Status.Conditions {
			if c.Type == v1.PodScheduled {
				if c.Status != v1.ConditionTrue {
					message = fmt.Sprintf("%s: %s", c.Reason, c.Message)
					break
				}
			}
		}
	}

	for _, cs := range pod.Status.ContainerStatuses {
		restartCount += cs.RestartCount
		if status == PodStatusNotReady && cs.State.Waiting != nil {
			message = fmt.Sprintf("%s: %s", cs.State.Waiting.Reason, cs.State.Waiting.Message)
		}
		if status == PodStatusNotReady && cs.State.Terminated != nil {
			message = fmt.Sprintf("%s: %s", cs.State.Terminated.Reason, cs.State.Terminated.Message)
		}
	}

	return
}

func isTrashPod(client kubernetes.Interface, pod v1.Pod) bool {
	ref := metav1.GetControllerOf(&pod)
	if ref == nil {
		return true
	}
	switch ref.Kind {
	case extensions.SchemeGroupVersion.WithKind("ReplicaSet").Kind:
		return NewDeploymentRes(client, pod.Namespace).GetOwnerForPod(pod, ref) == nil
	case apps.SchemeGroupVersion.WithKind("StatefulSet").Kind:
		return NewStatefulRes(client, pod.Namespace).GetOwnerForPod(pod, ref) == nil
	}
	return false
}
