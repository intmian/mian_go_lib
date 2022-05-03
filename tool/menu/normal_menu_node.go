package menu

type NormalMenuNode struct {
	NormalMenuFuncNode
	NormalMenuLogicNode
}

type NormalMenuLogicNode struct {
	root   MenuNode
	parent MenuNode
	child  []MenuNode
	name   string
	id     int
}

func (n *NormalMenuLogicNode) GetID() int {
	return n.id
}

func (n *NormalMenuLogicNode) SetID(ID int) {
	n.id = ID
}

func (n *NormalMenuLogicNode) GetName() string {
	return n.name
}

func (n *NormalMenuLogicNode) BindParent(parent MenuNode) {
	n.parent = parent
}

func (n *NormalMenuLogicNode) BindChild(child MenuNode) {
	if n.child == nil {
		n.child = make([]MenuNode, 0)
	}
	for _, v := range n.child {
		if (v).GetName() == (child).GetName() {
			return
		}
	}
}

func (n *NormalMenuLogicNode) BindRoot(root MenuNode) {
	n.root = root
}

func (n *NormalMenuLogicNode) SetName(name string) {
	n.name = name
}

func (n *NormalMenuLogicNode) GetParent() MenuNode {
	return n.parent
}

func (n *NormalMenuLogicNode) GetRoot() MenuNode {
	return n.root
}

func (n *NormalMenuLogicNode) GetAllChild() []MenuNode {
	return n.child
}

type NormalMenuFuncNode struct {
	f func()
}

func (n *NormalMenuFuncNode) IsCallAble() bool {
	if n.f == nil {
		return false
	}
	return true
}

func (n *NormalMenuFuncNode) Do() {
	n.f()
}

func (n *NormalMenuFuncNode) BindDo(f func()) bool {
	if n.f == nil {
		n.f = f
		return true
	}
	return false
}
