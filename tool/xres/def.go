package xres

//ColumnType 列类型
type ColumnType int

const (
	CtNone         ColumnType = iota //无效
	CtInt                            //整数
	CtString                         //文本
	CtFloat                          //小数
	CtEnum                           //枚举
	CtBitEnum                        //位枚举
	CtVecDataPKey                    //主枚举列
	CtVecDataCKey                    //子枚举列
	CtVecDataValue                   //数据列
	CtLogicBegin   = 100             //逻辑类型
	CtData                           // 压缩过后的逻辑数据
)

type DataCell struct {
	// 用来表示类似于 数据类型:数据 的格式
	valueType int
	value     int
}
type Data []DataCell

type ResType interface {
	int | string | float64 | Data
}
