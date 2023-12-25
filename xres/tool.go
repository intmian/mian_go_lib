package xres

/*
为方便使用，位枚举可以无限量拓展。例如<歼灭者:222>就是一个值为222的位枚举，同时对应一个(222/32,1<<(222%32))的属性attr
一个属性集attrs可以有多个位枚举，用逗号分隔，例如一个同时为歼灭者和精英的属性的单位，就拥有<歼灭者:222>,<精英:333>。填表时可以填写"歼灭者,精英"
存储时一个attr存储为[]int32,例如<歼灭者:222>存储在222%32的第1<<(222/32)位上

对于用户来说位枚举是不可见的，他们只能看见属性。与枚举的共同点是，都只能有32位长，不同点是属性可以叠加
*/

// 请注意这里使用int和uint没有区别，因为是位操作

type Attrs []int
type Attr int

func (a *Attrs) HasAttr(attr Attr) bool {
	enumIndex := int(attr / 32)
	bitIndex := 1 << (attr % 32)
	if enumIndex >= len(*a) {
		return false
	}
	return (*a)[enumIndex]&bitIndex != 0
}

func (a *Attrs) setAttr(attr Attr) {
	enumIndex := int(attr / 32)
	bitIndex := 1 << (attr % 32)
	for enumIndex >= len(*a) {
		*a = append(*a, 0)
	}
	(*a)[enumIndex] |= bitIndex
}
