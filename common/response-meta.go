package common

type ResponseMeta struct {
	Pagination
	ItemsTotal int `json:"items_total"`
	PagesTotal int `json:"pages_total"`
}
