package model

import (
	"time"
)

// Status Define a custom type for the enum
type Status string

// Define constants using string
const (
	Open     Status = "opn"
	Closed          = "cld"
	Assigned        = "asn"
)

type Ticket struct {
	Id        int64      `json:"id,omitempty"`
	Title     string     `json:"title"`
	Message   string     `json:"message"`
	UserId    int        `json:"user_id"`
	Status    string     `json:"status"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

var MapStatusTicket map[Status]string

func init() {
	MapStatusTicket = map[Status]string{
		Open:     "Open",
		Closed:   "Closed",
		Assigned: "Assigned",
	}
}
