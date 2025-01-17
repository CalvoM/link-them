package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
	"sync"

	"github.com/CalvoM/link-them/models"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MyJson map[string]any

var (
	dbClient *gorm.DB
	token    string
)

var (
	baseActorDetailsURL  string = "https://api.themoviedb.org/3/person/%s?language=en-US&append_to_response=credits"
	baseMovieDetailsURL  string = "https://api.themoviedb.org/3/movie/%s?language=en-US&append_to_response=credits"
	baseCreditDetailsURL string = "https://api.themoviedb.org/3/credit/%s"
)

func scrapActors() {
	personID := 70000
	for personID < 80000 {
		url := fmt.Sprintf(baseActorDetailsURL, strconv.Itoa(personID))
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
			details["credits"] = bufJson["credits"]
			actor.Details = details
			result := dbClient.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "tmdb_id"}}, DoUpdates: clause.Assignments(MyJson{"details": details})}).Create(&actor)
			if result.Error != nil {
				log.Error().Msg(fmt.Sprintf("Failed to create the actor. %s", result.Error.Error()))
			}
		}
		personID++
	}
}

func scrapMovies() {
	movieID := 30853
	for movieID < 50000 {
		url := fmt.Sprintf(baseMovieDetailsURL, strconv.Itoa(movieID))
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("accept", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}

		defer res.Body.Close()
		if res.Status == "200 OK" {
			body, _ := io.ReadAll(res.Body)
			var bufJson MyJson
			json.Unmarshal(body, &bufJson)
			var movie models.Movie
			movie.Title = bufJson["title"].(string)
			movie.MovieID = uint(bufJson["id"].(float64))
			log.Info().Msg(fmt.Sprintf("Found movie ID %d", movie.MovieID))
			details := make(map[string]any)
			details["credits"] = bufJson["credits"]
			if bufJson["status"] == nil {
				// Default to released
				bufJson["status"] = "Released"
			}
			details["status"] = bufJson["status"]
			if bufJson["budget"] == nil {
				bufJson["budget"] = 0
			}
			details["budget"] = bufJson["budget"]
			if bufJson["revenue"] == nil {
				bufJson["revenue"] = 0
			}
			details["revenue"] = bufJson["revenue"]
			if bufJson["poster_path"] == nil {
				bufJson["poster_path"] = ""
			}
			details["poster_picture"] = bufJson["poster_path"].(string)
			details["credits"] = bufJson["credits"]
			movie.Details = details
			result := dbClient.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "tmdb_id"}}, DoUpdates: clause.Assignments(MyJson{"details": details})}).Create(&movie)
			if result.Error != nil {
				log.Error().Msg(fmt.Sprintf("Failed to create the movie. %s", result.Error.Error()))
			}
		} else {
			log.Info().Msg(res.Status)
		}
		movieID++
	}
}

func scrapCreditsFromActors() {
	// CreditID is some random string
	var creditsIDs []string
	result := dbClient.Table("actors").Select("jsonb_path_query(details, '$.credits.cast[*].credit_id')").Scan(&creditsIDs)
	if result.Error != nil {
		log.Info().Msg(result.Error.Error())
	}
	creditsIDs = lo.Map(creditsIDs, func(id string, index int) string { return id[1 : len(id)-1] })

	var savedCreditIDs []string
	result = dbClient.Table("credits").Select([]string{"tmdb_id"}).Scan(&savedCreditIDs)
	if result.Error != nil {
		log.Info().Msg(result.Error.Error())
	}
	remainingCreditIDs := lo.Filter(creditsIDs, func(id string, index int) bool { return slices.Contains(savedCreditIDs, id) == false })
	creditsIDsChunks := slices.Chunk(remainingCreditIDs, 2000)
	log.Info().Msg("Data is setup")
	var wg sync.WaitGroup
	for chunk := range creditsIDsChunks {
		wg.Add(1)
		log.Info().Msg("Wg Added")
		go func() {
			defer wg.Done()
			for _, creditID := range chunk {
				url := fmt.Sprintf(baseCreditDetailsURL, creditID)
				req, _ := http.NewRequest("GET", url, nil)
				req.Header.Add("accept", "application/json")
				req.Header.Add("Authorization", "Bearer "+token)
				res, err := http.DefaultClient.Do(req)
				if err != nil {
					log.Fatal().Msg(err.Error())
				}

				defer res.Body.Close()
				if res.Status == "200 OK" {
					log.Info().Msg(fmt.Sprintf("Found %s", creditID))
					body, _ := io.ReadAll(res.Body)
					var bufJson MyJson
					json.Unmarshal(body, &bufJson)
					var credit models.Credit
					credit.CreditID = bufJson["id"].(string)
					credit.Details = bufJson
					result := dbClient.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "tmdb_id"}}, DoUpdates: clause.Assignments(MyJson{"details": credit.Details})}).Create(&credit)
					if result.Error != nil {
						log.Error().Msg(fmt.Sprintf("Failed to create the credit. %s", result.Error.Error()))
					}
				} else {
					log.Info().Msg(res.Status)
				}
			}
		}()
	}
}

func scrapCreditsFromMovies() {
	// CreditID is some random string
	var creditsIDs []string
	result := dbClient.Table("movies").Select("jsonb_path_query(details, '$.credits.cast[*].credit_id');").Scan(&creditsIDs)
	if result.Error != nil {
		log.Info().Msg(result.Error.Error())
	}
	for _, creditID := range creditsIDs {
		creditID = creditID[1 : len(creditID)-1]
		url := fmt.Sprintf(baseCreditDetailsURL, creditID)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("accept", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}

		defer res.Body.Close()
		if res.Status == "200 OK" {
			log.Info().Msg(fmt.Sprintf("Found %s", creditID))
			body, _ := io.ReadAll(res.Body)
			var bufJson MyJson
			json.Unmarshal(body, &bufJson)
			var credit models.Credit
			credit.CreditID = bufJson["id"].(string)
			credit.Details = bufJson
			result := dbClient.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "tmdb_id"}}, DoUpdates: clause.Assignments(MyJson{"details": credit.Details})}).Create(&credit)
			if result.Error != nil {
				log.Error().Msg(fmt.Sprintf("Failed to create the credit. %s", result.Error.Error()))
			}
		} else {
			log.Info().Msg(res.Status)
		}
	}
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	dbClient = DBInit()
	token = os.Getenv("TOKEN")
	// scrapActors()
	// scrapMovies()
	scrapCreditsFromActors()
}
