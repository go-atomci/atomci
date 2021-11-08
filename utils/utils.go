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

package utils

import (
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

const (
	KC_RAND_KIND_NUM   = 0
	KC_RAND_KIND_LOWER = 1
	KC_RAND_KIND_UPPER = 2
	KC_RAND_KIND_ALL   = 3
)

func Krand(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all {
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}

func MakeToken() string {
	token := string(Krand(16, KC_RAND_KIND_ALL))
	return token
}

// Contains ..
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// IntContains ..
func IntContains(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// EnsureAbs prepends the WorkDir to the given path if it is not an absolute path.
func EnsureAbs(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(WorkDir(), path)
}

var (
	workDir     string
	workDirOnce sync.Once
)

// WorkDir returns the absolute path of work directory. It reads the value of envrionment
// variable GOGS_WORK_DIR. When not set, it uses the directory where the application's
// binary is located.
func WorkDir() string {
	workDirOnce.Do(func() {
		workDir = filepath.Dir(AppPath())
	})

	return workDir
}

var (
	appPath     string
	appPathOnce sync.Once
)

// AppPath returns the absolute path of the application's binary.
func AppPath() string {
	appPathOnce.Do(func() {
		var err error
		appPath, err = exec.LookPath(os.Args[0])
		if err != nil {
			panic("look executable path: " + err.Error())
		}
		appPath, err = filepath.Abs(appPath)
		if err != nil {
			panic("get absolute executable path: " + err.Error())
		}
	})

	return appPath
}

func PathsIsEqual(path1, path2 string) bool {
	return GetRootPath(path1) == GetRootPath(path2)
}

func GetRootPath(path string) string {
	if path == "" {
		return "/"
	}
	return path
}

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

const (
	VersionMinLen = 0
	VersionMaxLen = 63
)

// 正则表达式验证字符合法性
func Restricted(s, regdata string) bool {
	validName := regexp.MustCompile(`^` + regdata + `$`)
	legal := validName.MatchString(s)
	return legal
}
