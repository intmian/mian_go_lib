package xstorage

import (
	"github.com/intmian/mian_go_lib/tool/misc"
	"os"
)

type JsonCore struct {
	addr string
}

func NewJsonCore(addr string) *JsonCore {
	return &JsonCore{addr: addr}
}

func (j JsonCore) GetAll() (map[string]*ValueUnit, error) {
	jt := &misc.JsonTool{}
	m := make(map[string]*ValueUnit)
	// 如果文件不存在，会返回空map
	_, err := os.Stat(j.addr)
	if err != nil {
		return m, nil
	}
	err = jt.Read(j.addr, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (j JsonCore) SaveAll(data map[string]*ValueUnit) error {
	jt := &misc.JsonTool{}
	return jt.Write(j.addr, data)
}
