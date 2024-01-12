package xnews

import (
	"container/list"
	"github.com/intmian/mian_go_lib/tool/misc"
	"sync"
	"time"
)

type TopicSetting struct {
	Limit         []TopicTimeLimit
	Clear         []TopicClearTime
	DefaultRemain *TopicTimeRemain // 默认的留存时间，如果不设置则为永久留存
}

// Topic 一个主题
type Topic struct {
	misc.InitTag
	topicName string
	TopicSetting
	messageList list.List
	rwLock      sync.RWMutex
	pool        sync.Pool
}

// Init 初始化一个主题, DefaultRemain 为默认的留存时间，如果不设置则为永久留存
func (t *Topic) Init(name string, setting TopicSetting) error {
	*t = Topic{}
	t.topicName = name
	t.TopicSetting = setting
	t.pool.New = func() interface{} {
		return &Message{}
	}
	t.SetInitialized()
	return nil
}

func (t *Topic) SetTopicSetting(setting TopicSetting) error {
	if !t.IsInitialized() {
		return misc.ErrNotInit
	}
	t.rwLock.Lock()
	defer t.rwLock.Unlock()
	t.TopicSetting = setting
	return nil
}

func (t *Topic) IsEmpty() bool {
	if !t.IsInitialized() {
		return true
	}
	t.rwLock.RLock()
	defer t.rwLock.RUnlock()
	return t.messageList.Len() == 0
}

func (t *Topic) Update() {
	if !t.IsInitialized() {
		return
	}
	t.rwLock.Lock()
	defer t.rwLock.Unlock()
	// 先清理
	isCleared := false
	for _, clearT := range t.Clear {
		if time.Now().Sub(clearT.lastTime) < clearT.duration {
			continue
		}
		clearT.lastTime = time.Now()
		if isCleared {
			continue
		}
		// 清空
		for i := t.messageList.Front(); i != nil; i = i.Next() {
			t.pool.Put(i.Value)
		}
		t.messageList.Init()
		isCleared = true
	}
	if isCleared {
		return
	}
	// 清除过期
	for i := t.messageList.Front(); i != nil; i = i.Next() {
		if i.Value.(*Message).CreateTime.Add(time.Duration(*t.DefaultRemain)).After(time.Now()) {
			continue
		}
		t.pool.Put(i.Value)
		t.messageList.Remove(i)
	}
}

func (t *Topic) AddMessage(content string, remain *TopicTimeRemain) error {
	if !t.IsInitialized() {
		return misc.ErrNotInit
	}
	t.rwLock.Lock()
	defer t.rwLock.Unlock()
	msg := t.pool.Get().(*Message)
	msg.Content = content
	msg.CreateTime = time.Now()
	t.messageList.PushBack(msg)
	if remain == nil {
		remain = t.DefaultRemain
	}
	for _, limitT := range t.Limit {
		out := limitT.Add(1)
		if out <= 0 {
			continue
		}
		for i := 0; i < out; i++ {
			t.pool.Put(t.messageList.Front().Value)
			t.messageList.Remove(t.messageList.Front())
		}
	}
	return nil
}

func (t *Topic) Get() ([]string, error) {
	if !t.IsInitialized() {
		return nil, misc.ErrNotInit
	}
	t.rwLock.RLock()
	defer t.rwLock.RUnlock()
	s := make([]string, 0, t.messageList.Len())
	for i := t.messageList.Front(); i != nil; i = i.Next() {
		s = append(s, i.Value.(*Message).Content)
	}
	return s, nil
}

func NewTopic(name string, setting TopicSetting) (*Topic, error) {
	t := &Topic{}
	err := t.Init(name, setting)
	if err != nil {
		return nil, err
	}
	return t, nil
}
