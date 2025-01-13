package models

import "gorm.io/gorm"

type Actor struct {
	gorm.Model
	ActorDetails
	Details map[string]any `gorm:"type:jsonb;null"`
}

type ActorDetails struct {
	Name    string `json:"name"`
	ActorID uint   `json:"id" gorm:"column:tmdb_id;unique"`
}
