package pushmod

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/intmian/mian_go_lib/tool/cipher"
	"github.com/intmian/mian_go_lib/tool/misc"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type DingSetting struct {
	Token             string
	Secret            string
	SendInterval      int32 // 每隔多少时间
	IntervalSendCount int32 // 有多少次发送机会
	Ctx               context.Context
}

type DingRobotToken struct {
	accessToken string
	secret      string
}

type DingRobotMgr struct {
	dingRobotToken DingRobotToken
	isInit         bool
	goMgr          misc.GoLimit
}

const ApiUrl = "https://oapi.dingtalk.com/robot/send"

// GetDingSign 获得签名
func GetDingSign(token string) (timestamp string, sign string) {
	timestamp = fmt.Sprintf("%d", time.Now().UnixMilli())
	/*
		把timestamp+"\n"+密钥当做签名字符串，使用HmacSHA256算法计算签名，然后进行Base64 encode，最后再把签名参数再进行urlEncode，得到最终的签名（需要使用UTF-8字符集）。
	*/
	s := timestamp + "\n" + token
	s2 := cipher.HmacSha256Sign(token, s)
	sign = base64.StdEncoding.EncodeToString(s2)
	sign = url.QueryEscape(sign)
	return
}

func (m *DingRobotMgr) Init(setting DingSetting) error {
	if m.isInit {
		return misc.ErrNotInit
	}
	m.dingRobotToken.accessToken = setting.Token
	m.dingRobotToken.secret = setting.Secret
	if setting.Ctx == nil {
		setting.Ctx = context.Background()
	}
	err := m.goMgr.Init(misc.GoLimitSetting{
		TimeInterval:         time.Duration(setting.SendInterval) * time.Second,
		EveryIntervalCallNum: setting.IntervalSendCount,
	}, setting.Ctx)
	if err != nil {
		return errors.WithMessage(err, "DingRobotMgr Init")
	}
	err = m.goMgr.Start()
	if err != nil {
		return errors.WithMessage(err, "DingRobotMgr Start")
	}

	m.isInit = true
	return nil
}

func NewDingRobotMgr(setting DingSetting) (*DingRobotMgr, error) {
	m := &DingRobotMgr{}
	err := m.Init(setting)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *DingRobotMgr) Push(title string, content string) error {
	return m.pushDing(title, content, false)
}

func (m *DingRobotMgr) PushMarkDown(title string, content string) error {
	return m.pushDing(title, content, true)
}

func (m *DingRobotMgr) SetSetting(setting interface{}) error {
	if !m.isInit {
		return misc.ErrNotInit
	}
	if setting == nil {
		return ErrTypeErr
	}
	settingT, ok := setting.(DingSetting)
	if !ok {
		return ErrTypeErr
	}
	m.dingRobotToken.accessToken = settingT.Token
	m.dingRobotToken.secret = settingT.Secret
	return nil
}

func (m *DingRobotMgr) pushDing(title string, content string, markDown bool) error {
	if !m.isInit {
		return misc.ErrNotInit
	}
	var mes DingMessage
	if markDown {
		mesT := NewDingMarkdown()
		mesT.Markdown.Title = title
		mesT.Markdown.Text = content
		mes = mesT
	} else {
		mesT := NewDingText()
		mesT.Text.Content = title + "\n" + content
		mes = mesT
	}
	return m.Send(mes)
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
	body, _ := ioutil.ReadAll(respond.Body)
	// {"errcode":0,"errmsg":"ok"}
	jsonData := struct {
		ErrCode int32  `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return fmt.Errorf("SendDingMessage: %v", err)
	}
	if jsonData.ErrCode != 0 {
		return fmt.Errorf("SendDingMessage: %d|%s", jsonData.ErrCode, jsonData.ErrMsg)
	}
	return nil
}

func (m *DingRobotMgr) Send(message DingMessage) error {
	if !m.isInit {
		return fmt.Errorf("DingRobotMgr not init")
	}
	err := make(chan error)
	err2 := m.goMgr.Call(func() {
		err2 := SendDingMessage(m.dingRobotToken.accessToken, m.dingRobotToken.secret, message)
		err <- err2
	})
	if err2 != nil {
		return err2
	}
	return <-err
}

type DingMessage interface {
	ToJson() string
}

type At struct {
	AtMobiles []string `json:"atMobiles"`
	AtUserIds []string `json:"atUserIds"`
	IsAtAll   bool     `json:"isAtAll"`
}

type DingText struct {
	At   At `json:"at"`
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
	MsgType string `json:"msgtype"`
}

func (m *DingText) ToJson() string {
	b, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(b)
}

func NewDingText() *DingText {
	return &DingText{
		MsgType: "text",
		Text: struct {
			Content string `json:"content"`
		}{},
		At: At{},
	}
}

type DingLink struct {
	MsgType string `json:"msgtype"`
	Link    struct {
		Text       string `json:"text"`
		Title      string `json:"title"`
		PicUrl     string `json:"picUrl"`
		MessageUrl string `json:"messageUrl"`
	} `json:"link"`
}

func (m *DingLink) ToJson() string {
	b, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(b)
}

func NewDingLink() *DingLink {
	return &DingLink{
		MsgType: "link",
	}
}

type DingMarkdown struct {
	MsgType  string `json:"msgtype"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
	At At `json:"at"`
}

func (m *DingMarkdown) ToJson() string {
	b, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(b)
}
func NewDingMarkdown() *DingMarkdown {
	return &DingMarkdown{
		MsgType: "markdown",
		Markdown: struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		}{},
		At: At{},
	}
}

type btn struct {
	Title     string `json:"title"`
	ActionURL string `json:"actionURL"`
}

type DingActionCard struct {
	MsgType    string `json:"msgtype"`
	ActionCard struct {
		Title          string `json:"title"`
		Text           string `json:"text"`
		BtnOrientation string `json:"btnOrientation"`
		Btns           []struct {
			Title     string `json:"title"`
			ActionURL string `json:"actionURL"`
		} `json:"btns"`
		SingleTitle string `json:"singleTitle"`
		SingleURL   string `json:"singleURL"`
	} `json:"actionCard"`
}

func (m *DingActionCard) ToJson() string {
	b, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(b)
}

func NewDingActionCard() *DingActionCard {
	return &DingActionCard{
		MsgType: "actionCard",
	}
}

type DingFeedCard struct {
	MsgType  string `json:"msgtype"`
	FeedCard struct {
		Links []struct {
			Title      string `json:"title"`
			MessageURL string `json:"messageURL"`
			PicURL     string `json:"picURL"`
		} `json:"links"`
	} `json:"feedCard"`
}
