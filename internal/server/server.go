package server

import (
	"log"
	"net/http"
	"time"

	"analytics-platform/internal/database"
	"analytics-platform/internal/handlers"
	"analytics-platform/internal/middleware"
	"analytics-platform/internal/repositories"
	"analytics-platform/internal/services"
)

func Start(db *database.Database) {

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Repositories

	projectRepo := repositories.NewProjectRepository(db)
	eventRepo := repositories.NewEventRepository(db)
	analyticsRepo := repositories.NewAnalyticsRepository(db)

	// Services

	cacheService := services.NewCacheService(
		db.Redis,
	)

	// Handlers

	eventHandler := handlers.NewEventHandler(eventRepo)

	analyticsHandler := handlers.NewAnalyticsHandler(
		analyticsRepo,
		cacheService,
	)

	// Protected Routes

	mux.HandleFunc(
		"/events",
		middleware.APIKeyMiddleware(
			projectRepo,
			eventHandler.CreateEvent,
		),
	)

	mux.HandleFunc(
		"/analytics/events",
		middleware.APIKeyMiddleware(
			projectRepo,
			analyticsHandler.GetTotalEvents,
		),
	)

	mux.HandleFunc(
		"/analytics/pages",
		middleware.APIKeyMiddleware(
			projectRepo,
			analyticsHandler.GetTopPages,
		),
	)

	log.Println("REGISTERED: /events")
	log.Println("REGISTERED: /analytics/events")
	log.Println("REGISTERED: /analytics/pages")
	log.Println("Server running on :8080")

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
