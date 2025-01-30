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

  1. actors
  2. movies
  3. actor_credits
  4. movie_credits

This will populate the DB with data so that it is possible to start the web server with data to display
and analyze.

## Starting the web server

Please run

  ```sh
  go run main.go
  ```

## Design Specification

Please read the specification [here](https://docs.google.com/document/d/1M2wU7mJmwb-g56dbN7Okx5ViCS9QEb0JLQU8T8oMuy4/edit?usp=sharing)

## Documentation

In addition to this **README**, we have **docs/** that has our d2 files that illustrate the various components and how they interact.
To run the d2 files:

1. Please install d2lang cli from the [instructions](https://d2lang.com/tour/install/)
2. Inside the **docs** folder please run

    ```sh
    d2 --watch <d2_file>
    ```

    Where **d2_file** is any of the d2 files in the **docs/** folder
3. Open the browser on the port used as stated in the logs.
