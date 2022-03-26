package misc

import (
	"encoding/json"
	"io/ioutil"
)

// TJsonTool 是对于json文件的抽象，可以用匿名组合的方式进行二次包装
type TJsonTool struct {
	addr  string
	pData interface{}
}

//Data 从数据结构中取出指针，此时应该用.(*TYPE)转换为某指针继续使用
func (t *TJsonTool) Data() interface{} {
	return t.pData
}

//ReadFile 将data对应的文件反序列化到string中，如果出现问题，会返回err
//切记切记一定要传进指针
func ReadFile(filename string, data interface{}) error {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, data); err != nil {
		return err
	}
	return nil
}

//NewTJsonTool 针对传入的数据初始化，在此时会载入数据。如果传入地址为空则不进行载入
//!!!!!!!!当传入值也能正常工作，但是传入指针可以避免浪费了右值，并且可以直接对结构体进行操作!!!!!!!
func NewTJsonTool(addr string, pData interface{}) *TJsonTool {
	r := &TJsonTool{pData: pData, addr: addr}
	if addr == "" {
		return r
	}
	if !r.Load(addr) {
		return nil
	}
	return r
}

//Load 从addr对应的文件中重新载入json配置
func (t *TJsonTool) Load(addr string) bool {
	if t.addr == "" {
		return false
	}
	err := ReadFile(addr, &t.pData)
	if err != nil {
		return false
	}
	return true
}

//Save 序列化数据结构到文件
func (t *TJsonTool) Save() bool {
	if t.addr == "" {
		return false
	}
	return t.SaveOther(t.addr)
}

//SaveOther 序列化数据结构到某个文件
func (t *TJsonTool) SaveOther(addr string) bool {
	s, err := json.Marshal(t.pData)
	if err != nil {
		return false
	}
	err = ioutil.WriteFile(addr, s, 0666)
	if err != nil {
		return false
	}
	return true
}

//CheckSetDefault 检查默认值是否存在，若否则置上。返回默认值是否存在
func CheckSetDefault(m *map[string]interface{}, key string, value interface{}) bool {
	_, ok := (*m)[key]
	if !ok {
		(*m)[key] = value
		return false
	}
	return true
}
