package models

import "time"

type Event struct {
	Event     string    `json:"event"`
	Page      string    `json:"page"`
	UserID    string    `json:"user_id"`
	ProjectID int       `json:"project_id"`
	Timestamp time.Time `json:"timestamp"`
}