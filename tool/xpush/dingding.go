package xpush

import (
	"fmt"
	"github.com/intmian/mian_go_lib/tool/cipher"
	"time"
)

type DingRobotToken struct {
	accessToken string
	secret      string
}

type DingRobotMgr struct {
	dingRobotToken DingRobotToken
	message        chan string
}

type DingMessage interface {
	QuickAddIdAt(id string)
	QuickAddMobileAt(mobile string)
	ToJson() string
}

const ApiUrl = "https://oapi.dingtalk.com/robot/send"

// GetDingSign 获得签名
func GetDingSign(token string) (timestamp string, sign string) {
	timestamp = fmt.Sprintf("%d", time.Now().UnixNano()/1e6)
	/*
		把timestamp+"\n"+密钥当做签名字符串，使用HmacSHA256算法计算签名，然后进行Base64 encode，最后再把签名参数再进行urlEncode，得到最终的签名（需要使用UTF-8字符集）。
	*/
	s := timestamp + "\n" + token
	s2 := cipher.Sha2562String(s)
	return timestamp, s2
}
