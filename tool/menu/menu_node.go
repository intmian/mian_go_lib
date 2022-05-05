package menu

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/intmian/mian_go_lib/tool/misc"
)

//MakeUntilPressFunc 创建一个按下一般都会返回的函数
func MakeUntilPressFunc(f func(chan bool)) func() {
	// 按下esc的话会触发退出，但是只有在真正结束函数的时候才会返回
	endChan := make(chan bool)
	finishChan := make(chan bool)
	return func() {
		go func() {
			f(endChan)
			finishChan <- true
		}()
		misc.BindKeyDown("esc", func() {
			endChan <- true
		})
		for {
			select {
			case <-finishChan:
				return
			}
		}
	}
}

//MakeUntilPressForFunc 创建一个按下esc才停止否则会一直循环调用的函数
func MakeUntilPressForFunc(f func()) func() {
	forFunc := func(endChan chan bool) {
		for {
			select {
			case <-endChan:
				return
			default:
				f()
			}
		}
	}
	return MakeUntilPressFunc(forFunc)
}

//MakeUntilPressForShowFunc 返回一个每隔waitSecond刷新一次printFunc，按下esc才停止的函数
func MakeUntilPressForShowFunc(printFunc func(), waitSecond int) func() {
	forShowFunc := func(endChan chan bool) {
		for {
			select {
			case <-endChan:
				return
			default:
				misc.Clear()
				println(misc.Green("esc.") + "退出")
				printFunc()
				time.Sleep(time.Duration(waitSecond) * time.Second)
			}
		}
	}
	return MakeUntilPressFunc(forShowFunc)
}

func makeValueStr(oriValue, nowValue string) string {
	lenOriValue := len(oriValue)
	lenNowValue := len(nowValue)
	if lenNowValue >= lenOriValue {
		return misc.Green(nowValue)
	} else {
		valueA := nowValue[:lenNowValue]
		valueB := oriValue[lenNowValue+1:]
		return misc.Green(valueA+"_") + misc.Red(valueB)
	}
}

func parseList2text(kv kvSlice, nowIndex int, nowInput string) string {
	var text string
	for _, v := range kv {
		s := "%-15s: %25s\n"
		valueStr := ""
		if v.key == kv[nowIndex].key {
			valueStr = makeValueStr(v.value, nowInput)
		} else {
			valueStr = v.value
		}
		text += fmt.Sprintf(s, v.key, valueStr)
	}
	return text
}

func findMapIndex(kv map[string]string, index int) string {
	if index < 0 {
		return ""
	} else if index >= len(kv) {
		return ""
	} else {
		// 逆序排列的
		for k := range kv {
			if index == 0 {
				return k
			}
			index--
		}
	}
	return ""
}

type kvSlice []struct {
	key   string
	value string
}

//MakeListInputFunc 创建一个输入列表的函数
func MakeListInputFunc(kv map[string]string, callBack func()) func() {
	return func() {
		// 获得kv的第一个key
		nowIndex := 0
		nowInput := ""
		var copySlice []struct {
			key   string
			value string
		}
		for k, v := range kv {
			// 注意一下显示的可能是无序的，是因为hash的机制
			copySlice = append(copySlice, struct {
				key   string
				value string
			}{k, v})
		}
		// 排序
		sort.Slice(copySlice, func(i, j int) bool {
			return copySlice[i].key < copySlice[j].key
		})

		for {
			misc.Clear()
			text := parseList2text(copySlice, nowIndex, nowInput)
			text += "\n" + misc.Green("[]") + "选择，" + misc.Green("esc") + "退出"
			println(text)
			input := misc.WaitKeyDown()
			if input == 27 /*esc*/ {
				copySlice[nowIndex].value = nowInput
				// 回写kv并回调
				for k := range kv {
					for _, vv := range copySlice {
						if vv.key == k {
							kv[k] = vv.value
						}
					}
				}
				callBack()
				return
			} else if input == '[' {
				if nowIndex > 0 {
					if nowInput != "" {
						copySlice[nowIndex].value = nowInput
					}
					nowInput = ""
					nowIndex--
				}

			} else if input == ']' {
				if nowIndex < len(kv)-1 {
					if nowInput != "" {
						copySlice[nowIndex].value = nowInput
					}
					nowInput = ""
					nowIndex++
				}
			} else if input == 8 /*backspace*/ {
				if nowInput != "" {
					nowInput = nowInput[:len(nowInput)-1]
				}
			} else {
				nowInput += string(input)
			}
		}
	}
}

type uniKVMap map[string]interface{}
type uniKVSlice []struct {
	key   string
	value interface{}
}

func interface2text(i interface{}) string {
	switch i.(type) {
	case string:
		return i.(string)
	case int:
		return strconv.Itoa(i.(int))
	case float64:
		return strconv.FormatFloat(i.(float64), 'f', -1, 64)
	case bool:
		if i.(bool) {
			return " √ " // 不加上空格的话配合颜色显示会有问题，不知道为什么
		} else {
			return " × "
		}
	default:
		return ""
	}
}

func text2interface(s string, oriInterface interface{}) interface{} {
	switch oriInterface.(type) {
	case string:
		return s
	case int:
		i, _ := strconv.Atoi(s)
		return i
	case float64:
		f, _ := strconv.ParseFloat(s, 64)
		return f
	case bool:
		if s == "√" {
			s = "true"
		} else if s == "×" {
			s = "false"
		}
	default:
		return ""
	}
	return ""
}

