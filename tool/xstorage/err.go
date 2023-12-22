package xstorage

import "github.com/intmian/mian_go_lib/tool/misc"

const (
	SqliteDBFileAddrEmptyErr                   = misc.ErrStr("sqlite db file addr is empty")
	SqliteDBFileAddrNotExistErr                = misc.ErrStr("sqlite db file addr not exist")
	MgrNotInitErr                              = misc.ErrStr("mgr not init")
	SqliteCoreNotInitErr                       = misc.ErrStr("sqlite core not init")
	SetValueErr                                = misc.ErrStr("set value error")
	KeyIsEmptyErr                              = misc.ErrStr("key is empty")
	NotUseCacheErr                             = misc.ErrStr("not use cache")
	KeyNotExistErr                             = misc.ErrStr("key not exist")
	ValueIsNilErr                              = misc.ErrStr("value is nil")
	SliceButValueIntIsNilErr                   = misc.ErrStr("slice but ValueInt is nil")
	ValueTypeErr                               = misc.ErrStr("value type error")
	OpenSqliteErr                              = misc.ErrStr("open sqlite error")
	AutoMigrateErr                             = misc.ErrStr("auto migrate error")
	RecIsNilErr                                = misc.ErrStr("rec is nil")
	NotUseDbErr                                = misc.ErrStr("not use db")
	RecordToMapErr                             = misc.ErrStr("record to map error")
	PoolTypeErr                                = misc.ErrStr("pool type error")
	NotUseCacheAndNotUseDbErr                  = misc.ErrStr("not use cache and not use db")
	GetAllValueErr                             = misc.ErrStr("get all value error")
	GinEngineRunErr                            = misc.ErrStr("gin engine run error")
	NotUseCacheOrNotUseDbAndFullInitLoadErr    = misc.ErrStr("not use cache or not use db and full init load")
	UseJsonButNotUseCacheAndNotFullInitLoadErr = misc.ErrStr("use json, but not use cache and not full init load")
	NewSqliteCoreErr                           = misc.ErrStr("new sqlite core error")
	NotUseDbAndFullInitLoadErr                 = misc.ErrStr("not use db and full init load")
	ValueUnitIsNilErr                          = misc.ErrStr("valueUnit is nil")
	ValueIsDirtyErr                            = misc.ErrStr("value is dirty")
	RemoveFromMapErr                           = misc.ErrStr("remove from map error")
	DeleteValueErr                             = misc.ErrStr("delete value error")
	KeyCanNotContainSquareBracketsErr          = misc.ErrStr("Key can not contain []")
	SqliteData2ModelErr                        = misc.ErrStr("sqliteData2Model")
	CreateValueErr                             = misc.ErrStr("create value error")
	RemoveValueErr                             = misc.ErrStr("remove value error")
	GetSliceValueErr                           = misc.ErrStr("get slice value error")
	sqliteModel2DataErr                        = misc.ErrStr("sqliteModel2Data error")
	GetErr                                     = misc.ErrStr("get error")
)
