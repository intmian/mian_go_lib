package blog_list

import "mio-blog/modules/blog"

//priListCore 优先队列核心
type priListCore []blog.TBlog

func (p *priListCore) Len() int {
	return len(*p)
}

func (p *priListCore) Less(i, j int) bool {
	// 暂时按照ID划分
	return (*p)[i].ID < (*p)[j].ID
}

func (p *priListCore) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}

func (p *priListCore) Push(x interface{}) {
	*p = append(*p, x.(blog.TBlog))
}

func (p *priListCore) Pop() interface{} {
	old := *p
	n := len(old)
	x := old[n-1]
	*p = old[0 : n-1]
	return x
}

func PriListCoreTest() {
	core := priListCore{}
	print(core.Len())
	core.Push(blog.TBlog{ID: 1})
	core.Push(blog.TBlog{ID: 2})
	core.Push(blog.TBlog{ID: 3})
	print(core.Len())
}
