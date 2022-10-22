package syncmap

import (
	"fmt"
	"sync"
)

type SyncMap struct {
	sync.RWMutex
	m map[string]interface{}
}

func NewSyncMap() *SyncMap {
	return &SyncMap{
		m: make(map[string]interface{}),
	}
}

func (sm *SyncMap) Lenght() int {
	return len(sm.m)
}

func (sm *SyncMap) Get(k string) interface{} {
	sm.RLock()
	defer sm.RUnlock()
	if val, ok := sm.m[k]; ok {
		return val
	}
	return nil
}

func (sm *SyncMap) Set(k string, v interface{}) bool {
	sm.Lock()
	defer sm.Unlock()
	if val, ok := sm.m[k]; !ok {
		sm.m[k] = v
	} else if val != v {
		sm.m[k] = v
	} else {
		return false
	}
	return true
}

//原子操作,请谨慎重入锁
func (sm *SyncMap) Atomic(atomic func()) {
	sm.Lock()
	defer sm.Unlock()
	atomic()
}

// key是否存在
func (sm *SyncMap) Check(k string) bool {
	sm.RLock()
	defer sm.RUnlock()
	if _, ok := sm.m[k]; !ok {
		return false
	}
	return true
}

func (sm *SyncMap) Delete(k string) {
	sm.Lock()
	defer sm.Unlock()
	delete(sm.m, k)
}

func (sm *SyncMap) GetAll() ([]string, []string) {
	keys := []string{}
	vals := []string{}
	sm.RLock()
	defer sm.RUnlock()
	for k, v := range sm.m {
		keys = append(keys, k)
		vals = append(vals, fmt.Sprintf("%v", v))
	}
	return keys, vals
}

//getset为对key对应的value的操作函数，例如：func(v uint8)uint8{return v++},则key对应的v值+1
func (sm *SyncMap) GetAndSetUint8(key string, getSet func(oldval uint8) (newval uint8)) (newVal uint8) {
	sm.Lock()
	defer sm.Unlock()
	v, ok := sm.m[key].(uint8)
	if !ok {
		v = 0
	}
	newval := getSet(v)
	sm.m[key] = newval
	return newVal
}
func (sm *SyncMap) GetAndSetUint16(key string, getSet func(oldval uint16) (newval uint16)) (newVal uint16) {
	sm.Lock()
	defer sm.Unlock()
	v, ok := sm.m[key].(uint16)
	if !ok {
		v = 0
	}
	newval := getSet(v)
	sm.m[key] = newval
	return newVal
}
