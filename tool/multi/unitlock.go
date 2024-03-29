package multi

import "sync"

// UnitLock 支持锁定某个key，可以融合进一些复杂的逻辑中
type UnitLock[KeyType comparable] struct {
	m map[KeyType]*sync.RWMutex
	// 锁住map的操作
	mapLock sync.Mutex
	// 用来锁定所有数据的操作
	allLock sync.RWMutex
}

func NewUnitLock[KeyType comparable]() *UnitLock[KeyType] {
	return &UnitLock[KeyType]{m: map[KeyType]*sync.RWMutex{}}
}

func (u *UnitLock[KeyType]) Lock(key KeyType) {
	u.allLock.RLock()
	u.mapLock.Lock()
	if _, ok := u.m[key]; !ok {
		u.m[key] = &sync.RWMutex{}
	}
	u.mapLock.Unlock()
	u.m[key].Lock()
}

func (u *UnitLock[KeyType]) Unlock(key KeyType) {
	u.allLock.RUnlock()
	u.mapLock.Lock()
	u.m[key].Unlock()
	u.mapLock.Unlock()
}

func (u *UnitLock[KeyType]) RLock(key KeyType) {
	u.allLock.RLock()
	u.mapLock.Lock()
	if _, ok := u.m[key]; !ok {
		u.m[key] = &sync.RWMutex{}
	}
	u.mapLock.Unlock()
	u.m[key].RLock()
}

func (u *UnitLock[KeyType]) RUnlock(key KeyType) {
	u.allLock.RUnlock()
	u.mapLock.Lock()
	u.m[key].RUnlock()
	u.mapLock.Unlock()
}

func (u *UnitLock[KeyType]) SafeRun(key KeyType, f func()) {
	u.allLock.RLock()
	defer u.allLock.RUnlock()
	u.Lock(key)
	defer u.Unlock(key)
	f()
}

func (u *UnitLock[KeyType]) SafeRunR(key KeyType, f func()) {
	u.allLock.RLock()
	defer u.allLock.RUnlock()
	u.RLock(key)
	defer u.RUnlock(key)
	f()
}

func (u *UnitLock[KeyType]) LockAll() {
	u.allLock.Lock()
}

func (u *UnitLock[KeyType]) UnlockAll() {
	u.allLock.Unlock()
}

func (u *UnitLock[KeyType]) SafeRunAll(f func()) {
	u.LockAll()
	defer u.UnlockAll()
	f()
}
