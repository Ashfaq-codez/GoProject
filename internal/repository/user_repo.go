package repository

import (
	"database/sql"
	db "user-api/db/sqlc"
)

type UserRepository struct {
	q *db.Queries
}

func NewUserRepository(dbConn *sql.DB) *UserRepository {
	return &UserRepository{
		q: db.New(dbConn),
	}
}

func (r *UserRepository) Queries() *db.Queries {
	return r.q
}
