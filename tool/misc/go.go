package misc

import (
	"sync/atomic"
	"time"
)

type LimitMCoCallFuncMgrSetting struct {
	TimeInterval         int32
	EveryIntervalCallNum int32
}

type LimitMcoCallFuncMgr struct {
	funcC chan func()
	end   chan bool
}

func (m *LimitMcoCallFuncMgr) Init(setting LimitMCoCallFuncMgrSetting) {
	m.funcC = make(chan func())
	m.end = make(chan bool)
	go m.run(setting)
}

func (m *LimitMcoCallFuncMgr) run(setting LimitMCoCallFuncMgrSetting) {
	NowCallNum := int32(0)
	fumC2 := make(chan func())
	go func() {
		cancallnum := setting.EveryIntervalCallNum - atomic.LoadInt32(&NowCallNum)
		nextTime := time.After(time.Duration(setting.TimeInterval) * time.Second)
		for {
			select {
			case <-m.end:
				return
			case <-nextTime:
				nextTime = time.After(time.Duration(setting.TimeInterval) * time.Second)
				cancallnum = setting.EveryIntervalCallNum - atomic.LoadInt32(&NowCallNum)
			default:
				break
			}
			if cancallnum > 0 {
				fumC2 <- <-m.funcC
				cancallnum--
			} else {
				<-nextTime
				nextTime = time.After(time.Duration(setting.TimeInterval) * time.Second)
				cancallnum = setting.EveryIntervalCallNum - atomic.LoadInt32(&NowCallNum)
			}
		}
	}()
	for {
		select {
		case <-m.end:
			return
		case f := <-fumC2:
			go func() {
				atomic.AddInt32(&NowCallNum, 1)
				f()
				atomic.AddInt32(&NowCallNum, -1)
			}()
		}
	}
}

func (m *LimitMcoCallFuncMgr) Call(f func()) {
	go func() {
		m.funcC <- f
	}()
}

func (m *LimitMcoCallFuncMgr) Exit() {
	m.end <- true
}

var gMgr LimitMcoCallFuncMgr

func GetMCoCallDefault() *LimitMcoCallFuncMgr {
	return &gMgr
}
