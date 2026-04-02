.PHONY: run build test migrate-create migrate-up migrate-down docker-up docker-down

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

test:
	go test -v ./...

# Database migrations
migrate-create:
	migrate create -ext sql -dir migrations/sql -seq $(name)

migrate-up:
	migrate -database "postgres://postgres:postgres@localhost:5432/tottho_vandar?sslmode=disable" -path migrations/sql up

migrate-down:
	migrate -database "postgres://postgres:postgres@localhost:5432/tottho_vandar?sslmode=disable" -path migrations/sql down

# Docker
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Tools
tidy:
	go mod tidy

vendor:
	go mod vendor