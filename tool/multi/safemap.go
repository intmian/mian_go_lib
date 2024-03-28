package multi

import "sync"

// SafeMap 是对sync.Map 的封装，支持泛型，避免繁琐的类型转换
type SafeMap[KeyType any, ValueType any] struct {
	lock sync.Map
}

func (m *SafeMap[KeyType, ValueType]) Load(key any) (value ValueType, ok bool) {
	v, ok := m.lock.Load(key)
	if !ok {
		return
	}
	value, ok = v.(ValueType)
	return
}

func (m *SafeMap[KeyType, ValueType]) Store(key any, value ValueType) {
	m.lock.Store(key, value)
}

func (m *SafeMap[KeyType, ValueType]) Delete(key any) {
	m.lock.Delete(key)
}

func (m *SafeMap[KeyType, ValueType]) Range(f func(key any, value ValueType) bool) {
	m.lock.Range(func(key, value interface{}) bool {
		return f(key, value.(ValueType))
	})
}
