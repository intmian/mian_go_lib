package multi

import "sync"

// UnitLock 支持锁定某个key，可以融合进一些复杂的逻辑中
type UnitLock[KeyType comparable] struct {
	m map[KeyType]*sync.Mutex
	l sync.Mutex
}

func NewUnitLock[KeyType comparable]() *UnitLock[KeyType] {
	return &UnitLock[KeyType]{m: make(map[KeyType]*sync.Mutex)}
}

func (u *UnitLock[KeyType]) Lock(key KeyType) {
	u.l.Lock()
	if _, ok := u.m[key]; !ok {
		u.m[key] = &sync.Mutex{}
	}
	u.l.Unlock()
	u.m[key].Lock()
}

func (u *UnitLock[KeyType]) Unlock(key KeyType) {
	u.m[key].Unlock()
}

func (u *UnitLock[KeyType]) SafeRun(key KeyType, f func()) {
	u.Lock(key)
	defer u.Unlock(key)
	f()
}
