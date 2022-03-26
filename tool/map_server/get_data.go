package map_server

import (
	"mio-blog/libs/maps"
	"mio-blog/libs/tool/cmd_server"
)

// 命令头 EGetData

type tGetData struct {
	key1   interface{} // map的key
	key2   interface{} // 嵌套map的内层的key
	target maps.Target
}

func (receiver *TMapServer) SendGetData(key1 interface{}, key2 interface{}, target maps.Target, c chan cmd_server.AnyData) {
	receiver.c.Send(cmd_server.Plt{
		Cmd: EGetData,
		Data: tGetData{
			key1:   key1,
			key2:   key2,
			target: target,
		},
		Ret: c,
	})
}

func (receiver *TMapServer) procGetData(v cmd_server.AnyData) cmd_server.AnyData {
	d := v.(tGetData)
	return receiver.m.GetData(d.key1, d.key2, d.target)
}
