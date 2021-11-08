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

package controllers

import (
	"fmt"
	"path"

	"github.com/go-atomci/atomci/core/podexec"
	"github.com/go-atomci/atomci/middleware/log"
	"github.com/go-atomci/atomci/pkg/kube"

	"github.com/astaxie/beego"
	"k8s.io/client-go/tools/clientcmd"
)

type TerminalController struct {
	BaseController
}

func (t *TerminalController) PodTerminal() {
	cluster := t.Ctx.Input.Param(":cluster")
	namespace := t.Ctx.Input.Param(":namespace")
	podName := t.Ctx.Input.Param(":podname")
	containerName := t.Ctx.Input.Param(":containername")

	if cluster == "" || namespace == "" || podName == "" || containerName == "" {
		log.Log.Error("args missing, cluster: %s, naemspace: %s, podName: %s, containerName: %s", cluster, namespace, podName, containerName)
		t.HandleInternalServerError(fmt.Sprintf("args missing, cluster: %s, naemspace: %s, podName: %s, containerName: %s", cluster, namespace, podName, containerName))
		return
	}
	log.Log.Info("exec containerName: %s, pod: %s, namespace: %s", containerName, podName, namespace)

	pty, err := podexec.NewTerminalSession(t.Ctx.ResponseWriter, t.Ctx.Request, nil)
	if err != nil {
		log.Log.Error("get pty failed: %v", err.Error())
		t.HandleInternalServerError(fmt.Sprintf("get pty failed: %v", err.Error()))
		return
	}

	defer func() {
		log.Log.Info("close session.")
		_ = pty.Close()
	}()

	kubeCli, err := kube.GetClientset(cluster)
	if err != nil {
		msg := fmt.Sprintf("get kubecli err :%v", err)
		log.Log.Error(msg)
		_, _ = pty.Write([]byte(msg))
		pty.Done()

		t.HandleInternalServerError(msg)
		return
	}

	ok, err := podexec.ValidatePod(kubeCli, namespace, podName, containerName)
	if !ok {
		msg := fmt.Sprintf("Validate pod error! err: %v", err)
		log.Log.Error(msg)
		_, _ = pty.Write([]byte(msg))
		pty.Done()

		t.HandleInternalServerError(msg)
		return
	}

	configFile := path.Join(beego.AppConfig.String("k8s::configPath"), cluster)
	cfg, err := clientcmd.BuildConfigFromFlags("", configFile)
	if err != nil {
		msg := fmt.Sprintf("build config occur error: %s", err.Error())
		log.Log.Error(msg)
		t.HandleInternalServerError(msg)
		return
	}
	err = podexec.ExecPod(kubeCli, cfg, []string{"/bin/sh"}, pty, namespace, podName, containerName)
	if err != nil {
		msg := fmt.Sprintf("Exec to pod error! err: %v", err)
		log.Log.Error(msg)
		_, _ = pty.Write([]byte(msg))
		pty.Done()

		t.HandleInternalServerError(msg)
		return
	}
}
