package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"analytics-platform/internal/middleware"
	"analytics-platform/internal/models"
)

type stubEventWriter struct {
	events []models.Event
}

func (s *stubEventWriter) CreateEvent(event models.Event) error {
	s.events = append(s.events, event)
	return nil
}

func TestCreateEventRequiresSingleJSONValue(t *testing.T) {
	const event = `{"event":"page_view","page":"/pricing","user_id":"u123","project_id":1}`
	tests := []struct {
		name, body string
		wantStatus int
	}{
		{"single event", event, http.StatusCreated},
		{"trailing whitespace", event + " \n\t", http.StatusCreated},
		{"second event", event + "\n" + event, http.StatusBadRequest},
		{"trailing null", event + " null", http.StatusBadRequest},
		{"trailing garbage", event + " garbage", http.StatusBadRequest},
		{"malformed first event", `{"event":`, http.StatusBadRequest},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &stubEventWriter{}
			handler := NewEventHandler(repo)
			req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(tc.body))
			req.Header.Set("X-API-Key", "project-one-key")
			res := httptest.NewRecorder()
			middleware.APIKeyMiddleware(testProjectLookup{}, handler.CreateEvent)(res, req)
			if res.Code != tc.wantStatus {
				t.Fatalf("got status %d, want %d: %s", res.Code, tc.wantStatus, res.Body.String())
			}
			if tc.wantStatus == http.StatusCreated {
				if got := res.Header().Get("Content-Type"); got != "application/json" {
					t.Errorf("got Content-Type %q, want application/json", got)
				}
				want := models.Event{Event: "page_view", Page: "/pricing", UserID: "u123", ProjectID: 1}
				if len(repo.events) != 1 || repo.events[0] != want {
					t.Fatalf("got persisted events %#v, want one %#v", repo.events, want)
				}
			} else if len(repo.events) != 0 {
				t.Fatalf("invalid body persisted %d event(s)", len(repo.events))
			}
		})
	}
}
