package xbi

import (
	"regexp"
	"strings"

	"github.com/intmian/mian_go_lib/tool/misc"
	"github.com/pkg/errors"
)

var (
	// 用于校验包含 ? 的复杂条件，例如 "Data.Duration > ?"
	// 允许的格式: 字段名 + 空格 + 操作符 + 空格(可选) + ?
	// 操作符白名单: =, <, >, <=, >=, <>, !=, LIKE, NOT LIKE
	// (?i) 表示不区分大小写
	regexComplexCondition = regexp.MustCompile(`(?i)^[a-zA-Z0-9_.]+\s+(?:=|!=|<>|>|<|>=|<=|LIKE|NOT\s+LIKE)\s*\?$`)

	// 用于校验普通字段名，例如 "Data.Rows"
	regexSimpleColumn = regexp.MustCompile(`^[a-zA-Z0-9_.]+$`)
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

func ReadLogWithFilter[T any](x *XBi, tableName string, filter *ReadLogFilter) ([]DbLogData[T], int64, error) {
	if x == nil {
		return nil, 0, errors.New("XBi instance is nil")
	}
	if !x.ini.IsInitialized() {
		return nil, 0, errors.New("XBi file not initialized")
	}
	if tableName == "" {
		return nil, 0, errors.New("table name is empty")
	}
	if filter == nil {
		return nil, 0, errors.New("filter is nil")
	}

	tx := x.setting.Db.WithContext(x.setting.Ctx).Table(tableName)

	if filter.conditions != nil {
		for k, v := range filter.conditions {
			if strings.Contains(k, "?") {
				// 安全检查：防止 SQL 注入
				if !regexComplexCondition.MatchString(k) {
					return nil, 0, errors.Errorf("invalid query condition key: %s", k)
				}
				tx = tx.Where(k, v)
			} else {
				// 安全检查：防止作为列名传入奇怪的 SQL 片段
				if !regexSimpleColumn.MatchString(k) {
					return nil, 0, errors.Errorf("invalid query column name: %s", k)
				}
				tx = tx.Where(k+" = ?", v)
			}
		}
	}

	var count int64
	if err := tx.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrap(err, "count failed")
	}

	if filter.orderBy != nil {
		field := filter.orderBy.Field
		// 安全检查：防止 Order By 注入
		if !regexSimpleColumn.MatchString(field) {
			return nil, 0, errors.Errorf("invalid order by field: %s", field)
		}

		orderStr := field
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
		return nil, 0, errors.Wrap(err, "ReadLogWithFilter failed")
	}
	return results, count, nil
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
