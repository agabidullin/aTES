package db

import (
	"github.com/agabidullin/aTES/tasks/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&model.Account{})
}
