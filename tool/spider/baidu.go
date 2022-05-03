package spider

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type BaiduNew struct {
	title   string
	content string
	source  string
	time    string
}

func timeValid(timeStr string) bool {
	if strings.Contains(timeStr, "分钟前") {
		return true
	} else if strings.Contains(timeStr, "小时前") {
		return true
	} else if strings.Contains(timeStr, "今天") {
		return true
	} else if strings.Contains(timeStr, "昨天") {
		if strings.Contains(timeStr, ":") {
			// 筛选出小时，仅显示24小时内的
			t := strings.Split(timeStr, "昨天")
			t2 := strings.Split(t[1], ":")
			hour, _ := strconv.Atoi(t2[0])
			nowHour := time.Now().Hour()
			if nowHour <= hour {
				return true
			}
		}
	}
	return false
}

func ParseNewToString(new BaiduNew) string {
	return new.title + "\r\n" + new.content + "\r\n" + new.source + "\r\n" + new.time
}

func ParseNewToMarkdown(keywords []string, news [][]BaiduNew) string {
	if len(news) == 0 {
		return ""
	}
	if len(news) != len(keywords) {
		return ""
	}
	s := ""
	for i, keyword := range keywords {
		s += "- " + keyword + "\r\n"
		for _, baiduNew := range news[i] {
			s += "  - " + baiduNew.title + "\r\n"
			s += "    - " + baiduNew.content + "\r\n"
			s += "    - " + baiduNew.source + "\r\n"
			s += "    - " + baiduNew.time + "\r\n"
			s += "\r\n"
		}
	}
	return s
}

func GetBaiduNews(keyword string, limitHour bool) (newsReturn []BaiduNew, reErrorExist bool, noNews bool) {
	newsReturn = make([]BaiduNew, 0)
	reErrorExist = false
	noNews = false

	header := http.Header{"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36"}}
	u := "https://www.baidu.com/s?tn=news&rtt=4&bsst=1&cl=2&wd=" + keyword
	httpUrl, _ := url.Parse(u)
	req := &http.Request{
		Method: "GET",
		URL:    httpUrl,
		Header: header,
	}
	client := &http.Client{}
	response, _ := client.Do(req)

	text, _ := ioutil.ReadAll(response.Body)
	reStr := `"accessibilityData":{"titleAriaLabel":"标题[： ](.*)","absAriaLabel":"摘要[： ](.*)","sourceAriaLabel":"新闻来源[： ](.*)","timeAriaLabel":"发布于[： ](.{0,20})"}`
	reg1 := regexp.MustCompile(reStr)
	if reg1 == nil {
		return
	}
	//根据规则提取关键信息
	results := reg1.FindAllStringSubmatch(string(text), -1)
	if len(results) == 0 {
		noNews = true
	}
	for _, result := range results {
		bn := BaiduNew{}
		if len(result) != 5 {
			reErrorExist = true
			continue
		}
		bn.title = result[1]
		bn.content = result[2]
		bn.content = strings.Replace(bn.content, " 摘要结束，点击查看详情", "...", -1)
		bn.source = result[3]
		bn.time = result[4]
		if limitHour {
			if !timeValid(bn.time) {
				continue
			}
		}
		newsReturn = append(newsReturn, bn)
	}
	return
}
