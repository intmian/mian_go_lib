package misc

import (
	"os"
	"os/exec"
	"strconv"
)

type SingleMenu struct {
	Name    string // 项目名
	F       func() // 可以被call的节点
	SubMenu []*SingleMenu
}

type CmdMenu struct {
	root         *SingleMenu
	now          *SingleMenu
	HisList      []*SingleMenu // 方便完成返回到上一步
	HisListIndex int
}

func (m *CmdMenu) Init(root *SingleMenu) {
	m.root = root
	m.now = root
	m.HisList = make([]*SingleMenu, 100)
	m.HisListIndex = 0
}

func (m *CmdMenu) returnToLast() {
	index := m.HisListIndex
	if index == 0 {
		return
	}
	if m.HisList[index-1] == nil {
		return
	}
	m.HisListIndex -= 1
	m.now = m.HisList[m.HisListIndex]
}

func (m *CmdMenu) gotoSub(index int) bool {
	if m.now == nil {
		return false
	}
	if len(m.now.SubMenu)-1 < index {
		return false
	}
	m.HisList[m.HisListIndex] = m.now
	m.HisListIndex++
	m.now = m.now.SubMenu[index]
	return true
}

func (m *CmdMenu) gotoRoot() {
	m.now = m.root
	m.HisListIndex = 0
}

func (m *CmdMenu) clear() {
	cmd := exec.Command("cmd.exe", "/c", "cls")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		print("clear fail")
	}
}

func (m *CmdMenu) do() (exit bool) {
	m.clear()
	lenSub := len(m.now.SubMenu)
	if lenSub == 0 {
		// 跑到叶节点了，就执行这个的逻辑，停一下，再接着跑
		println(m.now.Name)
		m.now.F()
		Stop()
		m.returnToLast()
		return false
	}
	if m.now == nil {
		return
	}
	for i, s := range m.now.SubMenu {
		println(strconv.Itoa(i+1) + ":" + s.Name)
	}
	CanReturn := false
	if m.HisListIndex > 0 {
		CanReturn = true
	}
	nowIndex := lenSub + 1
	homeIndex := -1
	backIndex := -1
	exitIndex := -1

	if m.now != m.root {
		homeIndex = nowIndex
		nowIndex += 1
		println(strconv.Itoa(homeIndex) + ":Home")
	}
	if CanReturn {
		backIndex = nowIndex
		nowIndex += 1
		println(strconv.Itoa(backIndex) + ":Back")
	}
	exitIndex = nowIndex

	println(strconv.Itoa(exitIndex) + ":Exit")
	inputIndex := 0
	err := Input("请输入下一步", 2, &inputIndex)
	if err != nil {
		return false
	}
	switch {
	case inputIndex < 0:
		return false
	case inputIndex < lenSub+1:
		m.gotoSub(inputIndex - 1)
		return false
	case inputIndex == homeIndex:
		m.gotoRoot()
		return false
	case inputIndex == backIndex:
		m.returnToLast()
		return false
	case inputIndex == exitIndex:
		return true
	default:
		return false
	}
}

func (m *CmdMenu) Run() {
	for {
		if m.do() {
			return
		}
	}
}
