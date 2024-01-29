package xstorage

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/intmian/mian_go_lib/tool/misc"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/*
初始化时传入数据库地址和表名
列Key用来存储键
列ValueInt用来存储整数值
列ValueString用来存储字符串值
...
请注意假如value类型为slice，会被存储于key[0]、Key[1]、Key[2]...列中，Key[0]、Key[1]、Key[2]...列的值为value的每个元素
*/

type SqliteCore struct {
	db *gorm.DB
	misc.InitTag
	rwLock sync.RWMutex
}

type EmptyLogger struct {
}

func (e EmptyLogger) LogMode(level logger.LogLevel) logger.Interface {
	return e
}

func (e EmptyLogger) Info(ctx context.Context, s string, i ...interface{}) {}

func (e EmptyLogger) Warn(ctx context.Context, s string, i ...interface{}) {}

func (e EmptyLogger) Error(ctx context.Context, s string, i ...interface{}) {}

func (e EmptyLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
}

func NewSqliteCore(DbFileAddr string) (*SqliteCore, error) {
	// 依靠外层进行日志交互，为了避免本地打印日志过多，这里不使用日志库
	db, err := gorm.Open(sqlite.Open(DbFileAddr), &gorm.Config{
		Logger: EmptyLogger{},
	})
	if err != nil {
		return nil, errors.Join(ErrOpenSqlite, err)
	}
	err = db.AutoMigrate(&KeyValueModel{})
	if err != nil {
		return nil, errors.Join(ErrAutoMigrate, err)
	}
	sqliteCore := &SqliteCore{
		db: db,
	}
	sqliteCore.SetInitialized()
	return sqliteCore, nil
}

func (m *SqliteCore) Get(key string, rec *ValueUnit) (bool, error) {
	if rec == nil {
		return false, ErrRecIsNil
	}
	return m.getInner(key, rec, true)
}

func (m *SqliteCore) getInner(key string, rec *ValueUnit, needLock bool) (exist bool, retErr error) {
	if !m.IsInitialized() {
		return false, ErrSqliteCoreNotInit
	}
	if needLock {
		m.rwLock.RLock()
		defer m.rwLock.RUnlock()
	}

	var keyValueModel KeyValueModel
	result := m.db.Where("Key = ?", key).First(&keyValueModel)
	// 如果没有这个

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, errors.Join(ErrSqliteDBFileAddrNotExist, result.Error)
		}
	}
	sliceNum, err := sqliteModel2Data(keyValueModel, rec)
	if err != nil {
		return false, err
	}

	if sliceNum == 0 {
		return true, nil
	}

	// slice 的内容存放在 Key[0]、Key[1]、Key[2]...Key[sliceNum-1] 列中
	switch ValueType(keyValueModel.ValueType) {
	case ValueTypeSliceInt:
		rec.Data = make([]int, sliceNum)
		rec.Type = ValueTypeSliceInt
		for i := 0; i < sliceNum; i++ {
			var keyValueModel2 KeyValueModel
			result := m.db.Where("Key = ?", key+"["+strconv.Itoa(i)+"]").First(&keyValueModel2)
			if result.Error != nil {
				return false, errors.Join(ErrSqliteDBFileAddrNotExist, result.Error)
			}
			if keyValueModel2.ValueInt == nil {
				return false, fmt.Errorf("slice but ValueInt is nil, Key: %s[%d]", key, i)
			}
			rec.Data.([]int)[i] = *keyValueModel2.ValueInt
		}
	case ValueTypeSliceString:
		rec.Data = make([]string, sliceNum)
		rec.Type = ValueTypeSliceString
		for i := 0; i < sliceNum; i++ {
			var keyValueModel2 KeyValueModel
			result := m.db.Where("Key = ?", key+"["+strconv.Itoa(i)+"]").First(&keyValueModel2)
			if result.Error != nil {
				return false, errors.Join(ErrSqliteDBFileAddrNotExist, result.Error)
			}
			if keyValueModel2.ValueString == nil {
				return false, fmt.Errorf("slice but ValueString is nil, Key: %s[%d]", key, i)
			}
			rec.Data.([]string)[i] = *keyValueModel2.ValueString
		}
	case ValueTypeSliceFloat:
		rec.Data = make([]float32, sliceNum)
		rec.Type = ValueTypeSliceFloat
		for i := 0; i < sliceNum; i++ {
			var keyValueModel2 KeyValueModel
			result := m.db.Where("Key = ?", key+"["+strconv.Itoa(i)+"]").First(&keyValueModel2)
			if result.Error != nil {
				return false, errors.Join(ErrSqliteDBFileAddrNotExist, result.Error)
			}
			if keyValueModel2.ValueFloat == nil {
				return false, fmt.Errorf("slice but ValueFloat is nil, Key: %s[%d]", key, i)
			}
			rec.Data.([]float32)[i] = *keyValueModel2.ValueFloat
		}
	case ValueTypeSliceBool:
		rec.Data = make([]bool, sliceNum)
		rec.Type = ValueTypeSliceBool
		for i := 0; i < sliceNum; i++ {
			var keyValueModel2 KeyValueModel
			result := m.db.Where("Key = ?", key+"["+strconv.Itoa(i)+"]").First(&keyValueModel2)
			if result.Error != nil {
				return false, errors.Join(ErrSqliteDBFileAddrNotExist, result.Error)
			}
			if keyValueModel2.ValueInt == nil {
				return false, fmt.Errorf("slice but ValueInt is nil, Key: %s[%d]", key, i)
			}
			if *keyValueModel2.ValueInt == 0 {
				rec.Data.([]bool)[i] = false
			} else {
				rec.Data.([]bool)[i] = true
			}
		}
	default:
		return false, ErrValueType
	}
	return true, nil
}

