package xstorage

import (
	"errors"
	"github.com/intmian/mian_go_lib/tool/misc"
	"strings"
	"sync"
)

type CfgParam struct {
	Key       string // 与前端协调的key
	ValueType ValueType
	CanUser   bool      // 是否可以用户配置
	RealKey   string    // storage里面的key
	Default   ValueUnit // 如果storage里面没有这个值，就会使用这个值。
}

type ParamMap struct {
	paramMapLock sync.RWMutex
	paramMap     map[string]*CfgParam
}

func (p *ParamMap) AddParam(param *CfgParam) error {
	p.paramMapLock.Lock()
	defer p.paramMapLock.Unlock()
	if param == nil {
		return ErrParamIsNil
	}
	if param.Key == "" || param.ValueType == 0 {
		return ErrParamIsInvalid
	}
	if param.RealKey == "" {
		param.RealKey = param.Key
	}
	_, ok := p.paramMap[param.Key]
	if ok {
		return ErrKeyAlreadyExist
	}
	p.paramMap[param.Key] = param
	return nil
}

func (p *ParamMap) GetParam(key string) *CfgParam {
	p.paramMapLock.RLock()
	defer p.paramMapLock.RUnlock()
	return p.paramMap[key]
}

// CfgExt 为xstorage增加一个通用的配置机制
type CfgExt struct {
	core     *XStorage
	paramMap ParamMap
	initTag  misc.InitTag
}

func (c *CfgExt) Init(core *XStorage) error {
	if core == nil {
		return ErrCoreIsNil
	}
	c.paramMap.paramMap = make(map[string]*CfgParam)
	c.core = core
	c.initTag.SetInitialized()
	return nil
}

func NewCfgExt(core *XStorage) (*CfgExt, error) {
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
	err := c.paramMap.AddParam(param)
	if err != nil {
		return err
	}
	return nil
}

func (c *CfgExt) Set(key string, value string) error {
	if !c.initTag.IsInitialized() {
		return ErrNotInitialized
	}
	if key == "" || value == "" {
		return ErrParamIsEmpty
	}

	param := c.paramMap.GetParam(key)
	if param == nil {
		return ErrKeyNotFound
	}
	v := StringToUnit(value, param.ValueType)
	if v == nil {
		return ErrParamIsInvalid
	}

	return c.core.Set(param.RealKey, v)
}

func (c *CfgExt) SetUser(user, key string, value string) error {
	if !c.initTag.IsInitialized() {
		return ErrNotInitialized
	}
	if user == "" || key == "" || value == "" {
		return ErrParamIsEmpty
	}

	param := c.paramMap.GetParam(key)
	if param == nil {
		return ErrKeyNotFound
	}

	v := StringToUnit(value, param.ValueType)
	if v == nil {
		return ErrParamIsInvalid
	}

	return c.core.Set(Join(param.RealKey, user), v)
}

func (c *CfgExt) GetAll() (map[string]ValueUnit, error) {
	if !c.initTag.IsInitialized() {
		return nil, ErrNotInitialized
	}
	ret := make(map[string]ValueUnit)
	for k, v := range c.paramMap.paramMap {
		value, err := c.core.Get(v.RealKey)
		if err != nil {
			return nil, err
		}
		if value == nil {
			continue
		}
		ret[k] = *value
	}
	return ret, nil
}

func (c *CfgExt) GetWithFilter(prefix, user string) (map[string]ValueUnit, error) {
	if !c.initTag.IsInitialized() {
		return nil, ErrNotInitialized
	}
	ret := make(map[string]ValueUnit)
	realPrefix := prefix + "."
	for logicKey, params := range c.paramMap.paramMap {
		if prefix != "" && !strings.HasPrefix(logicKey, realPrefix) {
			continue
		}
		if user == "" {
			value, err := c.core.Get(params.RealKey)
			if err != nil {
				return nil, err
			}
			if value == nil {
				return nil, ErrKeyNotFound
			}
			ret[logicKey] = *value
		} else {
			value, err := c.core.Get(Join(params.RealKey, user))
			if err != nil {
				return nil, err
			}
			if value == nil {
				return nil, ErrKeyNotFound
			}
			ret[logicKey] = *value
		}
	}
	return ret, nil
}

func (c *CfgExt) GetUser(user string, keys ...string) (*ValueUnit, error) {
	if !c.initTag.IsInitialized() {
		return nil, ErrNotInitialized
	}
	if len(keys) == 0 {
		return nil, ErrParamIsEmpty
	}
	realLogicKey := Join(keys...)
	param := c.paramMap.GetParam(realLogicKey)
	if param == nil {
		return nil, ErrParamIsInvalid
	}
	if !param.CanUser {
		return nil, ErrParamIsInvalid
	}
	v, err := c.core.Get(Join(param.RealKey, user))
	if err != nil {
		if errors.Is(err, ErrKeyNotFound) {
			return &param.Default, nil
		} else {
			return nil, err
		}
	}
	return v, err
}

func (c *CfgExt) Get(keys ...string) (*ValueUnit, error) {
	if !c.initTag.IsInitialized() {
		return nil, ErrNotInitialized
	}
	if len(keys) == 0 {
		return nil, ErrParamIsEmpty
	}
	realLogicKey := Join(keys...)
	param := c.paramMap.GetParam(realLogicKey)
	if param == nil {
		return nil, ErrParamIsInvalid
	}

	v, err := c.core.Get(param.RealKey)
	if err != nil {
		if errors.Is(err, ErrKeyNotFound) {
			return &param.Default, nil
		} else {
			return nil, err
		}
	}
	return v, err
}
