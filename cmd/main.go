package main

import (
	"log"
	"net/http"
	"subscriptions/configs"
	"subscriptions/internal/db"
	"subscriptions/internal/subscriptions"
)

func main() {
	log.Println("INFO service starting")

	router := http.NewServeMux()
	config := configs.LoadConfig()
	database := db.NewDB(config)
	log.Println("INFO database connected")

	// repositories
	subscriptionsRepository := subscriptions.NewSubscriptionsRepository(database)

	// services
	subscriptionsService := subscriptions.NewSubscriptionsService(subscriptionsRepository)

	// handlers
	subscriptions.NewSubscriptionsHandler(router, subscriptions.SubscriptionsHandlerDeps{
		SubscriptionsRepository: subscriptionsRepository,
		SubscriptionsService:    subscriptionsService,
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	log.Printf("INFO http server listening on %s", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Printf("ERROR http server failed: %v", err)
		panic(err)
	}
}
