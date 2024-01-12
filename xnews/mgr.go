package xnews

import (
	"errors"
	"github.com/intmian/mian_go_lib/tool/misc"
	"sync"
)

type XNews struct {
	topics map[string]*Topic
	misc.InitTag
	l sync.RWMutex
}

func (x *XNews) Init() error {
	*x = XNews{}
	x.topics = make(map[string]*Topic)
	x.SetInitialized()
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
	err := t.Init(name, setting)
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
