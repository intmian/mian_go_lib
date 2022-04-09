package menu

import (
	// "github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

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
			hook.End()
		})
	}
}
