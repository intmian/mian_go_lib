package menu

// MenuNode 菜单节点
type MenuNode interface {
	MenuFuncNode
	MenuLogicNode
}

// MenuFuncNode 菜单功能节点
type MenuFuncNode interface {
	do()
	stop() <-chan bool // 返回一个通道，用于显示什么时候停止
}

// MenuLogicNode 菜单逻辑节点
type MenuLogicNode interface {
	goRoot()
	goChild(int)
	goParent()
	isCallAble() bool
	getText()
	getParent() MenuLogicNode
	getRoot() MenuLogicNode
	getAllChild() []MenuLogicNode
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
	outInput(string) string
}
