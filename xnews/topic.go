package xnews

import (
	"container/list"
	"context"
	"github.com/intmian/mian_go_lib/tool/misc"
	"sync"
	"time"
)

type TopicSetting struct {
	Limit         []TopicTimeLimit
	Clear         []TopicClearTime
	DefaultRemain *TopicTimeRemain // 默认的留存时间，如果不设置则为永久留存
	ctx           context.Context
}

func (t *TopicSetting) AddForeverLimit(num int) {
	t.Limit = append(t.Limit, TopicTimeLimit{
		Num: num,
	})
}

func (t *TopicSetting) AddNowLimit(duration time.Duration, num int) {
	t.AddLimit(time.Now(), duration, num)
}

// AddLimit 用于限制某个topic的时间间隔内做多保留，如果添加多条limit，则任意一条达到上限，则都会删除最旧的信息，同时触发多条则会删除多条
func (t *TopicSetting) AddLimit(startTime time.Time, duration time.Duration, num int) {
	t.Limit = append(t.Limit, TopicTimeLimit{
		Duration:        &duration,
		Num:             num,
		LastResetTime:   startTime,
		ThisDurationNum: 0,
	})
}

func (t *TopicSetting) AddNowClear(duration time.Duration) {
	t.AddClear(time.Now(), duration)
}

func (t *TopicSetting) AddClear(startTime time.Time, duration time.Duration) {
	t.Clear = append(t.Clear, TopicClearTime{
		lastTime: startTime,
		duration: duration,
	})
}

func (t *TopicSetting) SetDefaultRemain(duration time.Duration) {
	t.DefaultRemain = (*TopicTimeRemain)(&duration)
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
func (t *Topic) Init(name string, setting TopicSetting, ctx context.Context) error {
	*t = Topic{}
	t.topicName = name
	t.TopicSetting = setting
	t.pool.New = func() interface{} {
		return &Message{}
	}
	t.ctx = ctx
	go t.startGo()
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

func (t *Topic) startGo() {
	// 后面如果有性能上的问题，需要改成协程+时序处理的形式，不要轮询，而是排队 wait时间唤醒
	for {
		select {
		case <-t.ctx.Done():
			return
		default:
			t.update()
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func (t *Topic) update() {
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
		if t.DefaultRemain == nil {
			continue
		}
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
	needDelete := false
	for i, _ := range t.Limit {
		out := t.Limit[i].Add(1)
		if out <= 0 {
			continue
		}
		if out > 0 {
			needDelete = true
		}
	}
	if needDelete {
		t.pool.Put(t.messageList.Front().Value)
		t.messageList.Remove(t.messageList.Front())
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

func NewTopic(name string, setting TopicSetting, ctx context.Context) (*Topic, error) {
	t := &Topic{}
	err := t.Init(name, setting, ctx)
	if err != nil {
		return nil, err
	}
	return t, nil
}
