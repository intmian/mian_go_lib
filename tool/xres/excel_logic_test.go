package xres

import (
	"github.com/intmian/mian_go_lib/tool/misc"
	"testing"
)

func TestExcelOri2Logic(t *testing.T) {
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
}
