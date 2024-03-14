package spider

import "testing"

func TestGetWeather(t *testing.T) {
	s, err := GetWeather()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)
}
