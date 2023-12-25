package misc

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"reflect"
	"sync"
)

type IFileRwUnit interface {
	Read(addr string, pData interface{}) error
	Write(addr string, pData interface{}) error
}

// FileUnit 是一个可以用于任意类型的非一次性序列化工具
// 可以将变量的指针与文件中的实际值关联
// NewFileUnit 会将结构体绑定某一个文件地址，之后可以通过Load和Save来进行读写
type FileUnit struct {
	addr       string
	pData      interface{}
	fileRWUnit IFileRwUnit
	needLock   bool
	sync.RWMutex
}

// NewFileUnit 初始化文件单元
func NewFileUnit(pData interface{}, kFileRWUnit IFileRwUnit, addr string, needLock bool) *FileUnit {
	// 使用反射判断是否为指针
	if reflect.TypeOf(pData).Kind() != reflect.Ptr {
		return nil
	}
	if addr == "" {
		return nil
	}
	return &FileUnit{pData: pData, fileRWUnit: kFileRWUnit, addr: addr, needLock: needLock}
}

// Load 从addr对应的文件中载入数据
func (t *FileUnit) Load() error {
	t.Lock()
	defer t.Unlock()
	if t.fileRWUnit == nil {
		return fmt.Errorf("t.fileRWUnit is nil")
	}
	var err error
	if err = t.fileRWUnit.Read(t.addr, t.pData); err != nil {
		return err
	}
	return nil
}

// save2Addr 序列化数据结构到文件
func (t *FileUnit) save2Addr(addr string) error {
	if t.fileRWUnit == nil {
		return fmt.Errorf("t.fileRWUnit is nil")
	}
	if addr == "" {
		return fmt.Errorf("addr is empty")
	}
	var err error
	if err = t.fileRWUnit.Write(addr, t.pData); err != nil {
		return err
	}
	return nil
}

// Save 序列化数据结构到文件
func (t *FileUnit) Save() error {
	t.RLock()
	defer t.RUnlock()
	return t.save2Addr(t.addr)
}

// SaveOther 序列化数据结构到某个文件
func (t *FileUnit) SaveOther(addr string) error {
	t.RLock()
	defer t.RUnlock()
	return t.save2Addr(addr)
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
