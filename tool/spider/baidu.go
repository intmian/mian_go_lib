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
	same    float64
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
	s := ""
	var noNews []string
	for i, keyword := range keywords {
		newsNum := len(news[i])
		if newsNum != 0 {
			continue
		}
		noNews = append(noNews, keyword)
	}
	if len(noNews) != 0 {
		s += "#### 无新动态\r\n"
		str := "- "
		for i, noNew := range noNews {
			if i != 0 {
				str += "、"
			}
			str += noNew
		}
		s += str + "\r\n"
	}

	for i, keyword := range keywords {
		newsNum := len(news[i])
		if newsNum == 0 {
			continue
		}
		s += "#### " + keyword + " " + strconv.Itoa(newsNum) + "条新闻\r\n"
		for _, baiduNew := range news[i] {
			baiduNew.time = strings.Replace(baiduNew.time, "undefined", "近期", -1)
			// 来源 时间：标题（链接）
			s += fmt.Sprintf("- %s %s：[%s](%s) \r\n", baiduNew.source, baiduNew.time, baiduNew.title, baiduNew.href)
		}
	}
	return s
}

func getBaiduNewsPage(keyword string, page int) (result []BaiduNew, err error) {
	result = make([]BaiduNew, 0)
	/*
		服务器和本机读取的结果可能不一样，本机的类似于直接搜索百度新闻，服务器上的就是精简后的，而且去除重复新闻。
		服务器请求可能会出现网络问题
	*/
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
	webText, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.WithMessage(err, "ioutil.ReadAll error")
	}
	//f, _ := os.Create(fmt.Sprintf("baidu_%s_%d_%s.html", keyword, page, time.Now().Format("2006-01-02_15:04:05")))
	//f.WriteString(string(webText))
	if strings.Contains(string(webText), "网络不给力，请稍后重试") {
		return nil, errors.New("网络不给力，请稍后重试")
	}
	reStr := `\{"titleAriaLabel":"标题[： ](.*)","absAriaLabel":"摘要[： ](.*)","sourceAriaLabel":"新闻来源[： ](.*)","timeAriaLabel":"发布于[： ](.{0,20})"\}.*href="(.*)" target`
	reg1 := regexp.MustCompile(reStr)
	if reg1 == nil {
		return
	}
	//根据规则提取关键信息
	results := reg1.FindAllStringSubmatch(string(webText), -1)
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
	if len(result) == 0 && page <= 1 {
		//加入全量打印用于调试。报错+打印内容方便debug。
		f, _ := os.Create(fmt.Sprintf("baidu_%s_%d_%s.html", keyword, page, time.Now().Format("2006-01-02_15:04:05")))
		f.WriteString(string(webText))
		return nil, errors.New("no news")
	}
	return
}

