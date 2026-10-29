package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"analytics-platform/internal/middleware"
	"analytics-platform/internal/repositories"
)

type stubAnalytics struct {
	total uint64
	calls int
}

func (s *stubAnalytics) GetTotalEvents(int) (uint64, error) {
	s.calls++
	return s.total, nil
}

func (s *stubAnalytics) GetTopPages(int) ([]repositories.PageStat, error) {
	return nil, nil
}

type stubAnalyticsCache struct {
	value            string
	err              error
	setKey, setValue string
}

func (s *stubAnalyticsCache) Get(string) (string, error) { return s.value, s.err }
func (s *stubAnalyticsCache) Set(key, value string) error {
	s.setKey, s.setValue = key, value
	return nil
}

func TestTotalEventsResponseIsConsistentAcrossCacheStates(t *testing.T) {
	tests := []struct {
		name, cached  string
		cacheErr      error
		wantRepoCalls int
	}{
		{"cache hit", "42", nil, 0},
		{"cache miss", "", errors.New("cache miss"), 1},
		{"invalid cache value", "not-a-number", nil, 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &stubAnalytics{total: 42}
			cache := &stubAnalyticsCache{value: tc.cached, err: tc.cacheErr}
			h := NewAnalyticsHandler(repo, cache)
			req := httptest.NewRequest(http.MethodGet, "/analytics/events?project_id=1", nil)
			req.Header.Set("X-API-Key", "project-one-key")
			res := httptest.NewRecorder()
			middleware.APIKeyMiddleware(testProjectLookup{}, h.GetTotalEvents)(res, req)
			if res.Code != http.StatusOK {
				t.Fatalf("got status %d, want 200", res.Code)
			}
			var body map[string]interface{}
			if err := json.Unmarshal(res.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if len(body) != 1 || body["total_events"] != float64(42) {
				t.Errorf("got response %v, want {total_events: 42}", body)
			}
			if repo.calls != tc.wantRepoCalls {
				t.Errorf("got %d repository calls, want %d", repo.calls, tc.wantRepoCalls)
			}
			if tc.wantRepoCalls > 0 && (cache.setKey != "analytics:events:1" || cache.setValue != "42") {
				t.Errorf("got cache write %q=%q, want analytics:events:1=42", cache.setKey, cache.setValue)
			}
		})
	}
}
