package xbi

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/intmian/mian_go_lib/fork/d1_gorm_adapter/gormd1"
	"github.com/intmian/mian_go_lib/tool/misc"
	"gorm.io/gorm"
)

type testLog struct {
	A string
	B int
}

type testLogEntity struct {
	data testLog
}

func (t *testLogEntity) TableName() string {
	return "log_test_db"
}

func (t *testLogEntity) GetWriteableData() *testLog {
	return &t.data
}

func TestNewXBi(t *testing.T) {
	account := misc.InputWithFile("account")
	token := misc.InputWithFile("token")
	dbid := misc.InputWithFile("dbid")
	str := "d1://%s:%s@%s"
	str = fmt.Sprintf(str, account, token, dbid)

	db, err := gorm.Open(gormd1.Open(str), &gorm.Config{})

	if err != nil {
		t.Fatal("连接数据库失败:", err)
	}

	ctx := context.Background()
	errChan := make(chan error, 1)
	xbi, err := NewXBi(Setting{
		Db:        db,
		errorChan: errChan,
		ctx:       ctx,
	})
	if err != nil {
		t.Fatal("创建XBi失败:", err)
	}
	err = RegisterLogEntity[testLog](xbi, &testLogEntity{})
	if err != nil {
		t.Fatal("注册日志实体失败:", err)
	}

	if xbi == nil {
		t.Fatal("创建XBi失败")
	}

	testLog1 := &testLogEntity{}
	testLog1Data := testLog1.GetWriteableData()
	testLog1Data.A = "testB"
	testLog1Data.B = 123

	err = WriteLog[testLog](xbi, testLog1)
	if err != nil {
		t.Fatal("写入日志失败:", err)
	}

	select {
	case err = <-errChan:
		if err != nil {
			t.Fatal("写入日志失败:", err)
		}
	case <-time.After(time.Second * 3):
		t.Fatal("写入日志超时")
	default:
	}

	entity := testLogEntity{}
	conditions := map[string]interface{}{
		"A": "testB",
	}
	data, err := ReadLog[testLog](xbi, string(entity.TableName()), conditions, 0, 0)

	if err != nil {
		t.Fatal("读取日志失败:", err)
	}

	if len(data) == 0 {
		t.Fatal("读取日志为空")
	}

	for _, v := range data {
		t.Log("读取日志:", v)
	}

	t.Log("XBi测试通过")
}
