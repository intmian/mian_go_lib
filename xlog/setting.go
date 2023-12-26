package xlog

import "github.com/intmian/mian_go_lib/xpush"

type LogSetting struct {
	LogInfo
	LogPrint
	LogStrategy
	LogRecordStrategy
	PushInfo
}

type Printer func(string) bool

type LogStrategy struct {
	IfMisc  bool
	IfDebug bool
}

type LogRecordStrategy struct {
	IfPrint bool
	IfPush  bool
	IfFile  bool
}

type PushInfo struct {
	EmailTargetAddr string // 用;分割
	EmailFromAddr   string
	PushMgr         *xpush.XPush
	PushStyle       []xpush.PushType
}

type LogPrint struct {
	Printer    Printer
	IfUseColor bool
}

type LogInfo struct {
	LogAddr string
	LogTag  string // 如果不为空则push时默认加上tag作为来源。用于区分来源
}

// DefaultSetting 默认设置
// 不包含推送相关
func DefaultSetting() LogSetting {
	l := LogSetting{}
	l.LogAddr = "./log"
	l.IfPrint = true
	l.IfUseColor = true
	l.IfFile = true
	l.Printer = func(s string) bool {
		println(s)
		return true
	}
	return l
}
