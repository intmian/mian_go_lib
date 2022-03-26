package map_server

import (
	"mio-blog/libs/maps"
	"mio-blog/libs/tool/cmd_server"
)

// 命令头 EUpdateData

type tUpdateData struct {
	key    interface{}
	value1 interface{}
	value2 interface{}
	mapT   maps.TMapType
	target maps.Target
}

// SendUpdateData 将数据进行插入。
//如果被插入的是个集合则在集合中插入。如果集合未被建立则建立。如果被插入的是个meta则调用meta的插入。
//当插入的是 字符串-集合：
//key : 键 value1 : 值 value2 : 为字符串的指令。
//当出入的是 字符串-Meta 本质上是[int][string]string这个map：
//key ：键 value1 ：Meta的键或对此meta的指令 value1 ：Meta的值或对此值的指令。
func (receiver *TMapServer) SendUpdateData(key interface{}, value1 interface{}, value2 interface{}, mapT maps.TMapType, target maps.Target) {
	receiver.c.Send(cmd_server.Plt{
		Cmd: EUpdateData,
		Data: tUpdateData{
			key:    key,
			value1: value1,
			value2: value2,
			mapT:   mapT,
			target: target,
		},
		Ret: nil,
	})
}

func (receiver *TMapServer) procUpdateData(v cmd_server.AnyData) cmd_server.AnyData {
	data := v.(tUpdateData)
	receiver.m.UpdateData(data.key, data.value1, data.value2, data.mapT, data.target)
	return nil
}
