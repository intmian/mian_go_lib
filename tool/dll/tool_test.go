package dll

import (
	"fmt"
	"testing"
)

func TestPtr(t *testing.T) {
	p := Bytes2Ptr([]byte("1111"))
	s := string(Ptr2Bytes(p))
	if s != "1111" {
		t.Error("TestPtr failed")
	}
}

func TestPlayerKey(t *testing.T) {
	str := GetPlayerKey(1, 78910)
	areaID, playerID := ParsePlayerKey(str)

	fmt.Println(str, areaID, playerID)
}
