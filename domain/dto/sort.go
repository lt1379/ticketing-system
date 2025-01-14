package dto

type Sort struct {
	Name string `json:"sort_name"`
	Dir  string `json:"sort_dir" binding:"required,oneofci=asc desc"`
}
