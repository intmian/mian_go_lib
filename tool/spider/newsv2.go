package spider

import (
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"net/http"
	"sort"
	"time"
)

type BBCRssItem struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	PubDate     time.Time `json:"pubDate"`
}

func GetBBCRss(client *http.Client) ([]BBCRssItem, error) {
	const bbcRssUrl = "https://feeds.bbci.co.uk/zhongwen/trad/rss.xml"
	fp := gofeed.NewParser()
	// 配置代理 7890
	if client != nil {
		fp.Client = client
	}
	feed, err := fp.ParseURL(bbcRssUrl)
	if err != nil {
		return nil, errors.WithMessage(err, "GetBBCRss")
	}
	var items []BBCRssItem
	for _, item := range feed.Items {
		//<lastBuildDate>Wed, 25 Dec 2024 08:20:58 GMT</lastBuildDate>
		newsTime, err := time.Parse(time.RFC1123, item.Published)
		if err != nil {
			return nil, errors.WithMessage(err, "GetBBCRss")
		}
		items = append(items, BBCRssItem{
			Title:       item.Title,
			Description: item.Description,
			Link:        item.Link,
			PubDate:     newsTime,
		})
	}

	// 根据时间排序
	sort.Slice(items, func(i, j int) bool {
		return items[i].PubDate.After(items[j].PubDate)
	})

	return items, nil
}

func GetBBCRssWithDay(day time.Time, client *http.Client) ([]BBCRssItem, error) {
	items, err := GetBBCRss(client)
	if err != nil {
		return nil, errors.WithMessage(err, "GetBBCRssWithDay")
	}
	var res []BBCRssItem
	for _, item := range items {
		if item.PubDate.Day() == day.Day() {
			res = append(res, item)
		}
	}
	return res, nil
}

type GoogleRssItem struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	PubDate     time.Time `json:"pubDate"`
}

func GetGoogleRss(keyWord string, client *http.Client) ([]GoogleRssItem, error) {
	// 先支持中文搜索再说
	const googleRssUrl = "https://news.google.com/rss/search?q=%s&hl=zh-CN&gl=CN&ceid=CN%3Azh-Hans"
	fp := gofeed.NewParser()
	if client != nil {
		fp.Client = client
	}
	feed, err := fp.ParseURL(googleRssUrl)
	if err != nil {
		return nil, errors.WithMessage(err, "GetGoogleRss")
	}
	var items []GoogleRssItem
	for _, item := range feed.Items {
		newsTime, err := time.Parse(time.RFC1123, item.Published)
		if err != nil {
			return nil, errors.WithMessage(err, "GetGoogleRss")
		}
		items = append(items, GoogleRssItem{
			Title:       item.Title,
			Description: item.Description,
			Link:        item.Link,
			PubDate:     newsTime,
		})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].PubDate.After(items[j].PubDate)
	})
	return items, nil
}

func GetGoogleRssWithDay(keyWord string, day time.Time, client *http.Client) ([]GoogleRssItem, error) {
	items, err := GetGoogleRss(keyWord, client)
	if err != nil {
		return nil, errors.WithMessage(err, "GetGoogleRssWithDay")
	}
	var res []GoogleRssItem
	for _, item := range items {
		if item.PubDate.Day() == day.Day() {
			res = append(res, item)
		}
	}
	return res, nil
}

type NYTimesRssItem struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	PubDate     time.Time `json:"pubDate"`
}

func GetNYTimesRss(client *http.Client) ([]NYTimesRssItem, error) {
	const nyTimesRssWorldUrl = "https://rss.nytimes.com/services/xml/rss/nyt/World.xml"
	const nyTimesRssAsiaUrl = "https://rss.nytimes.com/services/xml/rss/nyt/AsiaPacific.xml"

	titleMap := make(map[string]bool)

	fp := gofeed.NewParser()
	if client != nil {
		fp.Client = client
	}
	feed, err := fp.ParseURL(nyTimesRssWorldUrl)
	if err != nil {
		return nil, errors.WithMessage(err, "GetNYTimesRss")
	}
	var items []NYTimesRssItem
	for _, item := range feed.Items {
		newsTime, err := time.Parse(time.RFC1123Z, item.Published)
		if err != nil {
			return nil, errors.WithMessage(err, "GetNYTimesRss")
		}
		items = append(items, NYTimesRssItem{
			Title:       item.Title,
			Description: item.Description,
			Link:        item.Link,
			PubDate:     newsTime,
		})
		titleMap[item.Title] = true
	}
	feed, err = fp.ParseURL(nyTimesRssAsiaUrl)
	if err != nil {
		return nil, errors.WithMessage(err, "GetNYTimesRss")
	}
	for _, item := range feed.Items {
		newsTime, err := time.Parse(time.RFC1123Z, item.Published)
		if err != nil {
			return nil, errors.WithMessage(err, "GetNYTimesRss")
		}
		if _, ok := titleMap[item.Title]; ok {
			continue
		}
		items = append(items, NYTimesRssItem{
			Title:       item.Title,
			Description: item.Description,
			Link:        item.Link,
			PubDate:     newsTime,
		})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].PubDate.After(items[j].PubDate)
	})

	return items, nil
}

func GetNYTimesRssWithDay(day time.Time, client *http.Client) ([]NYTimesRssItem, error) {
	items, err := GetNYTimesRss(client)
	if err != nil {
		return nil, errors.WithMessage(err, "GetNYTimesRssWithDay")
	}
	var res []NYTimesRssItem
	for _, item := range items {
		if item.PubDate.Day() == day.Day() {
			res = append(res, item)
		}
	}
	return res, nil
}
