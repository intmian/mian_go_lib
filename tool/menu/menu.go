package menu

import (
	"encoding/json"
	"fmt"
	"github.com/intmian/mian_go_lib/tool/misc"
	"strconv"
)

const PAGE_NUM = 10

type Menu struct {
	root    MenuNode
	now     MenuNode
	ID2Node map[int]MenuNode
}

func (m *Menu) Do() {
	m.now = m.root
nodeCircle:
	for {
		misc.Clear()
		canRoot := m.now != m.root
		canParent := m.now.GetParent() != nil
		iPage := 1
		for {
			print(getText(m.root, m.now.GetParent(), m.now.GetAllChild(), iPage, canRoot, canParent))
			c := misc.WaitKeyDown()
			switch c {
			case 'r':
				if canRoot {
					m.now = m.root
					continue nodeCircle
				} else {
					continue
				}
			case 'p':
				if canParent {
					m.now = m.now.GetParent()
					continue nodeCircle
				} else {
					continue
				}
			case 'e':
				return
			case '←':
				if iPage > 1 {
					iPage--
				}
			case '→':
				if misc.GetMaxPage(len(m.now.GetAllChild()), PAGE_NUM, true) > iPage {
					iPage++
				}
			}
			if c >= '0' && c <= '9' {
				pageIndex, _ := strconv.ParseInt(string(c), 10, 32)
				index := misc.GetPageIndexOriIndex(int(pageIndex), iPage, PAGE_NUM, true)
				if index >= 0 && index < len(m.now.GetAllChild()) {
					next := m.now.GetAllChild()[index]
					if next.IsCallAble() {
						misc.Clear()
						next.Do()
						if next.GetAllChild() != nil && len(next.GetAllChild()) > 0 {
							m.now = next
						}
						continue nodeCircle
					} else {
						m.now = next
						continue nodeCircle
					}
				}
			}

		}
	}
}

func getText(root MenuNode, parent MenuNode, children []MenuNode, page int, canRoot bool, canParent bool) string {
	head := ""
	content := ""
	foot := ""
	if canRoot {
		head += "r." + root.GetName() + " "
	}
	if canParent {
		head += "p." + parent.GetName() + " "
	}
	head += "e.exit"
	if head != "" {
		head += "\n"
	}
	begin, end := misc.GetPageStartEnd(page, PAGE_NUM, len(children), true)
	for i := begin; i < end; i++ {
		realIndex := misc.GetPageIndexOriIndex(i, PAGE_NUM, len(children), true)
		content += strconv.Itoa(i) + "." + children[realIndex].GetName() + "\n"
	}

	if len(children) > PAGE_NUM {
		perm := ""
		if page > 1 {
			perm += "1..<-"
		}
		max := misc.GetMaxPage(len(children), PAGE_NUM, true)
		perm += "%d"
		if page < max {
			t := "->..%d"
			t = fmt.Sprintf(t, max)
			perm += t
		}
		perm = fmt.Sprintf(perm, page)
		foot += perm
	}
	return ""
}

// BindJson 用来确认菜单结构
/*
{
    "nodes" : [
        {
            "id":0,
            "name" : "0号节点",
            "father" : -1,
            "child" : [1,2,3]
        }
    ],
    "root" : 0
}
*/
type BindJson struct {
	Nodes []struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		Father int    `json:"father"`
		Child  []int  `json:"child"`
	} `json:"nodes"`
	Root int `json:"root"`
}

type FuncBind struct {
	ID   int
	FUNC func()
}

type BindInfo struct {
	/*
		类似于
		nodes:
		节点编号,节点名称,父节点,[子节点1,子节点2,子节点3]
		...
		root:
		节点编号
	*/
	LogicBindText string
	FuncBindList  []FuncBind
}

func (m *Menu) Init(info BindInfo) bool {
	b := BindJson{}
	err := json.Unmarshal([]byte(info.LogicBindText), &b)
	if err != nil {
		return false
	}

	// 添加所有逻辑节点
	for _, v := range b.Nodes {
		n := NormalMenuNode{}
		n.SetID(v.Id)
		n.SetName(v.Name)
		m.ID2Node[v.Id] = &n
	}

	// 添加所有功能节点
	for _, v := range info.FuncBindList {
		m.ID2Node[v.ID].BindDo(v.FUNC)
	}

	// 添加祖先
	m.root = m.ID2Node[b.Root]

	// 添加所有父子关系
	for _, v := range b.Nodes {
		for _, c := range v.Child {
			m.ID2Node[v.Id].BindChild(m.ID2Node[c])
			m.ID2Node[c].BindParent(m.ID2Node[v.Id])
		}
	}

	m.now = m.root
	return true
}
