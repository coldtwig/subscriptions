package main

import (
	"net/http"
	"subscriptions/configs"
	"subscriptions/internal/subscriptions"
)

func main() {
	router := http.NewServeMux()
	_ = configs.LoadConfig()

	// handlers
	subscriptions.NewSubscriptionsHandler(router)

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	server.ListenAndServe()
}
