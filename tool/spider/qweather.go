package spider

import (
	"encoding/json"
	"fmt"
	"github.com/intmian/mian_go_lib/tool/misc"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strings"
)

const qWeatherApi = "https://devapi.qweather.com"

type IndexReturn struct {
	Code       string `json:"code"`
	UpdateTime string `json:"updateTime"`
	FxLink     string `json:"fxLink"`
	Daily      []struct {
		Date     string `json:"date"`
		Type     string `json:"type"`
		Name     string `json:"name"`
		Level    string `json:"level"`
		Category string `json:"category"`
		Text     string `json:"text"`
	} `json:"daily"`
	Refer struct {
		Sources []string `json:"sources"`
		License []string `json:"license"`
	} `json:"refer"`
}

type WeatherReturn struct {
	Code       string `json:"code"`
	UpdateTime string `json:"updateTime"`
	FxLink     string `json:"fxLink"`
	Daily      []struct {
		FxDate         string `json:"fxDate"`
		Sunrise        string `json:"sunrise"`
		Sunset         string `json:"sunset"`
		Moonrise       string `json:"moonrise"`
		Moonset        string `json:"moonset"`
		MoonPhase      string `json:"moonPhase"`
		MoonPhaseIcon  string `json:"moonPhaseIcon"`
		TempMax        string `json:"tempMax"`
		TempMin        string `json:"tempMin"`
		IconDay        string `json:"iconDay"`
		TextDay        string `json:"textDay"`
		IconNight      string `json:"iconNight"`
		TextNight      string `json:"textNight"`
		Wind360Day     string `json:"wind360Day"`
		WindDirDay     string `json:"windDirDay"`
		WindScaleDay   string `json:"windScaleDay"`
		WindSpeedDay   string `json:"windSpeedDay"`
		Wind360Night   string `json:"wind360Night"`
		WindDirNight   string `json:"windDirNight"`
		WindScaleNight string `json:"windScaleNight"`
		WindSpeedNight string `json:"windSpeedNight"`
		Humidity       string `json:"humidity"`
		Precip         string `json:"precip"`
		Pressure       string `json:"pressure"`
		Vis            string `json:"vis"`
		Cloud          string `json:"cloud"`
		UvIndex        string `json:"uvIndex"`
	} `json:"daily"`
	Refer struct {
		Sources []string `json:"sources"`
		License []string `json:"license"`
	} `json:"refer"`
}

func QueryTodayIndex(key string, location string) (IndexReturn, error) {
	queryUrl := qWeatherApi + "/v7/indices/1d?type=1,3,6,8,9,11,14,16&location=" + url.QueryEscape(location) + "&key=" + key
	resp, err := http.Get(queryUrl)
	if err != nil {
		return IndexReturn{}, errors.WithMessage(err, "http get error")
	}
	defer resp.Body.Close()
	var indexReturn IndexReturn
	err = json.NewDecoder(resp.Body).Decode(&indexReturn)
	if err != nil {
		return IndexReturn{}, errors.WithMessage(err, "json decode error")
	}
	return indexReturn, nil
}

func QueryTodayWeather(key string, location string) (WeatherReturn, error) {
	queryUrl := qWeatherApi + "/v7/weather/3d?location=" + url.QueryEscape(location) + "&key=" + key
	resp, err := http.Get(queryUrl)
	if err != nil {
		return WeatherReturn{}, errors.WithMessage(err, "http get error")
	}
	defer resp.Body.Close()
	var weatherReturn WeatherReturn
	err = json.NewDecoder(resp.Body).Decode(&weatherReturn)
	if err != nil {
		return WeatherReturn{}, errors.WithMessage(err, "json decode error")
	}
	return weatherReturn, nil
}

type CityReturn struct {
	Code     string `json:"code"`
	Location []struct {
		Name      string `json:"name"`
		Id        string `json:"id"`
		Lat       string `json:"lat"`
		Lon       string `json:"lon"`
		Adm2      string `json:"adm2"`
		Adm1      string `json:"adm1"`
		Country   string `json:"country"`
		Tz        string `json:"tz"`
		UtcOffset string `json:"utcOffset"`
		IsDst     string `json:"isDst"`
		Type      string `json:"type"`
		Rank      string `json:"rank"`
		FxLink    string `json:"fxLink"`
	} `json:"location"`
	Refer struct {
		Sources []string `json:"sources"`
		License []string `json:"license"`
	} `json:"refer"`
}

func QueryCity(locationName string, key string) string {
	u := "https://geoapi.qweather.com/v2/city/lookup?" + "location=" + url.QueryEscape(locationName) + "&key=" + key
	resp, err := http.Get(u)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	var cityReturn CityReturn
	err = json.NewDecoder(resp.Body).Decode(&cityReturn)
	if err != nil {
		return ""
	}
	if len(cityReturn.Location) == 0 {
		return ""
	}
	return cityReturn.Location[0].Id
}

func MakeTodayWeatherMD(cityName string, index IndexReturn, w WeatherReturn) (string, error) {
	if len(index.Daily) == 0 || len(w.Daily) == 0 {
		return "", errors.New("no data")
	}
	todayWeather := w.Daily[0]
	md := misc.MarkdownTool{}
	md.AddTitle(cityName+"今日天气", 3)
	s := ""
	s = `白天%s->晚上%s，温度%s℃-%s℃，湿度%s%%，日出%s，日落%s`
	s = fmt.Sprintf(s, todayWeather.TextDay, todayWeather.TextNight, todayWeather.TempMin, todayWeather.TempMax, todayWeather.Humidity, todayWeather.Sunrise, todayWeather.Sunset)
	md.AddContent(s)
	for i := 0; i < len(index.Daily); i++ {
		name := index.Daily[i].Name
		name = strings.Replace(name, "指数", "", -1)
		md.AddList(name+":"+index.Daily[i].Text, 1)
	}
	return md.ToStr(), nil
}

func GetTodayWeatherMD(cityName string, key string) (string, error) {
	location := QueryCity(cityName, key)
	if location == "" {
		return "", errors.New("city not found")
	}
	indexReturn, err := QueryTodayIndex(key, location)
	if err != nil {
		return "", errors.WithMessage(err, "index error")
	}
	weatherReturn, err := QueryTodayWeather(key, location)
	if err != nil {
		return "", errors.WithMessage(err, "weather error")
	}
	md, err := MakeTodayWeatherMD(cityName, indexReturn, weatherReturn)
	if err != nil {
		return "", errors.WithMessage(err, "make md error")
	}
	return md, nil
}
