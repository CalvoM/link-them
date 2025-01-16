package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

type ActorsConnectionQuery struct {
	SrcActor    string `json:"srcActor"`
	SrcActorID  uint   `json:"srcActorID"`
	DestActor   string `json:"destActor"`
	DestActorID uint   `json:"destActorID"`
}

func (h handler) GetActorConnection(w http.ResponseWriter, r *http.Request) {
	var actorQuery ActorsConnectionQuery
	err := json.NewDecoder(r.Body).Decode(&actorQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var srcCreditsIDs []string
	CreditIDSelectQuery := []string{"jsonb_path_query(details, '$.credits.cast[*].credit_id')"}
	CreditIDWhereQuery := "tmdb_id = ? and name = ?"
	result := h.dbClient.Table("actors").Select(CreditIDSelectQuery).Where(CreditIDWhereQuery, actorQuery.SrcActorID, actorQuery.SrcActor).Scan(&srcCreditsIDs)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Server Error"))
		log.Fatal().Msg(result.Error.Error())
		return
	}
	var destCreditsIDs []string
	result = h.dbClient.Table("actors").Select(CreditIDSelectQuery).Where(CreditIDWhereQuery, actorQuery.DestActorID, actorQuery.DestActor).Scan(&destCreditsIDs)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Server Error"))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Query => %v", actorQuery)
}
