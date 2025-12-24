package xbi

import (
	"time"
)

func toRealEntity[T any](entity LogEntity[T]) RealLogEntity[T] {
	realEntity := RealLogEntity[T]{}
	realEntity.RecordTime = time.Now()
	realEntity.Data = entity.GetWriteableData()

	return realEntity
}
