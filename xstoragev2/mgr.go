package xstorage

import (
	"errors"
	"github.com/intmian/mian_go_lib/tool/misc"
	"sync"
)

type XStorage struct {
	dbCore   IDBCore
	fileCore IFileCore
	setting  XstorageSetting
	rwLock   sync.RWMutex
	initTag  misc.InitTag
	//map尽量不要包非pool指针，不然可能在频繁调用的情况下出现大量的内存垃圾，影响内存，gc也无法快速回收，如果低峰期依然有访问可能会出现同访问量、数据量的情况下，每天内存占用越来越高，直到内存耗尽才频繁gc，性能会有问题，特别是在单机多进程的情况下。
	kvMap map[string]*ValueUnit // 后面不放指针，避免影响gc，此为唯一数据，取出时取指针
	pool  sync.Pool
}

func (m *XStorage) Init(setting XstorageSetting) error {
	// 检查路径
	if setting.SaveType > DBBegin && setting.SaveType < FileBegin && setting.DBAddr == "" {
		return ErrSqliteDBFileAddrEmpty
	}
	if setting.SaveType > FileBegin && setting.FileAddr == "" {
		return ErrSqliteDBFileAddrEmpty
	}

	if !misc.HasOneProperty(setting.Property, UseCache, UseDisk) {
		return ErrNotUseCacheAndNotUseDb
	}
	if misc.HasProperty(setting.Property, FullInitLoad) && !misc.HasProperty(setting.Property, UseCache, UseDisk) {
		return ErrNotUseCacheOrNotUseDbAndFullInitLoad
	}
	if setting.SaveType == Toml && !misc.HasProperty(setting.Property, UseCache, FullInitLoad) {
		return ErrUseJsonButNotUseCacheAndNotFullInitLoad
	}
	switch setting.SaveType {
	case SqlLiteDB:
		dbCore, err := NewSqliteCore(setting.DBAddr)
		if err != nil {
			return errors.Join(ErrNewSqliteCore, err)
		}
		m.dbCore = dbCore
	case Toml:
		fileCore := NewTomlCore(setting.FileAddr)
		m.fileCore = fileCore
	}
	if misc.HasProperty(setting.Property, UseCache) {
		m.kvMap = make(map[string]*ValueUnit)
		m.pool.New = func() interface{} {
			return &ValueUnit{}
		}
	}
	if misc.HasProperty(setting.Property, FullInitLoad) {
		if misc.HasProperty(setting.Property, UseDisk) {
			kvMap, err := m.FromDiskGetAll()
			if err != nil {
				return errors.Join(ErrGetAllValue, err)
			}
			for key, valueUnit := range kvMap {
				if misc.HasProperty(setting.Property, UseCache) {
					err := m.recordToMap(key, valueUnit)
					if err != nil {
						return errors.Join(ErrRecordToMap, err)
					}
				}
			}
		} else {
			return ErrNotUseDbAndFullInitLoad
		}
	}
	m.initTag.SetInitialized()
	return nil
}

func NewXStorage(setting XstorageSetting) (*XStorage, error) {
	mgr := &XStorage{
		setting: setting,
	}
	err := mgr.Init(setting)
	if err != nil {
		return nil, err
	}
	return mgr, nil
}

func (m *XStorage) FromDiskGetAll() (map[string]*ValueUnit, error) {
	t := m.setting.SaveType
	var kvMap map[string]*ValueUnit
	var err error
	switch {
	case t > DBBegin && t < FileBegin:
		kvMap, err = m.dbCore.GetAll()
	case t > FileBegin:
		kvMap, err = m.fileCore.GetAll()
	}
	return kvMap, err
}

func (m *XStorage) Get(key string) (*ValueUnit, error) {
	var valueUnit ValueUnit
	ok, err := m.GetHP(key, &valueUnit)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &valueUnit, nil
}

func (m *XStorage) GetHP(key string, valueUnit *ValueUnit) (bool, error) {
	if !m.initTag.IsInitialized() {
		return false, ErrMgrNotInit
	}
	if valueUnit == nil {
		return false, ErrValueUnitIsNil
	}
	valueUnit.Reset()
	if misc.HasProperty(m.setting.Property, MultiSafe) {
		m.rwLock.RLock()
		defer m.rwLock.RUnlock()
	}
	if misc.HasProperty(m.setting.Property, UseCache) {
		if p, ok := m.kvMap[key]; ok {
			if valueUnit.dirty {
				return false, ErrValueIsDirty
			}
			*valueUnit = *p
			return true, nil
		}
	}
	if misc.HasProperty(m.setting.Property, UseDisk) {
		ok, err := m.onGetFromDisk(key, valueUnit)
		if err != nil {
			return false, errors.Join(ErrSqliteDBFileAddrNotExist, err)
		}
		if !ok {
			return false, nil
		}
		if misc.HasProperty(m.setting.Property, UseCache) {
			err := m.recordToMap(key, valueUnit)
			if err != nil {
				return false, errors.Join(ErrRecordToMap, err)
			}
		}
		return true, nil
	}
	if !misc.HasProperty(m.setting.Property, UseCache) && !misc.HasProperty(m.setting.Property, UseDisk) {
		return false, ErrNotUseCacheAndNotUseDb
	}
	return false, nil
}

