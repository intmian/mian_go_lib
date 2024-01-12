package xnews

type ErrStr string

const (
	ErrNil               = ErrStr("nil")
	ErrTopicAlreadyExist = ErrStr("topic already exist") // auto generated from .\mgr.go
	ErrInit              = ErrStr("init")                // auto generated from .\mgr.go
	ErrTopicNotExist     = ErrStr("topic not exist")     // auto generated from .\mgr.go
	ErrGetTopicFailed    = ErrStr("get topic failed")    // auto generated from .\mgr.go
	ErrAddMessageFailed  = ErrStr("add message failed")  // auto generated from .\mgr.go
)

func (e ErrStr) Error() string { return string(e) }
