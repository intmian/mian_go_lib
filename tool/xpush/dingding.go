package xpush

type DingRobotToken struct {
	accessToken string
	secret      string
}

type DingRobotMgr struct {
	dingRobotToken DingRobotToken
	message        chan string
}

type DingMessage interface {
	ToJson() string
}

const ApiUrl = "https://oapi.dingtalk.com/robot/send"
