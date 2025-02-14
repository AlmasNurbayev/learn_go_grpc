# Используем официальный образ Go
FROM golang:1.23-alpine AS builder
# Устанавливаем рабочую директорию в контейнере
WORKDIR /app

COPY go.mod ./
RUN go mod download && go mod verify && go mod tidy

# Копируем файлы проекта
COPY . .

# Собираем приложение
RUN go build -o SSO ./cmd/sso/main.go
RUN go build -o MIGRATOR ./cmd/migrator/main.go


# Используем stage 2: минимальный контейнер
FROM alpine:latest
WORKDIR /root/

# Копируем бинарники, миграции и конфиги из builder-образа
COPY --from=builder /app/config ./config
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/SSO .
COPY --from=builder /app/MIGRATOR .

# Запускаем приложение
CMD sh -c './MIGRATOR -typeTask up -dsn $DSN && ./SSO -config=./config/local.yaml'
