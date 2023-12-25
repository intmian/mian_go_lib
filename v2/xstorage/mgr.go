package xstorage

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/intmian/mian_go_lib/tool/misc"
	"github.com/intmian/mian_go_lib/tool/xlog"
	"sync"
)

type Mgr struct {
	dbCore   IDBCore
	fileCore IFileCore
	setting  KeyValueSetting
	rwLock   sync.RWMutex
	initTag  misc.InitTag
	//map尽量不要包非pool指针，不然可能在频繁调用的情况下出现大量的内存垃圾，影响内存，gc也无法快速回收，如果低峰期依然有访问可能会出现同访问量、数据量的情况下，每天内存占用越来越高，直到内存耗尽才频繁gc，性能会有问题，特别是在单机多进程的情况下。
	kvMap     map[string]*ValueUnit // 后面不放指针，避免影响gc，此为唯一数据，取出时取指针
	pool      sync.Pool
	ginEngine *gin.Engine
	log       *xlog.Mgr
	logFrom   string
}

func NewMgr(setting KeyValueSetting) (*Mgr, error) {
	// 检查路径
	if setting.SaveType > DBBegin && setting.SaveType < FileBegin && setting.DBAddr == "" {
		return nil, SqliteDBFileAddrEmptyErr
	}
	if setting.SaveType > FileBegin && setting.FileAddr == "" {
		return nil, SqliteDBFileAddrEmptyErr
	}

	if !misc.HasOneProperty(setting.Property, UseCache, UseDisk) {
		return nil, NotUseCacheAndNotUseDbErr
	}
	if misc.HasProperty(setting.Property, FullInitLoad) && !misc.HasProperty(setting.Property, UseCache, UseDisk) {
		return nil, NotUseCacheOrNotUseDbAndFullInitLoadErr
	}
	if setting.SaveType == Toml && !misc.HasProperty(setting.Property, UseCache, FullInitLoad) {
		return nil, UseJsonButNotUseCacheAndNotFullInitLoadErr
	}
	mgr := &Mgr{
		setting: setting,
	}
	switch setting.SaveType {
	case SqlLiteDB:
		dbCore, err := NewSqliteCore(setting.DBAddr)
		if err != nil {
			return nil, errors.Join(NewSqliteCoreErr, err)
		}
		mgr.dbCore = dbCore
	case Toml:
		fileCore := NewTomlCore(setting.FileAddr)
		mgr.fileCore = fileCore
	}
	if misc.HasProperty(setting.Property, UseCache) {
		mgr.kvMap = make(map[string]*ValueUnit)
		mgr.pool.New = func() interface{} {
			return &ValueUnit{}
		}
	}
	if misc.HasProperty(setting.Property, FullInitLoad) {
		if misc.HasProperty(setting.Property, UseDisk) {
			kvMap, err := mgr.FromDiskGetAll()
			if err != nil {
				return nil, errors.Join(GetAllValueErr, err)
			}
			for key, valueUnit := range kvMap {
				if misc.HasProperty(setting.Property, UseCache) {
					err := mgr.recordToMap(key, valueUnit)
					if err != nil {
						return nil, errors.Join(RecordToMapErr, err)
					}
				}
			}
		} else {
			return nil, NotUseDbAndFullInitLoadErr
		}
	}
	mgr.initTag.SetInitialized()
	return mgr, nil
}

