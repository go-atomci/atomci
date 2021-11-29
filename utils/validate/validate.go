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

package validate

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-atomci/atomci/constant"
	"github.com/go-atomci/atomci/utils/errors"

	kubevalidation "k8s.io/apimachinery/pkg/util/validation"
)

const (
	restrictedNameChars      = `[a-zA-Z0-9]+(?:[-|.][a-zA-Z0-9]+)*`
	restrictedEmailAddr      = `^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z0-9]{2,6}$`
	restrictedKubernetesName = `^[a-zA-Z0-9_.-]+$`
)

var LabelsStrMaxLenMap = map[string]int{
	constant.K8S_RESOURCE_TYPE_NODE: 2048,
	constant.K8S_RESOURCE_TYPE_APP:  64,
}

const (
	NameMinLen        = 1
	NameMaxLen        = 64
	DescriptionMinLen = 1
	DescriptionMaxLen = 64
)

const (
	NodePortMin = 30000
	NodePortMax = 32767
)

// 验证字符长度
func IsIllegalLength(s string, min int, max int) bool {
	if min == -1 {
		return (len(s) > max)
	}
	if max == -1 {
		return (len(s) <= min)
	}
	return (len(s) < min || len(s) > max)
}

// 正则表达式验证字符合法性
func Restricted(s, regdata string) bool {
	validName := regexp.MustCompile(`^` + regdata + `$`)
	legal := validName.MatchString(s)
	return legal
}

// 系统保留的关键字
func IsReservedBuName(s string) error {
	stringList := []string{
		"all", "default", "kube-system",
	}
	for _, str := range stringList {
		if s == str {
			return fmt.Errorf(fmt.Sprintf("string %v is system reserved keywords，can not be used", s))
		}
	}
	return nil
}

func FormatString(s string) string {
	s = strings.TrimSpace(s)
	return s
}

func ValidateKubernetesName(s string) error {
	if !Restricted(s, restrictedKubernetesName) {
		return fmt.Errorf("此类型的配置名称: \"%v\" 不允许有中文", s)
	}
	return nil
}
func ValidateName(s string) error {
	if err := IsReservedBuName(s); err != nil {
		return err
	}
	if IsIllegalLength(s, NameMinLen, NameMaxLen) {
		return fmt.Errorf("the name \"%v\": must be no more than %v characters", s, NameMaxLen)
	}
	if !Restricted(s, restrictedNameChars) {
		return fmt.Errorf("the name \"%v\" is invalid format", s)
	}
	return nil
}

func ValidateDescription(s string) error {
	if err := IsReservedBuName(s); err != nil {
		return err
	}
	if IsIllegalLength(s, DescriptionMinLen, DescriptionMaxLen) {
		return fmt.Errorf("\"%v\": must be no more than %v characters", s, DescriptionMaxLen)
	}
	return nil
}

func ValidateEmail(s string) error {
	if !Restricted(s, restrictedEmailAddr) {
		return fmt.Errorf("the email \"%v\" is invalid format", s)
	}
	return nil
}

// 字符串验证
func ValidateString(s string) error {
	if IsIllegalLength(s, NameMinLen, NameMaxLen) {
		return errors.NewBadRequest().
			SetCode("InvalidStringLength").
			SetMessage("invalid string length, acceptable range is [%v,%v]", NameMinLen, NameMaxLen)
	}
	if !Restricted(s, restrictedNameChars) {
		return errors.NewBadRequest().
			SetCode("InvalidString").
			SetMessage(`invalid string "%s", must start with lower alpha, number beginning and ending, and '-' can be used for separating`, s)
	}
	return nil
}

func ValidateLabels(object string, labels map[string]string) error {
	labelStr, err := json.Marshal(labels)
	if err != nil {
		return err
	}
	if len(labelStr) > LabelsStrMaxLenMap[object] {
		return errors.NewBadRequest().
			SetCode("InvalidLabelsLength").
			SetMessage("invalid labels length, acceptable range is [%v,%v]", 0, LabelsStrMaxLenMap[object])
	}
	for key, value := range labels {
		if errs := kubevalidation.IsQualifiedName(key); len(errs) != 0 {
			return errors.NewBadRequest().
				SetCode("InvalidLabelKey").
				SetMessage(`invalid label key "%v:%s"`, key, strings.Join(errs, "; "))
		}
		if errs := kubevalidation.IsValidLabelValue(value); len(errs) != 0 {
			return errors.NewBadRequest().
				SetCode("InvalidLabelValue").
				SetMessage(`invalid label value "%v:%s"`, value, strings.Join(errs, "; "))
		}
	}
	return nil
}

func ValidateNodePortNum(nodePort int32) error {
	if nodePort != 0 && (nodePort < NodePortMin || nodePort > NodePortMax) {
		return errors.NewBadRequest().
			SetCode("InvalidNodePort").
			SetMessage(`invalid node port "%v", acceptable range is [%v, %v], or 0`, nodePort, NodePortMin, NodePortMax)
	}
	return nil
}

func ValidatePortNum(port int32) error {
	if port < 1 || port > 65535 {
		return errors.NewBadRequest().
			SetCode("InvalidPort").
			SetMessage(`invalid port "%v", acceptable range is [1, 65535]`, port)
	}
	return nil
}
