package spider

import (
	"encoding/json"
	"errors"
	"github.com/antlabs/strsim"
	"net/http"
	url2 "net/url"
	"strconv"
	"time"
)

const GNewsApiUrl = "https://gnews.io/api/v4/search?"
const GNewsApiTopUrl = "https://gnews.io/api/v4/top-headlines?"

type GNewsLang string
type GNewsCountry string
type GNewsSortBy string

const (
	LanArabic     GNewsLang = "ar"
	LanChinese    GNewsLang = "zh"
	LanDutch      GNewsLang = "nl"
	LanEnglish    GNewsLang = "en"
	LanFrench     GNewsLang = "fr"
	LanGerman     GNewsLang = "de"
	LanGreek      GNewsLang = "el"
	LanHebrew     GNewsLang = "he"
	LanHindi      GNewsLang = "hi"
	LanItalian    GNewsLang = "it"
	LanJapanese   GNewsLang = "ja"
	LanMalayalam  GNewsLang = "ml"
	LanMarathi    GNewsLang = "mr"
	LanNorwegian  GNewsLang = "no"
	LanPortuguese GNewsLang = "pt"
	LanRomanian   GNewsLang = "ro"
	LanRussian    GNewsLang = "ru"
	LanSpanish    GNewsLang = "es"
	LanSwedish    GNewsLang = "sv"
	LanTamil      GNewsLang = "ta"
	LanTelugu     GNewsLang = "te"
	LanUkrainian  GNewsLang = "uk"

	CountryAustralia     GNewsCountry = "au"
	CountryBrazil        GNewsCountry = "br"
	CountryCanada        GNewsCountry = "ca"
	CountryChina         GNewsCountry = "cn"
	CountryEgypt         GNewsCountry = "eg"
	CountryFrance        GNewsCountry = "fr"
	CountryGermany       GNewsCountry = "de"
	CountryGreece        GNewsCountry = "gr"
	CountryHong          GNewsCountry = "Kong	hk"
	CountryIndia         GNewsCountry = "in"
	CountryIreland       GNewsCountry = "ie"
	CountryIsrael        GNewsCountry = "il"
	CountryItaly         GNewsCountry = "it"
	CountryJapan         GNewsCountry = "jp"
	CountryNetherlands   GNewsCountry = "nl"
	CountryNorway        GNewsCountry = "no"
	CountryPakistan      GNewsCountry = "pk"
	CountryPeru          GNewsCountry = "pe"
	CountryPhilippines   GNewsCountry = "ph"
	CountryPortugal      GNewsCountry = "pt"
	CountryRomania       GNewsCountry = "ro"
	CountryRussian       GNewsCountry = "ru"
	CountrySingapore     GNewsCountry = "sg"
	CountrySpain         GNewsCountry = "es"
	CountrySweden        GNewsCountry = "se"
	CountrySwitzerland   GNewsCountry = "ch"
	CountryTaiwan        GNewsCountry = "tw"
	CountryUkraine       GNewsCountry = "ua"
	CountryUnitedKingdom GNewsCountry = "gb"
	CountryUnitedStates  GNewsCountry = "us"

	/*
		// publishedAt = sort by publication date, the articles with the most recent publication date are returned first
			// relevance = sort by best match to keywords, the articles with the best match are returned first
	*/
	SortByPublishedAt GNewsSortBy = "publishedAt"
	SortByRelevance   GNewsSortBy = "relevance"
)

type GNewsSearch struct {
	q       string       // This parameter allows you to specify your search keywords to find the news articles you are looking for. The keywords will be used to return the most relevant articles. It is possible to use logical operators with keywords, see the section on query syntax.
	lang    GNewsLang    // This parameter allows you to specify the language of the news articles returned by the API. You have to set as value the 2 letters code of the language you want to filter.
	country GNewsCountry // This parameter allows you to specify the country where the news articles returned by the API were published, the contents of the articles are not necessarily related to the specified country. You have to set as value the 2 letters code of the country you want to filter.
	max     int          // This parameter allows you to specify the number of news articles returned by the API. The minimum value of this parameter is 1 and the maximum value is 100. The value you can set depends on your subscription.
	min     int          // This parameter allows you to specify the number of news articles returned by the API. The minimum value of this parameter is 1 and the maximum value is 100. The value you can set depends on your subscription.
	/*
		This parameter allows you to filter the articles that have a publication date greater than or equal to the specified value. The date must respect the following format:
		YYYY-MM-DDThh:mm:ssTZD
		TZD = time zone designator, its value must always be Z (universal time)
		e.g. 2022-08-21T16:27:09Z
	*/
	from   UniTimeStr  // This parameter allows you to filter the articles that have a publication date greater than or equal to the specified value. The date must respect the following format:
	to     UniTimeStr  // This parameter allows you to filter the articles that have a publication date smaller than or equal to the specified value. The date must respect the following format:
	sortby GNewsSortBy // This parameter allows you to choose with which type of sorting the articles should be returned. Two values are possible:
}

