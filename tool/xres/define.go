package xres

type ColumnType int

const (
	CtInt ColumnType = iota
	CtString
	CtFloat
	CtEnum
	CtVecDataPKey
	CtVecDataCKey
	CtVecDataValue
)

type ExcelCowMeta struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type ExcelCowsMeta struct {
	Columns map[string]ExcelCowMeta
}
