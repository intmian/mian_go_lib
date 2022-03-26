package blog_list

import (
	"container/heap"
	"mio-blog/modules/blog"
)

type IList interface {
	GetAll() []blog.TBlog
	GetBlogByIndex(begin uint32, end uint32) []blog.TBlog
	Find(ID uint32) (bool, blog.TBlog)
	Add(tBlog blog.TBlog)
	ReMove(id uint32)
}

//PriList 是一个优先队列
type PriList struct {
	core *priListCore
	init bool
}

func (p *PriList) Init() {
	p.core = &priListCore{}
	heap.Init(p.core)
	p.init = true
}

func (p *PriList) GetAll() []blog.TBlog {
	return *p.core
}

func (p *PriList) GetBlogByIndex(begin uint32, end uint32) []blog.TBlog {
	return (*p.core)[begin:end]
}

func (p *PriList) Find(ID uint32) (bool, blog.TBlog) {
	for _, v := range *p.core {
		if v.ID == ID {
			return true, v
		}
	}
	return false, blog.TBlog{}
}

func (p *PriList) Add(tBlog blog.TBlog) {
	heap.Push(p.core, tBlog)
}

func (p *PriList) ReMove(id uint32) {
	index := -1
	for i, v := range *p.core {
		if v.ID == id {
			index = i
			break
		}
	}
	heap.Remove(p.core, index)
}
