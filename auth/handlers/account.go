package handlers

import (
	"encoding/json"

	"github.com/agabidullin/aTES/auth/model"

	"net/http"

	"github.com/go-chi/render"
	"gorm.io/gorm"

	"github.com/agabidullin/aTES/common/events"
	"github.com/agabidullin/aTES/common/topics"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type AccountHandler struct {
	DB       *gorm.DB
	Producer *kafka.Producer
}

type RegisterAccountRequest struct {
	*model.Account
}

func (a *RegisterAccountRequest) Bind(r *http.Request) error {
	return nil
}

type RegisterAccountResponse struct {
	*model.Account
}

func (a *RegisterAccountResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

func (h AccountHandler) RegisterAccount(w http.ResponseWriter, r *http.Request) {
	data := &RegisterAccountRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	account := data.Account
	h.DB.Create(account)
	topic := topics.AccountsStream

	message := events.AccountCreatedPayload{PublicId: account.ID, Name: account.Name, Role: account.Role}
	ser, _ := json.Marshal(&message)

	h.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(events.AccountCreated),
		Value:          []byte(ser),
	}, nil)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &RegisterAccountResponse{Account: account})
}

type ChangeRoleRequest struct {
	ID   uint
	Role string
}

func (a *ChangeRoleRequest) Bind(r *http.Request) error {
	return nil
}

type ChangeRoleResponse struct {
	*model.Account
}

func (a *ChangeRoleResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

func (h AccountHandler) ChangeRole(w http.ResponseWriter, r *http.Request) {
	data := &ChangeRoleRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	var account model.Account
	h.DB.First(&account, data.ID)
	h.DB.Model(&account).Update("Role", data.Role)

	topic := topics.Accounts

	message := events.RoleChangedPayload{PublicId: account.ID, Role: account.Role}
	ser, _ := json.Marshal(&message)

	h.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(events.RoleChanged),
		Value:          []byte(ser),
	}, nil)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &RegisterAccountResponse{Account: &account})
}
