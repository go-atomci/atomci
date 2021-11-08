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
	"sync"
)

type locker struct {
	lock  *sync.Mutex
	count int
}

type SyncLocker struct {
	lockerList map[string]*locker
	mapLocker  *sync.Mutex
}

func NewSyncLocker() *SyncLocker {
	lock := &SyncLocker{
		mapLocker: &sync.Mutex{},
	}
	lock.lockerList = make(map[string]*locker)

	return lock
}

func (l *SyncLocker) Lock(key string) {
	var tmp *locker
	l.mapLocker.Lock()
	if l.lockerList[key] == nil {
		tmp = &locker{
			lock:  &sync.Mutex{},
			count: 0,
		}
		l.lockerList[key] = tmp
	} else {
		tmp = l.lockerList[key]
	}
	l.lockerList[key].count++
	l.mapLocker.Unlock()
	tmp.lock.Lock()
}

func (l *SyncLocker) Unlock(key string) {
	if l.lockerList[key] == nil {
		return
	}
	var tmp *locker
	l.mapLocker.Lock()
	tmp = l.lockerList[key]
	l.lockerList[key].count--
	if l.lockerList[key].count == 0 {
		delete(l.lockerList, key)
	}
	l.mapLocker.Unlock()
	tmp.lock.Unlock()
}
