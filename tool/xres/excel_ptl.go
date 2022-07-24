package xres

import "fmt"

//ExcelPtlRow 单行数据
type ExcelPtlRow struct {
	Data []interface{}
}

type Excel struct {
	Name string
	Rows []*ExcelPtlRow
}

//GetExcelFromLogic 将以列组织的数据重整为以行组织的数据
func GetExcelFromLogic(logic *ExcelLogic, meta *ExcelMeta) (*Excel, error) {
	if (logic == nil) || (meta == nil) {
		return nil, fmt.Errorf("%s:logic or meta is nil", logic.SheetName)
	}
	excel := Excel{}
	excel.Name = logic.SheetName
	excel.Rows = make([]*ExcelPtlRow, 0)

	// 进行原始检验
	// 检查原始数据第一列是否为序号列
	if logic.Columns[0].name != "ID" {
		return nil, fmt.Errorf("%s:first column is not 'ID'", logic.SheetName)
	}
	// 检查原始数据的每一列长度是否一致
	RowNum := -1
	for _, col := range logic.Columns {
		if RowNum == -1 {
			RowNum = len(col.data)
		} else if RowNum != len(col.data) {
			return nil, fmt.Errorf("%s:columns length is not equal", logic.SheetName)
		}
	}
	// 检查原始数据的列数与meta的列数是否一致
	if len(logic.Columns) != len(meta.ColumnMeta)+1 {
		return nil, fmt.Errorf("%s:columns length is not equal to meta", logic.SheetName)
	}

	// 将原始数据重整为以行组织的数据
	for i := 0; i < RowNum; i++ {
		row := ExcelPtlRow{}
		row.Data = make([]interface{}, 0)
		var tempData Data
		var tempDataCell DataCell
		lastCellType := CtNone
		idExistMap := make(map[int]bool)
		for i, col := range logic.Columns {
			// 校验ID
			if i == 0 {
				// 序号列
				if col.name != "ID" {
					return nil, fmt.Errorf("%s:first column is not 'ID'", logic.SheetName)
				}
				id, ok := col.data[i].(int)
				if !ok {
					return nil, fmt.Errorf("%s:first column is not int type", logic.SheetName)
				}
				_, ok = idExistMap[id]
				if ok {
					return nil, fmt.Errorf("%s:id is exist", logic.SheetName)
				}
				idExistMap[id] = true
			}
			// 读取此行的meta
			colMeta := meta.ColumnMeta[col.name]

			d := col.data[i] // 读取此列此行的数据
			if d == nil {
				return nil, fmt.Errorf("%s:column %s row %d is nil", logic.SheetName, col.name, i)
			}

			// 处理值
			// 如果是复杂类型的需要进行压缩
			switch colMeta.Type {
			case CtVecDataPKey:
				tempData = Data{}
				tempDataCell = DataCell{}
				tempDataCell.valueType = d.(int)
			case CtVecDataCKey:
				if lastCellType != CtVecDataValue {
					// 上一个key-value没有闭合
					return nil, fmt.Errorf("%s:column %s row %d last key do not have value", logic.SheetName, col.name, i)
				}
				tempData = append(tempData, tempDataCell)
				tempDataCell = DataCell{}
				tempDataCell.valueType = d.(int)
			case CtVecDataValue:
				if lastCellType != CtVecDataPKey && lastCellType != CtVecDataCKey {
					// 出现了孤立的value
					return nil, fmt.Errorf("%s:column %s row %d value has no key", logic.SheetName, col.name, i)
				}
				tempDataCell.value = d.(int)
				if i == len(logic.Columns)-1 || logic.Columns[i+1].ExcelType != CtVecDataCKey {
					row.Data = append(row.Data, tempData)
				}
			default:
				row.Data = append(row.Data, d)
			}
			lastCellType = colMeta.Type
		}
		excel.Rows = append(excel.Rows, &row)
	}
	return &excel, nil
}

// TODO: 写入写出以及万能的load to interface{}
