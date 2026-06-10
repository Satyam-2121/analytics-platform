package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"analytics-platform/internal/repositories"
	"analytics-platform/internal/services"
)

type AnalyticsHandler struct {
	AnalyticsRepo *repositories.AnalyticsRepository
	CacheService  *services.CacheService
}

func NewAnalyticsHandler(
	repo *repositories.AnalyticsRepository,
	cache *services.CacheService,
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

	cacheKey := fmt.Sprintf(
		"analytics:events:%d",
		projectID,
	)

	cached, err := h.CacheService.Get(cacheKey)

	if err == nil {
		json.NewEncoder(w).Encode(map[string]string{
			"source": "redis",
			"total":  cached,
		})
		return
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

	pages, err := h.AnalyticsRepo.GetTopPages(projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(pages)
}