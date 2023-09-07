package xpush

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/intmian/mian_go_lib/tool/cipher"
	"github.com/intmian/mian_go_lib/tool/misc"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
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

type DingSetting struct {
	Token             string
	Secret            string
	SendInterval      int32 // 每隔多少时间
	IntervalSendCount int32 // 有多少次发送机会
}

func (m *DingRobotMgr) Init(setting DingSetting) {
	m.dingRobotToken.accessToken = setting.Token
	m.dingRobotToken.secret = setting.Secret
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
	body, _ := ioutil.ReadAll(respond.Body)
	print(string(body))
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

type DingText struct {
	At struct {
		AtMobiles []string `json:"atMobiles"`
		AtUserIds []string `json:"atUserIds"`
		IsAtAll   bool     `json:"isAtAll"`
	} `json:"at"`
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
		}(struct{ Content string }{Content: ""}),
	}
}

type DingLink struct {
	MsgType    string `json:"msgtype"`
	Title      string `json:"title"`
	Text       string `json:"text"`
	MessageUrl string `json:"messageUrl"`
	PicUrl     string `json:"picUrl"`
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
	MsgType   string   `json:"msgtype"`
	Title     string   `json:"title"`
	Text      string   `json:"text"`
	AtMobiles []string `json:"atMobiles"`
	AtUserIds []string `json:"atUserIds"`
	IsAtAll   bool     `json:"isAtAll"`
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
	}
}

type btn struct {
	Title     string `json:"title"`
	ActionURL string `json:"actionURL"`
}

type DingActionCard struct {
	MsgType        string `json:"msgtype"`
	Title          string `json:"title"`
	Text           string `json:"text"`
	SingleTitle    string `json:"singleTitle"`
	SingleURL      string `json:"singleURL"`
	btns           []btn  `json:"btns"`
	BtnOrientation int32  `json:"btnOrientation"`
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

type FeedCardLink struct {
	Title      string `json:"title"`
	MessageURL string `json:"messageURL"`
	PicURL     string `json:"picURL"`
}

type FeedCard struct {
	Links []FeedCardLink `json:"links"`
}

type DingFeedCard struct {
	MsgType  string `json:"msgtype"`
	FeedCard `json:"feedCard"`
}

var dingMgr DingRobotMgr
var dingDefaultPushOnce sync.Once

func (m *Mgr) PushDing(title string, content string, markDown bool) error {
	dingDefaultPushOnce.Do(func() {
		dingMgr.Init(*m.pushDingSetting)
	})
	var mes DingMessage
	if markDown {
		mesT := NewDingMarkdown()
		mesT.Text = "##" + title + "\n" + content
		mes = mesT
	} else {
		mesT := NewDingText()
		mesT.Text.Content = title + "\n" + content
		mes = mesT
	}
	return dingMgr.Send(mes)
}
