package xlog

import (
	"testing"
)

func TestLog(t *testing.T) {
	setting := DefaultSetting()
	l, err := NewXLog(setting)
	if err != nil {
		t.Fatal(err)
	}
	l.Log(LogLevelError, "TEST", "testError")
	l.Log(LogLevelWarning, "TEST", "testWarning")
	l.Log(LogLevelInfo, "TEST", "testLog")
	l.Log(LogLevelMisc, "TEST", "testMisc")
	l.Log(LogLevelDebug, "TEST", "testDebug")
}
