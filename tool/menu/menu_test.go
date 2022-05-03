package menu

import "testing"

func TestMenu(t *testing.T) {
	logicText := `
{
    "nodes" : [
        {
            "id":0,
            "name" : "0号节点",
            "child" : [1,2,3]
        },
		{
			"id":1,
			"name" : "1号节点",
			"child" : [4,5]
		},
		{
			"id":2,
			"name" : "2号节点",
			"child" : [6,7]
		},
		{
			"id":3,
			"name" : "3号节点",
			"child" : []
		},
		{
			"id":4,
			"name" : "4号节点",
			"child" : []
		},
		{
			"id":5,
			"name" : "5号节点",
			"child" : []
		},
		{
			"id":6,
			"name" : "6号节点",
			"child" : []
		},
		{
			"id":7,
			"name" : "7号节点",
			"child" : []
		}
    ],
    "root" : 0
}
`
	m := Menu{}
	// 手动测试
	kv := make(map[string]string)
	kv["a"] = "1"
	kv["b"] = "2"
	kv["c"] = "3"
	f := MakeListInputFunc(kv, func() {
		for k, v := range kv {
			t.Logf("%s:%s", k, v)
		}
	})
	m.Init(BindInfo{
		LogicBindText: logicText,
		FuncBindList: []FuncBind{
			{
				ID:   7,
				FUNC: f,
			},
			{
				ID: 6,
				FUNC: func() {
					println("6号节点")
				},
			},
			{
				ID: 5,
				FUNC: func() {
					println("5号节点")
				},
			},
			{
				ID: 4,
				FUNC: func() {
					println("4号节点")
				},
			},
			{
				ID: 3,
				FUNC: func() {
					println("3号节点")
				},
			},
			{
				ID: 2,
				FUNC: func() {
					println("2号节点")
				},
			},
			{
				ID: 1,
				FUNC: func() {
					println("1号节点")
				},
			},
			{
				ID: 0,
				FUNC: func() {
					println("0号节点")
				},
			},
		},
	})
	m.Do()
}
