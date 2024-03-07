package main

import (
	"net/http"
	"time"

	"github.com/agabidullin/aTES/tasks/config"
	"github.com/agabidullin/aTES/tasks/db"
	"github.com/agabidullin/aTES/tasks/kafka"
	"github.com/agabidullin/aTES/tasks/oauth"
	"github.com/agabidullin/aTES/tasks/router"
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

	kafkaHandlers := kafka.KafkaHandlers{DB: database}

	go kafka.InitConsumer(kafkaHandlers.InitHandler)

	producer := kafka.InitProducer()
	defer producer.Close()

	// setup http server
	router := router.Init(service, database, producer)

	httpServer := &http.Server{
		Addr:              ":8082",
		ReadHeaderTimeout: 5 * time.Second,
		Handler:           router,
	}

	if err := httpServer.ListenAndServe(); err != nil {
		log.Printf("[PANIC] failed to start http server, %v", err)
	}
}
