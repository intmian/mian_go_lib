package xbi

import (
	"context"

	"gorm.io/gorm"
)

type Setting struct {
	// 权限配置
	Db        *gorm.DB
	ErrorChan chan<- error
	Ctx       context.Context
}

func GetDefaultSetting() Setting {
	return Setting{}
}
