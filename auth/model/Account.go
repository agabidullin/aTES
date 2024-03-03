package model

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Login    string
	Password string
	Name     string
	Role     string
}
