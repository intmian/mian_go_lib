package xbi

type Setting struct {
	// 权限配置
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	SecurityToken   string
	// 分组
	PpjName   PpjName
	DbName    DbName
	TableName TableName
}
