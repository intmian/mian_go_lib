package misc

func GetIndexPage(index int, pageSize int, pageBegin1 bool) (page int, pageIndex int) {
	page = index / pageSize
	pageIndex = index % pageSize
	if pageBegin1 {
		page = page + 1
	}
	return
}

func GetMaxPage(total int, pageSize int, pageBegin1 bool) (maxPage int) {
	if pageBegin1 {
		maxPage = total / pageSize
		if total%pageSize > 0 {
			maxPage = maxPage + 1
		}
	} else {
		maxPage = total / pageSize
	}
	return
}

func GetPageStartEnd(page int, pageSize int, total int, pageBegin1 bool) (pageStart int, pageEnd int) {
	if pageBegin1 {
		pageStart = (page-1)*pageSize + 1
		pageEnd = page * pageSize
	} else {
		pageStart = page * pageSize
		pageEnd = (page+1)*pageSize - 1
	}
	if pageEnd > total {
		if pageBegin1 {
			pageEnd = total
		} else {
			pageEnd = total - 1
		}
	}
	if pageStart > total {
		return -1, -1
	}
	return
}

func GetPageIndexOriIndex(index int, page int, pageSize int, pageBegin1 bool) int {
	if pageBegin1 {
		index = index - 1
		page = page - 1
	}
	return index + page*pageSize
}
