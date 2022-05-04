package spider

import (
	"github.com/intmian/mian_go_lib/tool/push"
	"testing"
)

func TestGetBaiduNews(t *testing.T) {
	params := []string{
		"nuc",
		"群晖",
		"macbook air",
		"扫地机器人 新品",
		"kindle",
	}
	keywords := []string{}
	newss := [][]BaiduNew{}
	for _, v := range params {
		news, b1, b2 := GetBaiduNews(v, true)
		if b1 == true || b2 == true {
			t.Error("get news error")
		}
		t.Logf("%s: %d", v, len(news))
		keywords = append(keywords, v)
		newss = append(newss, news)
	}

	s := ParseNewToMarkdown(keywords, newss)
	p := push.Mgr{}
	p.SetTag("auto")
	p.SetPushDeerToken("PDU10120Tp8PByEPFdrKiStSvMWeOdeFtwY7GuOmQ")
	p.PushPushDeer("新闻", s, true)
}
