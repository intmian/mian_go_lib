package misc

import (
	"context"
	"time"

	"gorm.io/gorm/logger"
)

type HookLogger struct {
	logger.Interface
	Hook SQLTraceHook
}

func (l *HookLogger) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (string, int64),
	err error,
) {
	// 保留原 logger 行为
	l.Interface.Trace(ctx, begin, fc, err)

	if l.Hook == nil {
		return
	}

	sql, rows := fc()
	duration := time.Since(begin)

	l.Hook(ctx, sql, rows, duration, err)
}

type SQLTraceHook func(
	ctx context.Context,
	sql string,
	rows int64,
	duration time.Duration,
	err error,
)
