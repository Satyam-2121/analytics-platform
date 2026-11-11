package repositories

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
)

func TestMalformedAPIKeyIsRejectedBeforeDatabaseLookup(t *testing.T) {
	repo := &ProjectRepository{}
	project, err := repo.GetByAPIKey("not-a-uuid")
	if project != nil || !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("got project %v, error %v; want no matching key", project, err)
	}
}
