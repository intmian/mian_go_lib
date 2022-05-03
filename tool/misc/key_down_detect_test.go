package misc

import (
	hook "github.com/robotn/gohook"
	"testing"
)

func TestWaitKeyDown(t *testing.T) {
	evChan := hook.Start()
	defer hook.End()
	for ev := range evChan {
		if ev.Kind == hook.KeyDown || ev.Kind == hook.KeyUp || ev.Kind == hook.KeyHold {
			println(ev.Kind, ev.Keychar, ev.Keycode)
		}

	}
}
