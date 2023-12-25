package token

import (
	"github.com/intmian/mian_go_lib/tool/cipher"
	"strconv"
	"strings"
	"sync"
	"time"
)

type JwtMgr struct {
	salt1 string
	salt2 string
	rwL   sync.RWMutex
}

func NewJwtMgr(salt1 string, salt2 string) *JwtMgr {
	m := &JwtMgr{salt1: salt1, salt2: salt2}
	m.rwL = sync.RWMutex{}
	return m
}

type Data struct {
	User       string
	Permission []string
	ValidTime  int64 // 时间戳
	token      string
}

func (m *JwtMgr) SetSalt(salt1, salt2 string) {
	// 上个锁。避免出现生成时一半是salt1，一半是salt2的情况
	m.rwL.Lock()
	defer m.rwL.Unlock()
	m.salt1 = salt1
	m.salt2 = salt2
}

func (m *JwtMgr) GenToken(user string, permission []string, validTime int64) string {
	m.rwL.RLock()
	defer m.rwL.RUnlock()
	s1 := m.salt1
	s2 := m.salt2
	key := user + strconv.FormatInt(validTime, 10) + strings.Join(permission, "")
	token := cipher.Sha2562String(s1 + key)
	token = cipher.Sha2562String(s2 + token)
	return token
}

func (m *JwtMgr) CheckSignature(data *Data, now time.Time, wantPermission string) bool {
	// 时间戳
	if data.ValidTime < now.Unix() {
		return false
	}
	// 检查有效性
	if m.GenToken(data.User, data.Permission, data.ValidTime) != data.token {
		return false
	}
	// 如果没有权限，则返回false
	for _, p := range data.Permission {
		if p == wantPermission {
			return true
		}
	}
	return false
}

func (m *JwtMgr) Signature(data *Data) {
	data.token = m.GenToken(data.User, data.Permission, data.ValidTime)
}
