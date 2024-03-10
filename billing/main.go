package main

import (
	"net/http"
	"time"

	"github.com/agabidullin/aTES/billing/config"
	"github.com/agabidullin/aTES/billing/db"
	"github.com/agabidullin/aTES/billing/kafka"
	"github.com/agabidullin/aTES/billing/oauth"
	"github.com/agabidullin/aTES/billing/router"
	"github.com/joho/godotenv"

	log "github.com/go-pkgz/lgr"
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
	service := oauth.Init()

	// setup http server
	router := router.Init(service, database)

	kafkaHandlers := kafka.KafkaHandlers{DB: database}

	go kafka.InitConsumer(kafkaHandlers.InitHandler)

	httpServer := &http.Server{
		Addr:              ":8083",
		ReadHeaderTimeout: 5 * time.Second,
		Handler:           router,
	}

	if err := httpServer.ListenAndServe(); err != nil {
		log.Printf("[PANIC] failed to start http server, %v", err)
	}
}
