package spider

import (
	"fmt"
	"github.com/antlabs/strsim"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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
	valid   float64
	href    string
}

func timeValid(timeStr string) bool {
	if timeStr == "undefined" {
		// 刚刚刷新出的新闻可能这样
		return true
	}
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
	dateTime := time.Now().Format("2006-01-02")
	s := "## " + dateTime + " 新闻汇总\r\n"
	for i, keyword := range keywords {
		newsNum := len(news[i])
		if newsNum == 0 {
			s += "### " + keyword + " 无新闻\r\n"
			continue
		}
		s += "### " + keyword + " " + strconv.Itoa(newsNum) + "条新闻\r\n"
		for _, baiduNew := range news[i] {
			baiduNew.time = strings.Replace(baiduNew.time, "undefined", "近期", -1)
			// 来源 时间：标题（链接）
			s += fmt.Sprintf("- %s %s：[%s](%s) \r\n", baiduNew.source, baiduNew.time, baiduNew.title, baiduNew.href)
		}
	}
	return s
}

func CutInvalidNews(news []BaiduNew, valid float64) []BaiduNew {
	newsReturn := make([]BaiduNew, 0)
	for _, baiduNew := range news {
		if baiduNew.valid < valid {
			newsReturn = append(newsReturn, baiduNew)
		}
	}
	return newsReturn
}

func getBaiduNewsPage(keyword string, page int) (result []BaiduNew, err error) {
	result = make([]BaiduNew, 0)
	header := http.Header{"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36"}}
	u := "https://www.baidu.com/s?tn=news&rtt=4&bsst=1&cl=2&wd=" + keyword
	if page > 1 {
		u += "&pn=" + strconv.Itoa((page-1)*10)
	}
	httpUrl, _ := url.Parse(u)
	req := &http.Request{
		Method: "GET",
		URL:    httpUrl,
		Header: header,
	}
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, "client.Do error")
	}
	text, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.WithMessage(err, "ioutil.ReadAll error")
	}
	reStr := `\{"titleAriaLabel":"标题[： ](.*)","absAriaLabel":"摘要[： ](.*)","sourceAriaLabel":"新闻来源[： ](.*)","timeAriaLabel":"发布于[： ](.{0,20})"\}.*href="(.*)" target`
	reg1 := regexp.MustCompile(reStr)
	if reg1 == nil {
		return
	}
	//根据规则提取关键信息
	results := reg1.FindAllStringSubmatch(string(text), -1)
	if len(results) == 0 {
		f, _ := os.Create(fmt.Sprintf("baidu_%s_%d_%s.html", keyword, page, time.Now().Format("2006-01-02_15:04:05")))
		f.WriteString(string(text))
		return nil, errors.New("no news")
	}
	for _, result2 := range results {
		bn := BaiduNew{}
		if len(result2) != 6 {
			continue
		}
		bn.title = result2[1]
		bn.content = result2[2]
		bn.content = strings.Replace(bn.content, " 摘要结束，点击查看详情", "...", -1)
		bn.source = result2[3]
		bn.time = result2[4]
		bn.href = result2[5]
		result = append(result, bn)
	}
	return
}

