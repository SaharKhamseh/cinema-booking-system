package models

type Theater struct {
	ID       uint     `json:"id" gorm:"primaryKey"`
	Name     string   `json:"name" gorm:"not null"`
	Capacity int      `json:"capacity" gorm:"not null"`
	Screens  []Screen `json:"screens" gorm:"foreignKey:TheaterID"`
}

type Screen struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"not null"`
	TheaterID uint   `json:"theater_id"`
	Capacity  int    `json:"capacity" gorm:"not null"`
	Seats     []Seat `json:"seats" gorm:"foreignKey:ScreenID"`
}

type Seat struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	ScreenID uint   `json:"screen_id"`
	Row      string `json:"row" gorm:"not null"`
	Number   int    `json:"number" gorm:"not null"`
	Category string `json:"category"` // e.g., standard, premium, vip
}
