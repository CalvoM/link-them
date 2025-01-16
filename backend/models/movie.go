package models

import "gorm.io/gorm"

type Movie struct {
	gorm.Model
	Title   string         `json:"title"`
	MovieID uint           `json:"id" gorm:"column:tmdb_id;uniqueIndex"`
	Details map[string]any `gorm:"type:jsonb;null"`
}
