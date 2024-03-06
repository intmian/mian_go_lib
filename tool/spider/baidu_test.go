package spider

import (
	"testing"
)

func TestGetBaiduNews(t *testing.T) {
	params := []string{
		//"gta6",
		//"iphone16 ",
		//"ps5pro ",
		//"apple glass",
		//"iphone16",
		"特斯拉",
	}
	keywords := []string{}
	newss := [][]BaiduNew{}
	for _, v := range params {
		news, err, _ := GetTodayBaiduNews(v)
		if err != nil {
			t.Error(err)
			return
		}
		t.Logf("%s: %d", v, len(news))
		keywords = append(keywords, v)
		newss = append(newss, news)
	}
	//for _, news := range newss {
	//	for _, baiduNew := range news {
	//		t.Log(baiduNew.title)
	//		t.Log(baiduNew.content)
	//		t.Log(baiduNew.source)
	//		t.Log(baiduNew.time)
	//		t.Log(baiduNew.valid)
	//		//println(baiduNew.title)
	//		//println(baiduNew.content)
	//		//println(baiduNew.source)
	//		//println(baiduNew.time)
	//		//println(baiduNew.valid)
	//	}
	//}
	_ = ParseNewToMarkdown(keywords, newss)
	//f, _ := os.Create("baidu.md")
	//f.WriteString(s)
	//p := xpush.Mgr{}
	//p.SetTag("auto")
	//p.SetPushDeerToken("PDU10120Tp8PByEPFdrKiStSvMWeOdeFtwY7GuOmQ")
	//p.PushPushDeer("新闻", s, true)
}
