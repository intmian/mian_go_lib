package consts

// CArr 提供对底层数组的保护，避免向上层返回内部数据被用户修改，所有可能会对原始数组修改的操作被屏蔽，必须复制一份。
// copy 会复制一个浅拷贝
// get、len、range 与原始数组一致
// section 复制一个切片
// 使用方法：
// var arr = []int{1, 2, 3}
// return arr => return ConstArr(arr)
type CArr[T any] struct {
	arr []T
}

// ConstArr 用于将切片转换为 CArr 结构体，提供保护机制
func ConstArr[T any](input []T) CArr[T] {
	return CArr[T]{arr: input}
}

// Get 返回指定索引处的元素，并指示是否存在该元素
func (m CArr[T]) Get(index int) (T, bool) {
	var empty T
	if index < 0 || index >= len(m.arr) {
		return empty, false
	}
	return m.arr[index], true
}

// Section 复制并返回从 start 到 end 之间的切片
func (m CArr[T]) Section(start, end int) []T {
	if start < 0 || start >= len(m.arr) || end < 0 || end > len(m.arr) || start >= end {
		return nil
	}
	var arr = make([]T, end-start)
	copy(arr, m.arr[start:end])
	return arr
}

// Len 返回数组的长度
func (m CArr[T]) Len() int {
	return len(m.arr)
}

// Copy 复制并返回整个数组的浅拷贝
func (m CArr[T]) Copy() []T {
	var arr = make([]T, len(m.arr))
	copy(arr, m.arr)
	return arr
}

// Range 遍历数组并调用指定的函数，直到函数返回 false 为止
func (m CArr[T]) Range(f func(index int, value T) bool) {
	for i, v := range m.arr {
		if !f(i, v) {
			break
		}
	}
}
