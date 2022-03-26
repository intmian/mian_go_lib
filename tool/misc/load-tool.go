package misc

type LoadTag struct {
	load bool
}

func (t *LoadTag) IsLoad() bool {
	return t.load
}

func (t *LoadTag) Loaded() {
	t.load = true
}
