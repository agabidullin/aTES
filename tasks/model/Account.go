package model

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	PublicId uint `gorm:"unique"`
	Name     string
	Role     string
	Balance  int
}
