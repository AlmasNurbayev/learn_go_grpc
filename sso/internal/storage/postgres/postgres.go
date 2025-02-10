package postgres

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func NewStorage(DSN string, log *slog.Logger) (*Storage, error) {
	const op = "postgres.NewStorage"
	log.With(slog.String("op", op)).Info("init storage " + DSN)

	db, err := sqlx.Open("postgres", DSN)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	db.SetMaxOpenConns(10) // Максимальное количество открытых соединений
	db.SetMaxIdleConns(5)  // Максимальное количество простаивающих соединений
	tx, err := db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db, tx: tx}, nil
}
