package models

import "database/sql"

type User struct {
	Id       int64          `db:"id"`
	Email    sql.NullString `db:"email"`
	Phone    sql.NullString `db:"phone"`
	PassHash []byte         `db:"pass_hash"`
}
