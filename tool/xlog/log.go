package xlog

import (
	"fmt"
	"github.com/intmian/mian_go_lib/tool/misc"
	"github.com/intmian/mian_go_lib/tool/xpush"
	"os"
	"strings"
	"time"
)

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
	PushMgr         *xpush.Mgr
	PushStyle       []xpush.PushType
}

type LogPrint struct {
	Printer    Printer
	IfUseColor bool
}

type LogInfo struct {
	LogAddr string
	LogTag  string // 标记日志的类型，用于在推送中区分不同的日志
}

type Mgr struct {
	LogInfo
	LogPrint
	LogStrategy
	LogRecordStrategy
	PushInfo
}

func SimpleNewMgr(pushMgr *xpush.Mgr, emailTargetAddr string, emailFromAddr string, logTag string) *Mgr {
	m := &Mgr{}
	m.PushMgr = pushMgr
	m.EmailTargetAddr = emailTargetAddr
	m.EmailFromAddr = emailFromAddr
	m.LogTag = logTag
	m.Printer = func(s string) bool {
		fmt.Println(s)
		return true
	}
	m.LogAddr = "\\log"
	m.PushStyle = []xpush.PushType{xpush.PushType_PUSH_PUSH_DEER}
	m.IfMisc = true
	m.IfPrint = true
	m.IfPush = true
	m.IfFile = true

	return m
}

func NewMgr(logAddr string, printer Printer, pushMgr *xpush.Mgr, pushStyle []xpush.PushType, ifMisc bool, ifDebug bool, ifPrint bool, ifPush bool, ifFile bool, emailTargetAddr string, emailFromAddr string, logTag string) *Mgr {
	m := &Mgr{}
	m.LogAddr = logAddr
	m.Printer = printer
	m.PushMgr = pushMgr
	m.PushStyle = pushStyle
	m.IfMisc = ifMisc
	m.IfDebug = ifDebug
	m.IfPrint = ifPrint
	m.IfPush = ifPush
	m.IfFile = ifFile
	m.EmailTargetAddr = emailTargetAddr
	m.EmailFromAddr = emailFromAddr
	m.LogTag = logTag
	return m
}

var logLevel2Str map[TLogLevel]string = map[TLogLevel]string{
	EError:   "ERROR",
	EWarning: "WARNING",
	ELog:     "LOG",
	EMisc:    "MISC",
	EDebug:   "DEBUG",
}

// detailLog 记录详细日志，并返回失败的模块
func (receiver *Mgr) detailLog(level TLogLevel, from string, info string, ifMisc, ifDebug, ifPrint, ifPush, ifFile bool) []error {
	errors := make([]error, 0)
	if !ifMisc && level == EMisc {
		return nil
	}
	if !ifDebug && level == EDebug {
		return nil
	}

	// log 格式为[级别]\t[日期]\t[发起人] 内容\n
	sLevel := logLevel2Str[level]

	t := time.Now()
	ts := t.Format("2006-01-02 15:04:05")

	content := parseLog(sLevel, ts, from, info)
	if ifPrint {
		var printContent string
		switch level {
		case EError:
			printContent = misc.Red(content)
		case EWarning:
			printContent = misc.Yellow(content)
		case EDebug:
			printContent = misc.Green(content)
		default:
			printContent = content
		}

		if !receiver.Printer(printContent) {
			errors = append(errors, fmt.Errorf("print failed"))
		}
	}

	if ifPush && level <= EWarning {
		for _, pushType := range receiver.PushStyle {
			switch pushType {
			case xpush.PushType_PUSH_EMAIL:
				if !receiver.PushMgr.PushEmail(receiver.EmailFromAddr, receiver.LogTag, receiver.EmailTargetAddr, receiver.LogTag+" "+sLevel+" log", content, false) {
					errors = append(errors, fmt.Errorf("push failed"))
				}
			case xpush.PushType_PUSH_PUSH_DEER:
				if _, suc := receiver.PushMgr.PushPushDeer(receiver.LogTag+" "+sLevel+" log", content, false); !suc {
					errors = append(errors, fmt.Errorf("push failed"))
				}
			case xpush.PushType_PUSH_DING:
				err := receiver.PushMgr.PushDing(receiver.LogTag+" "+sLevel+" log", content, false)
				if err != nil {
					errors = append(errors, fmt.Errorf("push failed"))
				}
			}
		}
	}

	if ifFile {
		fp, err := os.OpenFile(receiver.LogAddr+`/`+geneLogAddr(t),
			os.O_WRONLY|os.O_APPEND|os.O_CREATE,
			0666)
		isErr := false
		if err != nil {
			isErr = true
		}
		_, err = fp.Write([]byte(content))
		if err != nil {
			isErr = true
		}
		err = fp.Close()
		if err != nil {
			isErr = true
		}
		if isErr {
			errors = append(errors, fmt.Errorf("file failed"))
		}
	}

	return errors
}

// Log 记录一条日志， from 中应填入来源模块的大写
func (receiver *Mgr) Log(level TLogLevel, from string, info string) {
	errors := receiver.detailLog(level, from, info, receiver.IfMisc, receiver.IfDebug, receiver.IfPrint, receiver.IfPush, receiver.IfFile)

	// 如果有错误，则排除发生错误的那一种记录方式并将剩余的记录方式记录，同时发送一个日志错误日志
	canPrint := receiver.IfPrint
	canPush := receiver.IfPush
	canFile := receiver.IfFile
	errorReason := ""
	if errors != nil && len(errors) > 0 {
		for _, err := range errors {
			if err.Error() == "print failed" {
				canPrint = false
				errorReason += "print failed;"
			}
			if err.Error() == "push failed" {
				canPush = false
				errorReason += "push failed;"
			}
			if err.Error() == "file failed" {
				canFile = false
				errorReason += "file failed;"
			}
		}
		errorReason = strings.TrimRight(errorReason, ";")
		err := receiver.detailLog(EError, "LOG", errorReason, true, true, canPrint, canPush, canFile)
		if err != nil && len(err) > 0 {
			fmt.Println("救救我，我的日志记录有问题！")
		}
	}
}

func (receiver *Mgr) LogWithErr(level TLogLevel, from string, err error) {
	receiver.Log(EError, from, err.Error())
}

func parseLog(sLevel string, ts string, from string, info string) string {
	perm := "[%s]\t[%s]\t[%s]\t%s\n"
	perm = fmt.Sprintf(perm, sLevel, ts, from, info)
	return perm
}

func geneLogAddr(t time.Time) string {
	perm := `log_%d_%d_%d.txt`
	return fmt.Sprintf(perm, t.Year(), t.Month(), t.Day())
}

func GoWaitError(log *Mgr, c <-chan error, from string, s string) {
	go func() {
		err := <-c
		if err != nil {
			log.Log(EError, from, fmt.Sprintf("%s:%s", s, err.Error()))
		}
	}()
}
