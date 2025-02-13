package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	ctx context.Context
	db  *sqlx.DB
	log *slog.Logger
}

func NewStorage(DSN string, log *slog.Logger, timeout time.Duration) (*Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	const op = "postgres.NewStorage"
	log.With(slog.String("op", op)).Info("init storage " + DSN)

	db, err := sqlx.Connect("pgx", DSN)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	db.SetMaxIdleConns(5)  // Максимальное количество соединений
	db.SetMaxOpenConns(10) // Максимальное количество соединений

	// создание транзакции отключено, так как вызывает зависание при Close()
	// tx, err := db.Begin(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("%s: %w", op, err)
	// }

	return &Storage{ctx: ctx, db: db, log: log}, nil
}

func (s *Storage) Close() {
	const op = "postgres.Close"
	s.log.With(slog.String("op", op))

	if s.db != nil {
		s.log.Info("active Postgres conns", slog.Any("acquired_conns", s.db.Stats().OpenConnections))
		//s.tx.Rollback(context.Background())
		s.db.Close()
		s.log.Warn("DB connection closed")
	}
}
