package XBi

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
)

// XBi 是轻量业务日志结构体，封装了 SLS 客户端和项目、日志库信息
type XBi struct {
	client  *sls.Client
	project string // SLS 项目名
}

// NewXBi 初始化 SLS 客户端，传入 endpoint、accessKeyID、accessKeySecret、project、logstore
func NewXBi(endpoint, accessKeyID, accessKeySecret, project, logstore string) *XBi {
	client := sls.CreateNormalInterface(endpoint, accessKeyID, accessKeySecret)
	return &XBi{
		client:   client,
		project:  project,
		logstore: logstore,
	}
}

// LogEntry 业务日志实体，支持任意结构，实际写入会被序列化成 JSON
type LogEntry struct {
	Biz       string                 `json:"biz"`       // 业务名
	Timestamp int64                  `json:"timestamp"` // 时间戳，单位秒
	Data      map[string]interface{} `json:"data"`      // 业务数据，任意键值对
}

// WriteLog 发送日志到 SLS，参数是业务名和具体业务数据（map），自动填充时间戳
func (x *XBi) WriteLog(ctx context.Context, biz string, data map[string]interface{}) error {
	// 构造日志内容
	logEntry := LogEntry{
		Biz:       biz,
		Timestamp: time.Now().Unix(),
		Data:      data,
	}

	// 转为 JSON 字符串，作为单条日志内容
	contentBytes, err := json.Marshal(logEntry)
	if err != nil {
		return fmt.Errorf("序列化日志内容失败: %w", err)
	}

	// 构造 SLS Log 对象，SLS 的 Log 包含多个 KV，所以这里用单个 KV 保存整个 JSON 字符串
	slsLog := sls.Log{
		Contents: []*sls.LogContent{
			{
				Key:   "json",
				Value: string(contentBytes),
			},
		},
		Time: uint32(logEntry.Timestamp),
	}

	// 构造请求批量日志，SLS 支持批量写入，这里只写一条
	logGroup := &sls.LogGroup{
		Logs:   []*sls.Log{&slsLog},
		Topic:  biz,      // Topic 设置成业务名，方便 SLS 端按业务过滤
		Source: "XBi-go", // 日志来源，可自定义
	}

	// 发送日志，异步可自行封装，这里同步发送
	err = x.client.PutLogs(x.project, x.logstore, logGroup)
	if err != nil {
		return fmt.Errorf("发送日志失败: %w", err)
	}

	return nil
}
