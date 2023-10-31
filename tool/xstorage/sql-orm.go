package xstorage

import "gorm.io/gorm"

type KeyValueModel struct {
	gorm.Model
	// key 主键、索引
	key         *string `gorm:"primaryKey;index"`
	valueType   int
	valueInt    *int
	valueString *string
	valueFloat  *float32
}
