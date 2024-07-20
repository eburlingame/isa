# Isa: Online, multiplayer Uno game

## Local Development

- Start Redis:
  - Run `docker-compose -f docker-compose-services.yml up`
- Start the server
  - `cd api`
  - Run `go run src/*.go`
- Start the frontend
  - `cd frontend`
  - `yarn start`
