package repositories

import (
	"context"

	"analytics-platform/internal/database"
	"analytics-platform/internal/models"
)

type ProjectRepository struct {
	DB *database.Database
}

func NewProjectRepository(db *database.Database) *ProjectRepository {
	return &ProjectRepository{
		DB: db,
	}
}

func (r *ProjectRepository) CreateProject(name string, userID int) (*models.Project, error) {

	query := `
	INSERT INTO projects(name, user_id, api_key)
	VALUES($1, $2, gen_random_uuid())
	RETURNING id, name, api_key, user_id, created_at
	`

	var project models.Project

	err := r.DB.Postgres.QueryRow(
		context.Background(),
		query,
		name,
		userID,
	).Scan(
		&project.ID,
		&project.Name,
		&project.APIKey,
		&project.UserID,
		&project.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *ProjectRepository) GetByAPIKey(
	apiKey string,
) (*models.Project, error) {

	query := `
	SELECT id, name, api_key, user_id, created_at
	FROM projects
	WHERE api_key = $1
	`

	var project models.Project

	err := r.DB.Postgres.QueryRow(
		context.Background(),
		query,
		apiKey,
	).Scan(
		&project.ID,
		&project.Name,
		&project.APIKey,
		&project.UserID,
		&project.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &project, nil
}
