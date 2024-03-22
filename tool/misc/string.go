package misc

import (
	"github.com/yanyiwu/gojieba"
	"regexp"
	"sort"
	"unicode/utf8"
)

func CutToStrings(str string) []string {
	// 删除所有非中文字符
	//re := regexp.MustCompile("[^\u4e00-\u9fa5]")
	//str = re.ReplaceAllString(str, "")

	var words []string
	x := gojieba.NewJieba()
	defer x.Free()
	words = x.Cut(str, true)
	return words
}

func smallStrings(src []string) []string {
	// 删除所有非中文字符
	reg := regexp.MustCompile("[^\u4e00-\u9fa5]")
	var dst []string
	for _, v := range src {
		v = reg.ReplaceAllString(v, "")
		if v != "" {
			dst = append(dst, v)
		}
	}
	return dst
}

func stringTag(s string) map[string]int {
	words := CutToStrings(s)
	words = smallStrings(words)
	m := make(map[string]int)
	for _, v := range words {
		if utf8.RuneCountInString(v) < 2 {
			continue
		}
		m[v]++
	}
	return m
}

type StringTag struct {
	Tag string
	Num int
}

func string2Tag(s string) []StringTag {
	m := stringTag(s)
	var tags []StringTag
	for k, v := range m {
		tags = append(tags, StringTag{k, v})
	}
	// 排序
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Num > tags[j].Num
	})
	return tags
}