func GetBaiduNewsWithoutOld(keyword string, lastLinks []string, maxSame float64) (results []BaiduNew, newLinks []string, err error, retry int, folded int) {
	/*
		由于最近的新闻接口非常不稳定，所以需要传入上一次的最新链接，以便获取新的新闻。很多新闻最近没有时间了，而且有时会丢失.
		最多保存40条（对应4页）
		1.向后翻页一直翻到找到上一次的链接，再向后翻页直到，没有任何上次缓存的链接
		2.将所有最后一条缓存链接前的新闻，且不为缓存链接的新闻加入返回列表
		3.返回新的缓存链接（前40条）
	*/
	// 前置处理
	keyword = strings.Replace(keyword, " ", "+", -1)
	needInit := len(lastLinks) == 0

	// 初始化时，最多就读个20条
	if needInit {
		needNewsLen := 20
		for page := 1; page <= 5; page++ {
			news, PageRetry, err1 := getBaiduNewsPageRetry(keyword, page, 20)
			if err1 != nil {
				err = errors.WithMessage(err1, "getBaiduNewsPage error after retry")
				return
			}
			retry += PageRetry

			// 提取不到数据就要中断流程
			if len(news) == 0 {
				break
			}
			results = append(results, news...)
			if len(results) > needNewsLen {
				break
			}
		}
		if len(results) > needNewsLen {
			results = results[0:needNewsLen]
		}
		for _, news := range results {
			newLinks = append(newLinks, news.href)
		}
		results, folded = CutInvalidNews(results, maxSame)
		return
	}

	// 建立一个map，用于快速查找
	lastLinkMap := make(map[string]bool)
	for _, lastLink := range lastLinks {
		lastLinkMap[lastLink] = true
	}

	news1 := make([]BaiduNew, 0)
	findLastLink := false
	noLastLink := false
	for page := 1; page <= 20; page++ {
		pageNews, PageRetry, err1 := getBaiduNewsPageRetry(keyword, page, 20)
		if err1 != nil {
			err = errors.WithMessage(err1, "getBaiduNewsPage error after retry")
			return
		}
		retry += PageRetry
		if !findLastLink {
			for _, news := range pageNews {
				if lastLinkMap[news.href] {
					findLastLink = true
				}
			}
		}
		if findLastLink {
			allNotLastLink := true
			for _, news := range pageNews {
				if lastLinkMap[news.href] {
					allNotLastLink = false
					break
				}
			}
			if allNotLastLink {
				noLastLink = true
			}
		}
		// 向后翻页一直翻到没有之前链接的一页，但是去除最后一页的数据
		if !noLastLink {
			news1 = append(news1, pageNews...)
		} else {
			break
		}
	}

	// 找到最后一条缓存链接的位置，去除之后的
	// news1全量 news2去除了最后一条缓存链接之后的
	lastLinkIndex := -1
	for i, news := range news1 {
		if lastLinkMap[news.href] {
			lastLinkIndex = i
			break
		}
	}
	var news2 []BaiduNew
	if lastLinkIndex != -1 {
		news2 = news1[0:lastLinkIndex]
	}

	// 筛选出没有出现过的新闻，且不为缓存链接的新闻
	for _, news := range news2 {
		if !lastLinkMap[news.href] {
			results = append(results, news)
		}
	}

	// 将news1的链接加入newLinks
	for _, news := range news1 {
		newLinks = append(newLinks, news.href)
	}
	// 将lastLinks的链接加入newLinks。
	// 从前向后扫描一遍，如果出现lastLinks的链接就标记位置，之后将没有出现的lastLinks插入到上一个出现的lastlink的上方。（因为上一次有的这一次可能没有，如果删除了，下一次又冒出来了）
	newLinks = mergeLinks(lastLinks, newLinks)
	if len(newLinks) > 40 {
		newLinks = newLinks[0:40]
	}
	// 裁剪重复新闻
	if maxSame > 0 {
		results, folded = CutInvalidNews(results, maxSame)
	}
	return
}

func mergeLinks(old, new []string) []string {
	result := []string{}
	oldMap := make(map[string]int)
	newMap := make(map[string]bool)

	// 构建 oldMap 用于快速查找
	for i, val := range old {
		oldMap[val] = i
	}

	// 构建 newMap 用于快速查找
	for _, val := range new {
		newMap[val] = true
	}

	i := 0 // old 的指针
	for _, val := range new {
		result = append(result, val)
		// 如果当前元素在 old 中
		if idx, exists := oldMap[val]; exists {
			// 插入 old 中当前元素之后但不在 new 中的所有元素
			for i < len(old) && i < idx {
				i++
			}
			i++ // 跳过当前元素
			for i < len(old) && !newMap[old[i]] {
				result = append(result, old[i])
				i++
			}
		}
	}

	return result
}

func getBaiduNewsPageRetry(keyword string, page int, retryMax int) (results []BaiduNew, retry int, err error) {
	results = make([]BaiduNew, 0)
	for retry = 0; retry < retryMax; retry++ {
		results, err = getBaiduNewsPage(keyword, page)
		if err == nil {
			return
		}
		time.Sleep(time.Minute)
	}
	return
}

