package misc

import (
	"net"
	"testing"
	"time"
)

type test1 struct {
	RecallID *uint32
	TestId   uint32
}

func (m *test1) Reset() {
	//TODO implement me
	panic("implement me")
}

func (m *test1) String() string {
	//TODO implement me
	panic("implement me")
}

func (m *test1) ProtoMessage() {
	//TODO implement me
	panic("implement me")
}

func (m *test1) GetRecallID() uint32 {
	if m != nil && m.RecallID != nil {
		return *m.RecallID
	}
	return 0
}

type test2 struct {
	RecallID uint32
}

func (m *test2) Reset() {
	//TODO implement me
	panic("implement me")
}

func (m *test2) String() string {
	//TODO implement me
	panic("implement me")
}

func (m *test2) ProtoMessage() {
	//TODO implement me
	panic("implement me")
}

func (m *test2) GetRecallID() uint32 {
	return m.RecallID
}

type test3 struct {
}

func (t test3) Reset() {
	//TODO implement me
	panic("implement me")
}

func (t test3) String() string {
	//TODO implement me
	panic("implement me")
}

func (t test3) ProtoMessage() {
	//TODO implement me
	panic("implement me")
}

func (t test3) GetRecallID() uint32 {
	//TODO implement me
	panic("implement me")
}

func TestSync(t *testing.T) {
	sync := &Sync{}
	fEmpty := func(cConn net.Conn, areaID uint32, data ISyncProtoSend) error {
		return nil
	}
	_, err := sync.Wait(fEmpty, nil, 0, &test2{}, 1*time.Second)
	if err != nil && err.Error() != "not found RecallID ptr" {
		t.Error(err)
		return
	}
	_, err = sync.Wait(fEmpty, nil, 0, &test3{}, 1*time.Second)
	if err != nil && err.Error() != "not found RecallID ptr" {
		t.Error(err)
		return
	}
	rV := uint32(999)
	var recallID uint32
	s := &test1{TestId: 1}
	go func() {
		r, err := sync.Wait(fEmpty, nil, 0, s, 1*time.Second)
		if err != nil {
			t.Error(err)
			return
		}
		if recallID == 0 {
			t.Error("recallID == 0")
			return
		}
		if s.TestId != 1 {
			t.Error("test1.TestId != 1")
			return
		}
		if r.(*test1).TestId != rV {
			t.Error("test1.TestId != rV")
			return
		}
	}()
	time.Sleep(100 * time.Millisecond)
	recallID = s.GetRecallID()
	err = sync.OnRecResult(&test1{RecallID: &recallID, TestId: rV})
	if err != nil {
		t.Error(err)
		return
	}

	// 测试超时
	s = &test1{TestId: 1}
	go func() {
		_, err := sync.Wait(fEmpty, nil, 0, s, 1*time.Second)
		if err != nil && err.Error() != "timeout" {
			t.Error(err)
			return
		}
	}()
	time.Sleep(100 * time.Millisecond)
	recallID = s.GetRecallID()
	time.Sleep(2 * time.Second)
	err = sync.OnRecResult(&test1{RecallID: &recallID, TestId: rV})
	if err != nil && err.Error() != "not found" {
		t.Error(err)
		return
	}

}
