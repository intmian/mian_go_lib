package xbi

import (
	"context"

	"gorm.io/gorm"
)

type Setting struct {
	// 权限配置
	Db        *gorm.DB
	errorChan chan<- error
	ctx       context.Context
}
