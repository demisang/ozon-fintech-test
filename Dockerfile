FROM golang:1.20-alpine3.18 AS builder

ADD . /src
WORKDIR /src

RUN go mod download
RUN go build -o links ./cmd/links-service/main.go

FROM alpine:3.18

EXPOSE 8382

COPY --from=builder /src/links /app/
COPY --from=builder /src/config/main.example.yml /app/config/main.yml

ENTRYPOINT ["/app/links", "--config=/app/config/main.yml"]
