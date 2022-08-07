package xres

import (
	"fmt"
	"strconv"
	"strings"
)

type ExcelLogicCol struct {
	name      string
	ExcelType ColumnType
	data      []interface{}
}

type ExcelLogic struct {
	SheetName string
	Columns   []*ExcelLogicCol
}

func convertNormalStrFromOri(columnType ColumnType, str string, metaEnumMap map[string]int) interface{} {
	switch columnType {
	case CtInt, CtVecDataValue:
		v, err := strconv.Atoi(str)
		if err != nil {
			return nil
		}
		return v
	case CtString:
		return str
	case CtFloat:
		v, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil
		}
		return v
	case CtEnum, CtVecDataPKey, CtVecDataCKey:
		return metaEnumMap[str]
	case CtBitEnum:
		// 根据逗号分割,并将处理后的值压入一个数组
		strs := strings.Split(str, ",")
		bitEnums := make([]int, 0)
		for _, s := range strs {
			enum, ok := metaEnumMap[s]
			if !ok {
				return nil
			}
			realEnumIndex := enum % 32
			realEnum := 1 << uint(realEnumIndex)
			for len(bitEnums) <= realEnumIndex {
				bitEnums = append(bitEnums, 0)
			}
			bitEnums[realEnumIndex] |= realEnum
		}
		return bitEnums
	default:
		return nil
	}
}

//ExcelOri2Logic 将原始的以列组织的表，进行基础的校验，并将所有的文本翻译为逻辑值
func ExcelOri2Logic(ori *ExcelOri, meta *ExcelMeta) (*ExcelLogic, error) {
	if (ori == nil) || (meta == nil) {
		return nil, fmt.Errorf("%s:ori or meta is nil", ori.SheetName)
	}
	excel := ExcelLogic{}
	excel.SheetName = ori.SheetName
	excel.Columns = make([]*ExcelLogicCol, 0)

	// 进行原始检验
	// 检查原始数据第一列是否为序号列
	if ori.Columns[0].Name != "ID" {
		return nil, fmt.Errorf("%s:first column is not 'ID'", ori.SheetName)
	}
	// 检查原始数据的每一列长度是否一致
	cellsNum := -1
	for _, col := range ori.Columns {
		if cellsNum == -1 {
			cellsNum = len(col.CellStrings)
		} else if cellsNum != len(col.CellStrings) {
			return nil, fmt.Errorf("%s:columns(%s) num is not equal, this column num is %d, last column num is %d", ori.SheetName, col.Name, len(col.CellStrings), cellsNum)
		}
	}
	// 检查原始数据的列数与meta的列数是否一致
	if len(ori.Columns) != len(meta.ColumnMeta)+1 {
		return nil, fmt.Errorf("%s:columns length is not equal to meta", ori.SheetName)
	}

	// 对每一列的数据进行处理
	for _, colOri := range ori.Columns {
		colMeta, ok := meta.ColumnMeta[colOri.Name]
		if !ok {
			return nil, fmt.Errorf("%s:column %s not found in meta", ori.SheetName, colOri.Name)
		}
		col := ExcelLogicCol{}
		col.name = colOri.Name
		col.ExcelType = colMeta.Type
		col.data = make([]interface{}, cellsNum)
		for i, cellStr := range colOri.CellStrings {
			v := convertNormalStrFromOri(colMeta.Type, cellStr, colMeta.Data)
			if v == nil {
				return nil, fmt.Errorf("%s:column %s cell %d is not valid", ori.SheetName, colOri.Name, i)
			}
		}
	}

	return &excel, nil
}
