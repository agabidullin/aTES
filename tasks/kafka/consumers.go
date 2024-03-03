package kafka

import (
	"encoding/json"
	"fmt"

	messages "github.com/agabidullin/aTES/common/messages"
	topics "github.com/agabidullin/aTES/common/topics"
	"github.com/agabidullin/aTES/tasks/model"

	"gorm.io/gorm"
)

type KafkaHandlers struct {
	DB *gorm.DB
}

func (h *KafkaHandlers) InitHandler(topic string, key string, value string) {
	fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
		topic, key, value)
	switch topic {
	case topics.Accounts:
		switch key {
		case messages.RegisterAccountKey:
			{
				h.handleRegisterAccount(value)
			}
		case messages.ChangeRoleKey:
			{
				h.handleChangeRole(value)
			}
		}

	}
}

func (h *KafkaHandlers) handleRegisterAccount(value string) {
	message := messages.RegisterAccount{}
	err := json.Unmarshal([]byte(value), &message)
	if err == nil {
		account := &model.Account{PublicId: message.PublicId, Name: message.Name, Role: message.Role}
		h.DB.Create(account)
	}
}

func (h *KafkaHandlers) handleChangeRole(value string) {
	message := messages.ChangeRole{}
	err := json.Unmarshal([]byte(value), &message)
	if err == nil {
		h.DB.Model(&model.Account{}).Where("public_id = ?", message.PublicId).Update("role", message.Role)
	}
}
