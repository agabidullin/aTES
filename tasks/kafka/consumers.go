package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/agabidullin/aTES/common/events"
	"github.com/agabidullin/aTES/common/topics"
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
	case topics.AccountsStream:
		switch key {
		case events.AccountCreated:
			{
				h.handleAccountCreated(value)
			}
		}
	case topics.Accounts:
		switch key {
		case events.RoleChanged:
			{
				h.handleRoleChanged(value)
			}
		}

	}
}

func (h *KafkaHandlers) handleAccountCreated(value string) {
	message := events.AccountCreatedPayload{}
	err := json.Unmarshal([]byte(value), &message)
	if err == nil {
		account := &model.Account{PublicId: message.PublicId, Name: message.Name, Role: message.Role}
		h.DB.Create(account)
	}
}

func (h *KafkaHandlers) handleRoleChanged(value string) {
	message := events.RoleChangedPayload{}
	err := json.Unmarshal([]byte(value), &message)
	if err == nil {
		h.DB.Model(&model.Account{}).Where("public_id = ?", message.PublicId).Update("role", message.Role)
	}
}
