package postgres

import (
	"context"
	"database/sql"
	"errors"
	"sso/internal/models"
	"sso/internal/storage"

	"github.com/lib/pq"
)

func (s *Storage) SaveUser(ctx context.Context, email string, phone string, passHash []byte) (id int64, err error) {
	query := `INSERT INTO users (email, phone, pass_hash) VALUES ($1, $2, $3) RETURNING id`

	err = s.db.QueryRowContext(ctx, query, email, phone, passHash).Scan(&id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // Код ошибки уникальности PostgreSQL
				return 0, storage.ErrUserExists
			}
		}
		return 0, err
	}
	return id, nil
}

func (s *Storage) GetUserByPhone(ctx context.Context, phone string) (user models.User, err error) {
	query := `SELECT id, email, phone, pass_hash FROM users WHERE phone = $1`
	err = s.db.GetContext(ctx, &user, query, phone)
	if err != nil {
		return user, err
	}
	return user, nil
}
func (s *Storage) GetUserByEmail(ctx context.Context, email string) (user models.User, err error) {
	query := `SELECT id, email, phone, pass_hash FROM users WHERE email = $1`
	err = s.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, storage.ErrUserNotFound
		}
		return user, err
	}
	return user, nil
}
func (s *Storage) GetUserById(ctx context.Context, id int64) (user models.User, err error) {
	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, id int64) (isAdmin bool, err error) {
	return false, nil
}
