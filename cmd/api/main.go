package main

import (
	"fmt"
	"log"
	"time"
    
	"analytics-platform/internal/server"
	"analytics-platform/internal/database"
	"analytics-platform/internal/models"
	"analytics-platform/internal/repositories"
)

func main() {

	postgres, err := database.ConnectPostgres()
	if err != nil {
		log.Fatal(err)
	}

	redisClient, err := database.ConnectRedis()
	if err != nil {
		log.Fatal(err)
	}

	clickhouseConn, err := database.ConnectClickHouse()
	if err != nil {
		log.Fatal(err)
	}

	db := &database.Database{
		Postgres: postgres,
		Redis:    redisClient,
		CH:       clickhouseConn,
	}

	userRepo := repositories.NewUserRepository(db)

	email := fmt.Sprintf("test%d@example.com", time.Now().Unix())

	user, err := userRepo.CreateUser(email)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("User created: %+v\n", user)

	projectRepo := repositories.NewProjectRepository(db)

	project, err := projectRepo.CreateProject(
		"My Analytics App",
		user.ID,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Project created: %+v\n", project)

	eventRepo := repositories.NewEventRepository(db)

	err = eventRepo.CreateEvent(models.Event{
		Event:     "page_view",
		Page:      "/pricing",
		UserID:    "u123",
		ProjectID: project.ID,
		Timestamp: time.Now(),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Event inserted successfully")

	analyticsRepo := repositories.NewAnalyticsRepository(db)

	total, err := analyticsRepo.GetTotalEvents(project.ID)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Total Events: %d\n", total)

	_, err = analyticsRepo.GetTopPages(project.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("All services connected successfully")

	server.Start(db)
}