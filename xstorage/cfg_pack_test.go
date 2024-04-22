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
	err1 := cfg.AddParam(&CfgParam{
		Key:       "test1",
		ValueType: ValueTypeString,
		CanUser:   false,
		RealKey:   "test1",
	})
	err2 := cfg.AddParam(&CfgParam{
		Key:       "test2",
		ValueType: ValueTypeSliceString,
		CanUser:   true,
		RealKey:   "test2real",
	})
	err = misc.JoinErr(err1, err2)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("have not key", func(t *testing.T) {
		err := cfg.SetCfg()
	})
}
