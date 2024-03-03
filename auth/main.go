package main

import (
	"net/http"

	"github.com/agabidullin/aTES/auth/db"
	"github.com/agabidullin/aTES/auth/handlers"
	"github.com/agabidullin/aTES/auth/kafka"
	"github.com/agabidullin/aTES/auth/oauth"

	"github.com/go-chi/chi"
)

func main() {
	database := db.Init()

	oauth.Init(database)

	producer := kafka.Init()
	defer producer.Close()

	r := chi.NewRouter()

	accountHandler := handlers.AccountHandler{DB: database, Producer: producer}

	r.Route("/accounts", func(r chi.Router) {
		r.Post("/", accountHandler.RegisterAccount)      // POST /accounts
		r.Post("/changeRole", accountHandler.ChangeRole) // POST /accounts/changeRole
	})

	http.ListenAndServe(":8081", r)
}
