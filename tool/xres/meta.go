package xres

import (
	"strconv"
	"strings"
)

//ColumnType 列类型
type ColumnType int

type ExcelCowMetaOri struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type ExcelCowMeta struct {
	Type ColumnType
	Data map[string]int // 用于将填在excel中的枚举转换为int
	Child
}

type ExcelMeta struct {
	Columns map[string]*ExcelCowMeta
}

type ExcelMetaOri struct {
	Columns map[string]*ExcelCowMetaOri
}

func (m *ExcelCowMetaOri) GetColumnType() ColumnType {
	switch m.Type {
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

func (m *ExcelCowMetaOri) GetData() map[string]int {
	/*
		如果类型是枚举，则解析枚举数据
		枚举格式如下
		[枚举1:1]
		[枚举2:2]
		位枚举格式如下
		[枚举1:1]
		[枚举2:2]
	*/
	switch m.GetColumnType() {
	case CtEnum, CtBitEnum, CtVecDataPKey, CtVecDataCKey:
		data := make(map[string]int)
		for _, line := range strings.Split(m.Data, "\n") {
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

func (m *ExcelMetaOri) GetMeta() *ExcelMeta {
	meta := ExcelMeta{}
	meta.Columns = make(map[string]*ExcelCowMeta)
	for sheetName, cowMetaOri := range m.Columns {
		cowMeta := ExcelCowMeta{}
		cowMeta.Type = cowMetaOri.GetColumnType()
		cowMeta.Data = cowMetaOri.GetData()
		meta.Columns[sheetName] = &cowMeta
	}
	return &meta
}
