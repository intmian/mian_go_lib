package xstorage

import (
	"strconv"
	"strings"
)

// ValueUnit 加入类型，用于在非反射的情况下直接处理类型
type ValueUnit struct {
	Type  ValueType
	Data  interface{}
	dirty bool
}

func UnitToString(unit *ValueUnit) string {
	s := ""
	switch unit.Type {
	case VALUE_TYPE_STRING:
		s = ToBase[string](unit)
	case VALUE_TYPE_INT:
		s = strconv.Itoa(ToBase[int](unit))
	case VALUE_TYPE_FLOAT:
		s = strconv.FormatFloat(float64(ToBase[float32](unit)), 'f', -1, 32)
		// 如果不含小数点，就加上小数点
		if !strings.Contains(s, ".") {
			s += ".0"
		}
	case VALUE_TYPE_BOOL:
		s = strconv.FormatBool(ToBase[bool](unit))

	}
	return ""
}

// ToBase 直接转换为基础类型，一方面是为了避免频繁的类型转换，另一方面是为了限制类型
func ToBase[T IValueType](unit *ValueUnit) T {
	// 判断类型是否正确
	t, ok := unit.Data.(T)
	if !ok {
		var empty T
		return empty
	}
	return t
}

func ToUnit[T IValueType](value T, valueType ValueType) *ValueUnit {
	return &ValueUnit{
		Type:  valueType,
		Data:  value,
		dirty: false,
	}
}

func Compare(unit1 *ValueUnit, unit2 *ValueUnit) bool {
	if unit1 == nil || unit2 == nil {
		return true
	}
	if unit1 == nil || unit2 == nil {
		return false
	}
	if unit1.Type != unit2.Type {
		return false
	}
	switch unit1.Type {
	case VALUE_TYPE_STRING:
		if ToBase[string](unit1) != ToBase[string](unit1) {
			return false
		}
	case VALUE_TYPE_INT:
		if ToBase[int](unit1) != ToBase[int](unit1) {
			return false
		}
	case VALUE_TYPE_FLOAT:
		if ToBase[float32](unit1) != ToBase[float32](unit1) {
			return false
		}
	case VALUE_TYPE_BOOL:
		if ToBase[bool](unit1) != ToBase[bool](unit1) {
			return false
		}
	case VALUE_TYPE_SLICE_INT:
		if len(ToBase[[]int](unit1)) != len(ToBase[[]int](unit1)) {
			return false
		}
		for i, val := range ToBase[[]int](unit1) {
			if val != ToBase[[]int](unit1)[i] {
				return false
			}
		}
	case VALUE_TYPE_SLICE_STRING:
		if len(ToBase[[]string](unit1)) != len(ToBase[[]string](unit1)) {
			return false
		}
		for i, val := range ToBase[[]string](unit1) {
			if val != ToBase[[]string](unit1)[i] {
				return false
			}
		}
	case VALUE_TYPE_SLICE_FLOAT:
		if len(ToBase[[]float32](unit1)) != len(ToBase[[]float32](unit1)) {
			return false
		}
		for i, val := range ToBase[[]float32](unit1) {
			if val != ToBase[[]float32](unit1)[i] {
				return false
			}
		}
	case VALUE_TYPE_SLICE_BOOL:
		if len(ToBase[[]bool](unit1)) != len(ToBase[[]bool](unit1)) {
			return false
		}
		for i, val := range ToBase[[]bool](unit1) {
			if val != ToBase[[]bool](unit1)[i] {
				return false
			}
		}
	}
	return true
}

func Copy(srcValue *ValueUnit, newValue *ValueUnit) {
	switch srcValue.Type {
	case VALUE_TYPE_STRING, VALUE_TYPE_INT, VALUE_TYPE_FLOAT, VALUE_TYPE_BOOL:
		*newValue = *srcValue
	case VALUE_TYPE_SLICE_INT:
		newValue.Type = VALUE_TYPE_SLICE_INT
		_, ok := srcValue.Data.([]int)
		var newData []int
		if !ok {
			// 将类型强制转换为[]int
			newData = make([]int, len(srcValue.Data.([]interface{})))
			for i, val := range srcValue.Data.([]interface{}) {
				newData[i] = int(val.(int64))
			}
		} else {
			newData = make([]int, len(srcValue.Data.([]int)))
			copy(newData, srcValue.Data.([]int))
		}
		newValue.Data = newData
	case VALUE_TYPE_SLICE_STRING:
		newValue.Type = VALUE_TYPE_SLICE_STRING
		_, ok := srcValue.Data.([]string)
		var newData []string
		if !ok {
			// 将类型强制转换为[]string
			newData = make([]string, len(srcValue.Data.([]interface{})))
			for i, val := range srcValue.Data.([]interface{}) {
				newData[i] = val.(string)
			}
		} else {
			newData = make([]string, len(srcValue.Data.([]string)))
			copy(newData, srcValue.Data.([]string))
		}
		newValue.Data = newData
	case VALUE_TYPE_SLICE_FLOAT:
		newValue.Type = VALUE_TYPE_SLICE_FLOAT
		_, ok := srcValue.Data.([]float32)
		var newData []float32
		if !ok {
			// 将类型强制转换为[]float32
			newData = make([]float32, len(srcValue.Data.([]interface{})))
			for i, val := range srcValue.Data.([]interface{}) {
				newData[i] = float32(val.(float64))
			}
		} else {
			newData = make([]float32, len(srcValue.Data.([]float32)))
			copy(newData, srcValue.Data.([]float32))
		}
		newValue.Data = newData
	case VALUE_TYPE_SLICE_BOOL:
		newValue.Type = VALUE_TYPE_SLICE_BOOL
		_, ok := srcValue.Data.([]bool)
		var newData []bool
		if !ok {
			// 将类型强制转换为[]bool
			newData = make([]bool, len(srcValue.Data.([]interface{})))
			for i, val := range srcValue.Data.([]interface{}) {
				newData[i] = val.(bool)
			}
		} else {
			newData = make([]bool, len(srcValue.Data.([]bool)))
			copy(newData, srcValue.Data.([]bool))
		}
		newValue.Data = newData
	}
}

func Join(src ...string) string {
	var result string
	for _, s := range src {
		result += "." + s
	}
	return result
}

func Split(src string) []string {
	return strings.Split(src, ".")
}
