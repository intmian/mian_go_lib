package menu

const PAGE_NUM = 10

type Menu struct {
	root  *MenuNode
	now   *MenuNode
	input inputModel
}

func (m Menu) Do() {

}

func (m Menu) GetText(page int) string {
	result := ""

	minIndex := page * PAGE_NUM
	maxIndex := minIndex + PAGE_NUM
	nowLen := len(m.now.)
}
