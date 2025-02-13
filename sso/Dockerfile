# Используем официальный образ Go
FROM golang:1.23-alpine AS builder
# Устанавливаем рабочую директорию в контейнере
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Копируем файлы проекта
COPY . .

# Собираем приложение
RUN go build -o SSO ./cmd/sso/main.go


# Используем stage 2: минимальный контейнер
FROM alpine:latest
WORKDIR /root/

# Копируем только бинарник и конфиги из builder-образа
COPY --from=builder /app/config ./config
COPY --from=builder /app/SSO .

# Запускаем приложение
CMD ["./SSO", "-config=./config/local.yaml"]