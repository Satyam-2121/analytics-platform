package middleware

import (
	"net/http"

	"analytics-platform/internal/repositories"
)

func APIKeyMiddleware(
	projectRepo *repositories.ProjectRepository,
	next http.HandlerFunc,
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get("X-API-Key")

		if apiKey == "" {
			http.Error(w, "missing api key", http.StatusUnauthorized)
			return
		}

		_, err := projectRepo.GetByAPIKey(apiKey)
		if err != nil {
			http.Error(w, "invalid api key", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}