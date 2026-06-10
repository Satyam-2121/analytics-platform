package repositories

import (
	"context"

	"analytics-platform/internal/database"
)

type PageStat struct {
	Page  string `json:"page"`
	Count uint64 `json:"count"`
}

type AnalyticsRepository struct {
	DB *database.Database
}

func NewAnalyticsRepository(db *database.Database) *AnalyticsRepository {
	return &AnalyticsRepository{
		DB: db,
	}
}

func (r *AnalyticsRepository) GetTotalEvents(projectID int) (uint64, error) {
	var count uint64

	err := r.DB.CH.QueryRow(
		context.Background(),
		`
		SELECT count(*)
		FROM events
		WHERE project_id = ?
		`,
		projectID,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *AnalyticsRepository) GetTopPages(projectID int) ([]PageStat, error) {

	rows, err := r.DB.CH.Query(
		context.Background(),
		`
		SELECT page, count(*)
		FROM events
		WHERE project_id = ?
		GROUP BY page
		ORDER BY count(*) DESC
		`,
		projectID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var stats []PageStat

	for rows.Next() {

		var stat PageStat

		err := rows.Scan(&stat.Page, &stat.Count)
		if err != nil {
			return nil, err
		}

		stats = append(stats, stat)
	}

	return stats, nil
}
