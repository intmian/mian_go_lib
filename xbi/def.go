package xbi

import "time"

type (
	PpjName   string
	DbName    string
	TableName string
)

// LogEntity 日志接口
type LogEntity[T any] interface {
	TableName() TableName
	GetWriteableData() T // 获取可写入的数据
}

type RealLogEntity[T any] struct {
	// 一些公共字段
	RecordTime time.Time
	// 需要搜索的字段会使用 反射注入
	Data T `gorm:"embedded"`
}
