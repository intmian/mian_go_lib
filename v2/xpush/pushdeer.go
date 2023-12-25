package xpush

import (
	"net/http"
	"net/url"
)

type PushDeerToken struct {
	Token string
}

func (m *Mgr) PushPushDeer(title string, content string, markDown bool) (string, bool) {
	if m.tag != "" {
		title = m.tag + ":" + title
	}
	if m.pushDeerToken == nil {
		return "", false
	}
	baseUrl := "https://api2.pushdeer.com/message/push"
	t := ""

	if markDown {
		t = "markdown"
	}

	resp, err := http.PostForm(baseUrl, url.Values{"pushkey": {m.pushDeerToken.Token}, "text": {title}, "desp": {content}, "type": {t}})

	if err != nil {
		return "", false
	}

	return resp.Status, true
}
