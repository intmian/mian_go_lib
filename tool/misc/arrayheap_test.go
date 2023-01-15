package misc

import "testing"

type s struct {
	a int
}

func (receiver s) Less(o Comparable) bool {
	return receiver.a < o.(s).a
}

func TestArrayHeap(t *testing.T) {
	var h ArrayHeap
	Init(&h)
	Push(&h, s{1})
	Push(&h, s{11})
	Push(&h, s{3})
	//print(Top(&h))
	//print(Len(&h))
	//print(Pop(&h))
	//print(Pop(&h))
	//print(Pop(&h))

	t.Log(Top(&h))
	t.Log(Len(&h))
	t.Log(Pop(&h))
	t.Log(Pop(&h))
	t.Log(Pop(&h))
}
