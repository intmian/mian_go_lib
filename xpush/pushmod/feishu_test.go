package pushmod

import (
	"testing"

	"github.com/intmian/mian_go_lib/tool/misc"
)

func TestFeishuRobotMgr_Send(t *testing.T) {
	m := &FeishuRobotMgr{}
	webhookUrl := misc.InputWithFile("f_web_hook")
	err := m.Init(FeishuSetting{
		WebhookUrl:        webhookUrl,
		SendInterval:      60,
		IntervalSendCount: 20,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Test Text Message
	text := NewFeishuText("test content from unit test")
	err = m.Send(text)
	if err != nil {
		t.Fatal(err)
	}

	// Test Markdown Message
	md := NewFeishuCard("test title", "**test markdown**\n- item 1\n- item 2")
	err = m.Send(md)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFeishuRobotMgr_Push(t *testing.T) {
	webhookUrl := misc.InputWithFile("f_web_hook")

	if webhookUrl == "" {
		t.Skip("webhookUrl is empty, skip test")
		return
	}

	m, err := NewFeishuRobotMgr(FeishuSetting{
		WebhookUrl:        webhookUrl,
		SendInterval:      60,
		IntervalSendCount: 20,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = m.Push("test push title2", "test push content")
	if err != nil {
		t.Fatal(err)
	}

	err = m.PushMarkDown("test markdown title2", "**test markdown content**")
	if err != nil {
		t.Fatal(err)
	}
}
