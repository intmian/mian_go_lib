package xres

import (
	"testing"
)

func TestGetExcelData(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		o, err := GetExcelData("./test.xlsx", "正常")
		if err != nil {
			t.Error(err)
		}
		// 检查是否全为空
		for _, c := range o.Columns {
			if len(c.CellStrings) == 0 {
				t.Error("have empty col")
			}
		}
		err = CheckExcelOriLegal(o, len(o.Columns)-1)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("lack ID", func(t *testing.T) {
		o, err := GetExcelData("./test.xlsx", "缺ID")
		if err != nil {
			t.Error(err)
		}
		err = CheckExcelOriLegal(o, len(o.Columns)-1)
		if err == nil {
			t.Error("缺ID，应该报错")
		}
	})
	t.Run("raw len do not", func(t *testing.T) {
		o, err := GetExcelData("./test.xlsx", "表格行数不对")
		if err != nil {
			t.Error(err)
		}
		err = CheckExcelOriLegal(o, len(o.Columns)-1)
		if err == nil {
			t.Error("长度不一致，应该报错")
		}
	})
	t.Run("have empty raw", func(t *testing.T) {
		o, err := GetExcelData("./test.xlsx", "出现空行")
		if err != nil {
			t.Error(err)
		}
		err = CheckExcelOriLegal(o, len(o.Columns)-1)
		if err == nil {
			t.Error("空行，应该报错")
		}
	})

}
