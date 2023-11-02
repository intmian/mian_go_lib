package xstorage

import (
	"github.com/intmian/mian_go_lib/tool/misc"
	"os"
	"strconv"
	"testing"
)

func TestMgrSimple(t *testing.T) {
	// 删除test.db文件
	os.Remove("test.db")
	m, err := NewMgr(KeyValueSetting{
		Property: misc.CreateProperty(MultiSafe, UseCache, UseDB, FullInitLoad),
		SaveType: SqlLiteDB,
		DBAddr:   "test.db",
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

func TestMgrBase(t *testing.T) {
	// 删除test.db文件
	os.Remove("test.db")
	m, err := NewMgr(KeyValueSetting{
		Property: misc.CreateProperty(MultiSafe, UseCache, UseDB, FullInitLoad),
		SaveType: SqlLiteDB,
		DBAddr:   "test.db",
	})
	if err != nil {
		t.Error(err)
		return
	}
	type set struct {
		key string
		v   *ValueUnit
	}
	type get struct {
		key string
	}
	type remove struct {
		key string
	}

	v1 := ToUnit("1", VALUE_TYPE_STRING)
	v2 := ToUnit(2, VALUE_TYPE_INT)
	v3 := ToUnit(float32(3.0), VALUE_TYPE_FLOAT)
	v4 := ToUnit(true, VALUE_TYPE_BOOL)
	v5 := ToUnit([]int{1, 2, 3}, VALUE_TYPE_SLICE_INT)
	v6 := ToUnit([]string{"1", "2", "3"}, VALUE_TYPE_SLICE_STRING)
	v7 := ToUnit([]float32{1.0, 2.0, 3.0}, VALUE_TYPE_SLICE_FLOAT)
	v8 := ToUnit([]bool{true, false, true}, VALUE_TYPE_SLICE_BOOL)

	cases := []*ValueUnit{v1, v2, v3, v4, v5, v6, v7, v8}
	for i, v := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			// get-》set-》get-》remove -》get -》set-》get
			ok, _, err := m.Get(strconv.Itoa(i))
			if err != nil {
				t.Error(err)
				return
			}
			if ok {
				t.Error("get error")
				return
			}
			err = m.Set(strconv.Itoa(i), v)
			if err != nil {
				t.Error(err)
				return
			}
			ok, result, err := m.Get(strconv.Itoa(i))
			if err != nil {
				t.Error(err)
				return
			}
			if !ok || result == nil {
				t.Error("get error")
				return
			}
			if result.Type != v.Type {
				t.Error("type error")
				return
			}
			err = m.Delete(strconv.Itoa(i))
			if err != nil {
				t.Error(err)
				return
			}
			ok, _, err = m.Get(strconv.Itoa(i))
			if err != nil {
				t.Error(err)
				return
			}
			if ok {
				t.Error("get error")
				return
			}
			err = m.Set(strconv.Itoa(i), v)
			if err != nil {
				t.Error(err)
				return
			}
			ok, result, err = m.Get(strconv.Itoa(i))
			if err != nil {
				t.Error(err)
				return
			}
		})
	}
	os.Remove("test.db")
}

// test 多线程

// test slice
