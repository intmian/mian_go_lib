package xstorage

// ValueUnit 加入类型，用于在非反射的情况下直接处理类型
type ValueUnit struct {
	Type  ValueType
	Data  interface{}
	dirty bool
}

// Get 直接转换为基础类型，一方面是为了避免频繁的类型转换，另一方面是为了限制类型
func Get[T IValueType](unit *ValueUnit) T {
	// 判断类型是否正确
	t, ok := unit.Data.(T)
	if !ok {
		var empty T
		return empty
	}
	return t
}
