package xres

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

type ExcelColOri struct {
	Name        string
	cellStrings []string
}

type ExcelOri struct {
	SheetName string
	cols      []*ExcelColOri
}

func GetExcelData(addr, sheet string) (*ExcelOri, error) {
	excelOri := ExcelOri{}
	excelOri.SheetName = sheet
	excelOri.cols = make([]*ExcelColOri, 0)
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
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		excelcolOri := ExcelColOri{}
		excelcolOri.cellStrings = make([]string, 0)
		for _, colCell := range row {
			excelcolOri.cellStrings = append(excelcolOri.cellStrings, colCell)
		}
		excelOri.cols = append(excelOri.cols, &excelcolOri)
	}
	return &excelOri, nil
}
