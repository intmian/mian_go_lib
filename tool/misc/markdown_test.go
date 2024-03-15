package misc

import "testing"

func TestMd(t *testing.T) {
	m := MarkdownTool{}
	m.AddTitle("test1", 1)
	m.AddTitle("test2", 2)
	m.AddContent("test3")
	m.AddContent("test4")
	m.AddList("test5", 1)
	m.AddList("test5", 1)
	m.AddList("test6", 2)
	m.AddList("test5", 1)
	m.AddMd("# test7\n")
	m.AddTitle("test8", 3)
	m.AddTitle("test8", 3)
	m.AddContent("test9")
	if m.ToStr() != "## test1\n### test2\ntest3\n\ntest4\n- test5\n- test5\n  - test6\n- test5\n#test7\n#### test8\n#### test8\ntest9\n" {
		t.Error("markdown error")
	}
}
