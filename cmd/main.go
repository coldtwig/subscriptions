package main

import (
	"log"
	"net/http"
	"subscriptions/configs"
	_ "subscriptions/docs"
	"subscriptions/internal/db"
	"subscriptions/internal/subscriptions"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title Subscriptions API
// @version 1.0
// @description REST service for subscriptions aggregation.
// @BasePath /
// @schemes http
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
	router.Handle("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

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
