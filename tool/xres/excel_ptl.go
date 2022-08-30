package xres

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"reflect"
)

//ExcelPtlRow 单行数据
type ExcelPtlRow struct {
	Data []interface{}
}

type ExcelPtl struct {
	SheetName   string
	ColumnTypes []ColumnType
	Names       []string
	Rows        []*ExcelPtlRow
}

//GetExcelFromLogic 将以列组织的数据重整为以行组织的数据
func GetExcelFromLogic(logic *ExcelLogic, meta *ExcelMeta) (*ExcelPtl, error) {
	if (logic == nil) || (meta == nil) {
		return nil, fmt.Errorf("%s:logic or meta is nil", logic.SheetName)
	}
	excel := ExcelPtl{}
	excel.SheetName = logic.SheetName
	excel.Rows = make([]*ExcelPtlRow, 0)
	excel.ColumnTypes = make([]ColumnType, 0)

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
	if len(logic.Columns) != len(meta.ColumnMeta) {
		return nil, fmt.Errorf("%s:columns length is not equal to meta", logic.SheetName)
	}

	// 处理类型
	for i, col := range logic.Columns {
		// 非复杂类型
		switch col.ExcelType {
		case CtVecDataPKey, CtVecDataCKey:
			continue
		case CtVecDataValue:
			if i == len(logic.Columns)-1 ||
				(logic.Columns[i+1].ExcelType != CtVecDataPKey &&
					logic.Columns[i+1].ExcelType != CtVecDataCKey) {
				excel.ColumnTypes = append(excel.ColumnTypes, CtData)
			}
		case CtBitEnum:
			excel.ColumnTypes = append(excel.ColumnTypes, CtAttrs)
		default:
			excel.ColumnTypes = append(excel.ColumnTypes, col.ExcelType)
		}
	}

	// 处理名称
	excel.Names = make([]string, len(logic.Columns))
	for i, col := range logic.Columns {
		excel.Names[i] = col.name
	}

	// 将原始数据重整为以行组织的数据
	for i := 0; i < RowNum; i++ {
		row := ExcelPtlRow{}
		row.Data = make([]interface{}, 0)
		var tempData Data
		var tempDataCell DataCell
		lastCellType := CtNone
		idExistMap := make(map[int]bool)
		for j, col := range logic.Columns {
			// 校验ID
			if j == 0 {
				// 序号列
				if col.name != "ID" {
					return nil, fmt.Errorf("%s:first column is not 'ID'", logic.SheetName)
				}
				id, ok := col.data[j].(int)
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
				return nil, fmt.Errorf("%s:column %s row %d is nil", logic.SheetName, col.name, j)
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
					return nil, fmt.Errorf("%s:column %s row %d last key do not have value", logic.SheetName, col.name, j)
				}
				tempData = append(tempData, tempDataCell)
				tempDataCell = DataCell{}
				tempDataCell.valueType = d.(int)
			case CtVecDataValue:
				if lastCellType != CtVecDataPKey && lastCellType != CtVecDataCKey {
					// 出现了孤立的value
					return nil, fmt.Errorf("%s:column %s row %d value has no key", logic.SheetName, col.name, j)
				}
				tempDataCell.value = d.(int)
				if j == len(logic.Columns)-1 || logic.Columns[j+1].ExcelType != CtVecDataCKey {
					row.Data = append(row.Data, tempData)
				}
			case CtAttrs:
				row.Data = append(row.Data, Attrs(d.([]int)))
			default:
				row.Data = append(row.Data, d)
			}
			lastCellType = colMeta.Type
		}
		excel.Rows = append(excel.Rows, &row)
	}
	return &excel, nil
}

