package menu

type normalMenuLogicNode struct {
	root   *MenuNode
	parent *MenuNode
	child  []*MenuNode
}

func (n normalMenuLogicNode) GetParent() *MenuNode {
	return n.parent
}

func (n normalMenuLogicNode) GetRoot() *MenuNode {
	return n.root
}

func (n normalMenuLogicNode) GetAllChild() []*MenuNode {
	return n.child
}
