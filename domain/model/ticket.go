package model

import "time"

type Ticket struct {
	Id        int64      `json:"id"`
	Title     string     `json:"title"`
	Message   string     `json:"message"`
	UserId    int        `json:"user_id"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
