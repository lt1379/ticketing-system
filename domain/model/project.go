package model

import "time"

type Project struct {
	Id          int
	Name        string
	Description string
	CreatedAt   time.Time `gorm:"autoCreateTime;index"`
	UpdateAt    time.Time `gorm:"autoUpdateTime;index"`
}
