package misc

type InitTag struct {
	Init bool
}

func (receiver *InitTag) IsInitialized() bool {
	return receiver.Init
}

func (receiver *InitTag) SetInitialized() {
	receiver.Init = true
}
