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
	m, err := NewXStorage(XstorageSetting{
		Property: misc.CreateProperty(MultiSafe, UseCache, UseDisk, FullInitLoad),
		SaveType: SqlLiteDB,
		DBAddr:   "test.db",
	})
	if err != nil {
		t.Error(err)
		return
	}
	err = m.Set("1", ToUnit("1", ValueTypeString))
	if err != nil {
		t.Error(err)
		return
	}
	err = m.Set("2", ToUnit(2, ValueTypeInt))
	if err != nil {
		t.Error(err)
		return
	}
	err = m.Set("3", ToUnit(float32(3.0), ValueTypeFloat))
	if err != nil {
		t.Error(err)
		return
	}
	v := &ValueUnit{}
	ok, err := m.GetHP("1", v)
	if err != nil {
		return
	}
	if !ok {
		t.Error("not ok")
		return
	}
	if v.Type != ValueTypeString {
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
	m, err := NewXStorage(XstorageSetting{
		Property: misc.CreateProperty(MultiSafe, UseCache, UseDisk, FullInitLoad),
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

	v1 := ToUnit("1", ValueTypeString)
	v2 := ToUnit(2, ValueTypeInt)
	v3 := ToUnit(float32(3.0), ValueTypeFloat)
	v4 := ToUnit(true, ValueTypeBool)
	v5 := ToUnit([]int{1, 2, 3}, ValueTypeSliceInt)
	v6 := ToUnit([]string{"1", "2", "3"}, ValueTypeSliceString)
	v7 := ToUnit([]float32{1.0, 2.0, 3.0}, ValueTypeSliceFloat)
	v8 := ToUnit([]bool{true, false, true}, ValueTypeSliceBool)

	cases := []*ValueUnit{v1, v2, v3, v4, v5, v6, v7, v8}
	for i, v := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			// get-》set-》get-》remove -》get -》set-》get -》remove -》get_default -》set-》get_default
			temp := &ValueUnit{}
			ok, err := m.GetHP(strconv.Itoa(i), temp)
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
			result := &ValueUnit{}
			ok, err = m.GetHP(strconv.Itoa(i), result)
			if err != nil {
				t.Error(err)
				return
			}
			if !ok {
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
			ok, err = m.GetHP(strconv.Itoa(i), result)
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
			ok, err = m.GetHP(strconv.Itoa(i), result)
			if err != nil {
				t.Error(err)
				return
			}
			if !ok {
				t.Error("get error")
				return
			}
			err = m.Delete(strconv.Itoa(i))
			if err != nil {
				t.Error(err)
				return
			}
			ok, err = m.GetAndSetDefaultHP(strconv.Itoa(i), v, result)
			if err != nil {
				t.Error(err)
				return
			}
			if !ok {
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
			v2 := ToUnit("23333", ValueTypeString)
			ok, err = m.GetAndSetDefaultHP(strconv.Itoa(i), v2, result)
			if err != nil {
				t.Error(err)
				return
			}
			if !ok {
				t.Error("get error")
				return
			}
			if result.Type != v.Type {
				t.Error("type error")
				return
			}
			if result.Type == ValueTypeString && ToBase[string](result) == "23333" {
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
	m, err := NewXStorage(XstorageSetting{
		Property: misc.CreateProperty(MultiSafe, UseCache, UseDisk, FullInitLoad),
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
		datas[i] = ToUnit(i, ValueTypeInt)
	}
	t1 := time.Now()
	c := make(chan chan error, testNum)
	for i := 0; i < testNum*2; i++ {
		if i%2 == 0 {
			go func(i int) {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(10)))
				err, c2 := m.SetAsync(strconv.Itoa(i/2), datas[i/2])
				if err != nil {
					t.Error(err)
					return
				}
				c <- c2
			}(i)
		} else {
			go func(i int) {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(10)))
				temp := &ValueUnit{}
				_, err := m.GetHP(strconv.Itoa(i/2), temp)
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
	if (t2.Sub(t1) > time.Second*2) || (t3.Sub(t2) > time.Second*10) {
		t.Error("time error")
		return
	}
	//time.Sleep(time.Second * 2)
	for i := 0; i < testNum; i++ {
		// 取出来判断是否正确
		result := &ValueUnit{}
		ok, err := m.GetHP(strconv.Itoa(i), result)
		if err != nil {
			t.Error(err)
			return
		}
		if !ok {
			t.Errorf("get error %d", i)
			return
		}
		if result.Type != ValueTypeInt {
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
	m, err := NewXStorage(XstorageSetting{
		SaveType: SqlLiteDB,
		DBAddr:   "test4.db",
		Property: misc.CreateProperty(MultiSafe, UseCache, UseDisk, FullInitLoad),
	})
	if err != nil {
		t.Error(err)
		return
	}
	models := make([]*ValueUnit, 0)
	models = append(models, ToUnit("1", ValueTypeString))
	models = append(models, ToUnit(2, ValueTypeInt))
	models = append(models, ToUnit(float32(3.0), ValueTypeFloat))
	models = append(models, ToUnit(true, ValueTypeBool))
	models = append(models, ToUnit([]int{1, 2, 3}, ValueTypeSliceInt))
	models = append(models, ToUnit([]string{"1", "2", "3"}, ValueTypeSliceString))
	models = append(models, ToUnit([]float32{1.0, 2.0, 3.0}, ValueTypeSliceFloat))
	models = append(models, ToUnit([]bool{true, false, true}, ValueTypeSliceBool))
	for i, v := range models {
		err := m.Set(strconv.Itoa(i), v)
		if err != nil {
			t.Error(err)
			return
		}
	}

	m2, err := NewXStorage(XstorageSetting{
		SaveType: SqlLiteDB,
		DBAddr:   "test4.db",
		Property: misc.CreateProperty(MultiSafe, UseCache, UseDisk, FullInitLoad),
	})
	if err != nil {
		t.Error(err)
		return
	}
	for i, v := range models {
		result := &ValueUnit{}
		ok, err := m2.GetHP(strconv.Itoa(i), result)
		if err != nil {
			t.Error(err)
			return
		}
		if !ok {
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
	mgr1, _ := NewXStorage(XstorageSetting{
		SaveType: SqlLiteDB,
		DBAddr:   "test5.db",
		Property: misc.CreateProperty(MultiSafe, UseCache, UseDisk, FullInitLoad),
	})
	mgr2, _ := NewXStorage(XstorageSetting{
		SaveType: SqlLiteDB,
		DBAddr:   "test5.db",
		Property: misc.CreateProperty(MultiSafe, UseDisk),
	})
	mgr3, _ := NewXStorage(XstorageSetting{
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
		err := mgr1.Set(`testSlice`, ToUnit(a, ValueTypeSliceInt))
		if err != nil {
			t.Error(err)
			return
		}
		v2 := &ValueUnit{}
		ok, err := mgr1.GetHP(`testSlice`, v2)
		if err != nil {
			t.Error(err)
			return
		}
		if !ok {
			t.Error("get error")
			return
		}
		if !Compare(ToUnit(a, ValueTypeSliceInt), v2) {
			t.Error("value error")
			return
		}
		ok, err = mgr2.GetHP(`testSlice`, v2)
		if err != nil {
			t.Error(err)
			return
		}
		if !ok {
			t.Error("get error")
			return
		}
		if !Compare(ToUnit(a, ValueTypeSliceInt), v2) {
			t.Error("value error")
			return
		}
		err = mgr3.Set(`testSlice`, ToUnit(a, ValueTypeSliceInt))
		if err != nil {
			t.Error(err)
			return
		}
		ok, err = mgr3.GetHP(`testSlice`, v2)
		if err != nil {
			t.Error(err)
			return
		}
		if !ok {
			t.Error("get error")
			return
		}
		if !Compare(ToUnit(a, ValueTypeSliceInt), v2) {
			t.Error("value error")
			return
		}
	}
}

func TestMgrToml(t *testing.T) {
	// 删除test.db文件
	os.Remove("test6.toml")
	os.Remove("test7.toml")
	m, _ := NewXStorage(XstorageSetting{
		Property: misc.CreateProperty(MultiSafe, UseCache, UseDisk, FullInitLoad),
		SaveType: Toml,
		FileAddr: "test6.json",
	})
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

	v1 := ToUnit("1", ValueTypeString)
	v2 := ToUnit(2, ValueTypeInt)
	v3 := ToUnit(float32(3.0), ValueTypeFloat)
	v4 := ToUnit(true, ValueTypeBool)
	v5 := ToUnit([]int{1, 2, 3}, ValueTypeSliceInt)
	v6 := ToUnit([]string{"1", "2", "3"}, ValueTypeSliceString)
	v7 := ToUnit([]float32{1.0, 2.0, 3.0}, ValueTypeSliceFloat)
	v8 := ToUnit([]bool{true, false, true}, ValueTypeSliceBool)

	cases := []*ValueUnit{v1, v2, v3, v4, v5, v6, v7, v8}
	for i, v := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			// get-》set-》get-》remove -》get -》set-》get -》remove -》get_default -》set-》get_default
			temp := &ValueUnit{}
			ok, err := m.GetHP(strconv.Itoa(i), temp)
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
			result := &ValueUnit{}
			ok, err = m.GetHP(strconv.Itoa(i), result)
			if err != nil {
				t.Error(err)
				return
			}
			if !ok {
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
			ok, err = m.GetHP(strconv.Itoa(i), result)
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
			ok, err = m.GetHP(strconv.Itoa(i), result)
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
			ok, err = m.GetAndSetDefaultHP(strconv.Itoa(i), v, result)
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
			v2 := ToUnit("23333", ValueTypeString)
			ok, err = m.GetAndSetDefaultHP(strconv.Itoa(i), v2, result)
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
			if result.Type == ValueTypeString && ToBase[string](result) == "23333" {
				t.Error("value error")
				return
			}

			m2, _ := NewXStorage(XstorageSetting{
				Property: misc.CreateProperty(MultiSafe, UseCache, UseDisk, FullInitLoad),
				SaveType: Toml,
				FileAddr: "test7.json",
			})
			ok, err = m2.GetHP(strconv.Itoa(i), result)
			if err != nil {
				t.Error(err)
				return
			}
			if ok {
				t.Error("get error")
				return
			}
			m2, _ = NewXStorage(XstorageSetting{
				Property: misc.CreateProperty(MultiSafe, UseCache, UseDisk, FullInitLoad),
				SaveType: Toml,
				FileAddr: "test7.json",
			})
			err = m2.Set(strconv.Itoa(i), v)
			if err != nil {
				t.Error(err)
				return
			}
			m2, _ = NewXStorage(XstorageSetting{
				Property: misc.CreateProperty(MultiSafe, UseCache, UseDisk, FullInitLoad),
				SaveType: Toml,
				FileAddr: "test7.json",
			})
			ok, err = m2.GetHP(strconv.Itoa(i), result)
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
			m2, _ = NewXStorage(XstorageSetting{
				Property: misc.CreateProperty(MultiSafe, UseCache, UseDisk, FullInitLoad),
				SaveType: Toml,
				FileAddr: "test7.json",
			})
			err = m2.Delete(strconv.Itoa(i))
			if err != nil {
				t.Error(err)
				return
			}
			m2, _ = NewXStorage(XstorageSetting{
				Property: misc.CreateProperty(MultiSafe, UseCache, UseDisk, FullInitLoad),
				SaveType: Toml,
				FileAddr: "test7.json",
			})
			ok, err = m2.GetHP(strconv.Itoa(i), result)
			if err != nil {
				t.Error(err)
				return
			}
			if ok {
				t.Error("get error")
				return
			}
			m2, _ = NewXStorage(XstorageSetting{
				Property: misc.CreateProperty(MultiSafe, UseCache, UseDisk, FullInitLoad),
				SaveType: Toml,
				FileAddr: "test7.json",
			})
			err = m2.Set(strconv.Itoa(i), v)
			if err != nil {
				t.Error(err)
				return
			}
			m2, _ = NewXStorage(XstorageSetting{
				Property: misc.CreateProperty(MultiSafe, UseCache, UseDisk, FullInitLoad),
				SaveType: Toml,
				FileAddr: "test7.json",
			})
			ok, err = m2.GetHP(strconv.Itoa(i), result)
			if err != nil {
				t.Error(err)
				return
			}
			if !ok || result == nil {
				t.Error("get error")
				return
			}
			m2, _ = NewXStorage(XstorageSetting{
				Property: misc.CreateProperty(MultiSafe, UseCache, UseDisk, FullInitLoad),
				SaveType: Toml,
				FileAddr: "test7.json",
			})
			err = m2.Delete(strconv.Itoa(i))
			if err != nil {
				t.Error(err)
				return
			}
			m2, _ = NewXStorage(XstorageSetting{
				Property: misc.CreateProperty(MultiSafe, UseCache, UseDisk, FullInitLoad),
				SaveType: Toml,
				FileAddr: "test7.json",
			})
			ok, err = m2.GetAndSetDefaultHP(strconv.Itoa(i), v, result)
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
			m2, _ = NewXStorage(XstorageSetting{
				Property: misc.CreateProperty(MultiSafe, UseCache, UseDisk, FullInitLoad),
				SaveType: Toml,
				FileAddr: "test7.json",
			})
			err = m2.Set(strconv.Itoa(i), v)
			if err != nil {
				t.Error(err)
				return
			}
			v2 = ToUnit("23333", ValueTypeString)
			ok, err = m2.GetAndSetDefaultHP(strconv.Itoa(i), v2, result)
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
			if result.Type == ValueTypeString && ToBase[string](result) == "23333" {
				t.Error("value error")
				return
			}
		})
	}
	os.Remove("test6.toml")
	os.Remove("test7.toml")
}

func TestMem(t *testing.T) {
	m, err := NewXStorage(XstorageSetting{
		Property: misc.CreateProperty(UseCache),
	})
	if err != nil {
		t.Error(err)
		return
	}
	/*1.随机生成1000组数据，并行写入
	  2.读取所有数据是否都正常
	  3.读取内存数据
	  4.随机删除500组数据
	  5.读取内存数据
	  6. 随机生成1000组数据
	  7. 读取所有数据是否都正常
	  8. 删除所有数据，检查内存
	*/
	testNum1 := 100000
	testNum2 := 50000
	tesetNum3 := 100000
	type TestCase struct {
		key         string
		RealValue   *ValueUnit
		IsFirstAdd  bool
		IsDelete    bool
		IsSecondAdd bool
	}
	testCases := make([]*TestCase, 0)
	for i := 0; i < testNum1; i++ {
		testCases = append(testCases, &TestCase{
			key:         strconv.Itoa(i),
			RealValue:   ToUnit(i, ValueTypeInt),
			IsFirstAdd:  true,
			IsDelete:    false,
			IsSecondAdd: false,
		})
	}
	for _, testCase := range testCases {
		if testCase.IsFirstAdd {
			err := m.Set(testCase.key, testCase.RealValue)
			if err != nil {
				t.Error(err)
				return
			}
		}
	}
	for _, testCase := range testCases {
		if testCase.IsFirstAdd {
			result := &ValueUnit{}
			ok, err := m.GetHP(testCase.key, result)
			if err != nil {
				t.Error(err)
				return
			}
			if !ok {
				t.Error("get error")
				return
			}
			if !Compare(testCase.RealValue, result) {
				t.Error("value error")
				return
			}
		}
	}
	for i := 0; i < testNum2; i++ {
		index := rand.Intn(testNum1)
		for testCases[index].IsDelete {
			index = rand.Intn(testNum1)
		}
		testCases[index].IsDelete = true
		err := m.Delete(testCases[index].key)
		if err != nil {
			t.Error(err)
			return
		}
	}
	for _, testCase := range testCases {
		if testCase.IsDelete {
			result := &ValueUnit{}
			ok, err := m.GetHP(testCase.key, result)
			if ok {
				t.Error("get error")
				return
			} else {
				if err != nil {
					t.Error(err)
					return
				}
			}
		}
	}
	for i := testNum1; i < testNum1+tesetNum3; i++ {
		testCases = append(testCases, &TestCase{
			key:         strconv.Itoa(i),
			RealValue:   ToUnit(i, ValueTypeInt),
			IsFirstAdd:  false,
			IsDelete:    false,
			IsSecondAdd: true,
		})
	}
	for _, testCase := range testCases {
		if testCase.IsSecondAdd {
			err := m.Set(testCase.key, testCase.RealValue)
			if err != nil {
				t.Error(err)
				return
			}
		}
	}
	for _, testCase := range testCases {
		if testCase.IsSecondAdd {
			result := &ValueUnit{}
			ok, err := m.GetHP(testCase.key, result)
			if err != nil {
				t.Error(err)
				return
			}
			if !ok {
				t.Error("get error")
				return
			}
			if !Compare(testCase.RealValue, result) {
				t.Error("value error")
				return
			}
		}
	}
	for _, testCase := range testCases {
		if testCase.IsSecondAdd || (testCase.IsFirstAdd && !testCase.IsDelete) {
			err := m.Delete(testCase.key)
			if err != nil {
				t.Logf("delete error %v", testCase)
				t.Error(err)
				return
			}
		}
	}
	for _, testCase := range testCases {
		if testCase.IsSecondAdd || (testCase.IsFirstAdd && !testCase.IsDelete) {
			result := &ValueUnit{}
			ok, err := m.GetHP(testCase.key, result)
			if err != nil {
				t.Error(err)
				return
			}
			if ok {
				t.Error("get error")
				return
			}
		}
	}
}
