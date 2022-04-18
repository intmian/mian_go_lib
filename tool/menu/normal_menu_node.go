package menu

type normalMenuLogicNode struct {
	root   *MenuNode
	parent *MenuNode
	child  []*MenuNode
	name   string
}

func (n normalMenuLogicNode) GetName() string {
	return n.name
}

func (n normalMenuLogicNode) BindParent(parent *MenuNode) {
	n.parent = parent
}

func (n normalMenuLogicNode) BindChild(child *MenuNode) {
	if n.child == nil {
		n.child = make([]*MenuNode, 0)
	}
	for _, v := range n.child {
		if (*v).GetName() == (*child).GetName() {
			return
		}
	}
}

func (n normalMenuLogicNode) BindRoot(root *MenuNode) {
	n.root = root
}

func (n normalMenuLogicNode) SetName(name string) {
	n.name = name
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