func (m *XStorage) onGetFromDisk(key string, rec *ValueUnit) (bool, error) {
	t := m.setting.SaveType
	var ok bool
	var err error
	switch {
	case t > DBBegin && t < FileBegin:
		ok, err = m.dbCore.Get(key, rec)
	case t > FileBegin:
		return false, nil
	}
	return ok, err
}

func (m *XStorage) GetAndSetDefault(key string, defaultValue *ValueUnit) (*ValueUnit, error) {
	var valueUnit ValueUnit
	ok, err := m.GetHP(key, &valueUnit)
	if err != nil {
		return nil, err
	}
	if ok {
		return &valueUnit, nil
	} else {
		return defaultValue, nil
	}
}

// GetAndSetDefaultHP get值，如果没有就设置并返回默认值，返回 是否setdefault数据，数据，错误
func (m *XStorage) GetAndSetDefaultHP(key string, defaultValue *ValueUnit, rec *ValueUnit) (bool, error) {
	if !m.initTag.IsInitialized() {
		return false, ErrMgrNotInit
	}
	ok, err := m.GetHP(key, rec)
	if err != nil {
		return false, errors.Join(ErrGet, err)
	}
	if ok {
		return true, nil
	}
	*rec = *defaultValue
	err = m.Set(key, defaultValue)
	if err != nil {
		return false, errors.Join(ErrSetValue, err)
	}
	return true, nil

}

// GetAndSetDefaultAsync get值，如果没有就设置并返回默认值，返回 是否setdefault数据，数据，错误
func (m *XStorage) GetAndSetDefaultAsync(key string, defaultValue *ValueUnit, rec *ValueUnit) (bool, error, chan error) {
	if !m.initTag.IsInitialized() {
		return false, ErrMgrNotInit, nil
	}
	ok, err := m.GetHP(key, rec)
	if err != nil {
		return false, errors.Join(ErrSqliteDBFileAddrNotExist, err), nil
	}
	if ok {
		return true, nil, nil
	}
	*rec = *defaultValue
	err, c := m.SetAsync(key, defaultValue)
	if err != nil {
		return false, errors.Join(ErrSetValue, err), nil
	}
	return true, nil, c
}

// SetDefault 设置默认值，如果已经存在则不设置
func (m *XStorage) SetDefault(key string, defaultValue *ValueUnit) error {
	if !m.initTag.IsInitialized() {
		return ErrMgrNotInit
	}
	if key == "" {
		return ErrKeyIsEmpty
	}
	if defaultValue == nil {
		return ErrValueIsNil
	}
	var Rec ValueUnit
	ok, err := m.GetHP(key, &Rec)
	if err != nil {
		return errors.Join(ErrGet, err)
	}
	if ok {
		return nil
	}
	err = m.Set(key, defaultValue)
	if err != nil {
		return errors.Join(ErrSetValue, err)
	}
	return nil
}

// SetDefaultAsync 设置默认值，如果已经存在则不设置
func (m *XStorage) SetDefaultAsync(key string, defaultValue *ValueUnit) (error, chan error) {
	if !m.initTag.IsInitialized() {
		return ErrMgrNotInit, nil
	}
	if key == "" {
		return ErrKeyIsEmpty, nil
	}
	if defaultValue == nil {
		return ErrValueIsNil, nil
	}
	var Rec ValueUnit
	ok, err := m.GetHP(key, &Rec)
	if err != nil {
		return errors.Join(ErrGet, err), nil
	}
	if ok {
		return nil, nil
	}
	err, c := m.SetAsync(key, defaultValue)
	if err != nil {
		return errors.Join(ErrSetValue, err), nil
	}
	return nil, c
}

