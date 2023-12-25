package xstorage

import (
	"github.com/intmian/mian_go_lib/tool/misc"
	"os"
)

type TomlCore struct {
	addr string
}

func NewTomlCore(addr string) *TomlCore {
	return &TomlCore{addr: addr}
}

func (j TomlCore) GetAll() (map[string]*ValueUnit, error) {
	tt := misc.TomlTool{}
	m := make(map[string]*ValueUnit)
	// 如果文件不存在，会返回空map
	_, err := os.Stat(j.addr)
	if err != nil {
		return m, nil
	}
	err = tt.Read(j.addr, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (j TomlCore) SaveAll(data map[string]*ValueUnit) error {
	tt := &misc.TomlTool{}
	return tt.Write(j.addr, data)
}
