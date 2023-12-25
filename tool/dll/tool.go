package dll

import (
	"fmt"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

func Int2Ptr(n int) uintptr {
	return uintptr(n)
}

func Bytes2Ptr(s []byte) uintptr {
	if len(s) == 0 {
		s = []byte{0x00}
	}
	if s[len(s)-1] != 0x00 {
		s = append(s, 0x00)
	}
	return uintptr(unsafe.Pointer(&s[0]))
}

func Ptr2Bytes(p uintptr) []byte {
	// 将c++的char*转换为go的[]byte,从前向后遍历，遇到0x00停止
	var b []byte
	for {
		if *(*byte)(unsafe.Pointer(p)) == 0x00 {
			break
		}
		b = append(b, *(*byte)(unsafe.Pointer(p)))
		p++
	}
	return b
}

func UseFunc(helloDll *syscall.LazyDLL, funcName string, args ...uintptr) (uintptr, uintptr, error) {
	newProc := helloDll.NewProc(funcName)
	res1, res2, err := newProc.Call(args...)
	return res1, res2, err
}

type ArgType int

const (
	AT_INT    ArgType = 0
	AT_STRING ArgType = 1
)

type UniArgMeta struct {
	ArgType ArgType
	ArgName string
}

type UniArg struct {
	UniArgMeta
	Arg interface{}
}

type Args struct {
	Args []uintptr
}

func (a *Args) Get() []uintptr {
	return a.Args
}

func (a *Args) AddInt(n int) *Args {
	a.Args = append(a.Args, Int2Ptr(n))
	return a
}

func (a *Args) AddString(s string) *Args {
	a.Args = append(a.Args, Bytes2Ptr([]byte(s)))
	return a
}

func (a *Args) AddBytes(b []byte) *Args {
	a.Args = append(a.Args, Bytes2Ptr(b))
	return a
}

func (a *Args) AddCStruct(s *CStruct) *Args {
	a.Args = append(a.Args, Bytes2Ptr(s.buff))
	return a
}

func (a *Args) AddUniArg(arg UniArg) *Args {
	switch arg.ArgType {
	case AT_INT:
		a.AddInt(arg.Arg.(int))
	case AT_STRING:
		a.AddString(arg.Arg.(string))
	}
	return a
}

func Str2ID(s string) uint32 {
	if len(s) != 4 {
		return 0
	}
	var id uint32
	for i := 0; i < 4; i++ {
		id += uint32(s[i]) << uint32(8*i)
	}
	return id
}

type CStruct struct {
	buff []byte
}

func (s *CStruct) AddInt32(n int32) *CStruct {
	n1 := byte(n & 0xFF)
	n2 := byte((n >> 8) & 0xFF)
	n3 := byte((n >> 16) & 0xFF)
	n4 := byte((n >> 24) & 0xFF)
	s.buff = append(s.buff, n1, n2, n3, n4)
	return s
}

func (s *CStruct) AddUint32(n uint32) *CStruct {
	n1 := byte(n & 0xFF)
	n2 := byte((n >> 8) & 0xFF)
	n3 := byte((n >> 16) & 0xFF)
	n4 := byte((n >> 24) & 0xFF)
	s.buff = append(s.buff, n1, n2, n3, n4)
	return s
}

func (s *CStruct) AddString(str string, add0 bool, size int) *CStruct {
	strB := []byte(str)
	if add0 {
		if strB[len(strB)-1] != 0x00 {
			strB = append(strB, 0x00)
		}
	}
	if size > 0 {
		for i := len(strB); i < size; i++ {
			strB = append(strB, 0x00)
		}
	}
	s.buff = append(s.buff, strB...)
	return s
}

func (s *CStruct) AddStringPtr(str string) *CStruct {
	ptr := Bytes2Ptr([]byte(str))
	s.AddInt32(int32(ptr))
	return s
}

type PlayerKey struct {
	AreaID   uint32
	PlayerID uint32
}

func GetPlayerKey(areaID, playerID uint32) string {
	return fmt.Sprintf("%v:%v", areaID, playerID)
}

func ParsePlayerKey(playerStr string) (uint32, uint32) {
	playerKeyStrs := strings.Split(playerStr, ":")
	if len(playerKeyStrs) != 2 {
		return 0, 0
	}
	areaID, err := strconv.ParseInt(playerKeyStrs[0], 10, 0)
	if err != nil {
		return 0, 0
	}
	playerID, err := strconv.ParseInt(playerKeyStrs[1], 10, 0)
	if err != nil {
		return 0, 0
	}

	return uint32(areaID), uint32(playerID)
}
