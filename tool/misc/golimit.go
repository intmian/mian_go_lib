package misc

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// GoLimitSetting 并发执行限流器设置
type GoLimitSetting struct {
	TimeInterval         time.Duration // 每隔多少秒重置一次
	EveryIntervalCallNum int32         // 每个周期允许执行多少次
}

// GoLimit 并发执行限流器
// 如果超过限制上线，会阻塞
// 如果一个函数跨越两个周期，会阻塞
// init 后调用 start使用
// 停止通过ctx实现或直接调用
// 不支持stop后重新start
type GoLimit struct {
	setting     GoLimitSetting
	ctx         context.Context
	funcC       chan func()  // 待执行函数
	nowCallNum  atomic.Int32 // 当前周期正在执行的函数数量
	calledNum   atomic.Int32 // 当前周期已经执行的函数数量
	fullRunLock sync.Mutex
	fullRun     atomic.Bool
	stop        func()
	started     bool
	InitTag
}

const (
	ErrCtxNil       = ErrStr("ctx is nil")
	ErrParamInvalid = ErrStr("param invalid")
)

func (m *GoLimit) Init(setting GoLimitSetting, ctx context.Context) error {
	if m.IsInitialized() {
		return ErrRepeatInit
	}
	m.funcC = make(chan func())
	if ctx == nil {
		return ErrCtxNil
	}
	if setting.EveryIntervalCallNum <= 0 || setting.TimeInterval.Seconds() <= 0 {
		return ErrParamInvalid
	}
	m.ctx, m.stop = context.WithCancel(ctx)
	m.setting = setting
	m.SetInitialized()
	return nil
}

const (
	ErrStop     = ErrStr("stop")
	ErrNotStart = ErrStr("not start")
)

func (m *GoLimit) Call(f func()) error {
	if !m.started {
		return ErrNotStart
	}
	if !m.IsInitialized() {
		return ErrNotInit
	}
	select {
	case <-m.ctx.Done():
		return ErrStop
	default:
		m.funcC <- f
	}
	return nil
}

//TODO: 涉及免闭包实现
//func Call[T func](m *GoLimit, f T) {
//	m.Call(func() {
//		f()
//	})
//}

func (m *GoLimit) Start() error {
	m.fullRun.Store(true)
	m.fullRunLock.Lock()
	go m.timeLimit()
	go m.callLimit()
	m.started = true
	return nil
}

func (m *GoLimit) Stop() error {
	m.stop()
	return nil
}

func (m *GoLimit) timeLimit() {
	// 用timer的唯一原因是在写单元测试时，能够应该控制第一个周期的执行时间完整
	t := time.NewTimer(m.setting.TimeInterval)
	for {
		select {
		case <-m.ctx.Done():
			return
		case <-t.C:
			//fmt.Printf("%stime:nowCallNum:%d calledNum:%d\n", time.Now().Format("15:04:05"), m.nowCallNum.Load(), m.calledNum.Load())
			m.calledNum.Store(0)
			if m.fullRun.Load() {
				m.fullRun.Store(false)
				m.fullRunLock.Unlock()
			}
			t.Reset(m.setting.TimeInterval)
		}
	}
}

func (m *GoLimit) callLimit() {
	for {
		select {
		case <-m.ctx.Done():
			return
		case f := <-m.funcC:
			// 做限制
			//fmt.Printf("%scheck:nowCallNum:%d calledNum:%d\n", time.Now().Format("15:04:05"), m.nowCallNum.Load(), m.calledNum.Load())
			if m.nowCallNum.Load()+m.calledNum.Load() >= m.setting.EveryIntervalCallNum {
				m.fullRun.Store(true)
				m.fullRunLock.Lock()
			}
			// 实际执行
			//fmt.Printf("%scall:nowCallNum:%d calledNum:%d\n", time.Now().Format("15:04:05"), m.nowCallNum.Load(), m.calledNum.Load())
			m.nowCallNum.Add(1)
			go func() {
				f()
				defer m.nowCallNum.Add(-1)
				m.calledNum.Add(1)
			}()
		}
	}
}
