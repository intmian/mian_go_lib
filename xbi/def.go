package xbi

import "time"

type (
	PpjName      string
	DbName       string
	TableName    string
	KeyForSearch string
)

// LogEntity 日志接口
type LogEntity[T any] interface {
	PrjName() PpjName
	DbName() DbName
	TableName() TableName
	KeysForSearch() []KeyForSearch // 需要搜索的字段
	GetWriteableData() T           // 获取可写入的数据
}

type realLogEntity struct {
	RecordTime      time.Time
	OriginalJsonStr string
	// 需要搜索的字段会使用 反射注入
}
