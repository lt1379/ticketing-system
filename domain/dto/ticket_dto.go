package dto

type RequestTicketDto struct {
	Title   string `json:"ticket_title" validate:"required,min=10,max=100"`
	Message string `json:"ticket_msg" validate:"required,min=100"`
	UserId  int    `json:"user_id" validate:"required"`
}
