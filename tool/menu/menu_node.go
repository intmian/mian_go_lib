package menu

import (
	"fmt"
	"github.com/intmian/mian_go_lib/tool/misc"
	"sort"
	"time"
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
		s := "%-10s: %10s\n"
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
			} else {
				nowInput += string(input)
			}
		}
	}
}
