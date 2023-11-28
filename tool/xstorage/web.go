package xstorage

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"regexp"
	"strconv"
)

type WebCode int // web业务码

// 可能会被外部调用，所以复杂命名
const (
	WebCodeNull WebCode = iota
	WebCodeSuc
	WebCodeFail
)

type EWebFailReason int // web失败原因

const (
	WebFailReasonNull EWebFailReason = iota
	WebFailReasonNoLegalParam
	WebFailReasonInnerError
)

func (m *Mgr) StartWeb() error {
	if m.setting.webPort != 0 {
		m.ginEngine = gin.Default()
		m.ginEngine.GET("/get", m.WebGet)
		m.ginEngine.GET("/set", m.WebSet)
		m.ginEngine.GET("/get_all", m.WebGetAll)
		addr := fmt.Sprintf("127.0.0.1:%d", m.setting.webPort)
		err := m.ginEngine.Run(addr)
		if err != nil {
			return errors.Join(errors.New("gin engine run error"), err)
		}
	}
	return nil
}
func (m *Mgr) WebGet(c *gin.Context) {
	if !m.initTag.IsInitialized() {
		c.JSON(200, gin.H{
			"code": WebCodeFail,
			"msg":  WebFailReasonInnerError,
		})
		return
	}
	// 读取正则表达式
	useRe := c.Query("useRe")
	perm := c.Query("perm")
	var results []ValueUnit
	if useRe != "true" {
		ok, result, err := m.Get(perm)
		if err != nil {
			m.Error("xStorage:WebGet:get value error:" + err.Error())
			c.JSON(200, gin.H{
				"code": WebCodeFail,
				"msg":  WebFailReasonNoLegalParam,
			})
			return
		}
		if !ok {
			c.JSON(200, gin.H{
				"code": WebCodeFail,
				"msg":  WebFailReasonNoLegalParam,
			})
			return
		}
		results = append(results, *result)
	} else {
		// 遍历并且搜索正则
		all, err := m.GetAll()
		if err != nil {
			m.Error("xStorage:WebGet:get all value error:" + err.Error())
			c.JSON(200, gin.H{
				"code": WebCodeFail,
				"msg":  WebFailReasonInnerError,
			})
			return
		}
		found := false
		for k, _ := range all {
			// 使用正则
			matched, err := regexp.MatchString(perm, k)
			if err != nil {
				m.Error("xStorage:WebGet:match string error:" + err.Error())
				c.JSON(200, gin.H{
					"code": WebCodeFail,
					"msg":  WebFailReasonInnerError,
				})
				return
			}
			if matched {
				found = true
				results = append(results, *all[k])
			}
		}
		if !found {
			c.JSON(200, gin.H{
				"code": WebCodeFail,
				"msg":  WebFailReasonNoLegalParam,
			})
			return
		}
	}
	c.JSON(200, gin.H{
		"code":   WebCodeSuc,
		"result": results,
	})
}

func (m *Mgr) WebSet(c *gin.Context) {
	if !m.initTag.IsInitialized() {
		c.JSON(200, gin.H{
			"code": WebCodeFail,
			"msg":  WebFailReasonInnerError,
		})
		return
	}
	key := c.Query("key")
	valueType := c.Query("value_type")
	valueTypeInt, err := strconv.Atoi(valueType)
	if err != nil {
		c.JSON(200, gin.H{
			"code": WebCodeFail,
			"msg":  WebFailReasonNoLegalParam,
		})
		return
	}
	value := c.Query("value")

	err = m.Set(key, ToUnit(value, ValueType(valueTypeInt)))
	if err != nil {
		m.Error("xStorage:WebSet:set value error:" + err.Error())
		c.JSON(200, gin.H{
			"code": WebCodeFail,
			"msg":  WebFailReasonInnerError,
		})
		return
	}
	c.JSON(200, gin.H{
		"code": WebCodeSuc,
	})
}

func (m *Mgr) WebGetAll(c *gin.Context) {
	if !m.initTag.IsInitialized() {
		c.JSON(200, gin.H{
			"code": WebCodeFail,
			"msg":  WebFailReasonInnerError,
		})
		return
	}
	all, err := m.GetAll()
	if err != nil {
		m.Error("xStorage:WebGet:get all value error:" + err.Error())
		c.JSON(200, gin.H{
			"code": WebCodeFail,
			"msg":  WebFailReasonInnerError,
		})
		return
	}
	c.JSON(200, gin.H{
		"code":   WebCodeSuc,
		"result": all,
	})
}
