package map_server

import (
	"mio-blog/libs/maps"
	"mio-blog/libs/tool/cmd_server"
)

type TMapServer struct {
	c cmd_server.CMDServer
	m maps.TMaps
}

const (
	EUpdateData cmd_server.CMD = iota
	EGetData
)

func NewTMapServer() *TMapServer {
	d := &TMapServer{}
	d.Init() // 因为需要内部函数作为参数所以多增加一个初始化
	return d
}

func (receiver *TMapServer) Init() {
	// 缓存多给一点，削峰填谷
	receiver.c = cmd_server.MakeCmdServer(
		map[cmd_server.CMD]cmd_server.FUNC{
			EUpdateData: receiver.procUpdateData,
			EGetData:    receiver.procGetData,
		},
		1000,
		1000)
	receiver.m = *maps.NewKMaps()
}

// 下次用几名组合搞继承，别用这套，太麻烦了

func (receiver *TMapServer) Start() {
	receiver.c.Start()
}

func (receiver *TMapServer) Stop() {
	receiver.c.Stop()
}
