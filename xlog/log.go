package xlog

import (
	"errors"
	"fmt"
	"github.com/intmian/mian_go_lib/tool/misc"
	"os"
	"strings"
	"time"
)

// XLog 日志管理器，支持以对应策略记录日志。目前暂不支持复杂策略，可以考虑组合使用多个xLog来实现
type XLog struct {
	LogSetting
	misc.InitTag
}

// NewXlog 创建一个日志管理器
func NewXlog(setting LogSetting) (*XLog, error) {
	m := &XLog{}
	err := m.Init(setting)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (receiver *XLog) Init(setting LogSetting) error {
	if !setting.IfPrint && !setting.IfPush && !setting.IfFile {
		return ErrNoLogWay
	}
	receiver.LogSetting = setting
	// 如果文件夹不存在则创建
	_, err := os.Stat(receiver.LogAddr)
	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(receiver.LogAddr, os.ModePerm)
		if err != nil {
			return err
		}
	}
	receiver.SetInitialized()
	return nil
}

// Log 记录一条日志， from 中应填入来源模块的大写，
// 例如from = "TEST"，则日志中会显示[TEST]，用于区分来源
// 因为别的模块的error处理等都是通过日志模块来进行的，所以日志模块的错误处理只能通过print来进行
func (receiver *XLog) Log(level LogLevel, from string, info string) {
	if !receiver.IsInitialized() {
		fmt.Println("日志模块未初始化！")
		return
	}
	err := receiver.detailLog(level, from, info, receiver.IfMisc, receiver.IfDebug, receiver.IfPrint, receiver.IfPush, receiver.IfFile)

	// 如果有错误，则排除发生错误的那一种记录方式并将剩余的记录方式记录，同时发送一个日志错误日志
	canPrint := receiver.IfPrint
	canPush := receiver.IfPush
	canFile := receiver.IfFile
	errorReason := ""
	// 如果日志出现记录失败，则需要去除掉失败的方式重新记录记录失败
	if err != nil {
		if errors.Is(err, ErrPrintFail) {
			canPrint = false
			errorReason += "print failed;"
		}
		if errors.Is(err, ErrPushPushDeerFail) || errors.Is(err, ErrPushEmailFail) || errors.Is(err, ErrPushDingFail) {
			canPush = false
			errorReason += "push failed;"
		}
		if errors.Is(err, ErrFileFail) {
			canFile = false
			errorReason += "file failed;"
		}
		errorReason = strings.TrimRight(errorReason, ";")
		// 如果所有的记录方式都失败了，那么就直接print
		if !canPrint && !canPush && !canFile {
			fmt.Println("日志模块出现问题，无法记录日志！")
			return
		}
		err = receiver.detailLog(LogLevelError, "LOG", errorReason, true, true, canPrint, canPush, canFile)
		if err != nil {
			fmt.Println("日志模块出现问题，无法记录日志！")
		}
	}
}

func (receiver *XLog) Error(from string, info string) {
	receiver.Log(LogLevelError, from, info)
}

func (receiver *XLog) ErrorErr(from string, err error) {
	receiver.Log(LogLevelError, from, err.Error())
}

func (receiver *XLog) Warning(from string, info string) {
	receiver.Log(LogLevelWarning, from, info)
}

func (receiver *XLog) WarningErr(from string, err error) {
	receiver.Log(LogLevelWarning, from, err.Error())
}

func (receiver *XLog) Info(from string, info string) {
	receiver.Log(LogLevelInfo, from, info)
}

func (receiver *XLog) Misc(from string, info string) {
	receiver.Log(LogLevelMisc, from, info)
}

func (receiver *XLog) Debug(from string, info string) {
	receiver.Log(LogLevelDebug, from, info)
}

func GoWaitError(log *XLog, c <-chan error, from string, s string) {
	if c == nil {
		return
	}
	go func() {
		err := <-c
		if err != nil {
			log.Log(LogLevelError, from, fmt.Sprintf("%s:%s", s, err.Error()))
		}
	}()
}

var logLevel2Str = map[LogLevel]string{
	LogLevelError:   "ERROR",
	LogLevelWarning: "WARNING",
	LogLevelInfo:    "INFO",
	LogLevelMisc:    "MISC",
	LogLevelDebug:   "DEBUG",
}

// detailLog 根据日志配置，记录详细日志，并返回失败的模块
func (receiver *XLog) detailLog(level LogLevel, from string, info string, ifMisc, ifDebug, ifPrint, ifPush, ifFile bool) error {
	var err error
	if !ifMisc && level == LogLevelMisc {
		return nil
	}
	if !ifDebug && level == LogLevelDebug {
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
		case LogLevelError:
			printContent = misc.Red(content)
		case LogLevelWarning:
			printContent = misc.Yellow(content)
		case LogLevelDebug:
			printContent = misc.Green(content)
		default:
			printContent = content
		}

		if !receiver.Printer(printContent) {
			err = errors.Join(err, ErrPrintFail)
		}
	}

	if ifPush && level <= LogLevelWarning {
		err2 := receiver.PushMgr.Push(receiver.LogTag+" "+sLevel+" log", content, false)
		err = errors.Join(err, err2)
	}

	if ifFile {
		fp, err2 := os.OpenFile(receiver.LogAddr+`/`+geneLogAddr(t),
			os.O_WRONLY|os.O_APPEND|os.O_CREATE,
			0666)
		isErr := false
		if err2 != nil {
			isErr = true
		}
		_, err2 = fp.Write([]byte(content))
		if err2 != nil {
			isErr = true
		}
		err2 = fp.Close()
		if err2 != nil {
			isErr = true
		}
		if isErr {
			err = errors.Join(err, ErrFileFail)
		}
	}

	return err
}

func parseLog(sLevel string, ts string, from string, info string) string {
	perm := "[%s]\t[%s]\t[%s]\t%s\n"
	perm = fmt.Sprintf(perm, sLevel, ts, from, info)
	return perm
}

func geneLogAddr(t time.Time) string {
	perm := `%d_%d_%d.txt`
	return fmt.Sprintf(perm, t.Year(), t.Month(), t.Day())
}

const (
	ErrPrintFail        = misc.ErrStr("print failed")
	ErrPushPushDeerFail = misc.ErrStr("push failed")
	ErrPushEmailFail    = misc.ErrStr("push failed")
	ErrPushDingFail     = misc.ErrStr("push failed")
	ErrFileFail         = misc.ErrStr("file failed")
	ErrNoLogWay         = misc.ErrStr("no log way")
)
