package xpush

import (
	"github.com/intmian/mian_go_lib/tool/misc"
	"net/smtp"
	"strings"
)

type EmailToken struct {
	host  string
	User  string
	Token string
}

func (m *Mgr) PushEmail(fromEmailAddr string, fromName string, targetEmailAddr string, title string, content string, markDown bool) bool {
	if m.pushEmailToken == nil {
		return false
	}
	if markDown {
		content = misc.MarkdownToHTML(content)
	}
	// 发送邮件
	mailType := ""
	if markDown {
		mailType = "html"
	}
	err := sendToMail(m.pushEmailToken.host, m.pushEmailToken.User, m.pushEmailToken.Token, fromEmailAddr, fromName, targetEmailAddr, title, content, mailType)
	if err != nil {
		return false
	}

	return true
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
