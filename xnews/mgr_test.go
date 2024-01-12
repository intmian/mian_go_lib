package xnews

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"
)

func TestLimit(t *testing.T) {
	newsMgr1, err := NewXNews(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	var setting1 TopicSetting
	setting1.AddNowLimit(time.Second*2, 20)
	setting1.AddLimit(time.Now().Add(-1*time.Second), time.Second*2, 10)
	err = newsMgr1.AddTopic("topic1", setting1)
	if err != nil {
		t.Fatal(err)
	}
	msgIndex := 0
	genMsg := func() string {
		msgIndex++
		return strconv.Itoa(msgIndex)
	}
	// 第0秒
	for i := 0; i < 11; i++ {
		err = newsMgr1.AddMessage("topic1", genMsg())
		if err != nil {
			t.Fatal(err)
		}
	}
	msgs, err := newsMgr1.GetTopic("topic1")
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 10 {
		t.Fatal("len(msgs) != 10")
	}
	if msgs[0] != "2" {
		t.Fatal("msgs[0] != 1")
	}
	msgs, err = newsMgr1.GetTopic("topic2")
	if !(errors.Is(err, ErrTopicNotExist)) {
		t.Fatal("! errors.is(err, ErrTopicNotExist)")
	}
	if msgs != nil {
		t.Fatal("msgs != nil")
	}
	// 第1秒
	time.Sleep(time.Second + time.Millisecond*100)
	for i := 0; i < 11; i++ {
		err = newsMgr1.AddMessage("topic1", genMsg())
		if err != nil {
			t.Fatal(err)
		}
	}
	msgs, err = newsMgr1.GetTopic("topic1")
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 19 {
		t.Fatal("len(msgs) != 17")
	}
	if msgs[0] != "4" {
		t.Fatal("msgs[0] != 2")
	}

	//var setting2 TopicSetting
	//setting2.AddForeverLimit(100)
	//err = newsMgr1.AddTopic("topic2", setting2)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//for i := 0; i < 101; i++ {
	//	err = newsMgr1.AddMessage("topic2", "msg2")
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//}
	//msgs, err = newsMgr1.GetTopic("topic2")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//if len(msgs) != 100 {
	//	t.Fatal("len(msgs) != 100")
	//}
}

func TestClear(t *testing.T) {
	newsMgr1, err := NewXNews(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	var setting1 TopicSetting
	setting1.AddNowClear(time.Second * 2)
	setting1.AddClear(time.Now().Add(-1*time.Second), time.Second*2)
	err = newsMgr1.AddTopic("topic1", setting1)

	// 第0秒 10条
	for i := 0; i < 10; i++ {
		err = newsMgr1.AddMessage("topic1", "msg1")
		if err != nil {
			t.Fatal(err)
		}
	}
	msgs, err := newsMgr1.GetTopic("topic1")
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 10 {
		t.Fatal("len(msgs) != 10")
	}

	// 第1秒 10条
	time.Sleep(time.Second + time.Millisecond*100)
	msgs, err = newsMgr1.GetTopic("topic1")
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 0 {
		t.Fatal("len(msgs) != 0")
	}
	for i := 0; i < 10; i++ {
		err = newsMgr1.AddMessage("topic1", "msg1")
		if err != nil {
			t.Fatal(err)
		}
	}

	// 第2秒
	time.Sleep(time.Second + time.Millisecond*100)
	msgs, err = newsMgr1.GetTopic("topic1")
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 0 {
		t.Fatal("len(msgs) != 0")
	}

}
