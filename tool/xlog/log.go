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

type Mgr struct {
	logAddr         string
	printer         Printer
	pushMgr         *xpush.Mgr
	pushStyle       []xpush.PushType
	ifMisc          bool
	ifDebug         bool
	ifPrint         bool
	ifPush          bool
	ifFile          bool
	emailTargetAddr string // 用;分割
	emailFromAddr   string
	logTag          string // 标记日志的类型，用于在推送中区分不同的日志
}

func (receiver *Mgr) SetLogAddr(logAddr string) {
	receiver.logAddr = logAddr
}

func (receiver *Mgr) SetPrinter(printer Printer) {
	receiver.printer = printer
}

func (receiver *Mgr) SetPushMgr(pushMgr *xpush.Mgr) {
	receiver.pushMgr = pushMgr
}

func (receiver *Mgr) SetPushStyle(pushStyle []xpush.PushType) {
	receiver.pushStyle = pushStyle
}

func (receiver *Mgr) SetIfMisc(ifMisc bool) {
	receiver.ifMisc = ifMisc
}

func (receiver *Mgr) SetIfDebug(ifDebug bool) {
	receiver.ifDebug = ifDebug
}

func (receiver *Mgr) SetIfPrint(ifPrint bool) {
	receiver.ifPrint = ifPrint
}

func (receiver *Mgr) SetIfPush(ifPush bool) {
	receiver.ifPush = ifPush
}

func (receiver *Mgr) SetIfFile(ifFile bool) {
	receiver.ifFile = ifFile
}

func (receiver *Mgr) SetEmailTargetAddr(emailTargetAddr string) {
	receiver.emailTargetAddr = emailTargetAddr
}

func (receiver *Mgr) SetEmailFromAddr(emailFromAddr string) {
	receiver.emailFromAddr = emailFromAddr
}

func (receiver *Mgr) SetLogTag(logTag string) {
	receiver.logTag = logTag
}

func SimpleNewMgr(pushMgr *xpush.Mgr, emailTargetAddr string, emailFromAddr string, logTag string) *Mgr {
	m := &Mgr{pushMgr: pushMgr, emailTargetAddr: emailTargetAddr, emailFromAddr: emailFromAddr, logTag: logTag}
	m.printer = func(s string) bool {
		fmt.Println(s)
		return true
	}
	m.logAddr = "\\log"
	m.pushStyle = []xpush.PushType{xpush.PushType_PUSH_PUSH_DEER}
	m.ifMisc = true
	m.ifPrint = true
	m.ifPush = true
	m.ifFile = true

	return m
}

func NewMgr(logAddr string, printer Printer, pushMgr *xpush.Mgr, pushStyle []xpush.PushType, ifMisc bool, ifDebug bool, ifPrint bool, ifPush bool, ifFile bool, emailTargetAddr string, emailFromAddr string, logTag string) *Mgr {
	return &Mgr{logAddr: logAddr, printer: printer, pushMgr: pushMgr, pushStyle: pushStyle, ifMisc: ifMisc, ifDebug: ifDebug, ifPrint: ifPrint, ifPush: ifPush, ifFile: ifFile, emailTargetAddr: emailTargetAddr, emailFromAddr: emailFromAddr, logTag: logTag}
}

type Setting struct {
}

var logLevel2Str map[TLogLevel]string = map[TLogLevel]string{
	EError:   "ERROR",
	EWarning: "WARNING",
	ELog:     "LOG",
	EMisc:    "MISC",
	EDebug:   "DEBUG",
}

//detailLog 记录详细日志，并返回失败的模块
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

		if !receiver.printer(printContent) {
			errors = append(errors, fmt.Errorf("print failed"))
		}
	}

	if ifPush && level <= EWarning {
		for _, pushType := range receiver.pushStyle {
			switch pushType {
			case xpush.PushType_PUSH_EMAIL:
				if !receiver.pushMgr.PushEmail(receiver.emailFromAddr, receiver.logTag, receiver.emailTargetAddr, receiver.logTag+" "+sLevel+" log", content, false) {
					errors = append(errors, fmt.Errorf("push failed"))
				}
			case xpush.PushType_PUSH_PUSH_DEER:
				if _, suc := receiver.pushMgr.PushPushDeer(receiver.logTag+" "+sLevel+" log", content, false); !suc {
					errors = append(errors, fmt.Errorf("push failed"))
				}
			}
		}
	}

	if ifFile {
		fp, err := os.OpenFile(receiver.logAddr+`\`+geneLogAddr(t),
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

//Log 记录一条日志， from 中应填入来源模块的大写
func (receiver *Mgr) Log(level TLogLevel, from string, info string) {
	errors := receiver.detailLog(level, from, info, receiver.ifMisc, receiver.ifDebug, receiver.ifPrint, receiver.ifPush, receiver.ifFile)

	// 如果有错误，则排除发生错误的那一种记录方式并将剩余的记录方式记录，同时发送一个日志错误日志
	canPrint := receiver.ifPrint
	canPush := receiver.ifPush
	canFile := receiver.ifFile
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

func parseLog(sLevel string, ts string, from string, info string) string {
	perm := "[%s]\t[%s]\t[%s]\t%s\n"
	perm = fmt.Sprintf(perm, sLevel, ts, from, info)
	return perm
}

func geneLogAddr(t time.Time) string {
	perm := `log_%d_%d_%d.txt`
	return fmt.Sprintf(perm, t.Year(), t.Month(), t.Day())
}
