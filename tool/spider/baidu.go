package spider

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
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
		return true
	}
	return false
}

func getBaiduNews(keyword string, limitHour bool) (newsReturn []BaiduNew, reErrorExist bool, noNews bool) {
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
	response, err := client.Do(req)

	if err != nil {
		panic("failed to get : " + err.Error())
	}

	text, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic("failed to read : " + err.Error())
	}
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
