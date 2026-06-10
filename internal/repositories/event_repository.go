package repositories

import (
	"context"
	"time"

	"analytics-platform/internal/database"
	"analytics-platform/internal/models"
)

type EventRepository struct {
	DB *database.Database
}

func NewEventRepository(db *database.Database) *EventRepository {
	return &EventRepository{
		DB: db,
	}
}

func (r *EventRepository) CreateEvent(event models.Event) error {

	query := `
	INSERT INTO events
	(
		event,
		page,
		user_id,
		project_id,
		timestamp
	)
	VALUES
	(
		?,
		?,
		?,
		?,
		?
	)
	`

	return r.DB.CH.Exec(
		context.Background(),
		query,
		event.Event,
		event.Page,
		event.UserID,
		event.ProjectID,
		time.Now(),
	)
}