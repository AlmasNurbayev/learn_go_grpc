Тренировочный проект для реализации простого GRPC-сервера и клиента
Golang

# Варианты запуска

## Запуск через docker compose (рекомендуется для постоянной работы)

1. необходим .env согласно примера
2. docker compose up -d --build
   В этом случае миграции применяться автоматически

## Запуск через go run

если запускать или билдить Golang без docker-compose, то необходимо:

1. через docker compose запустить все контейнеры, чтобы стартовало DB
2. остановить контейнер SSO, чтобы освободить порт
3. создать .env в папку sso (а не в корне всех папок), иначе не прокинутся учетные данные для соединения с БД
4. после этого из sso запустить миграции go run cmd/migrator/main.go -typeTask "up" -dsn "postgres://..." (указать правильный DSN)
5. после миграции cd sso и оттуда запустить make run

# Прочие скрипты

## Github action

при git push tag Workflow настроено на билд ./sso/Dockergile и пуш image его в docker hub

## hook pre-push

Образец в sso (нужно установить golangci-lint, скопировать файл правильным образом, дать права на исполнение)

## make в sso - смотреть внутри

В папке SSO - makefile, для запуска команд нужно перейти в эту папку

# Подключение БД в SSO

Особенность БД - pgx/v5/stdlib в режиме совместимости с SQLX

(нативный режим PGX/v5, не совместимый со стандартной библиотекой pkg.go.dev/database/sql и несовместимый с github.com/jmoiron/sqlx)

# Тесты в SSO

Создают в указанной в конфиге БД реальные записи users и не удаляют за собой
