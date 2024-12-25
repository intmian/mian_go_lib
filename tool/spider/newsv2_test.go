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
	google, err = GetGoogleRssWithDay("特斯拉", time.Now().AddDate(0, 0, -1), &http.Client{
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
