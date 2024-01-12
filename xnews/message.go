package xnews

import "time"

type Message struct {
	CreateTime time.Time
	Content    string
	ExpireTime time.Time
}

func (m *Message) Reset() {
	*m = Message{}
}
