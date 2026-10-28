package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"analytics-platform/internal/middleware"
	"analytics-platform/internal/models"
)

type testProjectLookup struct{}

func (testProjectLookup) GetByAPIKey(key string) (*models.Project, error) {
	if key != "project-one-key" {
		return nil, errors.New("unknown key")
	}
	return &models.Project{ID: 1}, nil
}

func TestProjectAuthorization(t *testing.T) {
	events := NewEventHandler(nil)
	analytics := NewAnalyticsHandler(nil, nil)
	tests := []struct {
		name, method, target, body string
		handler                    http.HandlerFunc
	}{
		{"event", http.MethodPost, "/events", `{"event":"page_view","project_id":2}`, events.CreateEvent},
		{"total", http.MethodGet, "/analytics/events?project_id=2", "", analytics.GetTotalEvents},
		{"pages", http.MethodGet, "/analytics/pages?project_id=2", "", analytics.GetTopPages},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.target, strings.NewReader(tc.body))
			req.Header.Set("X-API-Key", "project-one-key")
			res := httptest.NewRecorder()
			middleware.APIKeyMiddleware(testProjectLookup{}, tc.handler)(res, req)
			if res.Code != http.StatusForbidden {
				t.Fatalf("got status %d, want %d", res.Code, http.StatusForbidden)
			}
		})
	}
}

func TestAPIKeyMiddlewarePassesAuthenticatedProject(t *testing.T) {
	handler := middleware.APIKeyMiddleware(testProjectLookup{}, func(w http.ResponseWriter, r *http.Request) {
		id, ok := middleware.ProjectID(r.Context())
		if !ok || id != 1 {
			t.Errorf("got project ID %d, present %v; want 1, true", id, ok)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-API-Key", "project-one-key")
	res := httptest.NewRecorder()
	handler(res, req)
	if res.Code != http.StatusNoContent {
		t.Fatalf("got status %d, want %d", res.Code, http.StatusNoContent)
	}
}
