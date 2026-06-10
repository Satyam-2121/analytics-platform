package repositories

import (
	"context"

	"analytics-platform/internal/database"
	"analytics-platform/internal/models"
)

type UserRepository struct {
	DB *database.Database
}

func NewUserRepository(db *database.Database) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) CreateUser(email string) (*models.User, error) {

	query := `
	INSERT INTO users(email)
	VALUES($1)
	RETURNING id, email, created_at
	`

	var user models.User

	err := r.DB.Postgres.QueryRow(
		context.Background(),
		query,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
