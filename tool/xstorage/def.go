package xstorage

type BindFileType int

const (
	JSON BindFileType = iota
	XML
)

type KeyValueProperty uint32

const (
	// MultiSafe 多线程安全，仅针对缓存
	MultiSafe KeyValueProperty = 1 << iota // TODO:
	// UseCache 缓存
	UseCache
	// UseDB 落盘
	UseDB
	// FullInitLoad 全量初始化加载，没有的话即懒加载，仅在同时使用缓存和数据库时有效，建议与MultiSafe一起使用，并使用setAsync方法进行set
	FullInitLoad
)

type keyValueSaveType uint32

const (
	null keyValueSaveType = iota
	SqlLiteDB
	Json
)

type KeyValueSetting struct {
	Property KeyValueProperty
	SaveType keyValueSaveType
	DBAddr   string
}

type ValueType int

const (
	VALUE_TYPE_NORMAL_BEGIN ValueType = iota
	VALUE_TYPE_STRING
	VALUE_TYPE_INT
	VALUE_TYPE_FLOAT
	VALUE_TYPE_BOOL
	VALUE_TYPE_NORMAL_END
	VALUE_TYPE_SLICE_BEGIN  = 100
	VALUE_TYPE_SLICE_STRING = iota + VALUE_TYPE_SLICE_BEGIN - VALUE_TYPE_NORMAL_END - 1
	VALUE_TYPE_SLICE_INT
	VALUE_TYPE_SLICE_FLOAT
	VALUE_TYPE_SLICE_BOOL
)

type IValueType interface {
	int | string | float32 | bool | []int | []string | []float32 | []bool
}
