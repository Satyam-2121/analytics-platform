package models

import "time"

type Project struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	APIKey    string    `json:"api_key"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}