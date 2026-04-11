package subscriptions

import "net/http"

type SubscriptionsHandler struct{}

func NewSubscriptionsHandler(router *http.ServeMux) *SubscriptionsHandler {
	handler := &SubscriptionsHandler{}

	router.HandleFunc("/create", GetSubscriptions())
	router.HandleFunc("/update", GetSubscriptions())
	router.HandleFunc("/delete", GetSubscriptions())
	router.HandleFunc("/subscriptions", GetSubscriptions())

	return handler
}

func GetSubscriptions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
func CreateSubscription() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
func DeleteSubscription() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
func UpdateSubscription() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
