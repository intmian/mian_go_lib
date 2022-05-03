package menu

import (
	"github.com/intmian/mian_go_lib/tool/misc"
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
				printFunc()
				time.Sleep(time.Duration(waitSecond) * time.Second)
			}
		}
	}
	return MakeUntilPressFunc(forShowFunc)
}
