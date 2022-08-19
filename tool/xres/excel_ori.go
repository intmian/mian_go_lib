package xres

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"path/filepath"
)

type ExcelColOri struct {
	Name        string
	CellStrings []string
}

type ExcelOri struct {
	SheetName string
	Columns   []*ExcelColOri
}

func GetExcelData(addr, sheet string) (*ExcelOri, error) {
	excelOri := ExcelOri{}
	// 以实际文件名为sheet名
	_, name := filepath.Split(addr)
	excelOri.SheetName = name
	excelOri.Columns = make([]*ExcelColOri, 0)
	f, err := excelize.OpenFile(addr)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// 获取工作表中指定单元格的值
	cols, err := f.GetCols(sheet)
	if err != nil {
		return nil, err
	}
	for _, col := range cols {
		excelColOri := ExcelColOri{}
		excelColOri.CellStrings = make([]string, 0)
		for i, colCell := range col {
			if i == 0 {
				excelColOri.Name = colCell
				continue
			}
			excelColOri.CellStrings = append(excelColOri.CellStrings, colCell)
		}
		excelOri.Columns = append(excelOri.Columns, &excelColOri)
	}
	return &excelOri, nil
}