type GNewsTop struct {
	category string
	Lang     GNewsLang    // This parameter allows you to specify the language of the news articles returned by the API. You have to set as value the 2 letters code of the language you want to filter.
	country  GNewsCountry // This parameter allows you To specify the country where the news articles returned by the API were published, the contents of the articles are not necessarily related To the specified country. You have To set as value the 2 letters code of the country you want To filter.
	From     UniTimeStr   // This parameter allows you to filter the articles that have a publication date greater than or equal to the specified value. The date must respect the following format:
	To       UniTimeStr   // This parameter allows you To filter the articles that have a publication date smaller than or equal To the specified value. The date must respect the following format:
}

type UniTimeStr string

func GetUniTimeStr(t time.Time) UniTimeStr {
	t = t.In(time.UTC)
	return UniTimeStr(t.Format("2006-01-02T15:04:05Z"))
}

type SearchResult struct {
	TotalArticles int `json:"totalArticles"`
	Articles      []struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Content     string    `json:"content"`
		Url         string    `json:"url"`
		Image       string    `json:"image"`
		PublishedAt time.Time `json:"publishedAt"`
		Source      struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"source"`
	} `json:"articles"`
}
type GNewsArticles []struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	Url         string    `json:"url"`
	Image       string    `json:"image"`
	PublishedAt time.Time `json:"publishedAt"`
	Source      struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"source"`
}

type TopResult struct {
	TotalArticles int           `json:"totalArticles"`
	Articles      GNewsArticles `json:"articles"`
}

func QueryGNews(search GNewsSearch, apikey string) (SearchResult, error) {
	url := GNewsApiUrl + "apikey=" + apikey
	if search.q != "" {
		q := url2.QueryEscape("\"" + search.q + "\"")
		url += "&q=" + q
	}
	if search.lang != "" {
		url += "&lang=" + string(search.lang)
	}
	if search.country != "" {
		url += "&country=" + string(search.country)
	}
	if search.max != 0 {
		url += "&max=" + strconv.Itoa(search.max)
	}
	if search.min != 0 {
		url += "&min=" + strconv.Itoa(search.min)
	}
	if search.from != "" {
		url += "&from=" + string(search.from)
	}
	if search.to != "" {
		url += "&to=" + string(search.to)
	}
	if search.sortby != "" {
		url += "&sortby=" + string(search.sortby)
	}
	// get
	resp, err := http.Get(url)
	if err != nil {
		return SearchResult{}, errors.Join(errors.New("http get error"), err)
	}
	// 解析json
	var result SearchResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return SearchResult{}, errors.Join(errors.New("json decode error"), err)
	}
	return result, nil
}

func QueryGNewsTop(top GNewsTop, apikey string) (TopResult, error) {
	url := GNewsApiTopUrl + "apikey=" + apikey
	if top.category != "" {
		url += "&category=" + top.category
	} else {
		url += "&category=general"
	}
	if top.Lang != "" {
		url += "&lang=" + string(top.Lang)
	}
	if top.country != "" {
		url += "&country=" + string(top.country)
	}
	if top.From != "" {
		url += "&from=" + string(top.From)
	}
	if top.To != "" {
		url += "&to=" + string(top.To)
	}
	// get
	resp, err := http.Get(url)
	if err != nil {
		return TopResult{}, errors.Join(errors.New("http get error"), err)
	}
	// 解析json
	var result TopResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return TopResult{}, errors.Join(errors.New("json decode error"), err)
	}
	return result, nil
}

func GetGNewsSumTop(token string, last time.Time) (TopResult, error) {
	// 将last到现在的时间拆分为一个小时段
	var times []time.Time
	for last.Before(time.Now()) {
		times = append(times, last)
		last = last.Add(time.Hour * 24)
	}
	times = append(times, time.Now())
	// 每个小时段查询一次
	var result GNewsArticles
	for i := 0; i < len(times)-1; i++ {
		req := GNewsTop{
			Lang: LanEnglish,
			From: GetUniTimeStr(times[i]),
			To:   GetUniTimeStr(times[i+1]),
		}
		r, err := QueryGNewsTop(req, token)
		if err != nil {
			return TopResult{}, err
		}
		result = append(result, r.Articles...)
	}
	var topResult TopResult
	// 进行去重
	for i := 0; i < len(result); i++ {
		isValid := true
		for j := i - 1; j >= 0; j-- {
			valid1 := strsim.Compare(result[i].Title, result[j].Title)
			valid2 := strsim.Compare(result[i].Content, result[j].Content)
			maxValid := valid1
			if valid2 > maxValid {
				maxValid = valid2
			}
			if maxValid > 0.7 {
				isValid = false
				continue
			}
		}
		if isValid {
			topResult.Articles = append(topResult.Articles, result[i])
		}
	}
	return topResult, nil
}
