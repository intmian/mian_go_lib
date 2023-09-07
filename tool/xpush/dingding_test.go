package xpush

import (
	"os"
	"strings"
	"testing"
)

func TestDingRobotMgr_Send(t *testing.T) {
	m := &DingRobotMgr{}
	token := ""
	secret := ""

	// 从本地文件 dingding_test.txt 读取测试内容token和secret
	file, _ := os.Open("dingding_test.txt")
	defer file.Close()
	buf := make([]byte, 1024)
	n, _ := file.Read(buf)
	str := string(buf[:n])
	strs := strings.Split(str, "\r\n")
	token = strs[0]
	secret = strs[1]

	m.Init(DingSetting{
		Token:             token,
		Secret:            secret,
		SendInterval:      60,
		IntervalSendCount: 20,
	})

	text := NewDingText()
	text.Text.Content = "test"
	err := m.Send(text)
	if err != nil {
		t.Error(err)
	}
}
