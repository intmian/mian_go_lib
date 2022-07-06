package xres

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

type ExcelCowOri struct {
	Name        string
	cellStrings []string
}

type ExcelOri struct {
	SheetName string
	Cows      []*ExcelCowOri
}

func GetExcelData(addr, sheet string) (*ExcelOri, error) {
	excelOri := ExcelOri{}
	excelOri.SheetName = sheet
	excelOri.Cows = make([]*ExcelCowOri, 0)
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
		excelCowOri := ExcelCowOri{}
		excelCowOri.cellStrings = make([]string, 0)
		for _, colCell := range row {
			excelCowOri.cellStrings = append(excelCowOri.cellStrings, colCell)
		}
		excelOri.Cows = append(excelOri.Cows, &excelCowOri)
	}
	return &excelOri, nil
}
