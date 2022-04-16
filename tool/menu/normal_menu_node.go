package menu

type normalMenuLogicNode struct {
	root   *MenuNode
	parent *MenuNode
	child  []*MenuNode
}

func (n normalMenuLogicNode) getParent() *MenuNode {
	return n.parent
}

func (n normalMenuLogicNode) getRoot() *MenuNode {
	return n.root
}

func (n normalMenuLogicNode) getAllChild() []*MenuNode {
	return n.child
}
