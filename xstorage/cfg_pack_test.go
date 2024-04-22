package xstorage

import (
	"github.com/intmian/mian_go_lib/tool/misc"
	"testing"
)

func TestCfg(t *testing.T) {
	m, err := NewXStorage(XStorageSetting{
		Property: misc.CreateProperty(MultiSafe, UseCache),
	})
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := NewCfgExt(m)
	if err != nil {
		t.Fatal(err)
	}
	cfg.AddParam(&CfgParam{
		Key:       "test1",
		ValueType: ValueTypeString,
		CanUser:   false,
		RealKey:   "test1",
	})

}
