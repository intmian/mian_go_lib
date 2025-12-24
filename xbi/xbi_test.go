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

func (t *testLogEntity) TableName() TableName {
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

	err = db.AutoMigrate(toRealEntity[*testLog](&testLogEntity{}))
	if err != nil {
		t.Fatal("创建数据表失败:", err)
	}

	xbi := NewXBi(Setting{
		Db: db,
	})

	if xbi == nil {
		t.Fatal("创建XBi失败")
	}

	testLog1 := &testLogEntity{}
	testLog1Data := testLog1.GetWriteableData()
	testLog1Data.A = "testA"
	testLog1Data.B = 123

	ctx := context.Background()
	errChan := make(chan error)

	err = WriteLog[*testLog](xbi, testLog1, ctx, errChan)
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

	logEntitie := testLogEntity{}
	data, err := ReadLog[testLog](xbi, string(logEntitie.TableName()), nil, 1, 0, ctx)

	if err != nil {
		t.Fatal("读取日志失败:", err)
	}

	if len(data) == 0 {
		t.Fatal("读取日志为空")
	}
}
