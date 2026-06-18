package models

import "time"

type Cage struct {
	ID              int       `json:"id" gorm:"primaryKey"`
	Name            string    `json:"name"`
	LocationX       float64   `json:"location_x"`
	LocationY       float64   `json:"location_y"`
	LocationZ       float64   `json:"location_z"`
	WaterTemp       float64   `json:"water_temp"`
	DissolvedOxygen float64   `json:"dissolved_oxygen"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (Cage) TableName() string {
	return "cages"
}

type FeedRecord struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	CageID     int       `json:"cage_id"`
	FeedAmount float64   `json:"feed_amount"`
	Operator   string    `json:"operator"`
	CreatedAt  time.Time `json:"created_at"`
}

func (FeedRecord) TableName() string {
	return "feed_records"
}

type FeedRequest struct {
	CageID     int     `json:"cage_id" binding:"required"`
	FeedAmount float64 `json:"feed_amount" binding:"required"`
	Operator   string  `json:"operator"`
}
