package models

type User struct {
	Id       int64  `db:"id"`
	Email    string `db:"email"`
	Phone    string `db:"phone"`
	PassHash []byte `db:"pass_hash"`
}
