package postgres

import (
	"context"
	"database/sql"
	"errors"
	"sso/internal/errorsPackage"
	"sso/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Storage) SaveUser(ctx context.Context, email string, phone string, passHash []byte, idRole int) (id int64, err error) {

	op := "postgres.SaveUser"
	s.log.With("op", op)

	// если передены пустые email или phone, то присваиваем им NULL, чтобы не нарушать уникальность
	emailOrNil := sql.NullString{String: email, Valid: email != ""}
	phoneOrNil := sql.NullString{String: phone, Valid: phone != ""}
	query := `INSERT INTO users (email, phone, pass_hash, role_id) VALUES ($1, $2, $3, $4) RETURNING id`

	err = s.db.QueryRowContext(ctx, query, emailOrNil, phoneOrNil, passHash, idRole).Scan(&id)
	if err != nil {
		if pqErr, ok := err.(*pgconn.PgError); ok {
			if pqErr.Code == "23505" { // Код ошибки уникальности PostgreSQL
				s.log.Error("user already exists", "error", err)
				return 0, errorsPackage.ErrUserExists
			}
		}
		return 0, err
	}
	return id, nil
}

func (s *Storage) GetUserByPhone(ctx context.Context, phone string) (user models.User, err error) {
	op := "postgres.GetUserByPhone"
	s.log.With("op", op)

	query := `SELECT id, email, phone, pass_hash FROM users WHERE phone = $1`
	err = s.db.QueryRowContext(ctx, query, phone).Scan(&user.Id, &user.Email, &user.Phone, &user.PassHash)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Error("user not found", "error", err)
			return user, errorsPackage.ErrUserNotFound
		}
		return user, err
	}
	return user, nil
}
func (s *Storage) GetUserByEmail(ctx context.Context, email string) (user models.User, err error) {

	// искусственное замедление запроса
	// var temp string
	// err = s.db.QueryRow(ctx, "SELECT pg_sleep(18)").Scan(&temp)
	// if err != nil {
	// 	s.log.Error("canceled query DB", "error", err)
	// 	return user, err
	// }
	op := "postgres.GetUserByEmail"
	s.log.With("op", op)

	//fmt.Println("storage", ctx.Value(middleware.TraceIDKey))

	query := `SELECT id, email, phone, pass_hash FROM users WHERE email = $1;`
	err = s.db.QueryRowContext(ctx, query, email).Scan(&user.Id, &user.Email, &user.Phone, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errorsPackage.ErrUserNotFound
		}
		return user, err
	}
	return user, nil
}
func (s *Storage) GetUserById(ctx context.Context, id int64) (user models.User, err error) {
	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, id int64) (isAdmin bool, err error) {
	op := "postgres.IsAdmin"
	s.log.With("op", op)

	query := `SELECT roles.name 
						FROM roles
						JOIN users ON roles.id = users.role_id
						WHERE users.id = $1`
	var name string
	err = s.db.QueryRowContext(ctx, query, id).Scan(&name)
	s.log.Info(name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Error("not found user")
			return false, errorsPackage.ErrUserNotFound
		}
		s.log.Error("error get role", "error", err.Error())
		return false, err
	}
	if name == "admin" { // если роль админ, то возвращаем true
		return true, nil
	} else {
		return false, nil
	}
}

func (s *Storage) GetRoleByName(ctx context.Context, name string) (id int, err error) {
	query := `SELECT id  FROM roles WHERE name = $1`
	err = s.db.QueryRowContext(ctx, query, name).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errorsPackage.ErrRoleNotFound
		}
		return 0, err
	}
	return id, nil
}
