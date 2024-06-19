package app

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/uchitsa/ubiquitous-octo-train/internal/domain"
	"github.com/uchitsa/ubiquitous-octo-train/internal/infrastructure"
)

type OrderService struct {
	orders       []domain.Order
	availability []domain.RoomAvailability
}

func NewOrderService() *OrderService {
	return &OrderService{
		orders: []domain.Order{},
		availability: []domain.RoomAvailability{
			{"reddison", "lux", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), 1},
			{"reddison", "lux", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), 1},
			{"reddison", "lux", time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC), 1},
			{"reddison", "lux", time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC), 1},
			{"reddison", "lux", time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC), 0},
		},
	}
}

func (s *OrderService) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder domain.Order
	err := json.NewDecoder(r.Body).Decode(&newOrder)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	daysToBook := daysBetween(newOrder.From, newOrder.To)
	if len(daysToBook) == 0 {
		http.Error(w, "Invalid date range", http.StatusBadRequest)
		return
	}

	unavailableDays := make(map[time.Time]struct{})
	for _, day := range daysToBook {
		unavailableDays[day] = struct{}{}
	}

	for _, dayToBook := range daysToBook {
		for i, availability := range s.availability {
			if availability.Date.Equal(dayToBook) && availability.Quota > 0 {
				availability.Quota -= 1
				s.availability[i] = availability
				delete(unavailableDays, dayToBook)
			}
		}
	}

	if len(unavailableDays) != 0 {
		http.Error(w, "Hotel room is not available for selected dates", http.StatusConflict)
		infrastructure.LogErrorf("Hotel room is not available for selected dates: %v", newOrder)
		return
	}

	s.orders = append(s.orders, newOrder)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newOrder)

	infrastructure.LogInfo("Order successfully created: %v", newOrder)
}

func daysBetween(from time.Time, to time.Time) []time.Time {
	if from.After(to) {
		return nil
	}

	days := make([]time.Time, 0)
	for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
		days = append(days, d)
	}

	return days
}
