package misc

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"regexp"
)

type MarkdownTool struct {
	s           string
	inContent   bool
	inReference bool
}

func (m *MarkdownTool) AddTitle(title string, lv int) {
	m.s += "#"
	for i := 0; i < lv; i++ {
		m.s += "#"
	}
	m.s += " " + title + "\n"
	m.inContent = false
	m.inReference = false
}

func (m *MarkdownTool) AddContent(content string) {
	if m.inContent {
		m.s += "\n"
	}
	m.s += content + "\n"
	m.inContent = true
	m.inReference = false
}

func (m *MarkdownTool) AddReference(reference string) {
	if m.inReference {
		m.s += ">  \n"
	}
	m.s += "> " + reference + "\n"
	m.inReference = true
	m.inContent = false
}

func (m *MarkdownTool) ToStr() string {
	return m.s
}

func (m *MarkdownTool) AddList(content string, lv int) {
	for i := 1; i < lv; i++ {
		m.s += "  "
	}
	if m.inReference {
		m.s += "> "
	}
	m.s += "- " + content + "\n"
}

func (m *MarkdownTool) AddMd(md string) {
	m.s += md
}

// MarkdownToHTML 提供markdown到html的功能
func MarkdownToHTML(md string) string {
	myHTMLFlags := 0 |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_DASHES |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES

	myExtensions := 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS |
		blackfriday.EXTENSION_HARD_LINE_BREAK

	renderer := blackfriday.HtmlRenderer(myHTMLFlags, "", "")
	bytes := blackfriday.MarkdownOptions([]byte(md), renderer, blackfriday.Options{
		Extensions: myExtensions,
	})
	theHTML := string(bytes)
	return bluemonday.UGCPolicy().Sanitize(theHTML)
}

// GetPicLinkFromStr 从字符串中提取所有的图片外链
func GetPicLinkFromStr(mdStr string) []string {
	reg := regexp.MustCompile(`![.*](.*)`)
	if reg == nil {
		return nil
	}
	//提取
	return reg.FindAllString(mdStr, -1)
}
