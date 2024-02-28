package xstorage

type ErrStr string

const(
    ErrNil = ErrStr("nil")
    ErrValueTypeNotMatch = ErrStr("value type not match")  // auto generated from .\sqlite.go
)

func (e ErrStr) Error() string { return string(e) }
