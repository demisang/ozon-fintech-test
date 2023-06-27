# Тестовое задание Ozon.Fintech

### Config:

```shell
cp config/main.example.yml config/main.yml
```

Switch In-Memory/Database storage in config-file line:

```yaml
storage: "db"
storage: "memory"
```

### Run:

```shell
docker-compose up -d db


docker-compose up -d db app

go run cmd/links-service/main.go --config=$PWD/config/main.yml
```

### Database:

`postgres://user:12345@localhost:5437/ozon_fintech?sslmode=disable`

```shell
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate create -ext sql -dir migrations -seq create_links_table
migrate -path migrations -database postgres://user:12345@localhost:5437/ozon_fintech?sslmode=disable up [N]
migrate -path migrations -database postgres://user:12345@localhost:5437/ozon_fintech?sslmode=disable down [N]
migrate -path migrations -database postgres://user:12345@localhost:5437/ozon_fintech?sslmode=disable goto [V]
migrate -path migrations -database postgres://user:12345@localhost:5437/ozon_fintech?sslmode=disable version
```

### Demo:

```shell
# Create new Link
curl --request POST \
  --url http://localhost:8382/create \
  --header 'Content-Type: application/json' \
  --data '{
	"url": "https://www.google.com/search?q=%D0%9A%D0%B0%D0%BA+%D0%B8%D1%81%D0%BA%D0%B0%D1%82%D1%8C+%D0%B2+%D0%B3%D1%83%D0%B3%D0%BB%D0%B5"
}'

# Get Link by Code
curl --request GET \
  --url http://localhost:8382/KiANRJctbf
```