func GetTodayBaiduNews(keyword string) (newsReturn []BaiduNew, err error) {
	newsReturn = make([]BaiduNew, 0)
	keyword = strings.Replace(keyword, " ", "+", -1)

	page := 1
	for {
		retryTimes := 10
		var news []BaiduNew
		for retryTimes > 0 {
			news, err = getBaiduNewsPage(keyword, page)
			if err == nil {
				break
			}
			retryTimes--
			// 休眠一分钟
			time.Sleep(time.Minute)
		}
		if err != nil {
			return nil, errors.WithMessage(err, "getBaiduNewsPage error after retry 10")
		}
		if len(news) == 0 {
			break
		}
		allValid := true
		for _, news1 := range news {
			if timeValid(news1.time) {
				newsReturn = append(newsReturn, news1)
			} else {
				allValid = false
			}
		}
		if !allValid {
			break
		}
		page++
		if page > 20 {
			break
		}
	}
	for i := 0; i < len(newsReturn); i++ {
		for j := i + 1; j < len(newsReturn); j++ {
			valid1 := strsim.Compare(newsReturn[i].title, newsReturn[j].title)
			valid2 := strsim.Compare(newsReturn[i].content, newsReturn[j].content)
			maxValid := valid1
			if valid2 > maxValid {
				maxValid = valid2
			}
			if maxValid > 0.1 && maxValid > newsReturn[j].valid {
				newsReturn[j].valid = maxValid
			}
		}
	}
	//// 写入到1.log
	//f, _ := os.Create("1.log")
	//for i := 0; i < len(newsReturn); i++ {
	//	f.WriteString(fmt.Sprintf("%s %f\n", newsReturn[i].title, newsReturn[i].valid))
	//}
	/*
		以下参数由以上数据参考决定
		如何才能进入特斯拉工作 0.000000
		特斯拉市值一夜蒸发3323亿元!他超越马斯克成为世界首富 0.000000
		特斯拉市值一夜蒸发3300亿元,马斯克“世界首富”跌没了! 0.551724
		特斯拉概念5日主力净流出31.77亿元,隆基绿能、通富微电居前 0.193548
		国产特斯拉2月份交付量同比下滑18% 环比也下滑超过10% 0.102564
		特斯拉位于柏林郊外的工厂及周边部分地区因环保人士在工厂纵火而... 0.000000
		马斯克“世界首富”跌没了!特斯拉市值蒸发超3300亿元 0.111111
		特斯拉市值一夜蒸发3323亿元!亚马逊创始人贝索斯反超马斯克,成为... 0.555556
		人形机器人新蓝海掀起巨浪,特斯拉/微美全息全面“狂奔”迈向新征程 0.138462
		连奔驰都不敢奋身搞电车了,全球只剩中国和特斯拉了?|热财经 0.103448
		港股异动丨隔夜特斯拉大跌,拖累汽车股普跌:恒大汽车跌超6%,零跑... 0.166667
		零百加速2秒,史上最强保时捷Taycan,轻松碾压特斯拉Model S Plaid? 0.000000
		特斯拉销量下滑7% 14家车企降价战短期需求受挑战 0.129032
		1秒破百?特斯拉全新一代Roadster或将年底发布 0.000000
		特斯拉股价今年跌幅达24% 贝佐斯反超马斯克再登富豪榜榜首 0.277778
		马斯克“世界首富”跌没了!贝佐斯反超马斯克成全球首富,特斯拉... 0.424242
		特斯拉股价周一大跌 7.16%:贝索斯替代马斯克再成世界首富 0.400000
		普通人可能没法驾驶?零百加速不到 1 秒!特斯拉新车性能爆表 0.133333
		特斯拉市值一夜蒸发超3000亿元,首富马斯克被贝佐斯超越 0.517241
		纳指跌0.4% 英伟达涨3.6%特斯拉跌7.2%苹果跌2.5% 0.137255
		特斯拉股价大跌 马斯克丢掉世界首富位子 0.433333
		特斯拉暴跌7%,马斯克痛失全球首富宝座 0.421053
		马斯克“世界首富”跌没了!特斯拉突然重挫,市值蒸发3300亿元 0.806452
		特斯拉股价大跌超7% 马斯克再度丢掉世界首富头衔 0.708333
		隔夜外围市场综述:特斯拉跌逾7% 0.187500
		美股三大指数集体收跌,特斯拉跌超7% 0.333333
		特斯拉收盘大跌超7%,年内已累跌超24%,2月中国销量创1年多新低 0.242424
		美股三大股指集体下跌 英伟达市值超越沙特阿美 特斯拉重挫逾7% 0.387097
		美股收评:三大指数集体下跌,特斯拉跌超7%,中概股指跌近4% 0.787879
		美股收盘:三大指数小幅收跌 特斯拉重挫逾7% 英伟达再创新高 0.903226
		特斯拉盘中重挫7.5%,市值蒸发超400亿美元,年内已累跌约25% 0.656863
		特斯拉跌幅扩大至6.7% 0.210526
		特斯拉汽车影院终于要来了!内设餐厅、超充和电影院 0.161290
		特斯拉盘前大跌! 0.333333
		特斯拉遭遇最差开年,市值蒸发逾 940 亿美元 0.333333
	*/
	newsReturn = CutInvalidNews(newsReturn, 0.2)
	return
}
