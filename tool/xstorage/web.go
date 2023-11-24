package xstorage

import "github.com/gin-gonic/gin"

func (m *Mgr) StartWeb() {
	if m.setting.webPort != 0 {
		m.ginEngine = gin.Default()
		m.
	}

}
func (m *Mgr) WebGet(c *gin.Context) {
	if !m.initTag.IsInitialized() {
		c.JSON(200, gin.H{
			"code": 1,/**/
			"msg":  "mgr not init",
		})
		return
	}
	// 读取正则表达式
	useRe := c.Query("useRe")
	perm := c.Query("perm")
	var results []ValueUnit
	if useRe == "true" {
		ok, result, err := m.Get(perm)
		if err != nil {
			m.Error("xStorage:WebGet:get value error:" + err.Error())
			c.JSON(200, gin.H{
				"code": 1,
				"msg":  "get value error",
			})
			return
		}
		if !ok {
			c.JSON(200, gin.H{
				"code": 1,
				"msg":  "no value",
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
				"code": 1,
				"msg":  "get all value error",
			})
			return
		}
	}
}

func (m *Mgr) WebSet(c *gin.Context) {

}