//Save2file 将excel数据保存到文件
//首先存入每一列的类型，然后存入每一行的数据
func (e *ExcelPtl) Save2file(addr string) error {
	// 创建文件
	file, err := os.Create(addr)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入表名称长度和表名称（写字符串竟然不会写\0...）
	err = binary.Write(file, binary.LittleEndian, int32(len(e.SheetName)))
	if err != nil {
		return err
	}
	_, err = file.WriteString(e.SheetName)
	if err != nil {
		return err
	}

	// 写入列类型
	err = binary.Write(file, binary.LittleEndian, int32(len(e.ColumnTypes)))
	if err != nil {
		return err
	}
	for _, colType := range e.ColumnTypes {
		err = binary.Write(file, binary.LittleEndian, int32(colType))
		if err != nil {
			return err
		}
	}

	// 写入列名称
	err = binary.Write(file, binary.LittleEndian, int32(len(e.Names)))
	if err != nil {
		return err
	}
	for _, name := range e.Names {
		err = binary.Write(file, binary.LittleEndian, int32(len(name)))
		_, err = file.WriteString(name)
		if err != nil {
			return err
		}
	}

	// 写入每一行的数据
	err = binary.Write(file, binary.LittleEndian, int32(len(e.Rows)))
	if err != nil {
		return err
	}
	for _, row := range e.Rows {
		var mark int32
		mark = 23333
		err = binary.Write(file, binary.LittleEndian, mark)
		for i, data := range row.Data {
			// 获得此列的类型
			colType := e.ColumnTypes[i]
			switch colType {
			case CtInt, CtEnum:
				err = binary.Write(file, binary.LittleEndian, int32(data.(int)))
			case CtAttrs:
				err = binary.Write(file, binary.LittleEndian, int32(len(data.([]int))))
				for _, attr := range data.([]int) {
					err = binary.Write(file, binary.LittleEndian, int32(attr))
				}
			case CtFloat:
				err = binary.Write(file, binary.LittleEndian, data.(float64))
			case CtString:
				err = binary.Write(file, binary.LittleEndian, int32(len(data.(string))))
				if err != nil {
					return err
				}
				_, err = file.WriteString(data.(string))
			case CtData:
				err = binary.Write(file, binary.LittleEndian, int32(len(data.([]DataCell))))
				for _, cell := range data.(Data) {
					err = binary.Write(file, binary.LittleEndian, int32(cell.valueType))
					if err != nil {
						return err
					}
					err = binary.Write(file, binary.LittleEndian, int32(cell.value))
					if err != nil {
						return err
					}

				}
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//LoadFromFile 从文件中加载excel数据
func (e *ExcelPtl) LoadFromFile(addr string) error {
	// 清空自身
	e.ColumnTypes = nil
	e.Names = nil
	e.Rows = nil

	// 打开文件
	file, err := os.Open(addr)
	if err != nil {
		return err
	}
	defer file.Close()

	// 读取表名称长度和表名称
	var sheetNameLen int32
	err = binary.Read(file, binary.LittleEndian, &sheetNameLen)
	if err != nil {
		return err
	}
	sheetName := make([]byte, sheetNameLen)
	_, err = file.Read(sheetName)
	if err != nil {
		return err
	}
	e.SheetName = string(sheetName)

	// 读取列类型
	var colTypeLen int32
	err = binary.Read(file, binary.LittleEndian, &colTypeLen)
	if err != nil {
		return err
	}
	e.ColumnTypes = make([]ColumnType, colTypeLen)
	for i := 0; i < int(colTypeLen); i++ {
		var colType int32
		err = binary.Read(file, binary.LittleEndian, &colType)
		e.ColumnTypes[i] = ColumnType(colType)
		if err != nil {
			return err
		}
	}

	// 读取列名称
	var nameLen int32
	err = binary.Read(file, binary.LittleEndian, &nameLen)
	if err != nil {
		return err
	}
	e.Names = make([]string, nameLen)
	for i := 0; i < int(nameLen); i++ {
		var nameLen int32
		err = binary.Read(file, binary.LittleEndian, &nameLen)
		if err != nil {
			return err
		}
		buf := make([]byte, nameLen)
		_, err = file.Read(buf)
		if err != nil {
			return err
		}
		e.Names[i] = string(buf)
	}

	// 读取每一行的数据
	var rowLen int32
	err = binary.Read(file, binary.LittleEndian, &rowLen)
	if err != nil {
		return err
	}
	e.Rows = make([]*ExcelPtlRow, rowLen)
	for i := 0; i < int(rowLen); i++ {
		e.Rows[i] = &ExcelPtlRow{}
		e.Rows[i].Data = make([]interface{}, len(e.ColumnTypes))
	}
	for i := 0; i < int(rowLen); i++ {
		var mark int32
		err = binary.Read(file, binary.LittleEndian, &mark)
		if err != nil {
			return err
		}
		if mark != 23333 {
			return errors.New("协议错位")
		}
		for j, colType := range e.ColumnTypes {
			var data interface{}
			switch colType {
			case CtInt, CtEnum:
				var d int32
				err = binary.Read(file, binary.LittleEndian, &d)
				if err != nil {
					return err
				}
				data = int(d)
			case CtAttrs:
				var d int32
				err = binary.Read(file, binary.LittleEndian, &d)
				if err != nil {
					return err
				}
				data = make([]int, d)
				for k := 0; k < int(d); k++ {
					var num int32
					err = binary.Read(file, binary.LittleEndian, &num)
					if err != nil {
						return err
					}
					data.([]int)[k] = int(num)
				}
			case CtFloat:
				var d float64
				err = binary.Read(file, binary.LittleEndian, &d)
				if err != nil {
					return err
				}
				data = d
			case CtString:
				var d int32
				err = binary.Read(file, binary.LittleEndian, &d)
				if err != nil {
					return err
				}
				buf := make([]byte, d)
				_, err = file.Read(buf)
				if err != nil {
					return err
				}
				data = string(buf)
			case CtData:
				var d int32
				err = binary.Read(file, binary.LittleEndian, &d)
				if err != nil {
					return err
				}
				data = make(Data, d)
				for i := 0; i < int(d); i++ {
					var cellDataType int
					err = binary.Read(file, binary.LittleEndian, &cellDataType)
					if err != nil {
						return err
					}
					var cellData int
					err = binary.Read(file, binary.LittleEndian, &cellData)
					if err != nil {
						return err
					}
					data.(Data)[i] = DataCell{valueType: cellDataType, value: cellData}
				}
			}
			if err != nil {
				return err
			}
			e.Rows[i].Data[j] = data
		}
	}
	return nil
}

//Convert 将excel数据转换为任意结构体，根据Tag进行转换
func (e *ExcelPtl) Convert(Rec []interface{}) error {
	Rec = make([]interface{}, len(e.Rows))
	for i := 0; i < len(e.Rows); i++ {
		d := Rec[i]
		typeOfD := reflect.TypeOf(d)
		for j := 0; j < len(e.Rows[i].Data); j++ {
			cellData := e.Rows[i].Data[j]
			colType := e.ColumnTypes[j]
			for k := 0; k < typeOfD.NumField(); k++ {
				if typeOfD.Field(k).Tag.Get("excel") == e.Names[j] {
					switch colType {
					case CtInt, CtEnum:
						d.(*reflect.Value).Elem().Field(k).SetInt(int64(cellData.(int)))
					case CtBitEnum:
						d.(*reflect.Value).Elem().Field(k).Set(reflect.ValueOf(cellData.([]int)))
					case CtFloat:
						d.(*reflect.Value).Elem().Field(k).SetFloat(float64(cellData.(float64)))
					case CtString:
						d.(*reflect.Value).Elem().Field(k).SetString(cellData.(string))
					case CtData:
						d.(*reflect.Value).Elem().Field(k).Set(reflect.ValueOf(cellData.(Data)))
					}
				}
			}

		}
	}
	return nil
}

func GetResFromExcelPtl[t any](ptl *ExcelPtl) (map[int]t, error) {
	if ptl == nil {
		return nil, errors.New("ptl is nil")
	}
	resultMap := make(map[int]t)
	var tempT t
	tType := reflect.TypeOf(tempT)
	name2Index := make(map[string]int)
	for i := 0; i < len(ptl.Names); i++ {
		name2Index[ptl.Names[i]] = i
	}
	excelIndex2fieldIndex := make(map[int]int)
	for i := 0; i < tType.NumField(); i++ {
		field := tType.Field(i)
		// Tag 与 Name 一一对应
		if field.Tag.Get("excel") != "" {
			excelIndex2fieldIndex[name2Index[field.Tag.Get("excel")]] = i
		}
	}
	// 将excel数据转换为任意结构体，根据Tag进行转换
	for i := 0; i < len(ptl.Rows); i++ {
		tempT = reflect.New(tType).Elem().Interface().(t)
		for j := 0; j < len(ptl.Rows[i].Data); j++ {
			cellData := ptl.Rows[i].Data[j]
			fieldIndex, ok := excelIndex2fieldIndex[j]
			if !ok {
				continue
			}
			field := tType.Field(fieldIndex)
			if reflect.TypeOf(cellData).Kind() != field.Type.Kind() {
				return nil, errors.New("cellData type is not equal to field type")
			}
			if !field.IsExported() {
				return nil, fmt.Errorf("field %s is not exported", field.Name)
			}
			reflect.ValueOf(&tempT).Elem().Field(fieldIndex).Set(reflect.ValueOf(cellData))
		}
		resultMap[ptl.Rows[i].Data[0].(int)] = tempT
	}
	return resultMap, nil
}
