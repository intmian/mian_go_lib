package misc

// Error 是基础Error类，仅包含一个失败原因字符串
type Error struct {
	Reason string
}

func (e Error) Error() string {
	return e.Reason
}

type LosePicUrlError struct {
	Error
	URLs []string
}
