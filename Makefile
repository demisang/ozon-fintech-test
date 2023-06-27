build:
	docker build \
 		--tag "orlov/ozon-fintech-demo"

run:
	docker run -it orlov/ozon-fintech-demo

migrate:
	migrate -path migrations -database postgres://user:12345@localhost:5437/ozon_fintech?sslmode=disable up
