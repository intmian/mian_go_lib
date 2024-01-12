package misc

const (
	ErrNil = ErrStr("nil")
)

type ErrStr string

func (e ErrStr) Error() string { return string(e) }
