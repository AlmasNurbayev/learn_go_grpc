# Description: Makefile for SSO


migrate_up:
# TODO - move DSN to env variable or flag
	go run cmd/migrator/main.go -typeTask "up" -dsn "postgres://postgres:postgres@localhost:5435/sso_db?sslmode=disable"

migrate_down:
# TODO - move DSN to env variable or flag
	go run cmd/migrator/main.go -typeTask "down" -dsn "postgres://postgres:postgres@localhost:5435/sso_db?sslmode=disable"

build:
	go build -o SSO cmd/sso/main.go

run:
	go run cmd/sso/main.go -config=./config/local.yaml

test:
	go test -v -count=1 ./tests/...

docker_run:
	cd ../ && docker-compose up -d --build