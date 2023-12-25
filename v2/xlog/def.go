package xlog

type TLogLevel uint8

const (
	EError TLogLevel = iota
	EWarning
	ELog
	EDebug
	EMisc
)
