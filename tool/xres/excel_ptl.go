package xres

// TODO: 将逻辑从实际列转换为逻辑行
import (
	"fmt"
	"strconv"
)

type ExcelCol struct {
	name      string
	ExcelType ColumnType
	data      []interface{}
}

type Excel struct {
	SheetName string
	cols      []*ExcelCol
}

func convertNormalStrFromOri(columnType ColumnType, str string, metaEnumMap map[string]int) interface{} {
	switch columnType {
	case CtInt:
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
	case CtEnum, CtBitEnum:
		return metaEnumMap[str]
	default:
		return nil
	}
}

func GetExcelFromOri(ori *ExcelOri, meta *ExcelMeta) (*Excel, error) {
	if (ori == nil) || (meta == nil) {
		return nil, fmt.Errorf("%s:ori or meta is nil", ori.SheetName)
	}
	excel := Excel{}
	excel.SheetName = ori.SheetName
	excel.cols = make([]*ExcelCol, 0)

	// 检查原始数据第一列是否为序号列
	if ori.cols[0].Name != "序号" {
		return nil, fmt.Errorf("%s:first column is not '序号'", ori.SheetName)
	}

	// 检查原始数据的每一列长度是否一致
	cellsNum := -1
	for _, col := range ori.cols {
		if cellsNum == -1 {
			cellsNum = len(col.cellStrings)
		} else if cellsNum != len(col.cellStrings) {
			return nil, fmt.Errorf("%s:columns(%s) num is not equal, this column num is %d, last column num is %d", ori.SheetName, col.Name, len(col.cellStrings), cellsNum)
		}
	}

	for _, colOri := range ori.cols {
		colMeta, ok := meta.Columns[colOri.Name]
		if !ok {
			return nil, fmt.Errorf("%s:column %s not found in meta", ori.SheetName, colOri.Name)
		}
		col := ExcelCol{}
		col.name = colOri.Name
		col.ExcelType = colMeta.Type
		col.data = make([]interface{}, 0)
		tempData := Data{}
		tempSuperValue := SuperValue{}
		for _, cellStr := range colOri.cellStrings {
			switch colMeta.Type {
			case CtVecDataPKey:
				v := convertNormalStrFromOri(CtEnum, cellStr, colMeta.Data)
				if v == nil {
					return nil, fmt.Errorf("%s:column %s cell %s convert error", ori.SheetName, colOri.Name, cellStr)
				}
				tempSuperValue.valueType = v

			default:
				v := convertNormalStrFromOri(colMeta.Type, cellStr, colMeta.Data)
				if v == nil {
					return nil, fmt.Errorf("%s:convert %s to %s failed", ori.SheetName, cellStr, colMeta.Type)
				}
			}
		}
		excel.cols = append(excel.cols, &col)
	}
}
