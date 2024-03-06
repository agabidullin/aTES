package handlers

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"

	"github.com/agabidullin/aTES/tasks/model"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TasksHandler struct {
	DB *gorm.DB
}

type TaskCreateInput struct {
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

	task := &model.Task{Description: input.Description, AssigneeId: input.AssigneeId}
	h.DB.Omit(clause.Associations).Create(task)
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}

func (h TasksHandler) Shuffle(w http.ResponseWriter, r *http.Request) {
	h.DB.Transaction(func(tx *gorm.DB) error {
		var openedTasks []model.Task
		h.DB.Where("is_closed = ?", false).Find(&openedTasks)

		var accountsToAssign []model.Account
		queryResult := h.DB.Where("role != ? AND role != ?", "admin", "manager").Find(&accountsToAssign)

		for _, t := range openedTasks {
			t.Assignee = accountsToAssign[rand.Intn(int(queryResult.RowsAffected))]
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

	err := json.NewEncoder(w).Encode(task)

	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}
