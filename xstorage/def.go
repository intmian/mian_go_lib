package xstorage

type KeyValueProperty uint32

const (
	// MultiSafe 多线程安全，仅针对缓存
	MultiSafe KeyValueProperty = 1 << iota // TODO:
	// UseCache 缓存
	UseCache
	// UseDisk 落盘
	UseDisk
	// FullInitLoad 全量初始化加载，没有的话即懒加载，仅在同时使用缓存和数据库时有效，建议与MultiSafe一起使用，并使用setAsync方法进行set
	FullInitLoad
	// OpenWeb 启用web端口，建议不要使用内置的web服务，而是自己启动一个web服务，然后调用GetGinEngine方法获取gin引擎，然后自己注册路由
	OpenWeb
)

type keyValueSaveType uint32

const (
	DBBegin keyValueSaveType = iota
	SqlLiteDB
	FileBegin
	Toml // 为保证效率，必须开启UseCache、FullInitLoad
)

type XStorageSetting struct {
	Property KeyValueProperty
	SaveType keyValueSaveType
	DBAddr   string
	FileAddr string
}

type ValueType int

// 可能会被外部调用，所以复杂命名
const (
	ValueTypeNormalBegin ValueType = iota
	ValueTypeString
	ValueTypeInt
	ValueTypeFloat
	ValueTypeBool
	ValueTypeNormalEnd
	ValueTypeSliceBegin  = 100
	ValueTypeSliceString = iota + ValueTypeSliceBegin - ValueTypeNormalEnd - 1
	ValueTypeSliceInt
	ValueTypeSliceFloat
	ValueTypeSliceBool
)

type IValueType interface {
	int | string | float32 | bool | []int | []string | []float32 | []bool
}
