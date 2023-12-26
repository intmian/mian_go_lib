package pushmod

import (
	"github.com/intmian/mian_go_lib/tool/misc"
	"net/http"
	"net/url"
)

type PushDeerSetting struct {
	Token string
}

type PushDeerMgr struct {
	setting *PushDeerSetting
	misc.InitTag
}

func (m *PushDeerMgr) SetSetting(setting interface{}) error {
	if !m.IsInitialized() {
		return misc.ErrNotInit
	}
	if setting == nil {
		return ErrTypeErr
	}
	settingT, ok := setting.(*PushDeerSetting)
	if !ok {
		return ErrTypeErr
	}
	m.setting = settingT
	return nil
}

func (m *PushDeerMgr) Init(setting *PushDeerSetting) error {
	m.setting = setting
	m.SetInitialized()
	return nil
}

func NewPushDeerMgr(setting *PushDeerSetting) (*PushDeerMgr, error) {
	m := &PushDeerMgr{}
	err := m.Init(setting)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *PushDeerMgr) Push(title string, content string) error {
	if !m.IsInitialized() {
		return misc.ErrNotInit
	}
	err := PushPushDeer(*m.setting, title, content, false)
	if err != nil {
		return err
	}
	return nil
}

func (m *PushDeerMgr) PushMarkDown(title string, content string) error {
	if !m.IsInitialized() {
		return misc.ErrNotInit
	}
	err := PushPushDeer(*m.setting, title, content, true)
	if err != nil {
		return err
	}
	return nil
}

func PushPushDeer(setting PushDeerSetting, title string, content string, markDown bool) error {
	baseUrl := "https://api2.pushdeer.com/message/push"
	t := ""

	if markDown {
		t = "markdown"
	}

	resp, err := http.PostForm(baseUrl, url.Values{"pushkey": {setting.Token}, "text": {title}, "desp": {content}, "type": {t}})

	if err != nil {
		return err
	}
	if resp.Status != "200 OK" {
		return ErrPushDeerPushFail
	}
	return nil
}
