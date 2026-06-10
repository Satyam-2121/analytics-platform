package handlers

import (
	"encoding/json"
	"net/http"

	"analytics-platform/internal/models"
	"analytics-platform/internal/repositories"
)

type EventHandler struct {
	EventRepo *repositories.EventRepository
}

func NewEventHandler(repo *repositories.EventRepository) *EventHandler {
	return &EventHandler{
		EventRepo: repo,
	}
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var event models.Event

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.EventRepo.CreateEvent(event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "event created",
	})
}