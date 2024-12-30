package spider

import (
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestGetBBCRss(t *testing.T) {
	bbc, err := GetBBCRss(&http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "http",
				Host:   "localhost:7890",
			}),
		},
	})
	if err != nil {
		t.Errorf("GetBBCRss() error = %v", err)
	}
	print("len(bbc):", len(bbc))

	// 昨天
	bbc, err = GetBBCRssWithDay(time.Now().AddDate(0, 0, -1), &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "http",
				Host:   "localhost:7890",
			}),
		},
	})
	if err != nil {
		t.Errorf("GetBBCRssWithDay() error = %v", err)
	}
	print("len(bbc):", len(bbc))
}
func TestGetGoogleRss(t *testing.T) {
	google, err := GetGoogleRss("特斯拉", &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "http",
				Host:   "localhost:7890",
			}),
		},
	})
	if err != nil {
		t.Errorf("GetGoogleRss() error = %v", err)
	}
	print("len(google):", len(google))

	// 昨天
	google, err = GetGoogleRssWithDay("model3", time.Now().AddDate(0, 0, -1), &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "http",
				Host:   "localhost:7890",
			}),
		},
	})
	if err != nil {
		t.Errorf("GetGoogleRssWithDay() error = %v", err)
	}
	print("len(google):", len(google))
}

func TestGetNYTimesRss(t *testing.T) {
	nytimes, err := GetNYTimesRss(&http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "http",
				Host:   "localhost:7890",
			}),
		},
	})
	if err != nil {
		t.Errorf("GetNYTimesRss() error = %v", err)
	}
	print("len(nytimes):", len(nytimes))

	// 昨天
	nytimes, err = GetNYTimesRssWithDay(time.Now().AddDate(0, 0, -1), &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "http",
				Host:   "localhost:7890",
			}),
		},
	})
	if err != nil {
		t.Errorf("GetNYTimesRssWithDay() error = %v", err)
	}
	print("len(nytimes):", len(nytimes))
}

func TestUTCTime(t *testing.T) {
	now := time.Now().Add(-time.Hour * 24)
	nowUTC := time.Now().UTC()
	for i := 0; i < 40; i++ {
		testTime := nowUTC.Add(time.Duration(-i) * time.Hour)
		if testTime.Year() == now.Year() && testTime.Month() == now.Month() && testTime.Day() == now.Day() {
			t.Logf("testTime:%v same day", testTime)
		} else {
			t.Logf("testTime:%v", testTime)
		}
	}
}
