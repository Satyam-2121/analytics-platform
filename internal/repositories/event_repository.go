package repositories

import (
	"context"
	"time"

	"analytics-platform/internal/database"
	"analytics-platform/internal/models"
)

type EventRepository struct {
	CH eventStore
}

type eventStore interface {
	Exec(context.Context, string, ...any) error
}

func NewEventRepository(db *database.Database) *EventRepository {
	return &EventRepository{
		CH: db.CH,
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

	timestamp := event.Timestamp
	if timestamp.IsZero() {
		timestamp = time.Now()
	}

	return r.CH.Exec(
		context.Background(),
		query,
		event.Event,
		event.Page,
		event.UserID,
		event.ProjectID,
		timestamp,
	)
}
