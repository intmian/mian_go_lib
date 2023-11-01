package xstorage

import (
	"errors"
	"github.com/intmian/mian_go_lib/tool/misc"
	"sync"
)

type Mgr struct {
	dbCore  IDBCore
	setting KeyValueSetting
	sync.RWMutex
	misc.InitTag
	kvMap map[string]*ValueUnit // 后面不放指针，避免影响gc，此为唯一数据，取出时取指针
	pool  sync.Pool
}

func NewMgr(setting KeyValueSetting) (*Mgr, error) {
	// 检查有效性
	if setting.SaveType != SqlLiteDB {
		return nil, errors.New("not support save type, only support sqlite")
	}
	// 如果使用sqlite存储，需要判断地址是否为空
	if setting.SaveType == SqlLiteDB && setting.DBAddr == "" {
		return nil, errors.New("sqlite db file addr is empty")
	}
	if setting.Property&UseCache == 0 && setting.Property&UseDB == 0 {
		return nil, errors.New("not use cache and not use db")
	}
	if setting.Property&FullInitLoad != 0 && (setting.Property&UseCache == 0 || setting.Property&UseDB == 0) {
		return nil, errors.New("not use cache or not use db and full init load")
	}
	mgr := &Mgr{
		setting: setting,
	}
	switch setting.SaveType {
	case SqlLiteDB:
		dbCore, err := NewSqliteCore(setting.DBAddr)
		if err != nil {
			return nil, errors.Join(errors.New("new sqlite core error"), err)
		}
		mgr.dbCore = dbCore
	}
	if setting.Property&UseCache != 0 {
		mgr.kvMap = make(map[string]*ValueUnit)
		mgr.pool.New = func() interface{} {
			return &ValueUnit{}
		}
	}
	if setting.Property&FullInitLoad != 0 {
		if setting.Property&UseDB != 0 {
			kvMap, err := mgr.dbCore.GetAll()
			if err != nil {
				return nil, errors.Join(errors.New("get all value error"), err)
			}
			for key, valueUnit := range kvMap {
				if setting.Property&UseCache != 0 {
					err := mgr.RecordToMap(key, valueUnit)
					if err != nil {
						return nil, errors.Join(errors.New("record to map error"), err)
					}
				}
			}
		} else {
			return nil, errors.New("not use db and full init load")
		}
	}
	mgr.SetInitialized()
	return mgr, nil
}

func (m *Mgr) Get(key string) (bool, *ValueUnit, error) {
	if !m.IsInitialized() {
		return false, nil, errors.New("mgr not init")
	}
	if m.setting.Property&MultiSafe != 0 {
		m.RLock()
		defer m.RUnlock()
	}
	if m.setting.Property&UseCache != 0 {
		if valueUnit, ok := m.kvMap[key]; ok {
			if valueUnit.dirty {
				return false, nil, errors.New("value is dirty")
			}
			return true, valueUnit, nil
		}
	}
	if m.setting.Property&UseDB != 0 {
		ok, valueUnit, err := m.dbCore.Get(key)
		if err != nil {
			return false, nil, errors.Join(errors.New("get value error"), err)
		}
		if !ok {
			return false, nil, nil
		}
		if m.setting.Property&UseCache != 0 {
			err := m.RecordToMap(key, valueUnit)
			if err != nil {
				return false, nil, errors.Join(errors.New("record to map error"), err)
			}
		}
		return true, valueUnit, nil
	}
	return false, nil, errors.New("not use cache and not use db")
}

func (m *Mgr) Set(key string, value *ValueUnit) error {
	if !m.IsInitialized() {
		return errors.New("mgr not init")
	}
	if key == "" {
		return errors.New("key is empty")
	}
	if value == nil {
		return errors.New("value is nil")
	}
	if m.setting.Property&MultiSafe != 0 {
		m.Lock()
		defer m.Unlock()
	}
	if m.setting.Property&UseCache != 0 {
		err := m.RecordToMap(key, value)
		if err != nil {
			return errors.Join(errors.New("record to map error"), err)
		}
	}
	if m.setting.Property&UseDB != 0 {
		err := m.dbCore.Set(key, value)
		if err != nil {
			return errors.Join(errors.New("set value error"), err)
		}
	}
	return nil
}

func (m *Mgr) SetAsyncDB(key string, value *ValueUnit) (error, chan error) {
	if !m.IsInitialized() {
		return errors.New("mgr not init"), nil
	}
	if key == "" {
		return errors.New("key is empty"), nil
	}
	if value == nil {
		return errors.New("value is nil"), nil
	}
	if m.setting.Property&UseDB == 0 {
		return errors.New("not use db"), nil
	}
	if m.setting.Property&MultiSafe != 0 {
		m.Lock()
		defer m.Unlock()
	}
	if m.setting.Property&UseCache != 0 {
		err := m.RecordToMap(key, value)
		if err != nil {
			return errors.Join(errors.New("record to map error"), err), nil
		}
	}
	errChan := make(chan error)
	go func() {
		err := m.dbCore.Set(key, value)
		if err != nil {
			errChan <- errors.Join(errors.New("set value error"), err)
		}
		errChan <- nil
	}()
	return nil, errChan
}

func (m *Mgr) Delete(key string) error {
	if !m.IsInitialized() {
		return errors.New("mgr not init")
	}
	if key == "" {
		return errors.New("key is empty")
	}
	if m.setting.Property&UseCache != 0 {
		err := m.RemoveFromMap(key)
		if err != nil {
			return errors.Join(errors.New("remove from map error"), err)
		}
	}
	if m.setting.Property&UseDB != 0 {
		err := m.dbCore.Delete(key)
		if err != nil {
			return errors.Join(errors.New("delete value error"), err)
		}
	}
	return nil
}

func (m *Mgr) RecordToMap(key string, value *ValueUnit) error {
	if !m.IsInitialized() {
		return errors.New("mgr not init")
	}
	if key == "" {
		return errors.New("key is empty")
	}
	if value == nil {
		return errors.New("value is nil")
	}
	if m.setting.Property&UseCache == 0 {
		return errors.New("not use cache")
	}
	newValue, ok := m.pool.Get().(*ValueUnit)
	if !ok {
		return errors.New("pool type error")
	}
	*newValue = *value
	m.kvMap[key] = newValue
	return nil
}

func (m *Mgr) RemoveFromMap(key string) error {
	if !m.IsInitialized() {
		return errors.New("mgr not init")
	}
	if key == "" {
		return errors.New("key is empty")
	}
	if m.setting.Property&UseCache == 0 {
		return errors.New("not use cache")
	}
	// 释放
	m.pool.Put(m.kvMap[key])
	delete(m.kvMap, key)
	return nil
}
