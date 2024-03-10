package db

import (
	"github.com/agabidullin/aTES/billing/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&model.Account{})
	db.AutoMigrate(&model.Task{})
}
