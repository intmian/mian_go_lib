package xbi

import (
	"context"

	"github.com/intmian/mian_go_lib/tool/misc"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// XBi 是轻量业务日志结构体，封装了 gorm 客户端和项目、日志库信息
// 因为数量暂时较少先都用数据库存
type XBi struct {
	ini      misc.InitTag
	dbClient *gorm.DB
	setting  Setting
}

// NewXBi 初始化 SLS 客户端，传入 Endpoint、AccessKeyID、AccessKeySecret、project、logstore
func NewXBi(setting Setting) *XBi {
	xbi := &XBi{
		setting:  setting,
		dbClient: setting.Db,
	}
	xbi.ini.SetInitialized()
	return xbi
}

func WriteLog[T any](x *XBi, log LogEntity[T], ctx context.Context, errChan chan<- error) error {
	if x == nil {
		return errors.New("XBi instance is nil")
	}
	if !x.ini.IsInitialized() {
		return errors.New("XBi file not initialized")
	}
	if log == nil {
		return errors.New("log entity is nil")
	}
	if x.dbClient == nil {
		return errors.New("db client is nil")
	}

	tableName := string(log.TableName())
	if tableName == "" {
		return errors.New("table name is empty")
	}

	go func() {
		realLog := toRealEntity(log)
		err := x.dbClient.WithContext(ctx).Table(tableName).Create(&realLog).Error
		select {
		case errChan <- err:
		case <-ctx.Done():
		}
	}()

	return nil
}

func ReadLog[T any](x *XBi, tableName string, conditions map[string]any, pageNum, page int, ctx context.Context) ([]RealLogEntity[T], error) {
	if x == nil {
		return nil, errors.New("XBi instance is nil")
	}
	if !x.ini.IsInitialized() {
		return nil, errors.New("XBi file not initialized")
	}
	if x.dbClient == nil {
		return nil, errors.New("db client is nil")
	}
	if tableName == "" {
		return nil, errors.New("table name is empty")
	}

	var results []RealLogEntity[T]
	err := x.dbClient.WithContext(ctx).Table(tableName).Where(conditions).Offset((pageNum - 1) * page).Limit(page).Find(&results).Error
	if err != nil {
		return nil, errors.Wrap(err, "ReadLog failed")
	}

	return results, nil
}
