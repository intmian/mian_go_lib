package spider

import (
	"github.com/intmian/mian_go_lib/tool/misc"
	"testing"
)

func TestGetWeather(t *testing.T) {
	s, err := GetWeatherDataOri("浙江", "杭州")
	if err != nil {
		t.Fatal(err)
	}
	s = misc.ReplaceUnicodeEscapes(s)
	weather, err := GetTodayWeather(s)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(weather)
}
