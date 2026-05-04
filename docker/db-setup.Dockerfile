FROM migrate/migrate:v4.17.0 AS migrate

FROM golang:1.25-alpine

WORKDIR /app

COPY --from=migrate /usr/local/bin/migrate /usr/local/bin/migrate

CMD ["sh", "-c", "migrate -path migrations -database \"$DATABASE_URL\" up && go run ./cmd/seed"]
