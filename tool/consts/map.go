package consts

// CMap 提供对底层 map 的保护，避免直接返回内部数据而被用户修改。
// copy 方法会复制一个浅拷贝
// get、len、range 与原始 map 一致
// delete、set 需要复制并返回一个新的 map
// 使用方法：
// var m = map[string]int{"a": 1, "b": 2}
// return m => return ConstMap(m)
type CMap[K comparable, V any] struct {
	m map[K]V
}

// ConstMap 用于将 map 转换为 CMap 结构体，提供保护机制
func ConstMap[K comparable, V any](input map[K]V) CMap[K, V] {
	return CMap[K, V]{m: input}
}

// Copy 复制并返回 map 的浅拷贝
func (m CMap[K, V]) Copy() map[K]V {
	var newMap = make(map[K]V)
	for k, v := range m.m {
		newMap[k] = v
	}
	return newMap
}

// Len 返回 map 的长度
func (m CMap[K, V]) Len() int {
	return len(m.m)
}

// Delete 复制 map 并删除指定的键，返回修改后的新 map
func (m CMap[K, V]) Delete(key K) map[K]V {
	newMap := m.Copy()
	delete(newMap, key)
	return newMap
}

// Get 获取指定键的值及其存在状态
func (m CMap[K, V]) Get(key K) (V, bool) {
	v, ok := m.m[key]
	return v, ok
}

// Set 复制 map 并设置指定键的值，返回修改后的新 map
func (m CMap[K, V]) Set(key K, value V) map[K]V {
	newMap := m.Copy()
	newMap[key] = value
	return newMap
}

// Range 遍历 map 并调用指定的函数，直到函数返回 false 为止
func (m CMap[K, V]) Range(f func(key K, value V) bool) {
	for k, v := range m.m {
		if !f(k, v) {
			break
		}
	}
}
