package spider

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

const lotteryRe = `.*<li><div class="zxkjc1"><span>(.*)<\/span><b>(.*)<\/b>(.*)<\/div><\/li>
.*<li><div class="zxkjc2">
((( *\n)* *<i.*>(.*)<\/i>\n)*)`

const ballRe = `\d{1,3}`

type Lottery struct {
	Issue  string
	name   string
	t      string
	Number []string
}

func GetLottery() []Lottery {
	header := http.Header{"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36"}}
	u := "http://www2.zgzcw.com"
	httpUrl, _ := url.Parse(u)
	req := &http.Request{
		Method: "GET",
		URL:    httpUrl,
		Header: header,
	}
	client := &http.Client{}
	response, _ := client.Do(req)

	text, _ := ioutil.ReadAll(response.Body)
	reg1 := regexp.MustCompile(lotteryRe)
	if reg1 == nil {
		return nil
	}
	//根据规则提取关键信息
	results := reg1.FindAllStringSubmatch(string(text), -1)
	if len(results) == 0 {
		return nil
	}
	var lotteries []Lottery
	for _, result := range results {
		l := Lottery{}

		l.name = result[1]
		l.Issue = result[2]
		l.t = result[3]
		l.Number = make([]string, 0)
		reg2 := regexp.MustCompile(ballRe)
		balls := reg2.FindAllStringSubmatch(result[4], -1)
		for _, ball := range balls {
			l.Number = append(l.Number, ball[0])
		}
		lotteries = append(lotteries, l)
	}
	return lotteries
}

func ParseLotteriesToMarkDown(lotteries []Lottery) string {
	var s string
	for _, l := range lotteries {
		s += "- " + l.name + l.Issue + "\r\n"
		s += "  - " + l.t + "\r\n"
		s += "  - "
		for _, n := range l.Number {
			s += n + " "
		}
		s += "\r\n"
	}
	return s
}
