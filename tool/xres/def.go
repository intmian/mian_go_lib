package xres

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
)

type Data map[int]int

type ResType interface {
	int | string | float64 | Data
}
