package dto

type RequestTicketDto struct {
	Title   string `json:"ticket_title" binding:"required,min=10,max=100"`
	Message string `json:"ticket_msg" binding:"required,min=100"`
	UserId  int    `json:"user_id" binding:"required"`
}
