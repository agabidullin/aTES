package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/agabidullin/aTES/common/events"
	"github.com/agabidullin/aTES/common/topics"
	"github.com/agabidullin/aTES/schemaregistry/validator"
	"github.com/agabidullin/aTES/tasks/model"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TasksHandler struct {
	DB       *gorm.DB
	Producer *kafka.Producer
}

type TaskCreateInput struct {
	Title       string
	Description string
	AssigneeId  uint
}

func (h TasksHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var input TaskCreateInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task := &model.Task{Title: input.Title, Description: input.Description, AssigneeId: input.AssigneeId}
	h.DB.Omit(clause.Associations).Create(task)
	err = json.NewEncoder(w).Encode(task)

	h.produceTaskCreatedEvent(task)
	h.produceTaskAssignedEvent(task)

	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}

func (h TasksHandler) produceTaskCreatedEvent(task *model.Task) {
	topic := topics.TasksStream
	event := events.TaskCreated
	message := events.TaskCreatedPayload{Event: &events.Event{Version: 1}, PublicId: task.ID, AssigneeId: task.AssigneeId, Title: task.Title, Description: task.Description}

	payload, err := validator.Validate(message, topic, event, 1)

	if err != nil {
		fmt.Print(err.Error())
		return
	}

	h.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(event),
		Value:          payload,
	}, nil)
}

func (h TasksHandler) produceTaskAssignedEvent(task *model.Task) {
	topic := topics.TasksLifecycle
	message := events.TaskAssignedPayload{PublicId: task.ID, AssigneeId: task.AssigneeId}
	ser, _ := json.Marshal(&message)

	h.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(events.TaskAssigned),
		Value:          []byte(ser),
	}, nil)
}

func (h TasksHandler) Shuffle(w http.ResponseWriter, r *http.Request) {
	h.DB.Transaction(func(tx *gorm.DB) error {
		var openedTasks []model.Task
		h.DB.Where("is_closed = ?", false).Find(&openedTasks)

		var accountsToAssign []model.Account
		queryResult := h.DB.Where("role != ? AND role != ?", "admin", "manager").Find(&accountsToAssign)

		for _, t := range openedTasks {
			t.AssigneeId = accountsToAssign[rand.Intn(int(queryResult.RowsAffected))].PublicId
			h.produceTaskAssignedEvent(&t)
			h.DB.Save(&t)
		}

		// return nil will commit the whole transaction
		return nil
	})
}

func (h TasksHandler) CompleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var task model.Task
	result := h.DB.Where("id = ? AND is_closed = ?", id, false).First(&task)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		http.Error(w, "Not found", http.StatusNotFound)
	}

	task.IsClosed = true
	h.DB.Save(&task)
	h.produceTaskCompletedEvent(&task)
	err := json.NewEncoder(w).Encode(task)

	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}

func (h TasksHandler) produceTaskCompletedEvent(task *model.Task) {
	topic := topics.TasksLifecycle
	message := events.TaskCompletedPayload{PublicId: task.ID, AssigneeId: task.AssigneeId}
	ser, _ := json.Marshal(&message)

	h.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(events.TaskCompleted),
		Value:          []byte(ser),
	}, nil)
}
