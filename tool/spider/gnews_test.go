package spider

import (
	"github.com/intmian/mian_go_lib/tool/ai"
	"os"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	now := time.Now()
	s := GetUniTimeStr(now)
	t.Log(s)
}

func TestGetGNews(t *testing.T) {
	req := GNewsSearch{
		q:      "tesla",
		lang:   LanChinese,
		sortby: SortByPublishedAt,
		from:   GetUniTimeStr(time.Now().AddDate(0, -1, -10)),
		to:     GetUniTimeStr(time.Now()),
	}
	result, err := QueryGNews(req, "ee54b7595ba81fc612c56689416abf6a")
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}

func TestGetGTop(t *testing.T) {
	req := GNewsTop{
		Lang: LanEnglish,
		From: GetUniTimeStr(time.Now().AddDate(0, 0, -1)),
		To:   GetUniTimeStr(time.Now()),
	}
	result, err := QueryGNewsTop(req, "ee54b7595ba81fc612c56689416abf6a")
	s := ""
	for _, v := range result.Articles {
		s += "title" + v.Title + "\n"
		s += "Description" + v.Description + "\n"
	}
	o := ai.NewOpenAI("https://api.openai-proxy.org/v1", "sk-7A4jfmtJ3QXhef0x9g9YIIxTLwK15C9T0vTsehdBlNLExMxk", false, "你是一台由mian研发的新闻机器人")
	re, err := o.Chat("" +
		"以下是一天内发生的热点新闻。" +
		"你是一位资深记者，请根据这些内容写一篇通信稿，要求输出markdown格式。首先使用中文在300字以内汇总以下新闻的内容，要求言语通顺、优美、专业，具备文学美感。最后用50字进行分析与评论：\n" + s)
	f, err := os.Create("openai_test.txt")
	f.WriteString(s)
	f.WriteString(re)
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}

func TestGetGNewsSumTop(t *testing.T) {
	result, err := GetGNewsSumTop("ee54b7595ba81fc612c56689416abf6a", time.Now().AddDate(0, 0, -7))
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}
