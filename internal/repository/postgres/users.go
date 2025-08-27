package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/Kofandr/To-do_list/internal/domain/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (postgres *PgxRepository) CreateUser(ctx context.Context, user *model.NewUser) (int, error) {
	const query = ` 
		INSERT INTO users (username, password) 
		VALUES ($1, $2)
		RETURNING user_id
	`

	var id int

	err := postgres.db.QueryRow(ctx, query, user.Username, user.Password).Scan(&id)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return id, fmt.Errorf("%w: user '%s' already exists", ErrDuplicate, user.Username)
		}

		return id, fmt.Errorf("failed to create user: %w", err)
	}

	return id, err
}

func (postgres *PgxRepository) GetUsers(ctx context.Context) (*[]model.User, error) {
	const query = ` 
		SELECT username, user_id 
		FROM users
	`

	rows, err := postgres.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var c model.User

		if err := rows.Scan(&c.Username, &c.UserID); err != nil {
			return nil, err
		}

		users = append(users, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &users, nil
}

func (postgres *PgxRepository) DeleteUser(ctx context.Context, id int) error {
	const query = ` 
		DELETE FROM users 
		WHERE user_id = $1
	`

	res, err := postgres.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (postgres *PgxRepository) UserExists(ctx context.Context, id int) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM users  WHERE user_id = $1)`

	var exists bool

	err := postgres.db.QueryRow(ctx, query, id).Scan(&exists)

	return exists, err
}

func (postgres *PgxRepository) GetUsersByName(ctx context.Context, username string) (*model.User, error) {
	const query = `
		SELECT user_id, username, password  
		FROM users 
		WHERE username = $1
		`

	var user model.User
	err := postgres.db.QueryRow(ctx, query, username).Scan(&user.UserID, &user.Username, &user.Password)

	return &user, err
}
