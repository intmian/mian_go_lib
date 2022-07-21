package xres

import "fmt"

type ExcelCow struct {
	name      string
	ExcelType ColumnType
	data      interface{}
}

type Excel struct {
	SheetName string
	Cows      []*ExcelCow
}

func Get[T any](from *ExcelCow) *T {
	return from.data.(*T)
}

func GetExcelFromOri(ori *ExcelOri, meta *ExcelMeta) (*Excel, error) {
	if (ori == nil) || (meta == nil) {
		return nil, fmt.Errorf("ori or meta is nil")
	}
	excel := Excel{}
	excel.SheetName = ori.SheetName
	excel.Cows = make([]*ExcelCow, 0)
	for _, cowOri := range ori.Cows {
		cowMeta, ok := meta.Columns[cowOri.Name]
		if !ok {
			return nil, fmt.Errorf("column %s not found in meta", cowOri.Name)
		}
		cow := ExcelCow{}
		cow.name = cowOri.Name
		cow.ExcelType = cowMeta.Type
		cow.data = cowOri.cellStrings
		excel.Cows = append(excel.Cows, &cow)
	}
}
