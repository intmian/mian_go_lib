package xbi

import (
	"encoding/json"
	"reflect"
	"time"
)

func toRealEntity[T any](entity LogEntity[T], logData interface{}) (interface{}, error) {
	realEntity := realLogEntity{}
	realEntity.RecordTime = time.Now()
	originalJsonBytes, err := json.Marshal(logData)
	if err != nil {
		return realEntity, err
	}
	realEntity.OriginalJsonStr = string(originalJsonBytes)

	// 使用反射获取需要搜索的字段，并注入到 realEntity 中
	v := reflect.ValueOf(&realEntity).Elem()
	keys := entity.KeysForSearch()
	for _, key := range keys {
		fieldVal := reflect.ValueOf(logData).Elem().FieldByName(string(key))
		if fieldVal.IsValid() {
			newField := reflect.New(fieldVal.Type()).Elem()
			newField.Set(fieldVal)
			v.FieldByName(string(key)).Set(newField)
		}
	}

	return realEntity, nil
}