// sqliteModel2Data 将从数据库取出来的model转化为ValueUnit，但是需要注意的是，如果是slice类型，只返回slice的长度，不返回具体的值
func sqliteModel2Data(keyValueModel KeyValueModel, rec *ValueUnit) (int, error) {
	value := &ValueUnit{}
	sliceNum := 0
	// 判断合法性
	switch ValueType(keyValueModel.ValueType) {
	case ValueTypeInt, ValueTypeBool:
		if keyValueModel.ValueInt == nil {
			return 0, ErrValueIsNil
		}
	case ValueTypeString:
		if keyValueModel.ValueString == nil {
			return 0, ErrValueIsNil
		}
	case ValueTypeFloat:
		if keyValueModel.ValueFloat == nil {
			return 0, ErrValueIsNil
		}
	case ValueTypeSliceInt, ValueTypeSliceString, ValueTypeSliceFloat, ValueTypeSliceBool:
		if keyValueModel.ValueInt == nil {
			return 0, ErrSliceButValueIntIsNil
		}
	default:
		return 0, ErrValueType
	}

	// 读取值
	switch ValueType(keyValueModel.ValueType) {
	case ValueTypeInt:
		value.Data = *keyValueModel.ValueInt
		value.Type = ValueTypeInt
	case ValueTypeBool:
		if (*keyValueModel.ValueInt) == 0 {
			value.Data = false
		} else {
			value.Data = true
		}
		value.Type = ValueTypeBool
	case ValueTypeString:
		value.Data = *keyValueModel.ValueString
		value.Type = ValueTypeString
	case ValueTypeFloat:
		value.Data = *keyValueModel.ValueFloat
		value.Type = ValueTypeFloat
	case ValueTypeSliceInt, ValueTypeSliceString, ValueTypeSliceFloat, ValueTypeSliceBool:
		sliceNum = *keyValueModel.ValueInt
		if sliceNum <= 0 {
			return 0, fmt.Errorf("slice but sliceNum is %d", sliceNum)
		}
	default:
		return 0, ErrValueType
	}
	*rec = *value
	return sliceNum, nil
}

