package xstorage

import "gorm.io/gorm"

type KeyValueModel struct {
	gorm.Model
	// Key 主键、索引
	Key         string `gorm:"primaryKey;index"`
	ValueType   int
	ValueInt    *int
	ValueString *string
	ValueFloat  *float32
}
