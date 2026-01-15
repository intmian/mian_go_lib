package pushmod

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/intmian/mian_go_lib/tool/misc"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"time"
)

type FeishuSetting struct {
	WebhookUrl        string
	SendInterval      int32 // 每隔多少时间
	IntervalSendCount int32 // 有多少次发送机会
	Ctx               context.Context
}

type FeishuRobotMgr struct {
	webhookUrl string
	isInit     bool
	goMgr      misc.GoLimit
}

func (m *FeishuRobotMgr) Init(setting FeishuSetting) error {
	if m.isInit {
		return misc.ErrNotInit
	}
	m.webhookUrl = setting.WebhookUrl
	if setting.Ctx == nil {
		setting.Ctx = context.Background()
	}
	err := m.goMgr.Init(misc.GoLimitSetting{
		TimeInterval:         time.Duration(setting.SendInterval) * time.Second,
		EveryIntervalCallNum: setting.IntervalSendCount,
	}, setting.Ctx)
	if err != nil {
		return errors.WithMessage(err, "FeishuRobotMgr Init")
	}
	err = m.goMgr.Start()
	if err != nil {
		return errors.WithMessage(err, "FeishuRobotMgr Start")
	}

	m.isInit = true
	return nil
}

func NewFeishuRobotMgr(setting FeishuSetting) (*FeishuRobotMgr, error) {
	m := &FeishuRobotMgr{}
	err := m.Init(setting)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *FeishuRobotMgr) Push(title string, content string) error {
	return m.pushFeishu(title, content, false)
}

func (m *FeishuRobotMgr) PushMarkDown(title string, content string) error {
	return m.pushFeishu(title, content, true)
}

func (m *FeishuRobotMgr) SetSetting(setting interface{}) error {
	if !m.isInit {
		return misc.ErrNotInit
	}
	if setting == nil {
		return ErrTypeErr
	}
	settingT, ok := setting.(FeishuSetting)
	if !ok {
		return ErrTypeErr
	}
	m.webhookUrl = settingT.WebhookUrl
	return nil
}

func (m *FeishuRobotMgr) pushFeishu(title string, content string, markDown bool) error {
	if !m.isInit {
		return misc.ErrNotInit
	}
	var mes FeishuMessage
	if markDown {
		mes = NewFeishuCard(title, content)
	} else {
		if title != "" {
			content = title + "\n" + content
		}
		mes = NewFeishuText(content)
	}
	return m.Send(mes)
}

func SendFeishuMessage(webhookUrl string, message FeishuMessage) error {
	messageStr := message.ToJson()
	// 以post的形式发送，header中Content-Type为application/json
	respond, err := http.Post(webhookUrl, "application/json", bytes.NewBufferString(messageStr))
	if err != nil {
		return fmt.Errorf("SendFeishuMessage: %v", err)
	}
	if respond.StatusCode != 200 {
		return fmt.Errorf("SendFeishuMessage: %d %v", respond.StatusCode, respond.Body)
	}
	body, _ := ioutil.ReadAll(respond.Body)
	// {"code":0,"msg":"success","data":{}}
	jsonData := struct {
		Code int32  `json:"code"`
		Msg  string `json:"msg"`
	}{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return fmt.Errorf("SendFeishuMessage: %v", err)
	}
	if jsonData.Code != 0 {
		return fmt.Errorf("SendFeishuMessage: %d|%s", jsonData.Code, jsonData.Msg)
	}
	return nil
}

func (m *FeishuRobotMgr) Send(message FeishuMessage) error {
	if !m.isInit {
		return fmt.Errorf("FeishuRobotMgr not init")
	}
	err := make(chan error)
	err2 := m.goMgr.Call(func() {
		err2 := SendFeishuMessage(m.webhookUrl, message)
		err <- err2
	})
	if err2 != nil {
		return err2
	}
	return <-err
}

type FeishuMessage interface {
	ToJson() string
}

type FeishuText struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

func (m *FeishuText) ToJson() string {
	b, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(b)
}

func NewFeishuText(content string) *FeishuText {
	return &FeishuText{
		MsgType: "text",
		Content: struct {
			Text string `json:"text"`
		}{
			Text: content,
		},
	}
}

type FeishuCard struct {
	MsgType string `json:"msg_type"`
	Card    struct {
		Header *struct {
			Title struct {
				Tag     string `json:"tag"`
				Content string `json:"content"`
			} `json:"title"`
			Template string `json:"template,omitempty"`
		} `json:"header,omitempty"`
		Elements []interface{} `json:"elements"`
	} `json:"card"`
}

type FeishuMarkdownElement struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

func (m *FeishuCard) ToJson() string {
	b, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(b)
}

func NewFeishuCard(title string, content string) *FeishuCard {
	c := &FeishuCard{
		MsgType: "interactive",
	}
	if title != "" {
		c.Card.Header = &struct {
			Title struct {
				Tag     string `json:"tag"`
				Content string `json:"content"`
			} `json:"title"`
			Template string `json:"template,omitempty"`
		}{
			Title: struct {
				Tag     string `json:"tag"`
				Content string `json:"content"`
			}{
				Tag:     "plain_text",
				Content: title,
			},
		}
	}
	c.Card.Elements = []interface{}{
		FeishuMarkdownElement{
			Tag:     "markdown",
			Content: content,
		},
	}
	return c
}
