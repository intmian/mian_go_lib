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
			realEnumIndex := enum / 32
			realEnum := 1 << uint(enum%32)
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

	err := CheckExcelOriLegal(ori, len(meta.ColumnMeta))
	if err != nil {
		return nil, err
	}
	// 列长度
	rawsLen := len(ori.Columns[0].CellStrings)

	// 对每一列的数据进行处理
	for _, colOri := range ori.Columns {
		colMeta, ok := meta.ColumnMeta[colOri.Name]
		if !ok {
			return nil, fmt.Errorf("%s:column %s not found in meta", ori.SheetName, colOri.Name)
		}
		col := ExcelLogicCol{}
		col.name = colOri.Name
		col.ExcelType = colMeta.Type
		col.data = make([]interface{}, rawsLen)
		for i, cellStr := range colOri.CellStrings {
			v := convertNormalStrFromOri(colMeta.Type, cellStr, colMeta.Data)
			if v == nil {
				return nil, fmt.Errorf("%s:column %s cell %d is not valid", ori.SheetName, colOri.Name, i)
			}
			col.data[i] = v
		}
		excel.Columns = append(excel.Columns, &col)
	}

	// 删除处于末尾的空行
	for i := rawsLen - 1; i >= 0; i-- {
		empty := true
		for _, col := range excel.Columns {
			if col.data[i] != nil || col.data[i] != "" {
				empty = false
				break
			}
		}
		if empty {
			for _, col := range excel.Columns {
				col.data = col.data[:i]
			}
		}
	}

	return &excel, nil
}

func CheckExcelOriLegal(ori *ExcelOri, metaLen int) error {
	// 进行原始检验
	// 检查原始数据第一列是否为序号列
	if ori.Columns[0].Name != "ID" {
		return fmt.Errorf("%s:first column is not 'ID'", ori.SheetName)
	}
	// 检查原始数据的每一列长度是否一致
	cellsNum := -1
	for _, col := range ori.Columns {
		if cellsNum == -1 {
			cellsNum = len(col.CellStrings)
		} else if cellsNum != len(col.CellStrings) {
			return fmt.Errorf("%s:columns(%s) num is not equal, this column num is %d, last column num is %d", ori.SheetName, col.Name, len(col.CellStrings), cellsNum)
		}
	}
	// 检查原始数据的列数与meta的列数是否一致
	if len(ori.Columns) != metaLen {
		return fmt.Errorf("%s:columns length is not equal to meta", ori.SheetName)
	}

	lastAllNullRawIndex := -1
	// 检查是否有非末行的空行
	for i := 0; i < len(ori.Columns[0].CellStrings); i++ {
		allNull := true
		for _, col := range ori.Columns {
			if col.CellStrings[i] != "" {
				allNull = false
				break
			}
		}
		if allNull {
			lastAllNullRawIndex = i
		}
		// 有效行上面不准出现空行
		if !allNull && (lastAllNullRawIndex != -1) {
			return fmt.Errorf("%s:have invalid row up to valid row", ori.SheetName)
		}
	}
	return nil
}
