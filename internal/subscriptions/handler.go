package subscriptions

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type SubscriptionsHandler struct {
	SubscriptionsRepository *SubscriptionsRepository
	SubscriptionsService    *SubscriptionsService
}

type SubscriptionsHandlerDeps struct {
	SubscriptionsRepository *SubscriptionsRepository
	SubscriptionsService    *SubscriptionsService
}

func NewSubscriptionsHandler(router *http.ServeMux, deps SubscriptionsHandlerDeps) *SubscriptionsHandler {
	handler := &SubscriptionsHandler{
		SubscriptionsRepository: deps.SubscriptionsRepository,
		SubscriptionsService:    deps.SubscriptionsService,
	}

	router.HandleFunc("POST /subscriptions", handler.Create())
	router.HandleFunc("GET /subscriptions", handler.GetAll())
	router.HandleFunc("PUT /subscriptions/{id}", handler.Update())
	router.HandleFunc("DELETE /subscriptions/{id}", handler.Delete())
	router.HandleFunc("POST /subscriptions/total", handler.SumAll())

	return handler
}

// Create godoc
// @Summary Create subscription
// @Description Creates a new subscription record. Dates are strings in MM-YYYY format.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body SubscriptionsCreateRequest true "Create subscription payload"
// @Success 201 {object} SubscriptionResponse
// @Failure 400 {string} string
// @Router /subscriptions [post]
func (handler *SubscriptionsHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body SubscriptionsCreateRequest
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			log.Printf("WARN invalid total request body: %v", err)
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

// GetAll godoc
// @Summary List subscriptions
// @Description Returns all subscriptions.
// @Tags subscriptions
// @Produce json
// @Success 200 {array} SubscriptionResponse
// @Failure 400 {string} string
// @Router /subscriptions [get]
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

// Update godoc
// @Summary Update subscription
// @Description Updates an existing subscription by ID. Dates are strings in MM-YYYY format.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Param request body SubscriptionsUpdateRequest true "Update subscription payload"
// @Success 200 {object} SubscriptionResponse
// @Failure 400 {string} string
// @Router /subscriptions/{id} [put]
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

// Delete godoc
// @Summary Delete subscription
// @Description Deletes a subscription by ID.
// @Tags subscriptions
// @Param id path int true "Subscription ID"
// @Success 204 {string} string
// @Failure 400 {string} string
// @Router /subscriptions/{id} [delete]
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

// SumAll godoc
// @Summary Calculate total cost
// @Description Calculates total subscription cost for a selected period with optional user_id and service_name filters. Fields from, to, start_date, end_date use MM-YYYY format.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body SubscriptionTotalRequest true "Total calculation filter"
// @Success 200 {integer} int
// @Failure 400 {string} string
// @Router /subscriptions/total [post]
func (handler *SubscriptionsHandler) SumAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body SubscriptionTotalRequest
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		from, err := time.Parse("01-2006", body.From)
		if err != nil {
			log.Printf("WARN invalid total from format: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		to, err := time.Parse("01-2006", body.To)
		if err != nil {
			log.Printf("WARN invalid total to format: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if from.After(to) {
			log.Printf("WARN invalid total range from=%s to=%s", body.From, body.To)
			http.Error(w, "from must be before or equal to to", http.StatusBadRequest)
			return
		}

		sum, err := handler.SubscriptionsService.SumAll(&SubscriptionTotalFilter{
			ServiceName: body.ServiceName,
			UserID:      body.UserID,
			From:        from,
			To:          to,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(sum)
	}
}
