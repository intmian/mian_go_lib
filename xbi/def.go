package xbi

// LogEntity 日志接口
type LogEntity[DataType any] interface {
	TableName() string
	GetWriteableData() *DataType // 获取可写入的数据
}

type DbLogData[DataType any] struct {
	// 一些公共字段
	RecordTime int64 `json:"record_time" gorm:"column:record_time"` // 记录时间戳ms
	// 内嵌真实的数据
	Data DataType `gorm:"embedded"`
}
