FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /bin/shortener cmd/shortener/main.go

FROM alpine AS runner

WORKDIR /app
COPY --from=builder /bin/shortener .

ENV BASE_URL=http://localhost:8080
ENV SERVER_ADDRESS=localhost:8080
ENV FILE_STORAGE_PATH=storage.json
ENV LOG_LEVEL=info
ENV LOG_ENCODING=json
ENV DATABASE_DSN=host=localhost port=5432 user=postgres dbname=url_shortener sslmode=disable

EXPOSE 8080
CMD ["./shortener"]