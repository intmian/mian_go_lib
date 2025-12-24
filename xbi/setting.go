package xbi

import "gorm.io/gorm"

type Setting struct {
	// 权限配置
	Db *gorm.DB
}
