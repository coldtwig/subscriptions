package subscriptions

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type SubscriptionsHandler struct {
	SubscriptionsRepository *SubscriptionsRepository
}

type SubscriptionsHandlerDeps struct {
	SubscriptionsRepository *SubscriptionsRepository
}

func NewSubscriptionsHandler(router *http.ServeMux, deps SubscriptionsHandlerDeps) *SubscriptionsHandler {
	handler := &SubscriptionsHandler{
		SubscriptionsRepository: deps.SubscriptionsRepository,
	}

	router.HandleFunc("POST /subscriptions", handler.Create())
	router.HandleFunc("GET /subscriptions", handler.GetAll())
	router.HandleFunc("PUT /subscriptions/{id}", handler.Update())
	router.HandleFunc("DELETE /subscriptions/{id}", handler.Delete())

	return handler
}

func (handler *SubscriptionsHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body SubscriptionsCreateRequest
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		startDate, err := time.Parse("01-2006", body.StartDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var endDate *time.Time
		if body.EndDate != nil {
			t, err := time.Parse("01-2006", *body.EndDate)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			endDate = &t
		}

		subscriptionBody := Subscription{
			ServiceName: body.ServiceName,
			Price:       body.Price,
			UserID:      body.UserID,
			StartDate:   startDate,
			EndDate:     endDate,
		}

		subscription, err := handler.SubscriptionsRepository.Create(&subscriptionBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(subscription)
	}
}

func (handler *SubscriptionsHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subscriptions, err := handler.SubscriptionsRepository.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(subscriptions)
	}
}

func (handler *SubscriptionsHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body SubscriptionsUpdateRequest
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		startDate, err := time.Parse("01-2006", body.StartDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var endDate *time.Time
		if body.EndDate != nil {
			t, err := time.Parse("01-2006", *body.EndDate)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			endDate = &t
		}

		subscription, err := handler.SubscriptionsRepository.Update(&Subscription{
			Model: gorm.Model{
				ID: uint(id),
			},
			ServiceName: body.ServiceName,
			Price:       body.Price,
			StartDate:   startDate,
			EndDate:     endDate,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(subscription)
	}
}

func (handler *SubscriptionsHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = handler.SubscriptionsRepository.Delete(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	}
}
