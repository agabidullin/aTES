package main

import (
	"net/http"
	"time"

	"github.com/agabidullin/aTES/tasks/db"
	"github.com/agabidullin/aTES/tasks/kafka"
	"github.com/agabidullin/aTES/tasks/oauth"
	"github.com/agabidullin/aTES/tasks/router"

	log "github.com/go-pkgz/lgr"
)

func main() {
	database := db.Init()
	service := oauth.Init()

	// setup http server
	router := router.Init(service)
	kafkaHandlers := kafka.KafkaHandlers{DB: database}

	go kafka.Init(kafkaHandlers.InitHandler)

	httpServer := &http.Server{
		Addr:              ":8082",
		ReadHeaderTimeout: 5 * time.Second,
		Handler:           router,
	}

	if err := httpServer.ListenAndServe(); err != nil {
		log.Printf("[PANIC] failed to start http server, %v", err)
	}
}
