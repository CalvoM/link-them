# Link Them Backend

## Getting started

1. Install the go dependencies in the *go.mod* file

  ```sh
  go mod download
  ```

2. Use the .env template to create a .env file .
3. Install the Docker and docker-compose.
4. Create the docker network

  ```sh
  docker network create link-them-net
  ```

5. Run the docker-compose services

  ```sh
  docker-compose up
  ```

## Getting data

The **scrapper/main.go** folder contains scrappers that get movies, actors and credits
from the TMDB Api. If you would like to scrap for any data please run the following:

  ```sh
  go run *.go -command <CMD>
```

The available options for **<CMD>** are:
  a. actors
  b. movies
  c. actor_credits
  d. movie_credits

This will populate the DB with data so that it is possible to start the web server with data to display
and analyze.

## Starting the web server

Please run

  ```sh
  go run main.go
  ```
