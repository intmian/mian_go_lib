package xpush

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/intmian/mian_go_lib/tool/cipher"
	"github.com/intmian/mian_go_lib/tool/misc"
	"net/http"
	"net/url"
	"time"
)

type DingRobotToken struct {
	accessToken string
	secret      string
}

type DingRobotMgr struct {
	dingRobotToken DingRobotToken
	isInit         bool
	goMgr          misc.LimitMcoCallFuncMgr
}

type DingMessage interface {
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
	sign = base64.StdEncoding.EncodeToString([]byte(s2))
	sign = url.QueryEscape(sign)
	return
}

type DingSetting struct {
	SendInterval      int32 // 每隔多少时间
	IntervalSendCount int32 // 有多少次发送机会
}

func (m *DingRobotMgr) Init(token string, secret string, setting DingSetting) {
	m.dingRobotToken.accessToken = token
	m.dingRobotToken.secret = secret
	m.isInit = true
	m.goMgr.Init(misc.LimitMCoCallFuncMgrSetting{
		TimeInterval:         setting.SendInterval,
		EveryIntervalCallNum: setting.IntervalSendCount,
	})
}

func SendDingMessage(accessToken, secret string, message DingMessage) error {
	timestamp, sign := GetDingSign(secret)
	messageStr := message.ToJson()
	// 以post的形式发送，header中Content-Type为application/json，access_token\timestamp\sign为url中的参数，请求body中为json格式的数据
	url1 := ApiUrl + "?access_token=" + accessToken + "&timestamp=" + timestamp + "&sign=" + sign
	respond, err := http.Post(url1, "application/json", bytes.NewBufferString(messageStr))
	if err != nil {
		return fmt.Errorf("SendDingMessage: %v", err)
	}
	if respond.StatusCode != 200 {
		return fmt.Errorf("SendDingMessage: %d %v", respond.StatusCode, respond.Body)
	}
	return nil
}

func (m *DingRobotMgr) Send(message DingMessage) error {
	if !m.isInit {
		return fmt.Errorf("DingRobotMgr not init")
	}
	err := make(chan error)
	m.goMgr.Call(func() {
		err2 := SendDingMessage(m.dingRobotToken.accessToken, m.dingRobotToken.secret, message)
		err <- err2
	})
	return <-err
}
