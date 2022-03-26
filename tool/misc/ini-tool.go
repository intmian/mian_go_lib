package misc

type InitTag struct {
	Init bool
}

func (receiver *InitTag) IsInit() bool {
	return receiver.Init
}

func (receiver *InitTag) Inited() {
	receiver.Init = true
}
