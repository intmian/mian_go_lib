package xres

import (
	"github.com/intmian/mian_go_lib/tool/misc"
	"os"
	"testing"
)

func TestPtl(t *testing.T) {
	o, err := GetExcelData("./test.xlsx", "正常")
	if err != nil {
		t.Error(err)
	}
	pMetaOri := &ExcelMetaOri{}
	err = misc.GTomlTool.Read(`testMeta.toml`, pMetaOri)
	if err != nil {
		t.Error(err)
	}

	if pMetaOri == nil {
		t.Error("meta is nil")
	}

	if pMetaOri.Columns == nil {
		t.Error("meta.columns is nil")
	}

	pMeta, err := pMetaOri.GetMeta()
	if err != nil {
		t.Error(err)
	}
	if pMeta == nil {
		t.Error("meta is nil")
	}

	pLogic, err := ExcelOri2Logic(o, pMeta)
	if err != nil {
		t.Error(err)
	}
	if pLogic == nil {
		t.Error("logic is nil")
	}
	if pLogic.SheetName == "" {
		t.Error("SheetName is empty")
	}
	if len(pLogic.Columns) != len(pMeta.ColumnMeta) {
		t.Error("Columns len not equal")
	}

	pPtl, err := GetExcelFromLogic(pLogic, pMeta)
	if err != nil {
		t.Error(err)
	}
	if pPtl == nil {
		t.Error("ptl is nil")
	}
	if pPtl.SheetName == "" {
		t.Error("SheetName is empty")
	}
	if len(pPtl.ColumnTypes) != len(pMeta.ColumnMeta) {
		t.Error("Columns len not equal")
	}
	if len(pPtl.Rows) == 0 {
		t.Error("Rows is empty")
	}

	// TODO: 可拓展数据列还没测试
	// TODO: python脚本检查待确认

	// 将ptl的列扩展10000行模拟真实情况
	for i := 0; i < 10000; i++ {
		pPtl.Rows = append(pPtl.Rows, pPtl.Rows[0])
	}

	err = pPtl.Save2file("./test.rxc")
	if err != nil {
		t.Error(err)
	}
	pPtl2 := &ExcelPtl{}
	err = pPtl2.LoadFromFile("./test.rxc")
	if err != nil {
		t.Error(err)
	}
	if pPtl2.SheetName == "" {
		t.Error("After WR sheetName is empty")
	}
	if len(pPtl2.ColumnTypes) != len(pMeta.ColumnMeta) {
		t.Error("After WR Columns len not equal")
	}
	if len(pPtl2.Rows) == 0 {
		t.Error("After WR Rows is empty")
	}
	for i, c := range pPtl2.ColumnTypes {
		if c != pPtl.ColumnTypes[i] {
			t.Error("After WR ColumnTypes not equal")
		}
	}
	err = os.Remove("./test.rxc")
	if err != nil {
		t.Error(err)
	}
	type Test struct {
		A int    `excel:"a"`
		B Attrs  `excel:"b"`
		V string `excel:"v"`
		K int    `excel:"k"`
	}
	tm, err := GetResFromExcelPtl[Test](pPtl)
	if err != nil {
		t.Error(err)
	}
	if tm == nil {
		t.Error("tm is nil")
	}
}
