package db

import (
	"fmt"
	"os"

	"github.com/CalvoM/link-them/models"
	"github.com/rs/zerolog/log"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	if godotenv.Load() != nil {
		log.Fatal().Msg("We could not load the environment setup")
	}
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("PG_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPasswd := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPasswd, dbHost, dbPort, dbName)
	dbClient, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	log.Info().Msg("DB Initialization successful.")
	dbClient.AutoMigrate(&models.Actor{})
	return dbClient
}
