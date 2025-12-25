package xbi

// LogEntity 日志接口
type LogEntity[DataType any] interface {
	TableName() string
	GetWriteableData() *DataType // 获取可写入的数据
}

type DbLogData[DataType any] struct {
	// 一些公共字段
	RecordTime int64 // ms 时间戳
	// 内嵌真实的数据
	Data DataType `gorm:"embedded"`
}