// GetBaiduNewsNew 获取百度新闻。由于最近的新闻接口非常不稳定，所以需要传入上一次的最新链接，以便获取新的新闻。很多新闻最近没有时间了，而且有时会丢失
func GetBaiduNewsNew(keyword string, lastLink string, maxSame float64) (results []BaiduNew, newestLink string, err error, retry int) {
	// 前置处理
	keyword = strings.Replace(keyword, " ", "+", -1)

	// 查找新闻
	page := 1
	for {
		retryTimes := 20
		var news []BaiduNew
		// 重试，至多20次
		for retryTimes > 0 {
			news, err = getBaiduNewsPage(keyword, page)
			if err == nil {
				break
			}
			retryTimes--
			retry++
			// 休眠一分钟
			time.Sleep(time.Minute)
		}
		if err != nil {
			err = errors.WithMessage(err, "getBaiduNewsPage error after retry")
			return
		}

		// 没有数据也返回
		if len(news) == 0 {
			break
		}

		results = append(results, news...)

		findLast := false
		for _, news1 := range news {
			if news1.href == lastLink {
				findLast = true
				break
			}
		}
		// 如果找到了上一次的链接，那么就返回，不然肯定在后面的页数
		if findLast {
			break
		}
		page++

		// 为空的情况下就读两页
		if lastLink == "" && page > 2 {
			break
		}

		if page > 20 {
			break
		}
	}

	// 删除lastLink及之前的数据
	lastLinkIndex := -1
	for i, news1 := range results {
		if news1.href == lastLink {
			lastLinkIndex = i
			break
		}
	}
	if lastLinkIndex != -1 {
		results = results[0:lastLinkIndex]
	}

	// 获得最新链接
	if len(results) > 0 {
		newestLink = results[0].href
	}

	results, _ = CutInvalidNews(results, maxSame)
	return
}

func CutInvalidNews(results []BaiduNew, maxSame float64) ([]BaiduNew, int) {
	// 计算有效性，重复度
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			valid1 := strsim.Compare(results[i].title, results[j].title)
			valid2 := strsim.Compare(results[i].content, results[j].content)
			maxValid := valid1
			if valid2 > maxValid {
				maxValid = valid2
			}
			if maxValid > 0.1 && maxValid > results[j].same {
				results[j].same = maxValid
			}
		}
	}
	var folded int
	results, folded = CutMoreSameNews(results, maxSame)
	return results, folded
}

func CutMoreSameNews(news []BaiduNew, maxSame float64) ([]BaiduNew, int) {
	folded := 0
	newsReturn := make([]BaiduNew, 0)
	for _, baiduNew := range news {
		if baiduNew.same < maxSame {
			newsReturn = append(newsReturn, baiduNew)
		} else {
			folded++
		}
	}
	return newsReturn, folded
}

func GetTodayBaiduNews(keyword string) (newsReturn []BaiduNew, err error, retry int) {
	newsReturn = make([]BaiduNew, 0)
	keyword = strings.Replace(keyword, " ", "+", -1)

	page := 1
	for {
		retryTimes := 20
		var news []BaiduNew
		for retryTimes > 0 {
			news, err = getBaiduNewsPage(keyword, page)
			if err == nil {
				break
			}
			retryTimes--
			retry++
			// 休眠一分钟
			time.Sleep(time.Minute)
		}
		if err != nil {
			return nil, errors.WithMessage(err, "getBaiduNewsPage error after retry"), retry
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

	// 如果时间全部都是undefined，那么就直接返回空，因为undefined的新闻可能是很古早的
	allUndefined := true
	for _, news1 := range newsReturn {
		if news1.time != "undefined" {
			allUndefined = false
			break
		}
	}
	if allUndefined {
		return []BaiduNew{}, nil, retry
	}

	for i := 0; i < len(newsReturn); i++ {
		for j := i + 1; j < len(newsReturn); j++ {
			valid1 := strsim.Compare(newsReturn[i].title, newsReturn[j].title)
			valid2 := strsim.Compare(newsReturn[i].content, newsReturn[j].content)
			maxValid := valid1
			if valid2 > maxValid {
				maxValid = valid2
			}
			if maxValid > 0.1 && maxValid > newsReturn[j].same {
				newsReturn[j].same = maxValid
			}
		}
	}
	//// 写入到1.log
	//f, _ := os.Create("1.log")
	//for i := 0; i < len(newsReturn); i++ {
	//	f.WriteString(fmt.Sprintf("%s %f\n", newsReturn[i].title, newsReturn[i].same))
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
	newsReturn, _ = CutMoreSameNews(newsReturn, 0.2)
	return
}
