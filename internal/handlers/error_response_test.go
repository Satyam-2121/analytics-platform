package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"analytics-platform/internal/middleware"
	"analytics-platform/internal/models"
	"analytics-platform/internal/repositories"
)

var storageFailure = errors.New("sensitive database connection detail")

type failingStore struct{}

func (failingStore) CreateEvent(models.Event) error                   { return storageFailure }
func (failingStore) GetTotalEvents(int) (uint64, error)               { return 0, storageFailure }
func (failingStore) GetTopPages(int) ([]repositories.PageStat, error) { return nil, storageFailure }

func TestStorageErrorsAreNotExposedToClients(t *testing.T) {
	events := NewEventHandler(failingStore{})
	analytics := NewAnalyticsHandler(failingStore{}, &stubAnalyticsCache{err: errors.New("cache miss")})
	tests := []struct {
		name, method, target, body string
		handler                    http.HandlerFunc
	}{
		{"event", http.MethodPost, "/events", `{"event":"view","project_id":1}`, events.CreateEvent},
		{"total", http.MethodGet, "/analytics/events?project_id=1", "", analytics.GetTotalEvents},
		{"pages", http.MethodGet, "/analytics/pages?project_id=1", "", analytics.GetTopPages},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.target, strings.NewReader(tc.body))
			req.Header.Set("X-API-Key", "project-one-key")
			res := httptest.NewRecorder()
			middleware.APIKeyMiddleware(testProjectLookup{}, tc.handler)(res, req)
			if res.Code != http.StatusInternalServerError || strings.TrimSpace(res.Body.String()) != "internal server error" {
				t.Fatalf("got status %d and body %q", res.Code, res.Body.String())
			}
		})
	}
}
