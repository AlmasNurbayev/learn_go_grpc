Тренировочный проект для реализации простого GRPC-сервера и клиента
Golang

## Запуск

если запускать или билдить Golang без docker-compose, то необходимо положить .env в корень SSO, иначе не прокинутся учетные данные для соединения с БД

## БД

Особенность БД - pgx/v5/stdlib в режиме совместимости с SQLX

(нативный режим PGX/v5, не совместимый со стандартной библиотекой pkg.go.dev/database/sql и несовместимый с github.com/jmoiron/sqlx)

В папке SSO - makefile, для запуска команд нужно перейти в эту папку

Тесты - создают в указанной в конфиге БД реальные записи users и не удаляют за собой
