package spider

import "testing"

func TestGetNews(t *testing.T) {
	params := []string{
		"nuc",
		"群晖",
		"macbook",
		"扫地机器人",
		"kindle",
	}
	for _, v := range params {
		news, b1, b2 := getBaiduNews(v, true)

		if b1 == true || b2 == true {
			t.Error("get news error")
		}
		t.Logf("%s: %d", v, len(news))
	}

	for _, v := range params {
		news, b1, b2 := getBaiduNews(v, false)

		if b1 == true || b2 == true {
			t.Error("get news error")
		}
		t.Logf("%s: %d", v, len(news))
	}
}
