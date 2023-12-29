package misc

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestGoCtx(t *testing.T) {
	var m GoLimit
	var runNum atomic.Int32
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Second*2))
	m.Init(GoLimitSetting{
		TimeInterval:         time.Second,
		EveryIntervalCallNum: 2,
	}, ctx)
	m.Start()
	for i := 0; i < 5; i++ {
		m.Call(func() {
			runNum.Add(1)
		})
	}
	time.Sleep(time.Second * 10)
	if runNum.Load() != 4 {
		t.Error("runNum != 4")
	}
}

func TestGoBase(t *testing.T) {
	var m GoLimit
	var runNum atomic.Int32
	var runTime []time.Duration
	m.Init(GoLimitSetting{
		TimeInterval:         time.Second,
		EveryIntervalCallNum: 2,
	}, context.Background())
	m.Start()
	now := time.Now()
	runTime = make([]time.Duration, 6, 6)
	for i := 0; i < 6; i++ {
		i2 := i
		m.Call(func() {
			if i2 == 2 {
				time.Sleep(time.Second)
			}
			fmt.Printf("%stest: %d\n", time.Now().Format("15:04:05"), i2)
			runNum.Add(1)
			runTime[i2] = time.Now().Sub(now)
		})
	}
	time.Sleep(time.Second * 5)
	if runNum.Load() != 6 {
		t.Error("runNum != 6")
	}
	if runTime[0] > time.Second {
		t.Error("runTime[0] > time.Second*2")
	}
	if runTime[1] > time.Second {
		t.Error("runTime[1] > time.Second*2")
	}
	if runTime[2] < time.Second*2 {
		t.Error("runTime[2] < time.Second*2")
	}
	if runTime[3] < time.Second {
		t.Error("runTime[3] < time.Second*2")
	}
	if runTime[4] < time.Second*2 {
		t.Error("runTime[4] < time.Second*2")
	}
	if runTime[5] < time.Second*3 {
		t.Error("runTime[5] < time.Second*2")
	}
}
