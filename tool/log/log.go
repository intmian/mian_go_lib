package log

import (
	"fmt"
	"mio-blog/libs/tool/misc"
	"mio-blog/setting"
	"os"
	"time"
)

// 以后局部全局变量的接口也这样暴露

type log struct {
	logAddr string
}

type TLogLevel uint8

const (
	EError TLogLevel = iota
	EWarning
	ELog
	EDebug
	EMisc
)

var logLevel2Str map[TLogLevel]string = map[TLogLevel]string{
	EError:   "ERROR",
	EWarning: "WARNING",
	ELog:     "LOG",
	EDebug:   "DEBUG",
	EMisc:    "MISC",
}

//log 记录一条日志， from 中应填入来源模块的大写
func (receiver *log) log(level TLogLevel, from string, info string) {
	// log 格式为[级别]\t[日期]\t[发起人] 内容\n
	if !setting.GSetting.Data().Debug && level > ELog {
		// 非debug模式，无视DEBUG与MISC级信息
		return
	}
	sLevel := logLevel2Str[level]

	t := time.Now()
	ts := t.Format("2006-01-02 15:04:05")

	perm := "[%s]\t[%s]\t[%s]\t%s\n"
	perm = fmt.Sprintf(perm, sLevel, ts, from, info)

	switch level {
	case EError:
		print(misc.Red(perm))
	case EWarning:
		print(misc.Yellow(perm))
	case EDebug:
		print(misc.Green(perm))
	default:
		print(perm)
	}

	// TODO 增加邮件
	if from == "LOG" {
		return // 对于自身发起的警告不写入文件，避免循环调用
	}

	fp, err := os.OpenFile(receiver.logAddr+`\`+geneLogAddr(t),
		os.O_WRONLY|os.O_APPEND,
		0666)
	if err != nil {
		receiver.log(EWarning, "LOG", "日志系统无法打开日志文件")
		return
	}
	_, err = fp.Write([]byte(perm))
	if err != nil {
		receiver.log(EWarning, "LOG", "日志系统无法写入日志文件")
		return
	}
	err = fp.Close()
	if err != nil {
		receiver.log(EWarning, "LOG", "日志系统无法关闭日志文件")
		return
	}
}

func geneLogAddr(t time.Time) string {
	perm := `log_%d_%d_%d.txt`
	return fmt.Sprintf(perm, t.Year(), t.Month(), t.Day())
}

var lLog log = log{setting.GSetting.Data().LogPath}

//Log 记录一条日志， from 中应填入来源模块的大写
func Log(level TLogLevel, from string, info string) {
	lLog.log(level, from, info)
}
