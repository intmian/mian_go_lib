package menu

type MenuFunc interface {
	do()
	stop() <-chan bool  // 返回一个通道，用于显示什么时候停止
}

type MenuLogicNode interface {
	MenuFunc
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

