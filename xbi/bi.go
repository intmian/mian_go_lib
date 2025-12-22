package xbi

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
)

// XBi 是轻量业务日志结构体，封装了 SLS 客户端和项目、日志库信息
type XBi struct {
	client  sls.ClientInterface
	project PpjName
	DbName  DbName
	table   TableName
}

// NewXBi 初始化 SLS 客户端，传入 Endpoint、AccessKeyID、AccessKeySecret、project、logstore
func NewXBi(setting Setting) *XBi {
	client := sls.CreateNormalInterface(setting.Endpoint, setting.AccessKeyID, setting.AccessKeySecret, setting.SecurityToken)
	return &XBi{
		client:  client,
		project: setting.PpjName,
		DbName:  setting.DbName,
		table:   setting.TableName,
	}
}

// LogEntry 业务日志实体，支持任意结构，实际写入会被序列化成 JSON
type LogEntry struct {
	Timestamp int64       `json:"timestamp"` // 时间戳，单位秒
	Data      interface{} `json:"data"`      // 业务数据，任意键值对
}

// WriteLog 发送日志到 SLS，参数是业务名和具体业务数据（map），自动填充时间戳
func (x *XBi) WriteLog(data interface{}) error {
	// 构造日志内容
	logEntry := LogEntry{
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
