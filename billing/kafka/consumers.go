package kafka

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/agabidullin/aTES/billing/model"
	"github.com/agabidullin/aTES/common/events"
	"github.com/agabidullin/aTES/common/topics"
	"github.com/agabidullin/aTES/schemaregistry/validator"

	"gorm.io/gorm"
)

type KafkaHandlers struct {
	DB *gorm.DB
}

func (h *KafkaHandlers) InitHandler(topic string, key string, value string) {
	fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
		topic, key, value)

	var message events.Event
	err := json.Unmarshal([]byte(value), &message)

	if err != nil {
		// TODO handle error
		return
	}

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
	case topics.TasksStream:
		switch key {
		case events.TaskCreated:
			{
				if _, err = validator.ValidateFromString(value, topics.TasksStream, events.TaskCreated, message.Version); err != nil {
					// TODO handle
					return
				}
				switch message.Version {
				case 1:
					{
						h.handleTaskCreatedV1(value)
					}
				case 2:
					{
						h.handleTaskCreatedV2(value)
					}
				}

			}
		}
	case topics.TasksLifecycle:
		switch key {
		case events.TaskAssigned:
			{
				h.handleTaskAssigned(value)
			}
		case events.TaskCompleted:
			{
				h.handleTaskCompleted(value)
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

const (
	AssignedTaskFeeMin     = 10
	AssignedTaskFeeMax     = 20
	CompletedTaskAmountMin = 20
	CompletedTaskAmountMax = 40
)

func (h *KafkaHandlers) handleTaskCreatedV1(value string) {
	message := events.TaskCreatedPayload{}
	err := json.Unmarshal([]byte(value), &message)
	if err == nil {
		assignedTaskFee := rand.Intn(AssignedTaskFeeMax-AssignedTaskFeeMin) + AssignedTaskFeeMin
		completedTaskAmount := rand.Intn(CompletedTaskAmountMax-CompletedTaskAmountMin) + CompletedTaskAmountMin
		task := &model.Task{
			PublicId:            message.PublicId,
			Title:               message.Title,
			Description:         message.Description,
			AssignedTaskFee:     assignedTaskFee,
			CompletedTaskAmount: completedTaskAmount,
		}
		h.DB.Create(task)
	}
}

func (h *KafkaHandlers) handleTaskCreatedV2(value string) {
	message := events.TaskCreatedPayload{}
	err := json.Unmarshal([]byte(value), &message)
	if err == nil {
		assignedTaskFee := rand.Intn(AssignedTaskFeeMax-AssignedTaskFeeMin) + AssignedTaskFeeMin
		completedTaskAmount := rand.Intn(CompletedTaskAmountMax-CompletedTaskAmountMin) + CompletedTaskAmountMin
		task := &model.Task{
			PublicId:            message.PublicId,
			Title:               message.Title,
			Description:         message.Description,
			AssignedTaskFee:     assignedTaskFee,
			CompletedTaskAmount: completedTaskAmount,
		}
		h.DB.Create(task)
	}
}

func (h *KafkaHandlers) handleTaskAssigned(value string) {
	message := events.TaskAssignedPayload{}
	err := json.Unmarshal([]byte(value), &message)
	if err == nil {
		var account model.Account
		h.DB.Where(&model.Account{PublicId: message.AssigneeId}).First(&account)

		var task model.Task
		h.DB.Where(&model.Task{PublicId: message.PublicId}).First(&task)

		// TODO Fin Transaction
		account.Balance -= task.AssignedTaskFee

		h.DB.Save(&account)
	}
}

func (h *KafkaHandlers) handleTaskCompleted(value string) {
	message := events.TaskCompletedPayload{}
	err := json.Unmarshal([]byte(value), &message)
	if err == nil {
		var task model.Task
		h.DB.Where(&model.Task{PublicId: message.PublicId}).First(&task)

		var account model.Account
		h.DB.Where(&model.Account{PublicId: message.AssigneeId}).First(&account)

		// TODO Fin Transaction
		account.Balance += task.CompletedTaskAmount

		h.DB.Save(&account)
	}
}
