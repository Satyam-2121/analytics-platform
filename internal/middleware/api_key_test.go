package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"analytics-platform/internal/models"
	"github.com/jackc/pgx/v5"
)

type lookupResult struct{ err error }

func (l lookupResult) GetByAPIKey(string) (*models.Project, error) { return nil, l.err }

func TestLookupFailuresReturnAppropriateStatus(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		status int
	}{
		{"unknown key", fmt.Errorf("lookup: %w", pgx.ErrNoRows), http.StatusUnauthorized},
		{"database unavailable", errors.New("private database detail"), http.StatusServiceUnavailable},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			h := APIKeyMiddleware(lookupResult{err: tc.err}, func(http.ResponseWriter, *http.Request) { called = true })
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.Header.Set("X-API-Key", "some-key")
			w := httptest.NewRecorder()
			h(w, r)
			if w.Code != tc.status || called || strings.Contains(w.Body.String(), "private database detail") {
				t.Fatalf("status %d, next called %v, body %q", w.Code, called, w.Body.String())
			}
		})
	}
}
