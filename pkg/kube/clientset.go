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
	"sync"

	"k8s.io/client-go/kubernetes"
)

var (
	clusterClientsetMapMutex sync.RWMutex
	clusterClientsetMap      = make(map[string]kubernetes.Interface)
)

func findClientset(cluster string) (client kubernetes.Interface, ok bool) {
	clusterClientsetMapMutex.RLock()
	defer clusterClientsetMapMutex.RUnlock()
	client, ok = clusterClientsetMap[cluster]
	return client, ok
}

func newClientset(cluster string) (client kubernetes.Interface, err error) {
	var ok bool
	clusterClientsetMapMutex.Lock()
	defer clusterClientsetMapMutex.Unlock()
	client, ok = clusterClientsetMap[cluster]
	if !ok {
		client, err = clientsetProvider(cluster)
		if err == nil {
			clusterClientsetMap[cluster] = client
		}
	}
	return client, err
}

func GetClientset(cluster string) (client kubernetes.Interface, err error) {
	var ok bool
	client, ok = findClientset(cluster)
	if !ok {
		client, err = newClientset(cluster)
	}
	return client, err
}

func UpdateClientset(cluster string) (client kubernetes.Interface, err error) {
	clusterClientsetMapMutex.Lock()
	defer clusterClientsetMapMutex.Unlock()
	client, err = clientsetProvider(cluster)
	if err != nil {
		return
	}
	clusterClientsetMap[cluster] = client
	return
}
