package xpush

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/intmian/mian_go_lib/tool/cipher"
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
	message        chan DingMessageCall
	isInit         bool
}

type DingMessage interface {
	ToJson() string
}

type DingMessageCall struct {
	DingMessage
	ret chan error
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
	SendInterval      int64 // 每隔多少时间
	IntervalSendCount int64 // 有多少次发送机会
}

func (m *DingRobotMgr) Init(token string, secret string) {
	m.dingRobotToken.accessToken = token
	m.dingRobotToken.secret = secret
	m.message = make(chan DingMessageCall)
	m.isInit = true
}

func SendDingMessage(accessToken, secret string, message DingMessage) error {
	timestamp, sign := GetDingSign(secret)
	messageStr := message.ToJson()
	// 以post的形式发送，header中Content-Type为application/json，access_token\timestamp\sign为url中的参数，请求body中为json格式的数据
	url := ApiUrl + "?access_token=" + accessToken + "&timestamp=" + timestamp + "&sign=" + sign
	respond, err := http.Post(url, "application/json", bytes.NewBufferString(messageStr))
	if err != nil {
		return fmt.Errorf("SendDingMessage: %v", err)
	}
	if respond.StatusCode != 200 {
		return fmt.Errorf("SendDingMessage: %d %v", respond.StatusCode, respond.Body)
	}
	return nil
}

func (m *DingRobotMgr) Run(setting DingSetting, end <-chan bool) error {
	if !m.isInit {
		return fmt.Errorf("DingRobotMgr not init")
	}
	if setting.SendInterval < 0 || setting.IntervalSendCount < 0 {
		return fmt.Errorf("DingRobotMgr setting error %v", setting)
	}
	if setting.SendInterval == 0 && setting.IntervalSendCount > 0 {
		return fmt.Errorf("DingRobotMgr setting error %v", setting)
	}
	go func() {
		for true {
			count := setting.IntervalSendCount
			select {
			case <-time.After(time.Duration(setting.SendInterval) * time.Second):
				continue
			case message := <-m.message:
				go func() {
					err := SendDingMessage(m.dingRobotToken.accessToken, m.dingRobotToken.secret, message)
					if err != nil {
						message.ret <- err
					} else {
						message.ret <- nil
					}
				}()
			case <-end:
				return
			}
		}
	}()
	return nil
}

func (m *DingRobotMgr) Send(message DingMessage) error {
	if !m.isInit {
		return fmt.Errorf("DingRobotMgr not init")
	}
	messageCall := DingMessageCall{DingMessage: message, ret: make(chan error)}
	m.message <- messageCall
	err := <-messageCall.ret
	return err
}
