package xpush

type PushType int8

const (
	PushTypeNull PushType = iota
	PushTypeEmail
	PushTypePushDeer
	PushTypeDing
	PushTypeMax
)
