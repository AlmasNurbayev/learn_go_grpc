package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewStorage(ctx context.Context, DSN string, log *slog.Logger) (*Storage, error) {
	const op = "postgres.NewStorage"
	log.With(slog.String("op", op)).Info("init storage " + DSN)

	db, err := pgxpool.New(ctx, DSN)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	db.Config().MaxConns = 10 // Максимальное количество простаивающих соединений

	// создание транзакции отключено, так как вызывает зависание при Close()
	// tx, err := db.Begin(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("%s: %w", op, err)
	// }

	return &Storage{db: db, log: log}, nil
}

func (s *Storage) Close() {
	const op = "postgres.Close"
	s.log.With(slog.String("op", op))

	if s.db != nil {
		s.log.Info("active Postgres conns", slog.Any("acquired_conns", s.db.Stat().AcquiredConns()))
		//s.tx.Rollback(context.Background())
		s.db.Close()
		s.log.Warn("DB connection closed")
	}
}
