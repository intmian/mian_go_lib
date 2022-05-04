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
	l.log(EError, "TEST", "testError")
	l.log(EWarning, "TEST", "testWarning")
	l.log(ELog, "TEST", "testLog")
	l.log(EMisc, "TEST", "testMisc")
	l.log(EDebug, "TEST", "testDebug")
}
