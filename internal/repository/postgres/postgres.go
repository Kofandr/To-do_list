package postgres

import (
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *PgxRepository {
	return &PgxRepository{
		db,
	}
}

var (
	ErrDuplicate = errors.New("duplicate entry")
)
