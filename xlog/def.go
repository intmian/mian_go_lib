package xlog

type LogLevel uint8

const (
	LogLevelError LogLevel = iota
	LogLevelWarning
	LogLevelInfo
	LogLevelDebug
	LogLevelMisc
	LogLevelBusiness // TODO: 业务日志，业务日志数量很大，需要后续单独处理，有topic，每个topic对应文件或者数据库，不进行常规处理，并可以进行单独的配置
)
