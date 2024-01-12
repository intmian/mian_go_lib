package xnews

import (
	"time"
)

// TopicTimeRemain 用于限制某个topic的单条信息的留存时间
type TopicTimeRemain time.Duration

// TopicTimeLimit 用于限制某个topic的时间间隔内做多保留的信息数量。
// 如果不填duration，则为永久上限。
// 可以将多条组合使用，但是请注意，任意一条如果达到上限，则都会删除最旧的信息。
type TopicTimeLimit struct {
	Duration        *time.Duration // 不填的话为永久上限
	Num             int
	LastResetTime   time.Time
	ThisDurationNum int
}

// Add 向limit当前周期内增加n个计数，返回淘汰的数量
func (l *TopicTimeLimit) Add(num int) int {
	if num < 0 {
		return 0
	}
	l.checkDuration()

	l.ThisDurationNum += num
	if l.ThisDurationNum > l.Num {
		outNum := l.ThisDurationNum - l.Num
		l.ThisDurationNum = l.Num
		return outNum
	}

	return 0
}

func (l *TopicTimeLimit) checkDuration() {
	if l.Duration == nil {
		return
	}
	if time.Now().Sub(l.LastResetTime) < *l.Duration {
		return
	}
	l.LastResetTime = time.Now()
	l.ThisDurationNum = 0
	return
}

// clearTime 用于在某个时间间隔内清理某个topic的信息
type TopicClearTime struct {
	lastTime time.Time
	duration time.Duration
}