// parseUniList2text 通用的列表解析
func parseUniList2text(kv uniKVSlice, nowIndex int, nowInput string, nowSearch string) string {
	var text string
	text += misc.Green("搜索:") + nowSearch
	if nowIndex == -1 {
		text += misc.Red("_") + "\n"
	} else {
		text += "\n"
	}
showLoop:
	for _, v := range kv {
		// 一个简单的模糊搜索

		if nowSearch != "" {
			nowSearchSub := strings.Split(nowSearch, " ")
			for _, vv := range nowSearchSub {
				if vv == "" {
					continue
				}
				if !strings.Contains(v.key, vv) {
					continue showLoop
				}
			}
		}

		s := "%-15s: %10s\n"
		valueStr := ""
		if nowIndex != -1 && v.key == kv[nowIndex].key {
			valueStr = makeValueStr(interface2text(v.value), nowInput)
		} else {
			valueStr = interface2text(v.value)
		}
		text += fmt.Sprintf(s, v.key, valueStr)
	}
	return text
}

func UniListFind(kv uniKVSlice, nowSearch string) []int {
	var indexs []int
	indexs = make([]int, 0)
	if nowSearch == "" {
		return indexs
	}
	nowSearchSub := strings.Split(nowSearch, " ")
findLoop:
	for i, v := range kv {
		for _, vv := range nowSearchSub {
			if vv == "" {
				continue
			}
			if !strings.Contains(v.key, vv) {
				continue findLoop
			}
			indexs = append(indexs, i)
		}
	}
	return indexs
}

func UniListFindPre(kv uniKVSlice, nowSearch string, nowIndex int) int {
	if nowSearch == "" {
		return -1
	}
	nowSearchSub := strings.Split(nowSearch, " ")
	for i := nowIndex - 1; i >= 0; i-- {
		for _, vv := range nowSearchSub {
			if vv == "" {
				continue
			}
			if !strings.Contains(kv[i].key, vv) {
				continue
			}
			return i
		}
	}
	return -1
}

func UniListFindNext(kv uniKVSlice, nowSearch string, nowIndex int) int {
	if nowSearch == "" {
		return -1
	}
	nowSearchSub := strings.Split(nowSearch, " ")
	for i := nowIndex + 1; i < len(kv); i++ {
		for _, vv := range nowSearchSub {
			if vv == "" {
				continue
			}
			if !strings.Contains(kv[i].key, vv) {
				continue
			}
			return i
		}
	}
	return -1
}

// MakeUniListInputFunc 创建一个通用输入列表的函数
func MakeUniListInputFunc(kv uniKVMap, callBack func()) func() {
	return func() {
		// 获得kv的第一个key
		nowIndex := -1
		nowInput := ""
		searchInput := ""
		var copySlice []struct {
			key   string
			value interface{}
		}
		for k, v := range kv {
			// 注意一下显示的可能是无序的，是因为hash的机制
			copySlice = append(copySlice, struct {
				key   string
				value interface{}
			}{k, v})
		}
		// 排序
		sort.Slice(copySlice, func(i, j int) bool {
			return copySlice[i].key < copySlice[j].key
		})

		for {
			text := parseUniList2text(copySlice, nowIndex, nowInput, searchInput)
			text += "\n" + misc.Green("[]") + "选择，" + misc.Green(".") + "反转，" + misc.Green("/") + "搜索，" + misc.Green("esc") + "退出"
			misc.Clear()
			println(text)
			input := misc.WaitKeyDown()
			if input == 27 /*esc*/ {
				// 保存并返回
				if nowIndex != -1 {
					copySlice[nowIndex].value = text2interface(nowInput, copySlice[nowIndex].value)
				}
				// 回写kv并回调
				for k := range kv {
					for _, vv := range copySlice {
						if vv.key == k {
							kv[k] = vv.value
						}
					}
				}
				callBack()
				return
			} else if input == '[' {
				// 保存并返回上一项
				if nowIndex == -1 {
					continue
				}
				pre := 0
				if searchInput != "" {
					preT := UniListFindPre(copySlice, searchInput, nowIndex)
					pre = preT
				} else {
					pre = nowIndex - 1
				}
				if nowInput != "" {
					copySlice[nowIndex].value = text2interface(nowInput, copySlice[nowIndex].value)
				}
				nowInput = ""
				nowIndex = pre
			} else if input == ']' {
				// 保存并返回下一项
				if nowIndex >= len(copySlice)-1 {
					continue
				}
				next := 0
				if searchInput != "" {
					nextT := UniListFindNext(copySlice, searchInput, nowIndex)
					if nextT == -1 {
						continue
					}
					next = nextT
				} else {
					next = nowIndex + 1
				}
				if nowInput != "" {
					copySlice[nowIndex].value = text2interface(nowInput, copySlice[nowIndex].value)
				}
				nowInput = ""
				nowIndex = next
			} else if input == '/' {
				if nowIndex != -1 {
					if nowInput != "" {
						copySlice[nowIndex].value = nowInput
					}
					nowInput = ""
					nowIndex = -1
				}
			} else if input == 8 /*backspace*/ {
				if nowIndex == -1 {
					if searchInput != "" {
						searchInput = searchInput[:len(searchInput)-1]
					}
				} else {
					if nowInput != "" {
						nowInput = nowInput[:len(nowInput)-1]
					}
				}
			} else if input == '.' {
				// 清空或者翻转
				if nowIndex == -1 {
					return
				}
				switch copySlice[nowIndex].value.(type) {
				case bool:
					copySlice[nowIndex].value = !copySlice[nowIndex].value.(bool)
				default:
					nowInput = ""
				}
			} else {
				if nowIndex == -1 {
					searchInput += string(input)
				} else {
					switch copySlice[nowIndex].value.(type) {
					case bool:
						//copySlice[nowIndex].value = !copySlice[nowIndex].value.(bool)
					default:
						nowInput += string(input)
					}
				}
			}
		}
	}
}
