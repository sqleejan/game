package models

type Pagination struct {
	Total     int `json:"total"`
	TotalPage int `json:"totalPage"`
}

func PageLocate(total, page, size int) (start, end int) {
	start = page * (size - 1)
	end = page * size
	if total < start {
		return 0, 0
	}
	if total < end {
		end = total
	}
	return
}
