package main

import (
	"fmt"
	"os"

	"github.com/CalvoM/link-them/models"
	"github.com/rs/zerolog/log"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBInit() *gorm.DB {
	if godotenv.Load("../.env") != nil {
		log.Fatal().Msg("We could not load the environment setup")
	}
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("PG_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPasswd := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPasswd, dbHost, dbPort, dbName)
	dbClient, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: false})
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	log.Info().Msg("DB Initialization successful.")
	dbClient.AutoMigrate(&models.Actor{})
	dbClient.AutoMigrate(&models.Movie{})
	dbClient.AutoMigrate(&models.Credit{})
	return dbClient
}
