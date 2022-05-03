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
	return

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
