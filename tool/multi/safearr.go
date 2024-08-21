package multi

import "sync"

type SafeArr[ValueType any] struct {
	arr  []ValueType
	lock sync.RWMutex
}

func NewSafeArr[ValueType any](arr []ValueType) *SafeArr[ValueType] {
	return &SafeArr[ValueType]{arr: arr}
}

func (m *SafeArr[ValueType]) Append(value ValueType) {
	m.lock.Lock()
	m.arr = append(m.arr, value)
	m.lock.Unlock()
}

func (m *SafeArr[ValueType]) Get(index int) (value ValueType, ok bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if index < 0 || index >= len(m.arr) {
		return
	}
	value = m.arr[index]
	ok = true
	return
}

func (m *SafeArr[ValueType]) Set(index int, value ValueType) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if index < 0 || index >= len(m.arr) {
		return
	}
	m.arr[index] = value
}

func (m *SafeArr[ValueType]) Len() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return len(m.arr)
}

func (m *SafeArr[ValueType]) Delete(index int) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if index < 0 || index >= len(m.arr) {
		return
	}
	m.arr = append(m.arr[:index], m.arr[index+1:]...)
}

func (m *SafeArr[ValueType]) SafeUse(f func(arr []ValueType)) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	f(m.arr)
}

func (m *SafeArr[ValueType]) Range(f func(index int, value ValueType) bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for i, v := range m.arr {
		if !f(i, v) {
			break
		}
	}
}

func (m *SafeArr[ValueType]) Copy() []ValueType {
	m.lock.RLock()
	defer m.lock.RUnlock()
	arr := make([]ValueType, len(m.arr))
	copy(arr, m.arr)
	return arr
}
