package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/CalvoM/link-them/models"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
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
	var storedIDs []uint
	res := dbClient.Table("actors").Select("tmdb_id").Scan(&storedIDs)
	if res.Error != nil {
		log.Fatal().Msg(res.Error.Error())
	}
	start := 500000
	var wg sync.WaitGroup
	for start < 1000000 {
		personID := start
		nextPersonID := start + 10000
		wg.Add(1)
		go func() {
			defer func() {
				defer wg.Done()
			}()
			for personID < nextPersonID {
				if slices.Contains(storedIDs, uint(personID)) {
					log.Info().Msg(fmt.Sprintf("Skipping %d", personID))
					personID++
					continue
				}
				url := fmt.Sprintf(baseActorDetailsURL, strconv.FormatUint(uint64(personID), 10))
				req, _ := http.NewRequest("GET", url, nil)

				req.Header.Add("accept", "application/json")
				req.Header.Add("Authorization", "Bearer "+token)

				res, err := http.DefaultClient.Do(req)
				if err != nil {
					log.Fatal().Msg(err.Error())
				}

				defer res.Body.Close()
				if res.StatusCode == http.StatusOK {
					body, _ := io.ReadAll(res.Body)
					var bufJson MyJson
					json.Unmarshal(body, &bufJson)
					var actor models.Actor
					actor.Name = bufJson["name"].(string)
					actor.ActorID = uint(bufJson["id"].(float64))
					log.Info().Msg(fmt.Sprintf("Found Actor ID %d", actor.ActorID))
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
					// result := dbClient.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "tmdb_id"}}, DoUpdates: clause.Assignments(MyJson{"details": details})}).Create(&actor)
					result := dbClient.Create(&actor)
					if result.Error != nil {
						log.Error().Msg(fmt.Sprintf("Failed to create the actor. %s", result.Error.Error()))
					}
				} else {
					log.Info().Msg(res.Status)
					if res.StatusCode == http.StatusTooManyRequests {
						// Consider using a backoff algorithm
						return
					}

				}
				personID++
			}
		}()
		start += 10000
	}
	wg.Wait()
}

func scrapMovies() {
	var storedIDs []uint
	res := dbClient.Table("movies").Select("tmdb_id").Scan(&storedIDs)
	if res.Error != nil {
		log.Fatal().Msg(res.Error.Error())
	}
	start := 100000
	var wg sync.WaitGroup
	for start < 500000 {
		movieID := start
		nextMovieID := start + 8000
		wg.Add(1)
		go func() {
			defer func() { defer wg.Done() }()
			for movieID < nextMovieID {
				if slices.Contains(storedIDs, uint(movieID)) {
					log.Info().Msg(fmt.Sprintf("Skipping %d", movieID))
					movieID++
					continue
				}

				url := fmt.Sprintf(baseMovieDetailsURL, strconv.Itoa(movieID))
				req, _ := http.NewRequest("GET", url, nil)
				req.Header.Add("accept", "application/json")
				req.Header.Add("Authorization", "Bearer "+token)

				res, err := http.DefaultClient.Do(req)
				if err != nil {
					log.Fatal().Msg(err.Error())
				}

				defer res.Body.Close()
				if res.StatusCode == http.StatusOK {
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
					// result := dbClient.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "tmdb_id"}}, DoUpdates: clause.Assignments(MyJson{"details": details})}).Create(&movie)
					result := dbClient.Create(&movie)
					if result.Error != nil {
						log.Error().Msg(fmt.Sprintf("Failed to create the movie. %s", result.Error.Error()))
					}
				} else {
					log.Info().Msg(res.Status)
					if res.StatusCode == http.StatusTooManyRequests {
						// Consider using a backoff algorithm
						return
					}
				}
				movieID++
			}
		}()
		start += 8000
	}
	wg.Wait()
}

func scrapCreditsFromActors() {
	var remainingCreditIDs []string
	query := "select jsonb_path_query(details, '$.credits.cast[*].credit_id')->>0 from actors except tmdb_id from credits;"
	result := dbClient.Raw(query).Scan(&remainingCreditIDs)
	if result.Error != nil {
		log.Fatal().Msg(result.Error.Error())
	}

	creditsIDsChunks := slices.Chunk(remainingCreditIDs, 3000)
	var wg sync.WaitGroup
	// TMDB Rate limiting is 50 requests/sec
	wgCountLimit := 50
	for chunk := range creditsIDsChunks {
		if wgCountLimit > 0 {
			wg.Add(1)
			wgCountLimit--
			go func() {
				defer func() {
					defer wg.Done()
					wgCountLimit++
				}()
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
					if res.StatusCode == http.StatusOK {
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
						if res.StatusCode == http.StatusTooManyRequests {
							// Consider using a backoff algorithm
							return
						}
					}
					time.Sleep(10 * time.Millisecond)
				}
			}()
		}
	}
	wg.Wait()
}

func scrapCreditsFromMovies() {
	var remainingCreditIDs []string
	query := "select jsonb_path_query(details, '$.credits.cast[*].credit_id')->>0 from movies except select tmdb_id from credits;"
	result := dbClient.Raw(query).Scan(&remainingCreditIDs)
	if result.Error != nil {
		log.Fatal().Msg(result.Error.Error())
	}

	creditsIDsChunks := slices.Chunk(remainingCreditIDs, 4000)
	var wg sync.WaitGroup
	// TMDB Rate limiting is 50 requests/sec
	wgCountLimit := 50
	for chunk := range creditsIDsChunks {
		if wgCountLimit > 0 {
			wg.Add(1)
			wgCountLimit--
			go func() {
				defer func() {
					defer wg.Done()
					wgCountLimit++
				}()
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
					if res.StatusCode == http.StatusOK {
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
						if res.StatusCode == http.StatusTooManyRequests {
							// Consider using a backoff algorithm
							return
						}
					}
					time.Sleep(1 * time.Millisecond)
				}
			}()
		}
	}
	wg.Wait()
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	dbClient = DBInit()
	token = os.Getenv("TOKEN")
	// availableCmds := []string{"actors", "movies", "actor_credits", "movie_credits"}
	command := flag.String("command", "actors", "resource to scrap")
	flag.Parse()
	switch *command {
	case "actors":
		scrapActors()
	case "movies":
		scrapMovies()
	case "actor_credits":
		scrapCreditsFromActors()
	case "movie_credits":
		scrapCreditsFromMovies()
	default:
		log.Error().Msg("We do not support that option.")
		fmt.Println("\u001b[31m We do not support that option.\u001b[0m")
	}
}
