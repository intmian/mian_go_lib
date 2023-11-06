package xstorage

// ValueUnit 加入类型，用于在非反射的情况下直接处理类型
type ValueUnit struct {
	Type  ValueType
	Data  interface{}
	dirty bool
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
		newValue.Data = make([]int, len(srcValue.Data.([]int)))
		copy(newValue.Data.([]int), srcValue.Data.([]int))
	case VALUE_TYPE_SLICE_STRING:
		newValue.Type = VALUE_TYPE_SLICE_STRING
		newValue.Data = make([]string, len(srcValue.Data.([]string)))
		copy(newValue.Data.([]string), srcValue.Data.([]string))
	case VALUE_TYPE_SLICE_FLOAT:
		newValue.Type = VALUE_TYPE_SLICE_FLOAT
		newValue.Data = make([]float32, len(srcValue.Data.([]float32)))
		copy(newValue.Data.([]float32), srcValue.Data.([]float32))
	case VALUE_TYPE_SLICE_BOOL:
		newValue.Type = VALUE_TYPE_SLICE_BOOL
		newValue.Data = make([]bool, len(srcValue.Data.([]bool)))
		copy(newValue.Data.([]bool), srcValue.Data.([]bool))
	}
}
