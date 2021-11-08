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

package dao

import (
	"sync"

	"github.com/astaxie/beego/orm"
)

var globalOrm orm.Ormer
var once sync.Once

// GetOrmer :set ormer singleton
func GetOrmer() orm.Ormer {
	once.Do(func() {
		globalOrm = orm.NewOrm()
	})
	return globalOrm
}

// Transactional invoke lambda function within transaction
func Transactional(ormer orm.Ormer, handle func() error) (err error) {
	err = ormer.Begin()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			ormer.Rollback()
			panic(p)
		} else if err != nil {
			ormer.Rollback()
		} else {
			err = ormer.Commit()
		}
	}()
	err = handle()
	return
}
