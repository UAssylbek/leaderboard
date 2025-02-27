# Базовый образ с Go 1.23
FROM golang:1.23-alpine AS builder

# Рабочая директория
WORKDIR /app

# Копируем go.mod и go.sum (сначала нужно создать их)
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Сборка приложения
RUN go build -o leaderboard ./main.go

# Финальный образ
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/leaderboard .
EXPOSE 8080
CMD ["./leaderboard"]