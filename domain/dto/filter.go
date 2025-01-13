package dto

type Filter struct {
	Name   string `json:"filter_name"`
	Type   string `json:"filter_type"`
	Value  string `json:"filter_value"`
	Value2 string `json:"filter_value2"`
}
