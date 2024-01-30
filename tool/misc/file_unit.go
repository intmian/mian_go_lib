package misc

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"sync"
)

type IFileRwUnit interface {
	Read(addr string, pData interface{}) error
	Write(addr string, pData interface{}) error
}

type FileUnitType int

const (
	FileUnitJson = iota
	FileUnitToml
)

// FileUnit 是一个可以用于任意类型的非一次性序列化工具
// 可以将变量的指针与文件中的实际值关联
// NewFileUnit 会将结构体绑定某一个文件地址，之后可以通过Load和Save来进行读写
type FileUnit[T any] struct {
	addr       string
	data       T
	fileRWUnit IFileRwUnit
	lock       sync.RWMutex
}

// NewFileUnit 初始化文件单元
func NewFileUnit[T any](unitType FileUnitType, addr string) *FileUnit[T] {
	if addr == "" {
		return nil
	}
	var t T
	var kFileRWUnit IFileRwUnit
	switch unitType {
	case FileUnitJson:
		kFileRWUnit = GJsonTool
	case FileUnitToml:
		kFileRWUnit = GTomlTool
	}
	return &FileUnit[T]{data: t, fileRWUnit: kFileRWUnit, addr: addr}
}

// Load 从addr对应的文件中载入数据
func (t *FileUnit[T]) Load() error {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.fileRWUnit == nil {
		return fmt.Errorf("t.fileRWUnit is nil")
	}
	var err error
	if err = t.fileRWUnit.Read(t.addr, &t.data); err != nil {
		return err
	}
	return nil
}

// save2Addr 序列化数据结构到文件
func (t *FileUnit[T]) save2Addr(addr string) error {
	if t.fileRWUnit == nil {
		return fmt.Errorf("t.fileRWUnit is nil")
	}
	if addr == "" {
		return fmt.Errorf("addr is empty")
	}
	var err error
	if err = t.fileRWUnit.Write(addr, &t.data); err != nil {
		return err
	}
	return nil
}

// Save 序列化数据结构到文件
func (t *FileUnit[T]) Save() error {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.save2Addr(t.addr)
}

// SaveOther 序列化数据结构到某个文件
func (t *FileUnit[T]) SaveOther(addr string) error {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.save2Addr(addr)
}

// SafeUseData 线程安全的使用数据的函数
func (t *FileUnit[T]) SafeUseData(f func(any), isWrite bool) {
	if isWrite {
		t.lock.Lock()
		defer t.lock.Unlock()
	} else {
		t.lock.RLock()
		defer t.lock.RUnlock()
	}
	f(t.data)
}

// Copy 复制一份向外传达，请注意类型内部是否存在指针
func (t *FileUnit[T]) Copy() any {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.data
}

type TomlTool struct{}

func (*TomlTool) Read(addr string, pData interface{}) error {
	_, err := toml.DecodeFile(addr, pData)
	return err
}

func (*TomlTool) Write(addr string, pData interface{}) error {
	f, err := os.Create(addr)
	if err != nil {
		return err
	}
	err = toml.NewEncoder(f).Encode(pData)
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

type JsonTool struct{}

func (*JsonTool) Read(addr string, pData interface{}) error {
	return readJsonFile(addr, pData)
}

var GJsonTool = &JsonTool{}

func (*JsonTool) Write(addr string, pData interface{}) error {
	return writeJsonFile(addr, pData)
}

var GTomlTool = &TomlTool{}
