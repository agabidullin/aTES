package model

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Description string
	Title       string
	IsClosed    bool
	AssigneeId  uint
	Assignee    Account `gorm:"not null; references: PublicId"`
}