func (m *SqliteCore) Set(key string, value *ValueUnit) error {
	if !m.IsInitialized() {
		return ErrSqliteCoreNotInit
	}
	// 为避免GetAll时，取出slice的成员作为单独的主键，这里不允许key中包含[]
	if strings.Contains(key, "[") || strings.Contains(key, "]") {
		return ErrKeyCanNotContainSquareBrackets
	}

	var needCreate []*KeyValueModel
	var needSet []*KeyValueModel
	var needRemove []*KeyValueModel // slice缩短的情况

	keyValueModels, err := sqliteData2Model(key, value)
	if err != nil {
		return errors.Join(ErrSqliteData2Model, err)
	}

	dbValue := &ValueUnit{}
	exist, err := m.getInner(key, dbValue, false)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Join(ErrSqliteDBFileAddrNotExist, err)
	}

	notList := dbValue.Type < ValueTypeSliceBegin
	if notList {
		if !exist {
			needCreate = keyValueModels
		}
		if exist {
			if !Compare(dbValue, value) {
				needSet = keyValueModels
			}
		}
	} else {
		sliceNum := 0
		switch dbValue.Type {
		case ValueTypeSliceInt:
			sliceNum = len(ToBase[[]int](dbValue))
		case ValueTypeSliceString:
			sliceNum = len(ToBase[[]string](dbValue))
		case ValueTypeSliceFloat:
			sliceNum = len(ToBase[[]float32](dbValue))
		case ValueTypeSliceBool:
			sliceNum = len(ToBase[[]bool](dbValue))
		}
		for i, keyValueModel := range keyValueModels {
			if i == 0 {
				// 长度节点
				if *keyValueModel.ValueInt != sliceNum {
					needSet = append(needSet, keyValueModel)
				}
				continue
			}
			if i > sliceNum {
				needCreate = append(needCreate, keyValueModel)
				continue
			}
			// 判断值
			switch ValueType(keyValueModel.ValueType) {
			case ValueTypeInt:
				if *keyValueModel.ValueInt != ToBase[[]int](dbValue)[i-1] {
					needSet = append(needSet, keyValueModel)
				}
			case ValueTypeBool:
				b := false
				if *keyValueModel.ValueInt != 0 {
					b = true
				}
				if b != ToBase[[]bool](dbValue)[i-1] {
					needSet = append(needSet, keyValueModel)
				}
			case ValueTypeString:
				if *keyValueModel.ValueString != ToBase[[]string](dbValue)[i-1] {
					needSet = append(needSet, keyValueModel)
				}
			case ValueTypeFloat:
				if *keyValueModel.ValueFloat != ToBase[[]float32](dbValue)[i-1] {
					needSet = append(needSet, keyValueModel)
				}
			}
		}
		for i := len(keyValueModels) - 1; i < sliceNum; i++ {
			needRemove = append(needRemove, &KeyValueModel{
				Key: key + "[" + strconv.Itoa(i) + "]",
			})
		}

	}

	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	for _, keyValueModel := range needCreate {
		result := m.db.Create(keyValueModel)
		if result.Error != nil {
			return errors.Join(ErrCreateValue, result.Error)
		}
	}

	for _, keyValueModel := range needSet {
		result := m.db.Where("Key = ?", keyValueModel.Key).Updates(keyValueModel)
		if result.Error != nil {
			return errors.Join(ErrSetValue, result.Error)
		}
	}

	for _, keyValueModel := range needRemove {
		result := m.db.Where("Key = ?", keyValueModel.Key).Delete(&KeyValueModel{})
		if result.Error != nil {
			return errors.Join(ErrRemoveValue, result.Error)
		}
	}

	return nil
}

