package xstorage

type ValueUnit struct {
	Type ValueType
	Data interface{}
}

type IValueType interface {
	int | string | float32 | bool | []int | []string | []float32 | []bool
}

func Get[T IValueType](unit *ValueUnit) T {
	return T(unit.Data)
}