func (m *XStorage) Set(key string, value *ValueUnit) error {
	if !m.initTag.IsInitialized() {
		return ErrMgrNotInit
	}
	if key == "" {
		return ErrKeyIsEmpty
	}
	if value == nil {
		return ErrValueIsNil
	}
	if misc.HasProperty(m.setting.Property, MultiSafe) {
		m.rwLock.Lock()
		defer m.rwLock.Unlock()
	}
	if misc.HasProperty(m.setting.Property, UseCache) {
		err := m.recordToMap(key, value)
		if err != nil {
			return errors.Join(ErrRecordToMap, err)
		}
	}
	if misc.HasProperty(m.setting.Property, UseDisk) {
		err := m.onSave2Disk(key, value)
		if err != nil {
			return errors.Join(ErrSetValue, err)
		}
	}
	return nil
}

func (m *XStorage) onSave2Disk(key string, value *ValueUnit) error {
	t := m.setting.SaveType
	var err error
	switch {
	case t > DBBegin && t < FileBegin:
		err = m.dbCore.Set(key, value)
	case t > FileBegin:
		err = m.fileCore.SaveAll(m.kvMap)
	}
	return err
}

func (m *XStorage) SetAsync(key string, value *ValueUnit) (error, chan error) {
	if !m.initTag.IsInitialized() {
		return ErrMgrNotInit, nil
	}
	if key == "" {
		return ErrKeyIsEmpty, nil
	}
	if value == nil {
		return ErrValueIsNil, nil
	}
	if misc.HasProperty(m.setting.Property, MultiSafe) {
		m.rwLock.Lock()
		defer m.rwLock.Unlock()
	}
	if misc.HasProperty(m.setting.Property, UseCache) {
		err := m.recordToMap(key, value)
		if err != nil {
			return errors.Join(ErrRecordToMap, err), nil
		}
	}
	if misc.HasProperty(m.setting.Property, UseDisk) {
		errChan := make(chan error)
		go func() {
			err := m.onSave2Disk(key, value)
			if err != nil {
				errChan <- errors.Join(ErrSetValue, err)
			}
			errChan <- nil
		}()
		return nil, errChan
	} else {
		return nil, nil
	}
}

func (m *XStorage) Delete(key string) error {
	if !m.initTag.IsInitialized() {
		return ErrMgrNotInit
	}
	if key == "" {
		return ErrKeyIsEmpty
	}
	if misc.HasProperty(m.setting.Property, UseCache) {
		err := m.removeFromMap(key)
		if err != nil {
			return errors.Join(ErrRemoveFromMap, err)
		}
	}
	if misc.HasProperty(m.setting.Property, UseDisk) {
		err := m.fromDiskDelete(key)
		if err != nil {
			return errors.Join(ErrDeleteValue, err)
		}
	}
	return nil
}

func (m *XStorage) fromDiskDelete(key string) error {
	t := m.setting.SaveType
	var err error
	switch {
	case t > DBBegin && t < FileBegin:
		err = m.dbCore.Delete(key)
	case t > FileBegin:
		err = m.fileCore.SaveAll(m.kvMap)
	}
	return err
}

func (m *XStorage) recordToMap(key string, value *ValueUnit) error {
	if key == "" {
		return ErrKeyIsEmpty
	}
	if value == nil {
		return ErrValueIsNil
	}
	if !misc.HasProperty(m.setting.Property, UseCache) {
		return ErrNotUseCache
	}
	newValue, ok := m.pool.Get().(*ValueUnit)
	if !ok {
		return ErrPoolType
	}
	Copy(value, newValue)
	m.kvMap[key] = newValue
	return nil
}

func (m *XStorage) removeFromMap(key string) error {
	if key == "" {
		return ErrKeyIsEmpty
	}
	if !misc.HasProperty(m.setting.Property, UseCache) {
		return ErrNotUseCache
	}
	// 释放
	if _, ok := m.kvMap[key]; !ok {
		return ErrKeyNotExist
	}
	m.kvMap[key].Reset()
	m.pool.Put(m.kvMap[key])
	delete(m.kvMap, key)
	return nil
}

func (m *XStorage) GetAll() (map[string]*ValueUnit, error) {
	if !m.initTag.IsInitialized() {
		return nil, ErrMgrNotInit
	}
	if misc.HasProperty(m.setting.Property, MultiSafe) {
		m.rwLock.RLock()
		defer m.rwLock.RUnlock()
	}
	if misc.HasProperty(m.setting.Property, UseCache) {
		return m.kvMap, nil
	}
	if misc.HasProperty(m.setting.Property, UseDisk) {
		kvMap, err := m.FromDiskGetAll()
		if err != nil {
			return nil, errors.Join(ErrGetAllValue, err)
		}
		return kvMap, nil
	}
	return nil, ErrNotUseCacheAndNotUseDb
}
