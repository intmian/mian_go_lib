package xstorage

import (
	"github.com/intmian/mian_go_lib/tool/misc"
)

type CfgParam struct {
	Key       string // 与前端协调的key
	RealKey   string // storage里面的key
	ValueType ValueType
}

// CfgExt 为xstorage增加一个通用的配置机制
type CfgExt struct {
	core     *XStorage
	paramMap map[string]*CfgParam
	initTag  misc.InitTag
}

func (c *CfgExt) Init(core *XStorage) error {
	if core == nil {
		return ErrCoreIsNil
	}
	c.paramMap = make(map[string]*CfgParam)
	c.core = core
	c.initTag.SetInitialized()
	return nil
}

func NewCfgExt(core *XStorage, serviceName string) (*CfgExt, error) {
	ret := &CfgExt{}
	err := ret.Init(core)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *CfgExt) AddParam(param *CfgParam) error {
	if !c.initTag.IsInitialized() {
		return ErrNotInitialized
	}
	// 校验key、以及各个参数的合法性
	if param == nil {
		return ErrParamIsNil
	}
	if param.Key == "" || param.RealKey == "" || param.ValueType == 0 {
		return ErrParamIsInvalid
	}
	_, ok := c.paramMap[param.Key]
	if ok {
		return ErrKeyAlreadyExist
	}
	c.paramMap[param.Key] = param
	return nil
}

func (c *CfgExt) SetCfg(key string, svr string, value ValueUnit) error {
	if !c.initTag.IsInitialized() {
		return ErrNotInitialized
	}
	if key == "" {
		return ErrParamIsEmpty
	}
	param, ok := c.paramMap[key]
	if !ok {
		return ErrKeyNotFound
	}
	if value.Type != param.ValueType {
		return ErrValueTypeNotMatch
	}
	if svr == "" {
		return ErrParamIsEmpty
	}
	// 避免奇怪的冲突
	if svr == "cfg" || svr == "plat" || svr == "user" {
		return ErrParamIsInvalid
	}
	return c.core.Set(Join("cfg", svr, param.RealKey), &value)
}

func (c *CfgExt) SetUserCfg(svr, user, key string, value ValueUnit) error {
	if !c.initTag.IsInitialized() {
		return ErrNotInitialized
	}
	if svr == "" {
		return ErrParamIsEmpty
	}
	// 避免奇怪的冲突
	if svr == "cfg" || svr == "plat" || svr == "user" {
		return ErrParamIsInvalid
	}
	if user == "" || key == "" {
		return ErrParamIsEmpty
	}
	param, ok := c.paramMap[key]
	if !ok {
		return ErrKeyNotFound
	}
	userKey := Join("cfg", svr, "user", user, param.RealKey)
	if value.Type != param.ValueType {
		return ErrValueTypeNotMatch
	}
	return c.core.Set(userKey, &value)
}

func (c *CfgExt) GetAllCfg(svr string) (map[string]ValueUnit, error) {
	if !c.initTag.IsInitialized() {
		return nil, ErrNotInitialized
	}
	ret := make(map[string]ValueUnit)
	for k, v := range c.paramMap {
		value, err := c.core.Get(Join("cfg", svr, v.RealKey))
		if err != nil {
			return nil, err
		}
		ret[k] = *value
	}
	return ret, nil
}

func (c *CfgExt) GetUserCfg(svr, user string) (map[string]ValueUnit, error) {
	if !c.initTag.IsInitialized() {
		return nil, ErrNotInitialized
	}
	ret := make(map[string]ValueUnit)
	for k, v := range c.paramMap {
		value, err := c.core.Get(Join("cfg", svr, "user", user, v.RealKey))
		if err != nil {
			return nil, err
		}
		ret[k] = *value
	}
	return ret, nil
}
