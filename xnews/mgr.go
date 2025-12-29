package xnews

import (
	"context"
	"errors"
	"sync"

	"github.com/intmian/mian_go_lib/tool/misc"
)

type XNews struct {
	topics map[string]*Topic
	misc.InitTag
	l   sync.RWMutex
	ctx context.Context
}

func NewXNews(ctx context.Context) (*XNews, error) {
	x := &XNews{}
	err := x.Init(ctx)
	if err != nil {
		return nil, err
	}
	return x, nil
}

func (x *XNews) Init(ctx context.Context) error {
	*x = XNews{}
	x.topics = make(map[string]*Topic)
	x.SetInitialized()
	x.ctx = ctx
	return nil
}

func (x *XNews) AddTopic(name string, setting TopicSetting) error {
	if !x.IsInitialized() {
		return misc.ErrNotInit
	}
	x.l.Lock()
	defer x.l.Unlock()
	if _, ok := x.topics[name]; ok {
		return ErrTopicAlreadyExist
	}
	t := &Topic{}
	ctx := x.ctx
	err := t.Init(name, setting, ctx)
	x.topics[name] = t
	if err != nil {
		return errors.Join(err, ErrInit)
	}
	return nil
}

func (x *XNews) DelTopic(name string) error {
	if !x.IsInitialized() {
		return misc.ErrNotInit
	}
	x.l.Lock()
	defer x.l.Unlock()
	if _, ok := x.topics[name]; !ok {
		return ErrTopicNotExist
	}
	delete(x.topics, name)
	return nil
}

func (x *XNews) GetTopic(name string) ([]string, error) {
	if !x.IsInitialized() {
		return nil, misc.ErrNotInit
	}
	x.l.RLock()
	defer x.l.RUnlock()
	if _, ok := x.topics[name]; !ok {
		return nil, ErrTopicNotExist
	}
	topic := x.topics[name]
	result, err := topic.Get()
	if err != nil {
		return nil, errors.Join(err, ErrGetTopicFailed)
	}
	return result, nil
}

func (x *XNews) AddMessage(topic string, message string) error {
	if !x.IsInitialized() {
		return misc.ErrNotInit
	}
	x.l.RLock()
	defer x.l.RUnlock()
	if _, ok := x.topics[topic]; !ok {
		return ErrTopicNotExist
	}
	err := x.topics[topic].AddMessage(message, nil)
	if err != nil {
		return errors.Join(err, ErrAddMessageFailed)
	}
	return nil
}

func (x *XNews) AddMessageWithExpire(topic string, message string, expire TopicTimeRemain) error {
	if !x.IsInitialized() {
		return misc.ErrNotInit
	}
	x.l.RLock()
	defer x.l.RUnlock()
	if _, ok := x.topics[topic]; !ok {
		return ErrTopicNotExist
	}
	err := x.topics[topic].AddMessage(message, &expire)
	if err != nil {
		return errors.Join(err, ErrAddMessageFailed)
	}
	return nil
}
