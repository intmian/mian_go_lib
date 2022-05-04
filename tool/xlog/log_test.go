package xlog

import (
	"fmt"
	"github.com/intmian/mian_go_lib/tool/push"
	"testing"
)

func TestLog(t *testing.T) {
	p := push.Mgr{}
	p.SetPushDeerToken("PDU10120Tp8PByEPFdrKiStSvMWeOdeFtwY7GuOmQ")
	pushStyle := []push.PushType{push.PushType_PUSH_PUSH_DEER}
	f := func(msg string) bool {
		fmt.Println(msg)
		return true
	}
	l := NewMgr("\\log", f, &p, pushStyle, true, true, true, true, true, "target@intmian.com", "from@intmian.com", "testlog")
	l.Log(EError, "TEST", "testError")
	l.Log(EWarning, "TEST", "testWarning")
	l.Log(ELog, "TEST", "testLog")
	l.Log(EMisc, "TEST", "testMisc")
	l.Log(EDebug, "TEST", "testDebug")
}
