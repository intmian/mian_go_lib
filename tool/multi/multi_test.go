package multi

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSafeMap(t *testing.T) {
	type 花里胡哨的类型 struct {
		花里胡哨的字段 int
	}
	m := SafeMap[string, 花里胡哨的类型]{}
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			m.Store(fmt.Sprintf("key%d", i), 花里胡哨的类型{花里胡哨的字段: i})
			wg.Done()
		}(i)
	}
	wg.Wait()
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			v, ok := m.Load(fmt.Sprintf("key%d", i))
			if !ok {
				t.Error("load error")
			}
			if v.花里胡哨的字段 != i {
				t.Error("load error")
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestUnitLock(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		u := NewUnitLock[string]()
		wg := sync.WaitGroup{}
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(i int) {
				u.Lock(fmt.Sprintf("key"))
				defer u.Unlock(fmt.Sprintf("key"))
				wg.Done()
			}(i)
		}
		c := make(chan any)
		go func() {
			wg.Wait()
			c <- nil
		}()
		select {
		case <-c:
		case <-time.After(time.Second):
			t.Error("lock error")
		}
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(i int) {
				u.SafeRun(fmt.Sprintf("key"), func() {
					wg.Done()
				})
			}(i)
		}
		go func() {
			wg.Wait()
			c <- nil
		}()
		select {
		case <-c:
		case <-time.After(time.Second):
			t.Error("lock error")
		}
	})
	t.Run("lock all", func(t *testing.T) {
		n := time.Now()
		u := NewUnitLock[string]()
		wg := sync.WaitGroup{}
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go u.SafeRun(fmt.Sprintf("key %d", i), func() {
				wg.Done()
			})
		}
		for i := 0; i < 100; i++ {
			wg.Add(1)
			i := i
			go func() {
				time.Sleep(time.Second)
				u.SafeRun(fmt.Sprintf("key %d", i), func() {
					if time.Since(n) < time.Second*10 {
						t.Error("lock error")
					}
					wg.Done()
				})
			}()
		}
		wg.Add(1)
		go func() {
			u.SafeRunAll(func() {
				time.Sleep(time.Second * 10)
				wg.Done()
			})
		}()
		timeout := time.After(time.Second * 12)
		c := make(chan any)
		go func() {
			wg.Wait()
			c <- nil
		}()
		select {
		case <-c:
		case <-timeout:
			t.Error("lock error")
		}
	})
}
