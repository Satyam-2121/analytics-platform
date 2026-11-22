package middleware

import (
	"context"
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

func (l lookupResult) GetByAPIKey(context.Context, string) (*models.Project, error) {
	return nil, l.err
}

type contextLookup struct {
	ctx context.Context
}

func (l *contextLookup) GetByAPIKey(ctx context.Context, _ string) (*models.Project, error) {
	l.ctx = ctx
	return nil, ctx.Err()
}

func TestAPIKeyLookupReceivesRequestCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
	r.Header.Set("X-API-Key", "some-key")
	lookup := &contextLookup{}
	h := APIKeyMiddleware(lookup, func(http.ResponseWriter, *http.Request) {
		t.Fatal("request continued after the lookup was canceled")
	})
	h(httptest.NewRecorder(), r)
	if lookup.ctx == nil || !errors.Is(lookup.ctx.Err(), context.Canceled) {
		t.Fatalf("lookup did not receive the canceled request context")
	}
}

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
