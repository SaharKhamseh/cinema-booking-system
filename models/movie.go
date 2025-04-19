package models

import (
	"time"
)

type Movie struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	Duration    int       `json:"duration" gorm:"not null"` // in minutes
	Genre       string    `json:"genre"`
	Language    string    `json:"language"`
	ReleaseDate time.Time `json:"release_date"`
	PosterURL   string    `json:"poster_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
