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

import "fmt"

// TriggerApplicationCreate ..
func TriggerApplicationCreate(clusterName, namespace, templateStr string, projectID, envID int64, force bool) error {
	tpl := NewTemplate()
	native := &NativeTemplate{
		Template: templateStr,
	}
	tpl = native
	if err := tpl.Validate(); err != nil {
		return fmt.Errorf("validate apps template occur error: %s, cluster: %s, namespace: %s", err.Error(), clusterName, namespace)
	}
	ar, err := NewAppRes(clusterName, envID, projectID)
	if err != nil {
		return fmt.Errorf("created app res occur error: %s, cluster: %s, namespace: %s", err.Error(), clusterName, namespace)
	}
	eparam := ExtensionParam{
		Force: force,
	}
	err = ar.InstallApp(namespace, "", tpl, &eparam)
	if err != nil {
		return fmt.Errorf("deploy application occur error: %s, cluster: %s, namespace: %s", err.Error(), clusterName, namespace)
	}
	return nil
}
