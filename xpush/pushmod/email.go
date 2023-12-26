package pushmod

import (
	"github.com/intmian/mian_go_lib/tool/misc"
	"net/smtp"
	"strings"
)

type EmailSetting struct {
	host            string
	User            string
	Token           string
	fromEmailAddr   string
	fromName        string
	targetEmailAddr string
}

type EmailMgr struct {
	misc.InitTag
	EmailSetting EmailSetting
}

func (m *EmailMgr) Push(title string, content string) error {
	if !m.IsInitialized() {
		return misc.ErrNotInit
	}
	return PushEmail(m.EmailSetting.host, m.EmailSetting.User, m.EmailSetting.Token, m.EmailSetting.fromEmailAddr, m.EmailSetting.fromName, m.EmailSetting.targetEmailAddr, title, content, false)
}

func (m *EmailMgr) PushMarkDown(title string, content string) error {
	if !m.IsInitialized() {
		return misc.ErrNotInit
	}
	return PushEmail(m.EmailSetting.host, m.EmailSetting.User, m.EmailSetting.Token, m.EmailSetting.fromEmailAddr, m.EmailSetting.fromName, m.EmailSetting.targetEmailAddr, title, content, true)
}

func (m *EmailMgr) SetSetting(setting interface{}) error {
	if !m.IsInitialized() {
		return misc.ErrNotInit
	}
	if setting == nil {
		return ErrTypeErr
	}
	settingT, ok := setting.(EmailSetting)
	if !ok {
		return ErrTypeErr
	}
	m.EmailSetting = settingT
	return nil
}

func (m *EmailMgr) Init(setting EmailSetting) error {
	m.EmailSetting = setting
	m.SetInitialized()
	return nil
}

func NewEmailMgr(setting EmailSetting) (*EmailMgr, error) {
	m := &EmailMgr{}
	m.EmailSetting = setting
	m.SetInitialized()
	return m, nil
}

func PushEmail(host, user, token, fromEmailAddr string, fromName string, targetEmailAddr string, title string, content string, markDown bool) error {
	if markDown {
		content = misc.MarkdownToHTML(content)
	}
	// 发送邮件
	mailType := ""
	if markDown {
		mailType = "html"
	}
	err := sendToMail(host, user, token, fromEmailAddr, fromName, targetEmailAddr, title, content, mailType)
	if err != nil {
		return err
	}
	return nil
}

func sendToMail(serverHost, serverUser, serverPwd, fromEmailAddr, fromName, targetEmailAddr, subject, body, mailType string) error {
	hp := strings.Split(serverHost, ":")
	auth := smtp.PlainAuth("", serverUser, serverPwd, hp[0])
	var content_type string
	if mailType == "html" {
		content_type = "Content-Type: text/" + mailType + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}
	msg := []byte("To: " + targetEmailAddr + "\r\nFrom: " + fromEmailAddr + "<" + fromName + ">" + "\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	sendTo := strings.Split(targetEmailAddr, ";")
	err := smtp.SendMail(serverHost, auth, serverUser, sendTo, msg)
	return err
}
