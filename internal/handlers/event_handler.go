package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"analytics-platform/internal/middleware"
	"analytics-platform/internal/models"
)

type eventWriter interface {
	CreateEvent(models.Event) error
}

type EventHandler struct {
	EventRepo eventWriter
}

func NewEventHandler(repo eventWriter) *EventHandler {
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

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		http.Error(w, "request body must contain a single JSON event", http.StatusBadRequest)
		return
	}

	projectID, ok := middleware.ProjectID(r.Context())
	if !ok || event.ProjectID != projectID {
		http.Error(w, "project does not match api key", http.StatusForbidden)
		return
	}
	if strings.TrimSpace(event.Event) == "" {
		http.Error(w, "event name is required", http.StatusBadRequest)
		return
	}

	err = h.EventRepo.CreateEvent(event)
	if err != nil {
		log.Printf("event storage failed: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "event created",
	})
}
