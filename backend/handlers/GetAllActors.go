package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/CalvoM/link-them/models"
)

func (h handler) GetAllActors(w http.ResponseWriter, r *http.Request) {
	var actors []models.ActorDetails
	result := h.dbClient.Table("actors").Select("name", "tmdb_id").Find(&actors)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Server Error"))
	}
	jsonResponse, err := json.Marshal(actors)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Server Error"))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
