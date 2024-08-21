package multi

import "sync"

// SafeMap 是对sync.Map 的封装，支持泛型，避免繁琐的类型转换
type SafeMap[KeyType comparable, ValueType any] struct {
	lock sync.Map
}

func NewSafeMap[KeyType comparable, ValueType any](m map[KeyType]ValueType) *SafeMap[KeyType, ValueType] {
	sm := &SafeMap[KeyType, ValueType]{}
	for k, v := range m {
		sm.Store(k, v)
	}
	return sm
}

func (m *SafeMap[KeyType, ValueType]) Load(key KeyType) (value ValueType, ok bool) {
	v, ok := m.lock.Load(key)
	if !ok {
		return
	}
	value, ok = v.(ValueType)
	return
}

func (m *SafeMap[KeyType, ValueType]) Store(key KeyType, value ValueType) {
	m.lock.Store(key, value)
}

func (m *SafeMap[KeyType, ValueType]) Delete(key KeyType) {
	m.lock.Delete(key)
}

func (m *SafeMap[KeyType, ValueType]) Range(f func(key KeyType, value ValueType) bool) {
	m.lock.Range(func(key, value interface{}) bool {
		return f(key.(KeyType), value.(ValueType))
	})
}
