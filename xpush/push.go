package xpush

import (
	"sync"

	"github.com/intmian/mian_go_lib/tool/misc"
	"github.com/intmian/mian_go_lib/xpush/pushmod"
	"github.com/pkg/errors"
)

type XPush struct {
	misc.InitTag
	pushMod  map[PushType]IPushMod
	needLock bool // 是否需要锁，用于在多线程中动态修改pushMod
	l        sync.RWMutex
}

func NewXPush(needLock bool) (*XPush, error) {
	m := &XPush{}
	err := m.Init(needLock)
	return m, err
}

func (m *XPush) Init(needLock bool) error {
	m.SetInitialized()
	m.pushMod = make(map[PushType]IPushMod)
	m.needLock = needLock
	return nil
}

func (m *XPush) OnLock() {
	if m.needLock {
		m.l.Lock()
	}
}

func (m *XPush) OnUnlock() {
	if m.needLock {
		m.l.Unlock()
	}
}

func (m *XPush) add(pushType PushType, pushMod IPushMod) error {
	m.OnLock()
	defer m.OnUnlock()
	if !m.IsInitialized() {
		return misc.ErrNotInit
	}
	if !m.IsTypeValid(pushType) {
		return ErrPushTypeInvalid
	}
	if _, ok := m.pushMod[pushType]; ok {
		return ErrPushTypeExist
	}
	m.pushMod[pushType] = pushMod
	return nil
}

func (m *XPush) IsTypeValid(pushType PushType) bool {
	if pushType >= PushTypeMax {
		return false
	}
	if pushType == PushTypeNull {
		return false
	}
	return true
}

func (m *XPush) AddDingDing(setting pushmod.DingSetting) error {
	var Ding pushmod.DingRobotMgr
	Ding.Init(setting)
	return m.add(PushTypeDing, &Ding)
}

func (m *XPush) AddFeishu(setting pushmod.FeishuSetting) error {
	var Feishu pushmod.FeishuRobotMgr
	Feishu.Init(setting)
	return m.add(PushTypeFeishu, &Feishu)
}

func (m *XPush) AddPushDeer(setting pushmod.PushDeerSetting) error {
	var PushDeer pushmod.PushDeerMgr
	err := PushDeer.Init(&setting)
	if err != nil {
		return errors.WithMessage(err, "PushDeer.Init")
	}
	return m.add(PushTypePushDeer, &PushDeer)
}

func (m *XPush) AddEmail(setting pushmod.EmailSetting) error {
	var Email pushmod.EmailMgr
	err := Email.Init(setting)
	if err != nil {
		return errors.WithMessage(err, "Email.Init")
	}
	return m.add(PushTypeEmail, &Email)
}

func (m *XPush) setting(pushType PushType, setting interface{}) error {
	m.OnLock()
	defer m.OnUnlock()
	if !m.IsInitialized() {
		return misc.ErrNotInit
	}
	if !m.IsTypeValid(pushType) {
		return ErrPushTypeInvalid
	}
	if _, ok := m.pushMod[pushType]; !ok {
		return ErrPushTypeNotExist
	}
	err := m.pushMod[pushType].SetSetting(setting)
	if err != nil {
		return err
	}
	return nil
}

func (m *XPush) SetDing(setting pushmod.DingSetting) error {
	return m.setting(PushTypeDing, setting)
}

func (m *XPush) SetPushDeer(setting pushmod.PushDeerSetting) error {
	return m.setting(PushTypePushDeer, setting)
}

func (m *XPush) SetEmail(setting pushmod.EmailSetting) error {
	return m.setting(PushTypeEmail, setting)
}

func (m *XPush) RemovePushType(pushType PushType) {
	m.OnLock()
	defer m.OnUnlock()
	delete(m.pushMod, pushType)
}

func (m *XPush) Push(title string, content string, markDown bool) error {
	m.l.RLock()
	defer m.l.RUnlock()
	if !m.IsInitialized() {
		return misc.ErrNotInit
	}
	var err error
	for _, v := range m.pushMod {
		if markDown {
			err = v.PushMarkDown(title, content)
		} else {
			err = v.Push(title, content)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
