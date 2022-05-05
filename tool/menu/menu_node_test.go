package menu

import (
	"testing"
)

func Test_makeValueStr(t *testing.T) {
	t.Logf(makeValueStr("aaa", ""))
	t.Logf(makeValueStr("aaa", "b"))
	t.Logf(makeValueStr("aaa", "bb"))
	t.Logf(makeValueStr("aaa", "bbb"))
	t.Logf(makeValueStr("aaa", "bbbb"))
}

func TestMakeListInputFunc(t *testing.T) {
	// 手动测试
	kv := make(map[string]string)
	kv["a"] = "1"
	kv["b"] = "2"
	kv["c"] = "3"
	MakeListInputFunc(kv, func() {
		for k, v := range kv {
			t.Logf("%s:%s", k, v)
		}
	})()
}

func TestMakeUniListInputFunc(t *testing.T) {
	kv := make(map[string]interface{})
	kv["test.title"] = "this is a test"
	kv["test.port"] = 1111
	kv["test.host"] = "192.168.0.1"
	kv["test.open"] = true
	kv["test.close"] = false
	kv["test.num"] = 222
	kv["test.float"] = 3.14
	MakeUniListInputFunc(kv, func() {
		for k, v := range kv {
			t.Logf("%s:%v", k, v)
		}
	})()
}
