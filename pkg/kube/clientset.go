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

package kube

import (
	"github.com/go-atomci/atomci/internal/core/settings"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetClientset(cluster string) (client kubernetes.Interface, cfg *rest.Config, err error) {

	pm := settings.NewSettingManager()
	resp, err := pm.GetIntegrateSettingByName(cluster, settings.KubernetesType)
	if err != nil {
		return nil, nil, err
	}
	return buildK8sClient(resp.IntegrateSettingReq.Config.(*settings.KubeConfig))
}

func buildK8sClient(kube *settings.KubeConfig) (client kubernetes.Interface, cfg *rest.Config, err error) {
	var k8sConfig *rest.Config
	switch kube.Type {
	case settings.KubernetesConfig:
		k8sConfig, err = clientcmd.RESTConfigFromKubeConfig([]byte(kube.Conf))
		if err != nil {
			return nil, nil, err
		}
	case settings.KubernetesToken:
		k8sConfig = &rest.Config{
			BearerToken:     kube.Conf,
			TLSClientConfig: rest.TLSClientConfig{Insecure: true},
			Host:            kube.URL,
		}
	}

	clientSet, err := kubernetes.NewForConfig(k8sConfig)
	return clientSet, k8sConfig, err
}
