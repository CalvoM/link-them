package main

import (
	"net/http"

	"github.com/CalvoM/link-them/db"
	"github.com/CalvoM/link-them/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Debug().Msg("Initializing Link Them Backend.")
	h := handlers.New(db.Init())
	r := mux.NewRouter()
	r.HandleFunc("/actors", h.GetAllActors).Methods(http.MethodGet)
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:4009",
	}
	log.Fatal().Msg(srv.ListenAndServe().Error())
}
