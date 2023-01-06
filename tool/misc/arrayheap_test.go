package misc

import "testing"

type s struct {
	Comparable
	a int
}

func (receiver s) Less(o Comparable) bool {
	return receiver.a < o.(s).a
}

func TestArrayHeap(t *testing.T) {
	var h ArrayHeap
	Init(&h)
	Push(&h, s{1})

}
