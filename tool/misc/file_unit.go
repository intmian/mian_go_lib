package misc

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"reflect"
)

type IFileRwUnit interface {
	Read(addr string, pData interface{}) error
	Write(addr string, pData interface{}) error
}

// TFileUnit 是一个可以用于任意类型的非一次性序列化工具，可以讲变量的指针与文件中的实际值关联，请注意这个容器不是线程安全的
type TFileUnit struct {
	addr       string
	pData      interface{}
	fileRWUnit IFileRwUnit
}

//NewFileUnit 初始化文件单元
func NewFileUnit(pData interface{}, kFileRWUnit IFileRwUnit, addr string) *TFileUnit {
	// 使用反射判断是否为指针
	if reflect.TypeOf(pData).Kind() != reflect.Ptr {
		return nil
	}
	if addr == "" {
		return nil
	}
	return &TFileUnit{pData: pData, fileRWUnit: kFileRWUnit, addr: addr}
}

//Load 从addr对应的文件中载入数据
func (t *TFileUnit) Load() error {
	if t.fileRWUnit == nil {
		return fmt.Errorf("t.fileRWUnit is nil")
	}
	var err error
	if err = t.fileRWUnit.Read(t.addr, t.pData); err != nil {
		return err
	}
	return nil
}

//save2Addr 序列化数据结构到文件
func (t *TFileUnit) save2Addr(addr string) error {
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

//Save 序列化数据结构到文件
func (t *TFileUnit) Save() error {
	return t.save2Addr(t.addr)
}

//SaveOther 序列化数据结构到某个文件
func (t *TFileUnit) SaveOther(addr string) error {
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
