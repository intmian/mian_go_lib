package pushmod

import "github.com/intmian/mian_go_lib/tool/misc"

// 前面的err是用来区分类型的，所以看上去有点扭曲
const (
	ErrTypeErr          = misc.ErrStr("type err")
	ErrPushDeerPushFail = misc.ErrStr("pushdeer push fail")
)
