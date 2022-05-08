package spider

import (
	"github.com/intmian/mian_go_lib/tool/xpush"
	"testing"
)

func TestDapan000001(t *testing.T) {
	price, inc, radio := GetDapan000001()
	if price == "" || inc == "" || radio == "" {
		t.Error("GetDapan000001 error")
	}
	s := ParseDapanToMarkdown("上证指数", price, inc, radio)
	p := xpush.Mgr{}
	p.SetTag("auto")
	p.SetPushDeerToken("PDU10120Tp8PByEPFdrKiStSvMWeOdeFtwY7GuOmQ")
	p.PushPushDeer("大盘", s, true)
}
