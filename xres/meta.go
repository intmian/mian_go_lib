package xres

import (
	"errors"
	"fmt"
	"github.com/Knetic/govaluate"
	"strconv"
	"strings"
)

/* template.toml
a.type = string
a.data = """
[["枚举1:1"]
["枚举2:2"]]
"""
*/

//ExcelColMetaOri 从toml中读取的原始数据
type ExcelColMetaOri struct {
	Type string `toml:"type"`
	Data string `toml:"data"`
}

type ExcelSheetMetaOri struct {
	expressions []string `toml:"limit"`
}

//ExcelColMeta 转化后的列元数据
type ExcelColMeta struct {
	Type ColumnType
	Data map[string]int // 用于将填在excel中的枚举转换为int
}

//ExcelMeta Excel元数据
type ExcelMeta struct {
	ColumnMeta  map[string]*ExcelColMeta
	expressions []*govaluate.EvaluableExpression
}

//ExcelMetaOri 从toml中读取的原始excel元数据
type ExcelMetaOri struct {
	Columns map[string]*ExcelColMetaOri `toml:"columns"`
	Sheet   ExcelSheetMetaOri           `toml:"sheet"`
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
		[枚举1:1]
		[枚举2:2]
		位枚举格式如下
		[枚举1:1]
		[枚举2:2]
	*/
	switch m.GetColumnType() {
	case CtEnum, CtVecDataPKey, CtVecDataCKey:
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
	case CtBitEnum:
		data := make(map[string]int)
		for _, line := range strings.Split(m.Data, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if !strings.HasPrefix(line, "<") || !strings.HasSuffix(line, ">") {
				continue
			}
			line = strings.TrimPrefix(line, "<")
			line = strings.TrimSuffix(line, ">")
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
func (m *ExcelMetaOri) GetMeta() (*ExcelMeta, error) {
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

	// ID列相关校验
	m.Sheet.expressions = append(m.Sheet.expressions, "ID>0")
	m.Sheet.expressions = append(m.Sheet.expressions, "unique(ID)")

	functions := map[string]govaluate.ExpressionFunction{
		"strlen": strLen,
		"unique": makeUniqueFunc(),
		"inc":    makeIncFunc(),
		"dec":    makeDecFunc(),
	}
	exprs := make([]*govaluate.EvaluableExpression, 0)
	for _, expr := range m.Sheet.expressions {
		exprL, err := govaluate.NewEvaluableExpressionWithFunctions(expr, functions)
		if err != nil {
			return nil, errors.New("GetMeta expr error：" + expr)
		}
		exprs = append(exprs, exprL)
	}
	meta.expressions = exprs
	return &meta, nil
}

func strLen(args ...interface{}) (interface{}, error) {
	v, ok := args[0].(string)
	if !ok {
		return nil, errors.New("strLen args error")
	}
	length := len(v)
	return (float64)(length), nil
}

func makeUniqueFunc() func(args ...interface{}) (interface{}, error) {
	uniqueMap := make(map[interface{}]bool)
	return func(args ...interface{}) (interface{}, error) {
		key := args[0]
		if _, ok := uniqueMap[key]; ok {
			return false, nil
		}
		uniqueMap[key] = true
		return true, nil
	}
}

func makeIncFunc() func(args ...interface{}) (interface{}, error) {
	lastF := 0.0
	return func(args ...interface{}) (interface{}, error) {
		v, ok := args[0].(float64)
		if ok {
			if v < lastF {
				return false, nil
			}
			lastF = v
			return true, nil
		}
		return false, fmt.Errorf("makeDecIntFunc type error")
	}
}

func makeDecFunc() func(args ...interface{}) (interface{}, error) {
	lastF := 0.0
	return func(args ...interface{}) (interface{}, error) {
		v, ok := args[0].(float64)
		if ok {
			if v > lastF {
				return false, nil
			}
			lastF = v
			return true, nil
		}
		return false, fmt.Errorf("makeDecIntFunc type error")
	}
}
