package main

import (
	"net/http"
	"subscriptions/configs"
	"subscriptions/internal/db"
	"subscriptions/internal/subscriptions"
)

func main() {
	router := http.NewServeMux()
	config := configs.LoadConfig()
	database := db.NewDB(config)

	// repositories
	subscriptionsRepository := subscriptions.NewSubscriptionsRepository(database)

	// handlers
	subscriptions.NewSubscriptionsHandler(router, subscriptions.SubscriptionsHandlerDeps{
		SubscriptionsRepository: subscriptionsRepository,
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
