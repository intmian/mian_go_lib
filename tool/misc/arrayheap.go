package misc

import "container/heap"

type Comparable interface {
	Less(other interface{}) bool
}

type ArrayHeap []Comparable

func (h *ArrayHeap) Len() int {
	return len(*h)
}

func (h *ArrayHeap) Less(i, j int) bool {
	return (*h)[i].Less((*h)[j])
}

func (h *ArrayHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *ArrayHeap) Push(x interface{}) {
	*h = append(*h, x.(Comparable))
}

func (h *ArrayHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func Top(h *ArrayHeap) Comparable {
	return (*h)[0]
}

func Pop(h *ArrayHeap) Comparable {
	x := Top(h)
	heap.Pop(h)
	return x
}

func Push(h *ArrayHeap, x Comparable) {
	heap.Push(h, x)
}

func Init(h *ArrayHeap) {
	heap.Init(h)
}

func Remove(h *ArrayHeap, i int) Comparable {
	x := (*h)[i]
	heap.Remove(h, i)
	return x
}

func Len(h *ArrayHeap) int {
	return len(*h)
}
