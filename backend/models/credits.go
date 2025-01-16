package models

import "gorm.io/gorm"

type Credit struct {
	gorm.Model
	CreditID string         `gorm:"column:tmdb_id;uniqueIndex"`
	Details  map[string]any `gorm:"type:jsonb;null"`
}
