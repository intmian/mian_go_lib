package spider

import (
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	now := time.Now()
	s := GetUniTimeStr(now)
	t.Log(s)
}

func TestGetGNews(t *testing.T) {
	req := GNewsSearch{
		q:      "tesla",
		lang:   LanChinese,
		sortby: SortByPublishedAt,
		from:   GetUniTimeStr(time.Now().AddDate(0, -1, -10)),
		to:     GetUniTimeStr(time.Now()),
	}
	result, err := QueryGNews(req, "ee54b7595ba81fc612c56689416abf6a")
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}

func TestGetGTop(t *testing.T) {
	req := GNewsTop{
		lang: LanEnglish,
		from: GetUniTimeStr(time.Now().AddDate(0, 0, -1)),
		to:   GetUniTimeStr(time.Now()),
	}
	result, err := QueryGNewsTop(req, "ee54b7595ba81fc612c56689416abf6a")
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}

func TestGetGNewsSumTop(t *testing.T) {
	result, err := GetGNewsSumTop("ee54b7595ba81fc612c56689416abf6a", time.Now().AddDate(0, 0, -7))
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}
