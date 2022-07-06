package xres

type ExcelCow struct {
	name      string
	ExcelType ColumnType
	data      interface{}
}

func Get[T any](from *ExcelCow) *T {
	return from.data.(*T)
}
