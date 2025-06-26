# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o webApp main.go

# Final stage
FROM alpine:3.20

WORKDIR /app

# Установим docker и docker-compose
RUN apk add --no-cache docker docker-compose

COPY --from=builder /app/webApp .
COPY .env ./
COPY assets ./assets

EXPOSE 8080

CMD ["./webApp"]