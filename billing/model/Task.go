package model

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	PublicId            uint `gorm:"unique"`
	Description         string
	Title               string
	AssignedTaskFee     int
	CompletedTaskAmount int
}