func (m *Mgr) FromDiskGetAll() (map[string]*ValueUnit, error) {
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

func (m *Mgr) Get(key string, valueUnit *ValueUnit) (bool, error) {
	if !m.initTag.IsInitialized() {
		return false, MgrNotInitErr
	}
	if valueUnit == nil {
		return false, ValueUnitIsNilErr
	}
	valueUnit.Reset()
	if misc.HasProperty(m.setting.Property, MultiSafe) {
		m.rwLock.RLock()
		defer m.rwLock.RUnlock()
	}
	if misc.HasProperty(m.setting.Property, UseCache) {
		if p, ok := m.kvMap[key]; ok {
			if valueUnit.dirty {
				return false, ValueIsDirtyErr
			}
			*valueUnit = *p
			return true, nil
		}
	}
	if misc.HasProperty(m.setting.Property, UseDisk) {
		ok, err := m.OnGetFromDisk(key, valueUnit)
		if err != nil {
			return false, errors.Join(SqliteDBFileAddrNotExistErr, err)
		}
		if !ok {
			return false, nil
		}
		if misc.HasProperty(m.setting.Property, UseCache) {
			err := m.recordToMap(key, valueUnit)
			if err != nil {
				return false, errors.Join(RecordToMapErr, err)
			}
		}
		return true, nil
	}
	if !misc.HasProperty(m.setting.Property, UseCache) && !misc.HasProperty(m.setting.Property, UseDisk) {
		return false, NotUseCacheAndNotUseDbErr
	}
	return false, nil
}

// Get 从mgr中提取数据并转换为对应类型。使用方便，但是性能逊色于GetHp
// 之所以保留这个函数，是因为相较于ValueUnit，对应的基础类型会小一点，直接返回值也还行
// 相较于库的内部接口数量有限，方便优化，外部引用的数量可能是无限的，所以给外部暴露这个接口
func Get[T IValueType](mgr *Mgr, key string) (bool, T, error) {
	var valueUnit ValueUnit
	var t T
	ok, err := mgr.Get(key, &valueUnit)
	if err != nil {
		return false, t, err
	}
	if !ok {
		return false, t, nil
	}
	t = ToBase[T](&valueUnit)
	return true, t, nil
}

// GetHp 从mgr中提取数据并转换为对应类型。性能优于Get
func GetHp[T IValueType](mgr *Mgr, key string, rec *T) (bool, error) {
	if rec == nil {
		return false, RecIsNilErr
	}
	var valueUnit ValueUnit
	ok, err := mgr.Get(key, &valueUnit)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	*rec = ToBase[T](&valueUnit)
	return true, nil
}

func (m *Mgr) OnGetFromDisk(key string, rec *ValueUnit) (bool, error) {
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

// GetAndSetDefault get值，如果没有就设置并返回默认值，返回 是否setdefault数据，数据，错误
func (m *Mgr) GetAndSetDefault(key string, defaultValue *ValueUnit, rec *ValueUnit) (bool, error) {
	if !m.initTag.IsInitialized() {
		return false, MgrNotInitErr
	}
	ok, err := m.Get(key, rec)
	if err != nil {
		return false, errors.Join(GetErr, err)
	}
	if ok {
		return true, nil
	}
	*rec = *defaultValue
	err = m.Set(key, defaultValue)
	if err != nil {
		return false, errors.Join(SetValueErr, err)
	}
	return true, nil

}

// GetAndSetDefaultAsync get值，如果没有就设置并返回默认值，返回 是否setdefault数据，数据，错误
func (m *Mgr) GetAndSetDefaultAsync(key string, defaultValue *ValueUnit, rec *ValueUnit) (bool, error, chan error) {
	if !m.initTag.IsInitialized() {
		return false, MgrNotInitErr, nil
	}
	ok, err := m.Get(key, rec)
	if err != nil {
		return false, errors.Join(SqliteDBFileAddrNotExistErr, err), nil
	}
	if ok {
		return true, nil, nil
	}
	*rec = *defaultValue
	err, c := m.SetAsync(key, defaultValue)
	if err != nil {
		return false, errors.Join(SetValueErr, err), nil
	}
	return true, nil, c
}

// SetDefault 设置默认值，如果已经存在则不设置
func (m *Mgr) SetDefault(key string, defaultValue *ValueUnit) error {
	if !m.initTag.IsInitialized() {
		return MgrNotInitErr
	}
	if key == "" {
		return KeyIsEmptyErr
	}
	if defaultValue == nil {
		return ValueIsNilErr
	}
	var Rec ValueUnit
	ok, err := m.Get(key, &Rec)
	if err != nil {
		return errors.Join(GetErr, err)
	}
	if ok {
		return nil
	}
	err = m.Set(key, defaultValue)
	if err != nil {
		return errors.Join(SetValueErr, err)
	}
	return nil
}

// SetDefaultAsync 设置默认值，如果已经存在则不设置
func (m *Mgr) SetDefaultAsync(key string, defaultValue *ValueUnit) (error, chan error) {
	if !m.initTag.IsInitialized() {
		return MgrNotInitErr, nil
	}
	if key == "" {
		return KeyIsEmptyErr, nil
	}
	if defaultValue == nil {
		return ValueIsNilErr, nil
	}
	var Rec ValueUnit
	ok, err := m.Get(key, &Rec)
	if err != nil {
		return errors.Join(GetErr, err), nil
	}
	if ok {
		return nil, nil
	}
	err, c := m.SetAsync(key, defaultValue)
	if err != nil {
		return errors.Join(SetValueErr, err), nil
	}
	return nil, c
}

func (m *Mgr) Set(key string, value *ValueUnit) error {
	if !m.initTag.IsInitialized() {
		return MgrNotInitErr
	}
	if key == "" {
		return KeyIsEmptyErr
	}
	if value == nil {
		return ValueIsNilErr
	}
	if misc.HasProperty(m.setting.Property, MultiSafe) {
		m.rwLock.Lock()
		defer m.rwLock.Unlock()
	}
	if misc.HasProperty(m.setting.Property, UseCache) {
		err := m.recordToMap(key, value)
		if err != nil {
			return errors.Join(RecordToMapErr, err)
		}
	}
	if misc.HasProperty(m.setting.Property, UseDisk) {
		err := m.OnSave2Disk(key, value)
		if err != nil {
			return errors.Join(SetValueErr, err)
		}
	}
	return nil
}

func (m *Mgr) OnSave2Disk(key string, value *ValueUnit) error {
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

func (m *Mgr) SetAsync(key string, value *ValueUnit) (error, chan error) {
	if !m.initTag.IsInitialized() {
		return MgrNotInitErr, nil
	}
	if key == "" {
		return KeyIsEmptyErr, nil
	}
	if value == nil {
		return ValueIsNilErr, nil
	}
	if misc.HasProperty(m.setting.Property, MultiSafe) {
		m.rwLock.Lock()
		defer m.rwLock.Unlock()
	}
	if misc.HasProperty(m.setting.Property, UseCache) {
		err := m.recordToMap(key, value)
		if err != nil {
			return errors.Join(RecordToMapErr, err), nil
		}
	}
	if misc.HasProperty(m.setting.Property, UseDisk) {
		errChan := make(chan error)
		go func() {
			err := m.OnSave2Disk(key, value)
			if err != nil {
				errChan <- errors.Join(SetValueErr, err)
			}
			errChan <- nil
		}()
		return nil, errChan
	} else {
		return nil, nil
	}
}

func (m *Mgr) Delete(key string) error {
	if !m.initTag.IsInitialized() {
		return MgrNotInitErr
	}
	if key == "" {
		return KeyIsEmptyErr
	}
	if misc.HasProperty(m.setting.Property, UseCache) {
		err := m.removeFromMap(key)
		if err != nil {
			return errors.Join(RemoveFromMapErr, err)
		}
	}
	if misc.HasProperty(m.setting.Property, UseDisk) {
		err := m.FromDiskDelete(key)
		if err != nil {
			return errors.Join(DeleteValueErr, err)
		}
	}
	return nil
}

func (m *Mgr) FromDiskDelete(key string) error {
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

func (m *Mgr) recordToMap(key string, value *ValueUnit) error {
	if key == "" {
		return KeyIsEmptyErr
	}
	if value == nil {
		return ValueIsNilErr
	}
	if !misc.HasProperty(m.setting.Property, UseCache) {
		return NotUseCacheErr
	}
	newValue, ok := m.pool.Get().(*ValueUnit)
	if !ok {
		return PoolTypeErr
	}
	Copy(value, newValue)
	m.kvMap[key] = newValue
	return nil
}

func (m *Mgr) removeFromMap(key string) error {
	if key == "" {
		return KeyIsEmptyErr
	}
	if !misc.HasProperty(m.setting.Property, UseCache) {
		return NotUseCacheErr
	}
	// 释放
	if _, ok := m.kvMap[key]; !ok {
		return KeyNotExistErr
	}
	m.kvMap[key].Reset()
	m.pool.Put(m.kvMap[key])
	delete(m.kvMap, key)
	return nil
}

func (m *Mgr) GetAll() (map[string]*ValueUnit, error) {
	if !m.initTag.IsInitialized() {
		return nil, MgrNotInitErr
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
			return nil, errors.Join(GetAllValueErr, err)
		}
		return kvMap, nil
	}
	return nil, NotUseCacheAndNotUseDbErr
}

func (m *Mgr) Log(level xlog.TLogLevel, info string) {
	if !m.initTag.IsInitialized() {
		return
	}
	if m.log != nil {
		m.log.Log(level, m.logFrom, info)
	}
}

func (m *Mgr) Error(info string) {
	if !m.initTag.IsInitialized() {
		return
	}
	if m.log != nil {
		m.log.Log(xlog.EError, m.logFrom, info)
	}
}
