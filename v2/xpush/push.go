package xpush

type PushType int8

const (
	PushType_PUSH_EMAIL PushType = iota
	PushType_PUSH_PUSH_DEER
	PushType_PUSH_DING
)

type Mgr struct {
	pushEmailToken  *EmailToken
	pushDeerToken   *PushDeerToken
	pushDingSetting *DingSetting
	dingMgr         *DingRobotMgr
	tag             string
}

func NewMgr(pushEmailToken *EmailToken, pushDeerToken *PushDeerToken, tag string) *Mgr {
	return &Mgr{pushEmailToken: pushEmailToken, pushDeerToken: pushDeerToken, tag: tag}
}

func NewDingMgr(pushDingSetting *DingSetting, tag string) *Mgr {
	return &Mgr{pushDingSetting: pushDingSetting, tag: tag}
}

func (m *Mgr) SetTag(tag string) {
	m.tag = tag
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
	m.pushDeerToken = &PushDeerToken{
		Token: token,
	}
}

// SetDingSetting 设置钉钉设置
func (m *Mgr) SetDingSetting(setting DingSetting) {
	m.pushDingSetting = &setting
}
