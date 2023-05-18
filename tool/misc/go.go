package misc

import (
	"sync/atomic"
	"time"
)

func LimitMCoCallFunc(f func())

type LimitMCoCallFuncMgrSetting struct {
	timeInterval         int32
	EveryIntervalCallNum int32
}

type LimitMcoCallFuncMgr struct {
	funcC chan func()
}

func (m *LimitMcoCallFuncMgr) Init(setting LimitMCoCallFuncMgrSetting, end <-chan bool) {
	go func() {
		NowCallNum := int32(0)
		CanCallNum := setting.EveryIntervalCallNum
		go func() {
			for {
				select {
				case <-end:
					return
				case <-time.After(time.Duration(setting.timeInterval) * time.Second):
					atomic.StoreInt32(&NowCallNum, 0)
				}
			}
		}()
		for {
			select {
			case <-end:
				return
			case f := <-m.funcC:
				if atomic.LoadInt32(&CanCallNum) > 0 {
					atomic.AddInt32(&NowCallNum, 1)
					ret := chan bool
					go func() {
						f()
						ret <- true
					}()

				}
			}
		}
	}()
}
