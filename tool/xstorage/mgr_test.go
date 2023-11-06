package xstorage

import (
	"github.com/intmian/mian_go_lib/tool/misc"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
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
	os.Remove("test.db")
}

func TestMgrBase(t *testing.T) {
	// 删除test.db文件
	os.Remove("test2.db")
	m, err := NewMgr(KeyValueSetting{
		Property: misc.CreateProperty(MultiSafe, UseCache, UseDB, FullInitLoad),
		SaveType: SqlLiteDB,
		DBAddr:   "test2.db",
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
			// get-》set-》get-》remove -》get -》set-》get -》remove -》get_default -》set-》get_default
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
			if !ok || result == nil {
				t.Error("get error")
				return
			}
			err = m.Delete(strconv.Itoa(i))
			if err != nil {
				t.Error(err)
				return
			}
			ok, result, err = m.GetAndSetDefault(strconv.Itoa(i), v)
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
			err = m.Set(strconv.Itoa(i), v)
			if err != nil {
				t.Error(err)
				return
			}
			v2 := ToUnit("23333", VALUE_TYPE_STRING)
			ok, result, err = m.GetAndSetDefault(strconv.Itoa(i), v2)
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
			if result.Type == VALUE_TYPE_STRING && ToBase[string](result) == "23333" {
				t.Error("value error")
				return
			}
		})
	}
	os.Remove("test2.db")
}

// test 多线程
func TestMgrMulti(t *testing.T) {
	os.Remove("test3.db")
	m, err := NewMgr(KeyValueSetting{
		Property: misc.CreateProperty(MultiSafe, UseCache, UseDB, FullInitLoad),
		SaveType: SqlLiteDB,
		DBAddr:   "test3.db",
	})
	if err != nil {
		t.Error(err)
		return
	}
	testNum := 1000
	datas := make([]*ValueUnit, testNum)
	for i := 0; i < testNum; i++ {
		datas[i] = ToUnit(i, VALUE_TYPE_INT)
	}
	t1 := time.Now()
	c := make(chan chan error, testNum)
	for i := 0; i < testNum*2; i++ {
		if i%2 == 0 {
			go func(i int) {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(10)))
				err, c2 := m.SetAsyncDB(strconv.Itoa(i/2), datas[i/2])
				if err != nil {
					t.Error(err)
					return
				}
				c <- c2
			}(i)
		} else {
			go func(i int) {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(10)))
				_, _, err := m.Get(strconv.Itoa(i / 2))
				if err != nil {
					t.Error(err)
					return
				}
			}(i)
		}
	}
	t2 := time.Now()
	chanErrors := make([]chan error, testNum)
	for i := 0; i < testNum; i++ {
		chanErrors[i] = <-c
	}
	for i := 0; i < testNum; i++ {
		err := <-chanErrors[i]
		if err != nil {
			t.Error(err)
			return
		}
		//if i%10 == 0 {
		//	t.Logf("set %d", i)
		//}
	}
	t3 := time.Now()
	t.Logf("set time %v, real set time %v", t2.Sub(t1), t3.Sub(t2))
	//time.Sleep(time.Second * 2)
	for i := 0; i < testNum; i++ {
		// 取出来判断是否正确
		ok, result, err := m.Get(strconv.Itoa(i))
		if err != nil {
			t.Error(err)
			return
		}
		if !ok || result == nil {
			t.Errorf("get error %d", i)
			return
		}
		if result.Type != VALUE_TYPE_INT {
			t.Error("type error")
			return
		}
		if ToBase[int](result) != i {
			t.Error("value error")
			return
		}
	}
	os.Remove("test3.db")
}

