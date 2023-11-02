package xstorage

import (
	"github.com/intmian/mian_go_lib/tool/misc"
	"os"
	"testing"
)

func TestMgr(t *testing.T) {
	// 删除test.db文件
	os.Remove("test.db")
	m, err := NewMgr(KeyValueSetting{
		Property: misc.CreateProperty(MultiSafe, UseCache, UseDB, FullInitLoad),
		SaveType: SqlLiteDB,
		DBAddr:   "tset.db",
	})
	if err != nil {
		t.Error(err)
		return
	}
	err = m.Set("1", ToUnit("1", VALUE_TYPE_STRING))
	if err != nil {
		t.Error(err)
		return
	}
	err = m.Set("2", ToUnit(2, VALUE_TYPE_INT))
	if err != nil {
		t.Error(err)
		return
	}
	err = m.Set("3", ToUnit(float32(3.0), VALUE_TYPE_FLOAT))
	if err != nil {
		t.Error(err)
		return
	}
	ok, v, err := m.Get("1")
	if err != nil {
		return
	}
	if !ok {
		t.Error("not ok")
		return
	}
	if v.Type != VALUE_TYPE_STRING {
		t.Error("type error")
		return
	}
	if ToBase[string](v) != "1" {
		t.Error("value error")
		return
	}
}
