package spider

import "testing"

func TestQueryTodayIndex(t *testing.T) {
	key := "cd893779afce417a8cd6eec2d23b4bc8"
	location := "101210101"
	indexReturn, err := QueryTodayIndex(key, location)
	if err != nil {
		t.Error(err)
	}
	t.Log(indexReturn)
}

func TestQueryCity(t *testing.T) {
	t.Log(QueryCity("杭州", "cd893779afce417a8cd6eec2d23b4bc8"))
}

func TestQueryTodayWeather(t *testing.T) {
	key := "cd893779afce417a8cd6eec2d23b4bc8"
	location := "101210101"
	weatherReturn, err := queryTodayWeather(key, location)
	if err != nil {
		t.Error(err)
	}
	t.Log(weatherReturn)
}

func TestGetTodayWeatherMD(t *testing.T) {
	key := "cd893779afce417a8cd6eec2d23b4bc8"
	location := "杭州"
	s, err := GetTodayWeatherMD(location, key)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)
}
