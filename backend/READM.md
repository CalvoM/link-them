# Link Them Backend

## Getting started

1. Install the go dependencies in the *go.mod* file

  ```sh
  go mod download
  ```

2. Use the .env template to create a .env file .
3. Install the docker and docker-compose.
4. Create the docker network

  ```sh
  docker network create link-them-net
  ```

5. Run the docker-compose services

  ```sh
  docker-compose up
  ```
