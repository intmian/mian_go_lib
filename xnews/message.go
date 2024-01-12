package xnews

import "time"

type Message struct {
	CreateTime time.Time
	Content    string
}

func (m *Message) Reset() {
	*m = Message{}
}
