package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"analytics-platform/internal/middleware"
	"analytics-platform/internal/repositories"
)

type analyticsReader interface {
	GetTotalEvents(projectID int) (uint64, error)
	GetTopPages(projectID int) ([]repositories.PageStat, error)
}

type analyticsCache interface {
	Get(key string) (string, error)
	Set(key, value string) error
}

type AnalyticsHandler struct {
	AnalyticsRepo analyticsReader
	CacheService  analyticsCache
}

func NewAnalyticsHandler(
	repo analyticsReader,
	cache analyticsCache,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		AnalyticsRepo: repo,
		CacheService:  cache,
	}
}

func (h *AnalyticsHandler) GetTotalEvents(
	w http.ResponseWriter,
	r *http.Request,
) {

	projectIDStr := r.URL.Query().Get("project_id")

	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		http.Error(w, "invalid project_id", http.StatusBadRequest)
		return
	}
	if authenticatedID, ok := middleware.ProjectID(r.Context()); !ok || projectID != authenticatedID {
		http.Error(w, "project does not match api key", http.StatusForbidden)
		return
	}

	cacheKey := fmt.Sprintf(
		"analytics:events:%d",
		projectID,
	)

	cached, err := h.CacheService.Get(cacheKey)

	if err == nil {
		if total, parseErr := strconv.ParseUint(cached, 10, 64); parseErr == nil {
			json.NewEncoder(w).Encode(map[string]uint64{"total_events": total})
			return
		}
	}

	total, err := h.AnalyticsRepo.GetTotalEvents(projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.CacheService.Set(
		cacheKey,
		strconv.FormatUint(total, 10),
	)

	json.NewEncoder(w).Encode(map[string]uint64{
		"total_events": total,
	})
}

func (h *AnalyticsHandler) GetTopPages(
	w http.ResponseWriter,
	r *http.Request,
) {

	projectIDStr := r.URL.Query().Get("project_id")

	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		http.Error(w, "invalid project_id", http.StatusBadRequest)
		return
	}
	if authenticatedID, ok := middleware.ProjectID(r.Context()); !ok || projectID != authenticatedID {
		http.Error(w, "project does not match api key", http.StatusForbidden)
		return
	}

	pages, err := h.AnalyticsRepo.GetTopPages(projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(pages)
}
