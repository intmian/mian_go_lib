package xres

import (
	"strconv"
	"strings"
)

//ExcelColMetaOri 从json中读取的原始数据
type ExcelColMetaOri struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

//ExcelColMeta 转化后的列元数据
type ExcelColMeta struct {
	Type ColumnType
	Data map[string]int // 用于将填在excel中的枚举转换为int
}

//ExcelMeta Excel元数据
type ExcelMeta struct {
	ColumnMeta map[string]*ExcelColMeta
}

//ExcelMetaOri 从json中读取的原始excel元数据
type ExcelMetaOri struct {
	Columns map[string]*ExcelColMetaOri
}

//GetColumnType 从原始文本中获得列类型
func (m *ExcelColMetaOri) GetColumnType() ColumnType {
	switch m.Type {
	case "int":
		return CtInt
	case "text":
		return CtString
	case "float":
		return CtFloat
	case "enum":
		return CtEnum
	case "bitEnum":
		return CtBitEnum
	case "HeadEnum":
		return CtVecDataPKey
	case "SubEnum":
		return CtVecDataCKey
	case "EnumValue":
		return CtVecDataValue
	default:
		return CtNone
	}
}

//GetData 将原始文本中的枚举文本转换为实际的map
func (m *ExcelColMetaOri) GetData() map[string]int {
	/*
		如果类型是枚举，则解析枚举数据
		枚举格式如下
		[枚举1:1],[枚举2:2]
		位枚举格式如下
		[枚举1:1],[枚举2:2]
	*/
	switch m.GetColumnType() {
	case CtEnum, CtBitEnum, CtVecDataPKey, CtVecDataCKey:
		data := make(map[string]int)
		for _, line := range strings.Split(m.Data, ",") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if !strings.HasPrefix(line, "[") || !strings.HasSuffix(line, "]") {
				continue
			}
			line = strings.TrimPrefix(line, "[")
			line = strings.TrimSuffix(line, "]")
			kv := strings.Split(line, ":")
			if len(kv) != 2 {
				continue
			}
			key := strings.TrimSpace(kv[0])
			value, err := strconv.Atoi(strings.TrimSpace(kv[1]))
			if err != nil {
				continue
			}
			data[key] = value
		}
		return data
	default:
		return nil
	}
}

//GetMeta 将原始的元数据转换为实际的元数据
func (m *ExcelMetaOri) GetMeta() *ExcelMeta {
	meta := ExcelMeta{}
	meta.ColumnMeta = make(map[string]*ExcelColMeta)
	for sheetName, colMetaOri := range m.Columns {
		colMeta := ExcelColMeta{}
		colMeta.Type = colMetaOri.GetColumnType()
		colMeta.Data = colMetaOri.GetData()
		meta.ColumnMeta[sheetName] = &colMeta
	}

	// 增加第一列为自增长ID列
	meta.ColumnMeta["ID"] = &ExcelColMeta{
		Type: CtInt,
		Data: nil,
	}

	return &meta
}
