package push

import (
	"github.com/intmian/mian_go_lib/tool/misc"
	"net/smtp"
	"strings"

	"github.com/Luoxin/go-pushdeer-sdk/psdk"
)

type PushType int8

const (
	PushType_PUSH_EMAIL PushType = iota
	PushType_PUSH_PUSH_DEER
)

type Mgr struct {
	pushEmailToken *EmailToken
	PushDeerToken  *PushDeerToken
}

type EmailToken struct {
	host  string
	User  string
	Token string
}

type PushDeerToken struct {
	Token string
}

func (m *Mgr) PushEmail(from string, to string, title string, content string, markDown bool) bool {
	if markDown {
		content = misc.MarkdownToHTML(content)
	}
	// 发送邮件
	mailType := ""
	if markDown {
		mailType = "html"
	}
	err := SendToMail(m.pushEmailToken.host, m.pushEmailToken.User, m.pushEmailToken.Token, from, to, title, content, mailType)
	if err != nil {
		return false
	}

	return true
}

func (m *Mgr) PushPushDeer(title string, content string, markDown bool) bool {
	p, err := psdk.New("api2.pushdeer.com", m.PushDeerToken.Token)
	if err != nil {
		return false
	}

	req := psdk.MessagePushReq{
		BaseReq: psdk.BaseReq{},
		PushKey: m.PushDeerToken.Token,
		Text:    title,
		Desp:    content,
	}

	if markDown {
		req.Type = "markdown"
	}

	_, err = p.MessagePush(&req)
	if err != nil {
		return false
	}
	return true
}

// SetEmailToken 设置邮件token
func (m *Mgr) SetEmailToken(host string, user string, token string) {
	m.pushEmailToken = &EmailToken{
		host:  host,
		User:  user,
		Token: token,
	}
}

// SetPushDeerToken 设置pushDeer token
func (m *Mgr) SetPushDeerToken(token string) {
	m.PushDeerToken = &PushDeerToken{
		Token: token,
	}
}

func SendToMail(host, user, password, from, to, subject, body, mailType string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailType == "html" {
		content_type = "Content-Type: text/" + mailType + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}
	msg := []byte("To: " + to + "\r\nFrom: " + from + "<" + user + ">" + "\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}
