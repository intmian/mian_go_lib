package xstorage

type ErrStr string

const(
    ErrNil = ErrStr("nil")
    ErrValueTypeNotMatch = ErrStr("value type not match")  // auto generated from .\sqlite.go
    ErrCoreIsNil = ErrStr("core is nil")  // auto generated from .\cfg_pack.go
    ErrNotInitialized = ErrStr("not initialized")  // auto generated from .\cfg_pack.go
    ErrKeyNotFound = ErrStr("key not found")  // auto generated from .\cfg_pack.go
    ErrParamIsNil = ErrStr("param is nil")  // auto generated from .\cfg_pack.go
    ErrParamIsInvalid = ErrStr("param is invalid")  // auto generated from .\cfg_pack.go
    ErrKeyAlreadyExist = ErrStr("key already exist")  // auto generated from .\cfg_pack.go
    ErrParamIsEmpty = ErrStr("param is empty")  // auto generated from .\cfg_pack.go
)

func (e ErrStr) Error() string { return string(e) }
