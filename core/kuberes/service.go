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
	"reflect"

	"github.com/go-atomci/atomci/middleware/log"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

type serviceAddress struct {
	Port         int32       `json:"port"`
	TargetPort   interface{} `json:"target_port"`
	NodePort     int32       `json:"node_port"`
	Protocol     string      `json:"protocol"`
	NodePortAddr string      `json:"node_port_addr"`
	ClusterAddr  string      `json:"cluster_addr"`
}

type podSvcAddress struct {
	TargetPort  int    `json:"target_port"`
	Protocol    string `json:"protocol"`
	ClusterAddr string `json:"cluster_addr"`
}

type ServiceDetail struct {
	Name           string           `json:"name"`
	Type           string           `json:"type"`
	ClusterIP      string           `json:"cluster_ip"`
	AddressList    []serviceAddress `json:"address_list"`
	PodsvcAddrList []podSvcAddress  `json:"podsvc_addr_list"`
}

func GetAppServices(client kubernetes.Interface, cluster, namespace, appDeploymentName, nodeIP string) ([]*ServiceDetail, error) {
	var selector map[string]string
	deployment, error := client.AppsV1().Deployments(namespace).Get(appDeploymentName, metav1.GetOptions{})
	if error != nil {
		log.Log.Error(fmt.Sprintf("Get %v deploy on namespace %v in cluster %v Error: %v", appDeploymentName, namespace, cluster, error.Error()))
	} else if deployment != nil {
		selector = deployment.Spec.Selector.MatchLabels

	}

	k8sSvcs, err := client.CoreV1().Services(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Log.Error(fmt.Sprintf("Get %v svc on namespace %v in cluster %v Error: %v", selector, namespace, cluster, err.Error()))
		return nil, err
	}
	svcs := []*ServiceDetail{}
	if len(selector) == 0 {
		return svcs, nil
	}

	for _, svcItem := range k8sSvcs.Items {
		log.Log.Debug("svc item labels: %v", svcItem.Spec.Selector)
		if len(selector) == len(svcItem.Spec.Selector) && reflect.DeepEqual(selector, svcItem.Spec.Selector) {
			svc := appServiceFormat(svcItem, nodeIP)
			svcs = append(svcs, svc)
		}
	}
	return svcs, nil
}

func appServiceFormat(svc apiv1.Service, nodeIP string) *ServiceDetail {
	svcDetail := &ServiceDetail{}
	svcDetail.Name = svc.Name
	svcDetail.Type = string(svc.Spec.Type)
	svcDetail.ClusterIP = svc.Spec.ClusterIP
	for _, item := range svc.Spec.Ports {
		var address serviceAddress
		address.Port = item.Port
		address.TargetPort = item.TargetPort
		address.Protocol = string(item.Protocol)
		address.NodePort = item.NodePort
		address.ClusterAddr = fmt.Sprintf("%s.%s:%v", svc.Name, svc.Namespace, address.Port)
		if apiv1.ServiceType(svcDetail.Type) == apiv1.ServiceTypeNodePort && nodeIP != "" {
			address.NodePortAddr = fmt.Sprintf("%s:%v", nodeIP, address.NodePort)
		} else {
			address.NodePortAddr = "<none>"
		}
		svcDetail.AddressList = append(svcDetail.AddressList, address)
	}
	return svcDetail
}
