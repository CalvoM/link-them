package handlers

import "gorm.io/gorm"

type handler struct {
	dbClient *gorm.DB
}

func New(db *gorm.DB) handler {
	return handler{db}
}
