package repositories

import (
	"context"
	"testing"
	"time"

	"analytics-platform/internal/models"
)

type fakeEventStore struct {
	args []any
}

func (s *fakeEventStore) Exec(_ context.Context, _ string, args ...any) error {
	s.args = args
	return nil
}

func TestCreateEventPreservesClientTimestamp(t *testing.T) {
	store := &fakeEventStore{}
	repo := &EventRepository{CH: store}
	want := time.Date(2026, time.November, 7, 9, 30, 0, 0, time.UTC)

	err := repo.CreateEvent(models.Event{
		Event: "page_view", Page: "/pricing", UserID: "u123",
		ProjectID: 1, Timestamp: want,
	})
	if err != nil {
		t.Fatalf("CreateEvent returned %v", err)
	}
	if got := store.args[4].(time.Time); !got.Equal(want) {
		t.Fatalf("stored timestamp %v, want %v", got, want)
	}
}

func TestCreateEventDefaultsMissingTimestamp(t *testing.T) {
	store := &fakeEventStore{}
	repo := &EventRepository{CH: store}
	before := time.Now()

	if err := repo.CreateEvent(models.Event{ProjectID: 1}); err != nil {
		t.Fatalf("CreateEvent returned %v", err)
	}
	after := time.Now()
	got := store.args[4].(time.Time)
	if got.Before(before) || got.After(after) {
		t.Fatalf("default timestamp %v is outside [%v, %v]", got, before, after)
	}
}
