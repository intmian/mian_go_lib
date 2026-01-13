package xbi

import (
	"github.com/intmian/mian_go_lib/tool/misc"
	"github.com/pkg/errors"
)

// XBi 是轻量业务日志结构体，封装了 gorm 客户端和项目、日志库信息
// 因为数量暂时较少先都用数据库存
type XBi struct {
	ini     misc.InitTag
	setting Setting
}

// NewXBi 初始化
func NewXBi(setting Setting) (*XBi, error) {
	xbi := &XBi{
		setting: setting,
	}

	if xbi.setting.Db == nil || setting.ErrorChan == nil || setting.Ctx == nil {
		return nil, errors.New("XBi setting is invalid")
	}

	xbi.ini.SetInitialized()
	return xbi, nil
}

func RegisterLogEntity[DataType any](x *XBi, entity LogEntity[DataType]) error {
	// 受限于 go 的泛型，目前只能用这种方式注册实体，以后如果泛型支持可以将entity LogEntity[DataType]去掉。因为泛型没办法实例化
	if x == nil {
		return errors.New("XBi instance is nil")
	}
	if !x.ini.IsInitialized() {
		return errors.New("XBi file not initialized")
	}
	err := x.setting.Db.Table(entity.TableName()).AutoMigrate(toDbData[DataType](entity))
	if err != nil {
		return errors.WithMessage(err, "RegisterLogEntity failed")
	}
	return nil
}

func WriteLog[T any](x *XBi, log LogEntity[T]) error {
	// 因为类方法无法使用泛型，所以只能用这种方式传递
	if x == nil {
		return errors.New("XBi instance is nil")
	}
	if !x.ini.IsInitialized() {
		return errors.New("XBi file not initialized")
	}
	if log == nil {
		return errors.New("log entity is nil")
	}

	tableName := string(log.TableName())
	if tableName == "" {
		return errors.New("table name is empty")
	}

	go func() {
		realLog := toDbData(log)
		err := x.setting.Db.WithContext(x.setting.Ctx).Table(tableName).Create(&realLog).Error
		if err != nil {
			select {
			case x.setting.ErrorChan <- err:
			case <-x.setting.Ctx.Done():
			}
		}
	}()

	return nil
}

type ReadLogFilter struct {
	conditions map[string]any
	page       *struct {
		Num  int
		Size int
	}
	orderBy *struct {
		Field string
		Desc  bool
	}
}

func (r *ReadLogFilter) SetConditions(conditions map[string]any) *ReadLogFilter {
	r.conditions = conditions
	return r
}

func (r *ReadLogFilter) SetPage(num, size int) *ReadLogFilter {
	r.page = &struct {
		Num  int
		Size int
	}{
		Num:  num,
		Size: size,
	}
	return r
}

func (r *ReadLogFilter) SetOrderBy(field string, desc bool) *ReadLogFilter {
	r.orderBy = &struct {
		Field string
		Desc  bool
	}{
		Field: field,
		Desc:  desc,
	}
	return r
}

func ReadLogWithFilter[T any](x *XBi, tableName string, filter ReadLogFilter) ([]DbLogData[T], error) {
	if x == nil {
		return nil, errors.New("XBi instance is nil")
	}
	if !x.ini.IsInitialized() {
		return nil, errors.New("XBi file not initialized")
	}
	if tableName == "" {
		return nil, errors.New("table name is empty")
	}

	tx := x.setting.Db.WithContext(x.setting.Ctx).Table(tableName)

	if filter.conditions != nil {
		tx = tx.Where(filter.conditions)
	}

	if filter.orderBy != nil {
		orderStr := filter.orderBy.Field
		if filter.orderBy.Desc {
			orderStr += " desc"
		}
		tx = tx.Order(orderStr)
	}

	if filter.page != nil {
		if filter.page.Num <= 0 {
			filter.page.Num = 1
		}
		tx = tx.Offset((filter.page.Num - 1) * filter.page.Size).Limit(filter.page.Size)
	}

	var results []DbLogData[T]
	err := tx.Find(&results).Error
	if err != nil {
		return nil, errors.Wrap(err, "ReadLogWithFilter failed")
	}
	return results, nil
}

func ReadLog[T any](x *XBi, tableName string, conditions map[string]any, pageNum, page int) ([]DbLogData[T], error) {
	if x == nil {
		return nil, errors.New("XBi instance is nil")
	}
	if !x.ini.IsInitialized() {
		return nil, errors.New("XBi file not initialized")
	}
	if tableName == "" {
		return nil, errors.New("table name is empty")
	}

	var results []DbLogData[T]
	if pageNum <= 0 {
		err := x.setting.Db.WithContext(x.setting.Ctx).Table(tableName).Where(conditions).Find(&results).Error
		if err != nil {
			return nil, errors.Wrap(err, "ReadLog failed")
		}
		return results, nil
	}

	err := x.setting.Db.WithContext(x.setting.Ctx).Table(tableName).Where(conditions).Offset((pageNum - 1) * page).Limit(pageNum).Find(&results).Error
	if err != nil {
		return nil, errors.Wrap(err, "ReadLog failed")
	}
	return results, nil
}
