package menu

// MenuNode 菜单节点
type MenuNode interface {
	MenuFuncNode
	MenuLogicNode
}

// MenuFuncNode 菜单功能节点
type MenuFuncNode interface {
	isCallAble() bool
	do()
	stop() <-chan bool // 返回一个通道，用于显示什么时候停止
}

// MenuLogicNode 菜单逻辑节点
type MenuLogicNode interface {
	getParent() *MenuNode
	getRoot() *MenuNode
	getAllChild() []*MenuNode
}

type UnCallableMenuLogicNode struct {
}

func (receiver UnCallableMenuLogicNode) stop() <-chan bool {
	return nil
}

func (receiver UnCallableMenuLogicNode) isCallAble() bool {
	return false
}

func (receiver UnCallableMenuLogicNode) do() {
	return
}

type inputModel interface {
	input() string
	inputWithLen(strLen int) string
	outInput(string) string
}