// test 重启getall
func TestMgrReBoot(t *testing.T) {
	/*
		先写入所有类型的数据，然后再新建一个mgr，看是否能够正确读取所有的数据
	*/
	os.Remove("test4.db")
	m, err := NewMgr(KeyValueSetting{
		SaveType: SqlLiteDB,
		DBAddr:   "test4.db",
		Property: misc.CreateProperty(MultiSafe, UseCache, UseDB, FullInitLoad),
	})
	if err != nil {
		t.Error(err)
		return
	}
	models := make([]*ValueUnit, 0)
	models = append(models, ToUnit("1", VALUE_TYPE_STRING))
	models = append(models, ToUnit(2, VALUE_TYPE_INT))
	models = append(models, ToUnit(float32(3.0), VALUE_TYPE_FLOAT))
	models = append(models, ToUnit(true, VALUE_TYPE_BOOL))
	models = append(models, ToUnit([]int{1, 2, 3}, VALUE_TYPE_SLICE_INT))
	models = append(models, ToUnit([]string{"1", "2", "3"}, VALUE_TYPE_SLICE_STRING))
	models = append(models, ToUnit([]float32{1.0, 2.0, 3.0}, VALUE_TYPE_SLICE_FLOAT))
	models = append(models, ToUnit([]bool{true, false, true}, VALUE_TYPE_SLICE_BOOL))
	for i, v := range models {
		err := m.Set(strconv.Itoa(i), v)
		if err != nil {
			t.Error(err)
			return
		}
	}

	m2, err := NewMgr(KeyValueSetting{
		SaveType: SqlLiteDB,
		DBAddr:   "test4.db",
		Property: misc.CreateProperty(MultiSafe, UseCache, UseDB, FullInitLoad),
	})
	if err != nil {
		t.Error(err)
		return
	}
	for i, v := range models {
		ok, result, err := m2.Get(strconv.Itoa(i))
		if err != nil {
			t.Error(err)
			return
		}
		if !ok || result == nil {
			t.Error("get error")
			return
		}
		if !Compare(v, result) {
			t.Errorf("compare error %v %v", v, result)
		}
	}
	os.Remove("test4.db")
}

// test slice
func TestMgrSlice(t *testing.T) {
	os.Remove("test5.db")
	defer os.Remove("test5.db")
	// 伸长 缩短 或改部分值
	mgr1, _ := NewMgr(KeyValueSetting{
		SaveType: SqlLiteDB,
		DBAddr:   "test5.db",
		Property: misc.CreateProperty(MultiSafe, UseCache, UseDB, FullInitLoad),
	})
	mgr2, _ := NewMgr(KeyValueSetting{
		SaveType: SqlLiteDB,
		DBAddr:   "test5.db",
		Property: misc.CreateProperty(MultiSafe, UseDB),
	})
	mgr3, _ := NewMgr(KeyValueSetting{
		Property: misc.CreateProperty(MultiSafe, UseCache),
	})
	a1 := []int{1, 2, 3}
	a2 := []int{1, 2, 3, 4}
	a3 := []int{1, 2, 3, 4, 5}
	a4 := []int{1, 2, 3, 4}
	a5 := []int{1, 6, 3}
	a6 := []int{1, 2, 3, 4, 5}
	a7 := []int{1, 2, 3, 3, 5}
	as := [][]int{a1, a2, a3, a4, a5, a6, a7}
	for _, a := range as {
		err := mgr1.Set(`testSlice`, ToUnit(a, VALUE_TYPE_SLICE_INT))
		if err != nil {
			t.Error(err)
			return
		}
		ok, v2, err := mgr1.Get(`testSlice`)
		if err != nil {
			t.Error(err)
			return
		}
		if !ok {
			t.Error("get error")
			return
		}
		if !Compare(ToUnit(a, VALUE_TYPE_SLICE_INT), v2) {
			t.Error("value error")
			return
		}
		ok, v2, err = mgr2.Get(`testSlice`)
		if err != nil {
			t.Error(err)
			return
		}
		if !ok {
			t.Error("get error")
			return
		}
		if !Compare(ToUnit(a, VALUE_TYPE_SLICE_INT), v2) {
			t.Error("value error")
			return
		}
		err = mgr3.Set(`testSlice`, ToUnit(a, VALUE_TYPE_SLICE_INT))
		if err != nil {
			t.Error(err)
			return
		}
		ok, v2, err = mgr3.Get(`testSlice`)
		if err != nil {
			t.Error(err)
			return
		}
		if !ok {
			t.Error("get error")
			return
		}
		if !Compare(ToUnit(a, VALUE_TYPE_SLICE_INT), v2) {
			t.Error("value error")
			return
		}
	}
}
