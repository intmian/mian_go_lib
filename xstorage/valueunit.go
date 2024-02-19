package xstorage

import (
	"strconv"
	"strings"
)

// TODO: 考虑用泛型改写一个v2的

// ValueUnit 加入类型，用于在非反射的情况下直接处理类型
type ValueUnit struct {
	Type  ValueType
	Data  interface{}
	dirty bool
}

func (v *ValueUnit) Reset() {
	*v = ValueUnit{}
}

// ToBase 直接转换为基础类型，一方面是为了避免频繁的类型转换，另一方面是为了限制类型
func ToBase[T IValueType](unit *ValueUnit) T {
	// 判断类型是否正确
	if unit == nil {
		var empty T
		return empty
	}
	t, ok := unit.Data.(T)
	if !ok {
		var empty T
		return empty
	}
	return t
}

func ToBaseF[T IValueType](unit *ValueUnit, err error) T {
	if err != nil {
		var empty T
		return empty
	}
	return ToBase[T](unit)
}

func ToUnit[T IValueType](value T, valueType ValueType) *ValueUnit {
	return &ValueUnit{
		Type:  valueType,
		Data:  value,
		dirty: false,
	}
}

func UnitToString(unit *ValueUnit) string {
	s := ""
	switch unit.Type {
	case ValueTypeString:
		s = ToBase[string](unit)
	case ValueTypeInt:
		s = strconv.Itoa(ToBase[int](unit))
	case ValueTypeFloat:
		s = strconv.FormatFloat(float64(ToBase[float32](unit)), 'f', -1, 32)
		// 如果不含小数点，就加上小数点
		if !strings.Contains(s, ".") {
			s += ".0"
		}
	case ValueTypeBool:
		s = strconv.FormatBool(ToBase[bool](unit))

	}
	return ""
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
	case ValueTypeString:
		if ToBase[string](unit1) != ToBase[string](unit1) {
			return false
		}
	case ValueTypeInt:
		if ToBase[int](unit1) != ToBase[int](unit1) {
			return false
		}
	case ValueTypeFloat:
		if ToBase[float32](unit1) != ToBase[float32](unit1) {
			return false
		}
	case ValueTypeBool:
		if ToBase[bool](unit1) != ToBase[bool](unit1) {
			return false
		}
	case ValueTypeSliceInt:
		if len(ToBase[[]int](unit1)) != len(ToBase[[]int](unit1)) {
			return false
		}
		for i, val := range ToBase[[]int](unit1) {
			if val != ToBase[[]int](unit1)[i] {
				return false
			}
		}
	case ValueTypeSliceString:
		if len(ToBase[[]string](unit1)) != len(ToBase[[]string](unit1)) {
			return false
		}
		for i, val := range ToBase[[]string](unit1) {
			if val != ToBase[[]string](unit1)[i] {
				return false
			}
		}
	case ValueTypeSliceFloat:
		if len(ToBase[[]float32](unit1)) != len(ToBase[[]float32](unit1)) {
			return false
		}
		for i, val := range ToBase[[]float32](unit1) {
			if val != ToBase[[]float32](unit1)[i] {
				return false
			}
		}
	case ValueTypeSliceBool:
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
	case ValueTypeString, ValueTypeInt, ValueTypeFloat, ValueTypeBool:
		*newValue = *srcValue
	case ValueTypeSliceInt:
		newValue.Type = ValueTypeSliceInt
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
	case ValueTypeSliceString:
		newValue.Type = ValueTypeSliceString
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
	case ValueTypeSliceFloat:
		newValue.Type = ValueTypeSliceFloat
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
	case ValueTypeSliceBool:
		newValue.Type = ValueTypeSliceBool
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
