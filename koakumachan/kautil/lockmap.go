package kautil

import "sync"

type LockMap[Key comparable, Value any] struct {
	m    map[Key]Value
	lock sync.RWMutex
}

func NewLockMap[Key comparable, Value any]() *LockMap[Key, Value] {
	return &LockMap[Key, Value]{
		m: make(map[Key]Value),
	}
}

func (lm *LockMap[Key, Value]) Set(key Key, value Value) {
	lm.lock.Lock()
	defer lm.lock.Unlock()
	lm.m[key] = value
}

func (lm *LockMap[Key, Value]) Get(key Key) (Value, bool) {
	lm.lock.RLock()
	defer lm.lock.RUnlock()
	v, ok := lm.m[key]
	return v, ok
}

func (lm *LockMap[Key, Value]) Delete(key Key) {
	lm.lock.Lock()
	defer lm.lock.Unlock()
	delete(lm.m, key)
}
