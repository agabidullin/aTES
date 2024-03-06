package main

import (
	"net/http"

	"github.com/agabidullin/aTES/auth/config"
	"github.com/agabidullin/aTES/auth/db"
	"github.com/agabidullin/aTES/auth/handlers"
	"github.com/agabidullin/aTES/auth/kafka"
	"github.com/agabidullin/aTES/auth/oauth"

	"github.com/go-chi/chi"

	"github.com/joho/godotenv"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}
}

func main() {
	conf := config.New()

	database := db.Init(conf.DSN)

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
