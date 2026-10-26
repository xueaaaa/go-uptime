migrate-up:
	migrate -path migrations -database "postgres://root:root@localhost:5432/go-uptime?sslmode=disable" up 1

migrate-down:
	migrate -path migrations -database "postgres://root:root@localhost:5432/go-uptime?sslmode=disable" down 1

build:
	go mod tidy && \
	go build -o go-uptime .

compose-build:
	docker compose build app

run:
	go mod tidy && \
	go run .

compose-run:
	docker compose down --remove-orphans && \
	docker compose run --service-ports app serve

test:
	go test -v ./...