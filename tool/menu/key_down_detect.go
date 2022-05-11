package menu

import (
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

//WaitKeyDown 阻塞式的等待按键按下
func WaitKeyDown() rune {
	evChan := hook.Start()
	defer func() {
		// ClearIOBuffer() // TODO: 没用
		hook.End()
	}()
	for ev := range evChan {
		if !robotgo.IsValid() {
			continue
		}
		if ev.Kind == hook.KeyDown {
			return ev.Keychar
		}
	}
	return ' '
}

//KeyDownSendChan 返回一个chan，第一次按下key时发送到chan中，后续按下key时不发送
func KeyDownSendChan(key rune) <-chan bool {
	result := make(chan bool)
	go func() {
		evChan := hook.Start()
		defer hook.End()
		for ev := range evChan {
			if ev.Kind == hook.KeyDown && ev.Keychar == key {
				result <- true
				return
			}
		}
	}()
	return result
}

func BindKeyDown(key string, callback func()) {
	hook.Register(hook.KeyDown, []string{key}, func(event hook.Event) {
		callback()
		hook.End()
	})
}

func BindKeysDown(keys []string, callback func()) {
	for _, key := range keys {
		hook.Register(hook.KeyDown, []string{key}, func(e hook.Event) {
			callback()
		})
	}
}
