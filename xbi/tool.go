package xbi

import (
	"time"
)

func toDbData[DataType any](entity LogEntity[DataType]) DbLogData[DataType] {
	realEntity := DbLogData[DataType]{}
	realEntity.RecordTime = time.Now().Unix()
	if data := entity.GetWriteableData(); data != nil {
		realEntity.Data = *data
	}

	return realEntity
}
