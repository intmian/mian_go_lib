package misc

import "fmt"

const (
	ErrNil = ErrStr("nil")
)

type ErrStr string

func (e ErrStr) Error() string { return string(e) }

func JoinErr(errs ...error) error {
	errNum := 0
	for _, err := range errs {
		if err != nil {
			errNum++
		}
	}
	if errNum == 0 {
		return nil
	}
	var errStr string
	errStr = "errors[%d]: "
	errStr = fmt.Sprintf(errStr, errNum)
	for i, err := range errs {
		if err != nil {
			errStr += fmt.Sprintf("%d: %s\n", i, err.Error())
		}
	}
	return ErrStr(errStr)
}
