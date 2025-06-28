# --- Stage 1: Build frontend ---
FROM node:20-alpine AS frontend-builder

WORKDIR /frontend
COPY ../frontend/package*.json ./
RUN npm install
COPY ../frontend .
RUN npm run build

# --- Stage 2: Build backend ---
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY --from=frontend-builder /frontend/dist ./frontend/dist

RUN CGO_ENABLED=0 GOOS=linux go build -o webApp main.go

# --- Stage 3: Final image ---
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache docker docker-compose

COPY --from=backend-builder /app/webApp .
COPY --from=backend-builder /app/frontend/dist ./frontend/dist
COPY .env ./
COPY assets ./assets

EXPOSE 8080

CMD ["./webApp"]