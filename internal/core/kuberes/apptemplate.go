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
	"strings"
	"sync"

	"github.com/go-atomci/atomci/internal/middleware/log"
	"github.com/go-atomci/atomci/internal/models"

	yamlencoder "github.com/ghodss/yaml"
)

const TemplateMaxSize = 1024000 //1000KB

//single app template interface
//simpleapp and nativeapp template support this interface
type AppTemplate interface {
	GenerateAppObject(cluster, namespace, tplname string, projectID int64) (*models.CaasApplication, error)
	UpdateAppObject(app *models.CaasApplication) error
	GenerateKubeObject(cluster, namespace string) (map[string]interface{}, error)
	GetAppName() string
	GetAppKind() string
	String() (string, error)
	Image(param []ContainerParam) AppTemplate
	DefaultLabel() AppTemplate
	Replicas(replicas int) AppTemplate
}

type WorkerResult struct {
	AppName string
	AppKind string
	Result  error
}

func CreateAppTemplateByApp(app models.CaasApplication) (AppTemplate, error) {
	return CreateNativeAppTemplate(app)
}

func DeployAppTemplates(appTplList []AppTemplate, projectid, envID int64, cluster, namespace, tname string, eparam *ExtensionParam) error {
	if len(appTplList) == 0 {
		return nil
	}
	errInfoList := []string{}
	ar, err := NewAppRes(cluster, envID, projectid)
	if err != nil {
		return err
	}
	workerResult := make(chan WorkerResult)
	var wg sync.WaitGroup
	for _, tpl := range appTplList {
		wg.Add(1)
		wk := NewDeployWorker(tpl.GetAppName(), namespace, tpl.GetAppKind(), ar, eparam, tpl)
		param := AppParam{Name: tpl.GetAppName()}
		go func(app AppTemplate) {
			defer wg.Done()
			err := wk.Start(tname, param)
			workerResult <- WorkerResult{
				AppName: app.GetAppName(),
				AppKind: app.GetAppKind(),
				Result:  err,
			}
		}(tpl)
	}
	go func() {
		wg.Wait()
		close(workerResult)
	}()
	for res := range workerResult {
		if res.Result != nil {
			errInfoList = append(errInfoList, res.AppName+":"+res.Result.Error())
			log.Log.Error("%v", res.Result)
		} else {
			log.Log.Info("deploy application " + res.AppName + " successfully!")
		}
	}
	if len(errInfoList) != 0 {
		return fmt.Errorf(strings.Join(errInfoList, ";"))
	}

	return nil
}

func AppTemplateToYamlString(tpl AppTemplate, cluster, namespace, podVersion string) (string, error) {
	objs, err := tpl.GenerateKubeObject(cluster, namespace)
	if err != nil && objs == nil {
		log.Log.Error("generate kubernetes object failed:", err)
		return "", err
	}
	ctx := []byte{}
	elems := []reflect.Value{}
	for _, obj := range objs {
		v := reflect.ValueOf(obj)
		switch v.Kind() {
		case reflect.Ptr:
			elems = append(elems, v)
		case reflect.Slice, reflect.Array:
			for i := 0; i < v.Len(); i++ {
				elems = append(elems, v.Index(i))
			}
		default:
			log.Log.Debug("object kind:", v.Kind())
		}
	}
	for _, elem := range elems {
		yamlBytes, err := yamlencoder.Marshal(elem.Interface())
		if err != nil {
			log.Log.Error("yaml marshal object failed:", err)
		}
		ctx = append(ctx, yamlBytes...)
		ctx = append(ctx, []byte(YamlSeparator)...)
	}

	return strings.TrimSuffix(string(ctx), YamlSeparator), nil
}
