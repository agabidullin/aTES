package handlers

import (
	"github.com/agabidullin/aTES/tasks/model"

	"gorm.io/gorm"
)

type AccountHandler struct {
	DB *gorm.DB
}

func (h *AccountHandler) AccountRegisterHandler(account *model.Account) {
	h.DB.Create(account)
}
