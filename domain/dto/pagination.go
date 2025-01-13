package dto

type Pagination struct {
	PageNumber  int `json:"page_number"`
	PerPage     int `json:"per_page"`
	TotalPage   int `json:"total_page"`
	TotalRecord int `json:"total_record"`
}

type RequestPagination struct {
	Filter   *Filter `json:"filter,omitempty"`
	Sort     Sort    `json:"sort"`
	PageSize int     `json:"page_size"`
}
