package misc

type FuncQueueCaller struct {
	queue  chan func()
	init   bool
	max    chan bool
	exit   chan bool
	exited chan bool
}

func (q *FuncQueueCaller) Init(num int, paraNum int) {
	q.queue = make(chan func(), num)
	q.max = make(chan bool, paraNum)
	q.exit = make(chan bool, 0)
	q.exited = make(chan bool, 0)
	go func() {
		// 按照顺序从队列中取出paraNum个函数并行运行，直到退出
		for {
			select {
			case <-q.exit:
				q.exited <- true
				return
			case f := <-q.queue:
				q.max <- true
				// 不一定能严格保证按照顺序执行，因为同时打开paraNum个协程，但是可以保证前面的比后面的优先被协程发射。同时发射的paraNum个不能保证先后，如果不用队列则只能保证同时只能执行paraNum个，所有协程一起竞争，
				// 同时发射只有当初始化时或者有协程同时返回时才可能触发
				go func() {
					f() // 暂不考虑函数炸掉、死等导致的泄漏
					<-q.max
				}()
			}
		}

	}()
	q.init = true
}

func (q *FuncQueueCaller) PushFunc(f func()) bool {
	if !q.init {
		return false
	}
	q.queue <- f
	return true
}

func (q *FuncQueueCaller) Exit() bool {
	if !q.init {
		return false
	}
	q.exit <- true
	<-q.exited
	return true
}
