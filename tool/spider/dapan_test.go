package spider

import (
	"testing"
)

func TestDapan000001(t *testing.T) {
	price, inc, radio := GetDapan000001()
	if price == "" || inc == "" || radio == "" {
		t.Error("GetDapan000001 error")
	}
	s := ParseDapanToMarkdown("上证指数", price, inc, radio)
	t.Log(s)
}
