package menu

// MenuNode 菜单节点
type MenuNode interface {
	MenuFuncNode
	MenuLogicNode
}

// MenuFuncNode 菜单功能节点
type MenuFuncNode interface {
	IsCallAble() bool
	Do()
	BindDo(func()) bool
}

// MenuLogicNode 菜单逻辑节点
type MenuLogicNode interface {
	GetParent() MenuNode
	GetRoot() MenuNode
	GetAllChild() []MenuNode
	GetName() string
	GetID() int

	BindParent(parent MenuNode)
	BindChild(child MenuNode)
	BindRoot(root MenuNode)
	SetName(name string)
	SetID(ID int)
}

type UnCallableMenuLogicNode struct {
}

func (receiver UnCallableMenuLogicNode) IsCallAble() bool {
	return false
}

func (receiver UnCallableMenuLogicNode) Do() {
	return
}
