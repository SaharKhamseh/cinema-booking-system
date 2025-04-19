package models

import (
	"time"
)

type ShowTime struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	MovieID   uint      `json:"movie_id"`
	Movie     Movie     `json:"movie" gorm:"foreignKey:MovieID"`
	ScreenID  uint      `json:"screen_id"`
	Screen    Screen    `json:"screen" gorm:"foreignKey:ScreenID"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Price     float64   `json:"price"`
}