// 将数据转换为model，如果是slice类型，会返回多个model
func sqliteData2Model(key string, value *ValueUnit) ([]*KeyValueModel, error) {
	keyValueModel := &KeyValueModel{
		Key:       key,
		ValueType: int(value.Type),
	}
	switch value.Type {
	case ValueTypeInt:
		valueInt := ToBase[int](value)
		keyValueModel.ValueInt = &valueInt
	case ValueTypeBool:
		valueBool := ToBase[bool](value)
		if valueBool {
			valueInt := 1
			keyValueModel.ValueInt = &valueInt
		} else {
			valueInt := 0
			keyValueModel.ValueInt = &valueInt
		}
	case ValueTypeString:
		valueString := ToBase[string](value)
		keyValueModel.ValueString = &valueString
	case ValueTypeFloat:
		valueFloat := ToBase[float32](value)
		keyValueModel.ValueFloat = &valueFloat
	case ValueTypeSliceInt:
		sliceNum := len(value.Data.([]int))
		if sliceNum <= 0 {
			return nil, fmt.Errorf("slice but sliceNum is %d", sliceNum)
		}
		valueInt := sliceNum
		keyValueModel.ValueInt = &valueInt
	case ValueTypeSliceString:
		sliceNum := len(value.Data.([]string))
		if sliceNum <= 0 {
			return nil, fmt.Errorf("slice but sliceNum is %d", sliceNum)
		}
		valueInt := sliceNum
		keyValueModel.ValueInt = &valueInt
	case ValueTypeSliceFloat:
		sliceNum := len(value.Data.([]float32))
		if sliceNum <= 0 {
			return nil, fmt.Errorf("slice but sliceNum is %d", sliceNum)
		}
		valueInt := sliceNum
		keyValueModel.ValueInt = &valueInt
	case ValueTypeSliceBool:
		sliceNum := len(value.Data.([]bool))
		if sliceNum <= 0 {
			return nil, fmt.Errorf("slice but sliceNum is %d", sliceNum)
		}
		valueInt := sliceNum
		keyValueModel.ValueInt = &valueInt
	}

	result := make([]*KeyValueModel, 1)
	result[0] = keyValueModel

	switch value.Type {
	case ValueTypeSliceInt:
		//for i, v := range ToBase[[]int](value) {
		//	Errslice = m.Set(Key+"["+strconv.Itoa(i)+"]", &ValueUnit{
		//		Type: ValueTypeInt,
		//		Data: v,
		//	})
		//}
		for i, v := range ToBase[[]int](value) {
			newV := v
			result = append(result, &KeyValueModel{
				Key:       key + "[" + strconv.Itoa(i) + "]",
				ValueType: int(ValueTypeInt),
				ValueInt:  &newV,
			})
		}
	case ValueTypeSliceString:
		for i, v := range ToBase[[]string](value) {
			newV := v
			result = append(result, &KeyValueModel{
				Key:         key + "[" + strconv.Itoa(i) + "]",
				ValueType:   int(ValueTypeString),
				ValueString: &newV,
			})
		}
	case ValueTypeSliceFloat:
		for i, v := range ToBase[[]float32](value) {
			newV := v
			result = append(result, &KeyValueModel{
				Key:        key + "[" + strconv.Itoa(i) + "]",
				ValueType:  int(ValueTypeFloat),
				ValueFloat: &newV,
			})
		}
	case ValueTypeSliceBool:
		for i, v := range ToBase[[]bool](value) {
			ti := 0
			if v {
				ti = 1
			}
			result = append(result, &KeyValueModel{
				Key:       key + "[" + strconv.Itoa(i) + "]",
				ValueType: int(ValueTypeBool),
				ValueInt:  &ti,
			})
		}
	}

	return result, nil
}

func (m *SqliteCore) Delete(key string) error {
	if !m.IsInitialized() {
		return ErrSqliteCoreNotInit
	}
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	return m.db.Where("Key = ?", key).Delete(&KeyValueModel{}).Error
}

func (m *SqliteCore) Have(key string) (bool, error) {
	if !m.IsInitialized() {
		return false, ErrSqliteCoreNotInit
	}
	m.rwLock.RLock()
	defer m.rwLock.RUnlock()
	var keyValueModel KeyValueModel
	result := m.db.Where("Key = ?", key).First(&keyValueModel)
	if result.Error != nil {
		return false, errors.Join(ErrSqliteDBFileAddrNotExist, result.Error)
	}
	return true, nil
}

func (m *SqliteCore) GetAll() (map[string]*ValueUnit, error) {
	if !m.IsInitialized() {
		return nil, ErrSqliteCoreNotInit
	}
	m.rwLock.RLock()
	defer m.rwLock.RUnlock()
	var keyValueModelList []KeyValueModel
	result := m.db.Find(&keyValueModelList)
	if result.Error != nil {
		return nil, errors.Join(ErrGetAllValue, result.Error)
	}
	keyValueModelMap := make(map[string]*ValueUnit)
	for _, keyValueModel := range keyValueModelList {
		// 跳过所有含有[]的key，因为这些key是slice的成员，不是真正的key
		if strings.Contains(keyValueModel.Key, "[") || strings.Contains(keyValueModel.Key, "]") {
			continue
		}
		unit := &ValueUnit{}
		sliceNum, err := sqliteModel2Data(keyValueModel, unit)
		if err != nil {
			return nil, errors.Join(ErrsqliteModel2Data, err)
		}

		if sliceNum != 0 {
			// slice 的内容存放在 Key[0]、Key[1]、Key[2]...Key[sliceNum-1] 列中
			_, err = m.Get(keyValueModel.Key, unit)
			if err != nil {
				return nil, errors.Join(ErrGetSliceValue, err)
			}
		}
		keyValueModelMap[keyValueModel.Key] = unit
	}
	return keyValueModelMap, nil
}
