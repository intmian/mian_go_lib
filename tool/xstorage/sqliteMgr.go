package xstorage

import (
	"context"
	"database/sql"
)

type SqliteMgr struct {
	KeyValueSetting
	db  *sql.DB
	ctx context.Context // 用于控制关闭系统时
}

func NewMgr(keyValueSetting KeyValueSetting) *SqliteMgr {
	return &SqliteMgr{KeyValueSetting: keyValueSetting}
}

func (m *SqliteMgr) Get(key string) (*ValueUnit, error) {
	// TODO
}

func (m *SqliteMgr) Set(key string, value *ValueUnit) error {
	// TODO
}

func (m *SqliteMgr) Have(key string) (bool, error) {
	// TODO
}
