package models

import (
	"time"
)

type Booking struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id"`
	User       User      `json:"user" gorm:"foreignKey:UserID"`
	ShowTimeID uint      `json:"show_time_id"`
	ShowTime   ShowTime  `json:"show_time" gorm:"foreignKey:ShowTimeID"`
	Seats      []Seat    `json:"seats" gorm:"many2many:booking_seats;"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"` // confirmed, cancelled, pending
	BookedAt   time.Time `json:"booked_at"`
}
