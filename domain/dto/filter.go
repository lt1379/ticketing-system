package dto

type Filter struct {
	Name   string `json:"filter_name"`
	Type   string `json:"filter_type" binding:"required,oneofci=before after between"`
	Value  string `json:"filter_value" binding:"required,datetime=2006-01-02"`
	Value2 string `json:"filter_value2" binding:"required_if=Type between"`
}
