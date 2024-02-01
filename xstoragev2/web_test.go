package xstorage

import (
	"fmt"
	"github.com/intmian/mian_go_lib/tool/misc"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestWeb(t *testing.T) {
	m, err := NewXStorage(XstorageSetting{
		Property: misc.CreateProperty(MultiSafe, UseCache),
	})
	w, err := NewWebPack(WebPackSetting{
		webPort: 11111,
	}, m)
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		err = w.StartWeb()
		if err != nil {
			t.Error(err)
			return
		}
	}()
	time.Sleep(time.Second)
	u1 := ToUnit("1", ValueTypeString)
	// 用http set设置
	url := "http://127.0.0.1:11111"
	urlSet := fmt.Sprintf("%s/set?key=%s&value=%s&value_type=%d", url, "test", ToBase[string](u1), u1.Type)
	resp, err := http.Get(urlSet)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatal("web set error")
	}
	// 打印返回的json
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "{\"code\":1}" {
		t.Fatal("web set error")
	}

	// 用http get获取
	urlGet := fmt.Sprintf("%s/get?perm=%s", url, "test")
	resp, err = http.Get(urlGet)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatal("web get error")
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "{\"code\":1,\"result\":[{\"Type\":1,\"Data\":\"1\"}]}" {
		t.Fatal("web get error")
	}
}

func TestMgr_WebMa(t *testing.T) {
	return
	m, err := NewXStorage(XstorageSetting{
		Property: misc.CreateProperty(MultiSafe, UseCache),
	})
	w, err := NewWebPack(WebPackSetting{
		webPort: 11111,
	}, m)
	if err != nil {
		t.Fatal(err)
	}
	err = w.StartWeb()
	if err != nil {
		t.Error(err)
		return
	}
}
