package menu

// MenuNode 菜单节点
type MenuNode interface {
	MenuFuncNode
	MenuLogicNode
}

// MenuFuncNode 菜单功能节点
type MenuFuncNode interface {
	IsCallAble() bool
	Do(end chan<- bool) // 阻塞执行、end用来强制结束
}

// MenuLogicNode 菜单逻辑节点
type MenuLogicNode interface {
	GetParent() *MenuNode
	GetRoot() *MenuNode
	GetAllChild() []*MenuNode
	GetName() string

	BindParent(parent *MenuNode)
	BindChild(child *MenuNode)
	BindRoot(root *MenuNode)
	SetName(name string)
}

type UnCallableMenuLogicNode struct {
}

func (receiver UnCallableMenuLogicNode) IsCallAble() bool {
	return false
}

func (receiver UnCallableMenuLogicNode) Do() {
	return
}

type inputModel interface {
	input() string
	inputWithLen(strLen int) string
	outInput(string) string
}
