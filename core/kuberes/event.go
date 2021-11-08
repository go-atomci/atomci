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
	"time"

	"github.com/go-atomci/atomci/middleware/log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type EventItem struct {
	EventUid        string    `orm:"column(event_uid);size(36)" json:"event_uid"`
	ActionType      string    `orm:"column(action_type);size(10)" json:"action_type"`
	EventType       string    `orm:"column(event_type);size(10)" json:"event_type"`
	Cluster         string    `orm:"column(cluster);size(20)" json:"cluster"`
	Namespace       string    `orm:"column(namespace);size(100)" json:"namespace"`
	SourceComponent string    `orm:"column(source_component);size(50)" json:"source_component"`
	SourceHost      string    `orm:"column(source_host);size(60)" json:"source_host"`
	ObjectKind      string    `orm:"column(object_kind);size(64)" json:"object_kind"`
	ObjectName      string    `orm:"column(object_name);size(100)" json:"object_name"`
	ObjectUid       string    `orm:"column(object_uid);size(36)" json:"object_uid"`
	FieldPath       string    `orm:"column(field_path);size(200)" json:"field_path"`
	Reason          string    `orm:"column(reason);size(100)" json:"reason"`
	Message         string    `orm:"column(message);type(text)" json:"message"`
	FirstTimestamp  time.Time `orm:"column(first_time)" json:"first_time"`
	LastTimestamp   time.Time `orm:"column(last_time);index" json:"last_time"`
}

func GetEventList(client kubernetes.Interface, cluster, namespace, appName string) ([]*EventItem, error) {
	labelSelector := ""
	deployment, error := client.AppsV1().Deployments(namespace).Get(appName, metav1.GetOptions{})
	if error != nil {
		log.Log.Error(fmt.Sprintf("Get %v deploy on namespace %v in cluster %v Error: %v", appName, namespace, cluster, error.Error()))
	} else if deployment != nil {
		selector := deployment.Spec.Selector.MatchLabels

		for k, v := range selector {
			if k == "version" {
				continue
			}
			labelSelector += k + "=" + v + ","
		}
		len := len(labelSelector) - 1
		labelSelector = labelSelector[:len]
	}

	if labelSelector == "" {
		labelSelector = LABLE_APPNAME_KEY + "=" + appName
	}
	log.Log.Debug("get app events labelSeletor: %v", labelSelector)
	k8sEvents, err := client.CoreV1().Events(namespace).List(metav1.ListOptions{
		// LabelSelector: labelSelector,
	})
	if err != nil {
		log.Log.Error(fmt.Sprintf("Get %v event on namespace %v in cluster %v Error: %v", labelSelector, namespace, cluster, err.Error()))
		return nil, err
	} else {
		log.Log.Debug("k8s events len: %v", len(k8sEvents.Items))
	}
	events := []*EventItem{}
	for _, eventItem := range k8sEvents.Items {
		event := &EventItem{}
		event = eventFormat(eventItem, event)
		events = append(events, event)
	}
	return events, nil
}

func eventFormat(item v1.Event, event *EventItem) *EventItem {
	firstTime, _ := time.Parse("2006-01-02 15:04:05", item.FirstTimestamp.Format("2006-01-02 15:04:05"))
	lastTime, _ := time.Parse("2006-01-02 15:04:05", item.LastTimestamp.Format("2006-01-02 15:04:05"))

	event.EventUid = string(item.ObjectMeta.UID)
	event.EventType = item.Type
	event.Cluster = item.ClusterName
	event.Namespace = item.ObjectMeta.Namespace
	event.SourceComponent = item.Source.Component
	event.SourceHost = item.Source.Host
	event.ObjectKind = item.InvolvedObject.Kind
	event.ObjectName = item.InvolvedObject.Name
	event.ObjectUid = string(item.InvolvedObject.UID)
	event.FieldPath = item.InvolvedObject.FieldPath
	event.Reason = item.Reason
	event.Message = item.Message
	event.FirstTimestamp = firstTime
	event.LastTimestamp = lastTime
	return event
}
