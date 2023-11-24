package misc

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"net"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

type ISyncProtoSend interface {
	proto.Message
	GetRecallID(uint32) // 为了反射修改RecallID，其proto文件中必须有RecallID字段，则必有GetRecallID方法
}
type ISyncProtoRec interface {
	proto.Message
	GetRecallID() uint32
}

type Sync struct {
	syncChanMap sync.Map
	syncChan    map[uint32]chan ISyncProtoRec
	ID          atomic.Uint32
}

func (s *Sync) Wait(send func(cConn net.Conn, areaID uint32, data ISyncProtoSend) error, cConn net.Conn, areaID uint32, data ISyncProtoSend, timeout time.Duration) (ISyncProtoRec, error) {
	// 生成一个id
	id := s.ID.Add(1)
	// 使用反射修改data的RecallID
	v := reflect.ValueOf(data)
	v.Elem().FieldByName("RecallID").SetUint(uint64(id))
	// 创建一个chan
	ch := make(chan ISyncProtoRec)
	s.syncChanMap.Store(id, ch)
	defer func() {
		s.syncChanMap.Delete(id)
		close(ch)
	}()
	err := send(cConn, areaID, data)
	if err != nil {
		return nil, err
	}
	// 等待
	select {
	case <-time.After(timeout):
		return nil, errors.New("timeout")
	case r := <-ch:
		return r, nil
	}
}

func (s *Sync) OnRecResult(data ISyncProtoRec) error {
	id := data.GetRecallID()
	v, ok := s.syncChanMap.Load(id)
	if !ok {
		return errors.New("not found")
	}
	ch, ok := v.(chan ISyncProtoRec)
	if !ok {
		return errors.New("not found")
	}
	ch <- data
	return nil
}
