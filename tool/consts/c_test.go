package consts

import "testing"

func TestArr(t *testing.T) {
	var arr = []int{1, 2, 3}
	carr := ConstArr(arr)
	// 判断所有读取操作是否够与原始数组一致
	if len(arr) != carr.Len() {
		t.Fatal("len error")
	}
	for i := 0; i < len(arr); i++ {
		v, ok := carr.Get(i)
		if !ok || v != arr[i] {
			t.Fatal("get error")
		}
	}

	arr2 := carr.Copy()
	arr2[0] = 4
	if arr[0] == arr2[0] {
		t.Fatal("copy error")
	}

	arr3 := carr.Section(0, 2)
	arr3[0] = 5
	arr3 = append(arr3, 6)
	if arr[0] == arr3[0] || arr[2] == 6 {
		t.Fatal("section error")
	}
}
