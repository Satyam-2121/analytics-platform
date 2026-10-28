package middleware

import (
	"context"
	"net/http"

	"analytics-platform/internal/models"
)

type projectIDKey struct{}

type projectLookup interface {
	GetByAPIKey(apiKey string) (*models.Project, error)
}

// ProjectID returns the project authenticated by APIKeyMiddleware.
func ProjectID(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(projectIDKey{}).(int)
	return id, ok
}

func APIKeyMiddleware(
	projectRepo projectLookup,
	next http.HandlerFunc,
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get("X-API-Key")

		if apiKey == "" {
			http.Error(w, "missing api key", http.StatusUnauthorized)
			return
		}

		project, err := projectRepo.GetByAPIKey(apiKey)
		if err != nil {
			http.Error(w, "invalid api key", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), projectIDKey{}, project.ID)
		next(w, r.WithContext(ctx))
	}
}
