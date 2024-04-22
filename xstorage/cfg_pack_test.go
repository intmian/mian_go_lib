package xstorage

import (
	"encoding/json"
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
	err3 := cfg.AddParam(&CfgParam{
		Key:       "test3",
		ValueType: ValueTypeInt,
		CanUser:   true,
		Default:   ToUnit(123, ValueTypeInt),
		RealKey:   "test3real",
	})
	err4 := cfg.AddParam(&CfgParam{
		Key:       "test4",
		ValueType: ValueTypeInt,
		CanUser:   false,
		Default:   ToUnit(123, ValueTypeInt),
		RealKey:   "test4real",
	})
	err = misc.JoinErr(err1, err2, err3, err4)
	if err != nil {
		t.Fatal(err)
	}

	// 测试不存在的key
	v, _ := cfg.Get("test1")
	if v != nil {
		t.Fatal("test1 should be nil")
	}

	// 测试默认
	v, _ = cfg.Get("test3")
	if v != nil {
		t.Fatal("test3 should be nil")
	}
	v, _ = cfg.Get("test4")
	if v == nil {
		t.Fatal("test4 should not be nil")
	}

	// 测非用户的
	bytes, _ := json.Marshal("hahaha")
	err = cfg.Set("test1", string(bytes))
	if err != nil {
		t.Fatal(err)
	}
	v, _ = cfg.Get("test1")
	if v == nil || ToBase[string](v) != "hahaha" {
		t.Fatal("test1 should be hahaha")
	}

	// 测试用户
	users := []string{"user1", "user2", ""}
	s := []string{"a", "b", "c"}
	bytes, _ = json.Marshal(s)
	for _, user := range users {
		err = cfg.SetUser(user, "test2", string(bytes))
		if user == "" {
			if err == nil {
				t.Fatal("user is empty, should return error")
			}
		} else {
			if err != nil {
				t.Fatal(err)
			}
		}
		if user == "" {
			continue
		}
		v, _ = cfg.GetUser(user, "test2")
		if v == nil || ToBase[[]string](v)[0] != "a" {
			t.Fatal("test2 should be a")
		}
	}
}
