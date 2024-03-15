package spider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

func GetWeatherDataOri(province, city string) (string, error) {
	pc := url.QueryEscape(province + city + "天气")
	url1 := "https://weathernew.pae.baidu.com/weathernew/pc?query=%s&srcid=4982&forecast=long_day_forecast"
	url1 = fmt.Sprintf(url1, pc)
	req, err := http.NewRequest("GET", url1, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("DNT", "1")
	req.Header.Set("Referer", "https://www.baidu.com/link?url=LOoW7nfxfB2350fjBnQho9KK8Q8Ohrk3zjDnkt5-ji2dYVikhoZM0eMLAh4n9zX9JVGtbVCjEWTkgvmPficS0lutwN8YMcIveqCrmGMqUwHSQ7gheKSPqJa3LUg9_6OV3Qe9jEyVbPedGbd9sfZhn3Pa41CWbxXZCfPkOePFfq1xroUjxSIr0DBtAEjutPQSB0QXgbifwjJl7mkWQ5ZS_a&wd=&eqid=f3b700cf000029540000000465f3e3c1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Set("sec-ch-ua", `"Chromium";v="122", "Not(A:Brand";v="24", "Google Chrome";v="122"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(body), nil
}

type Index string

const (
	ChuanYiIndex   Index = "穿衣"
	WuranIndex     Index = "污染"
	ChuyouIndex    Index = "出游"
	XicheIndex     Index = "洗车"
	ShaiBeiZiIndex Index = "晒被子"
)

type IndexInfo struct {
	Status string
	Why    string
}

type Weather struct {
	Date       time.Time
	IndexMap   map[Index]IndexInfo
	MinWeather int
	MaxWeather int
	Condition  string
}

func GetTodayWeather(s string) (Weather, error) {
	// 从原文中使用正则提取出对应的字符串
	ZhiShuPerm := `data\["zhishu"\] *= *(\{.*\})`
	Weather15DayData := `data\["weather15DayData"\] *= *(\[.*\])`
	// 从原文中使用正则提取出对应的字符串，然后使用json解析
	ZhiShuRe := regexp.MustCompile(ZhiShuPerm)
	Weather15DayDataRe := regexp.MustCompile(Weather15DayData)
	ZhiShuMatch := ZhiShuRe.FindStringSubmatch(s)
	Weather15DayDataMatch := Weather15DayDataRe.FindStringSubmatch(s)
	if len(ZhiShuMatch) < 2 || len(Weather15DayDataMatch) < 2 {
		return Weather{}, fmt.Errorf("error parsing origin string")
	}
	// 使用json解析
	var l LifeIndex
	err := json.Unmarshal([]byte(ZhiShuMatch[1]), &l)
	if err != nil {
		return Weather{}, fmt.Errorf("error parsing zhishu json: %v", err)
	}
	var d []Day15Weather
	err = json.Unmarshal([]byte(Weather15DayDataMatch[1]), &d)
	if err != nil {
		return Weather{}, fmt.Errorf("error parsing weather15daydata json: %v", err)
	}
	// 解析出来的数据转换为Weather结构体
	weather := Weather{
		Date: time.Now(),
	}
	todayStr := time.Now().Format("2006-01-02")
	for _, v := range d {
		if v.Date == todayStr {
			weather.Condition = v.WeatherText
			break
		}
	}
	weather.IndexMap = make(map[Index]IndexInfo)
	for _, v := range l.Item {
		weather.IndexMap[Index(v.ItemName)] = IndexInfo{
			Status: v.ItemTitle,
			Why:    v.ItemDesc,
		}
		if v.ItemName == "穿衣" {
			// 10℃~19℃ 中提取出最低温度和最高温度
			_, err := fmt.Sscanf(v.ItemDesc, "%d℃~%d℃", &weather.MinWeather, &weather.MaxWeather)
			if err != nil {
				return Weather{}, fmt.Errorf("error parsing temperature: %v", err)
			}
		}
	}
	// 有一个为空就返回错误
	if weather.MinWeather == 0 || weather.MaxWeather == 0 || weather.Condition == "" || len(weather.IndexMap) == 0 {
		return Weather{}, fmt.Errorf("error parsing weather data")
	}
	return weather, nil
}

type Day15Weather struct {
	FormatDate  string `json:"formatDate"`
	Date        string `json:"date"`
	FormatWeek  string `json:"formatWeek"`
	WeatherIcon string `json:"weatherIcon"`
	WeatherWind struct {
		WindDirectionDay   string `json:"windDirectionDay"`
		WindDirectionNight string `json:"windDirectionNight"`
		WindPowerDay       string `json:"windPowerDay"`
		WindPowerNight     string `json:"windPowerNight"`
	} `json:"weatherWind"`
	WeatherPm25 string `json:"weatherPm25"`
	WeatherText string `json:"weatherText"`
}

type LifeIndex struct {
	Url      string `json:"url"`
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	OtherUrl string `json:"other_url"`
	Item     []struct {
		ItemName      string `json:"item_name"`
		ItemTitle     string `json:"item_title"`
		ItemIcon      string `json:"item_icon"`
		ItemIconWhite string `json:"item_icon_white"`
		ItemDesc      string `json:"item_desc"`
		ItemUrl       string `json:"item_url"`
		ItemOtherUrl  string `json:"item_other_url"`
	} `json:"item"`
	StrategyLog struct {
		RecommendZhishuSort []string      `json:"recommend_zhishu_sort"`
		UserAttr            []interface{} `json:"user_attr"`
		ObserveWeather      struct {
			BodytempInfo      string `json:"bodytemp_info"`
			WindDirection     string `json:"wind_direction"`
			Site              string `json:"site"`
			Weather           string `json:"weather"`
			DewTemperature    string `json:"dew_temperature"`
			PrecipitationType string `json:"precipitation_type"`
			WindDirectionNum  string `json:"wind_direction_num"`
			Temperature       string `json:"temperature"`
			WindPower         string `json:"wind_power"`
			F1HInfo           []struct {
				PrecipitationProbability string `json:"precipitation_probability"`
				Temperature              string `json:"temperature"`
				Hour                     string `json:"hour"`
				WindDirection            string `json:"wind_direction"`
				Uv                       string `json:"uv"`
				UvNum                    string `json:"uv_num"`
				WindPower                string `json:"wind_power"`
				Weather                  string `json:"weather"`
				WindPowerNum             string `json:"wind_power_num"`
				Precipitation            string `json:"precipitation"`
			} `json:"f1hInfo"`
			UpdateTime          string `json:"update_time"`
			PublishTime         string `json:"publish_time"`
			Visibility          string `json:"visibility"`
			Pressure            string `json:"pressure"`
			PrecMonitorTime     string `json:"prec_monitor_time"`
			Precipitation       string `json:"precipitation"`
			RealFeelTemperature string `json:"real_feel_temperature"`
			UvInfo              string `json:"uv_info"`
			Uv                  string `json:"uv"`
			Humidity            string `json:"humidity"`
			UvNum               string `json:"uv_num"`
			WindPowerNum        string `json:"wind_power_num"`
			F1HInfoNumBaidu     int    `json:"f1hInfo#num#baidu"`
			PsPm25              string `json:"ps_pm25"`
		} `json:"observe_weather"`
	} `json:"strategy_log"`
}
