run:
	docker compose up -d

up:
	docker compose up -d db

down:
	docker compose down

migrate:
	migrate -path migrations -database postgres://user:12345@localhost:5437/ozon_fintech?sslmode=disable up

fmt:
	gofumpt -w .

tidy:
	go mod tidy

lint: fmt tidy
	golangci-lint run ./...

test: up
	go test -v -count=10 -race ./...
	make down
