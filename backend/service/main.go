package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/CalvoM/link-them/models"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type MyJson map[string]any

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	dbClient := DBInit()
	token := os.Getenv("TOKEN")
	personID := 1000
	for personID < 10000 {
		url := fmt.Sprintf("https://api.themoviedb.org/3/person/%s?language=en-US", strconv.Itoa(personID))
		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Add("accept", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}

		defer res.Body.Close()
		log.Info().Msg(res.Status)
		if res.Status == "200 OK" {
			body, _ := io.ReadAll(res.Body)
			var bufJson MyJson
			json.Unmarshal(body, &bufJson)
			var actor models.Actor
			actor.Name = bufJson["name"].(string)
			actor.ActorID = uint(bufJson["id"].(float64))
			details := make(map[string]any)
			details["also_known_as"] = bufJson["also_known_as"]
			if bufJson["birthday"] == nil {
				bufJson["birthday"] = ""
			}
			details["birthday"] = bufJson["birthday"].(string)
			if bufJson["deathday"] == nil {
				bufJson["deathday"] = ""
			}
			details["deathday"] = bufJson["deathday"].(string)
			if bufJson["profile_path"] == nil {
				bufJson["profile_path"] = ""
			}
			details["profile_picture"] = bufJson["profile_path"].(string)
			actor.Details = details
			result := dbClient.Create(&actor)
			if result.Error != nil {
				log.Error().Msg("Failed to create the actor.")
			}
		}
		personID++
	}
}
