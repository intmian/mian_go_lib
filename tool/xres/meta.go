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
)

type ResType interface {
	int | string | float64 |
}

type ExcelCowMeta struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type ExcelCowsMeta struct {
	Columns map[string]*ExcelCowMeta
}

func (m *ExcelCowsMeta) GetColumnType() ColumnType {
	switch m.Columns["type"].Data {
	case "整数":
		return CtInt
	case "文本":
		return CtString
	case "小数":
		return CtFloat
	case "枚举":
		return CtEnum
	case "位枚举":
		return CtBitEnum
	case "主枚举列":
		return CtVecDataPKey
	case "子枚举列":
		return CtVecDataCKey
	case "数据列":
		return CtVecDataValue
	default:
		return CtNone
	}
}
