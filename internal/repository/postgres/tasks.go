package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/Kofandr/To-do_list/internal/domain/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (postgres *PgxRepository) GetTask(ctx context.Context, id int) (*model.Task, error) {
	const query = `
		SELECT id, title, description, user_id, completed 
		FROM tasks 
		WHERE id = $1
		`

	var task model.Task
	err := postgres.db.QueryRow(ctx, query, id).Scan(&task.ID, &task.Title, &task.Description, &task.UserID, &task.Completed)

	return &task, err
}

func (postgres *PgxRepository) GetTasksUser(ctx context.Context, id int) (*[]model.Task, error) {
	const query = ` 
		SELECT id, title, description, user_id, completed 
		FROM tasks 
		WHERE user_id = $1
	`

	rows, err := postgres.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []model.Task

	for rows.Next() {
		var task model.Task

		if err := rows.Scan(); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return &tasks, nil
}

func (postgres *PgxRepository) CreateTask(ctx context.Context, task *model.RequestTask) (int, error) {
	const query = ` 
		INSERT INTO  (title, description, user_id) 
		VALUES ($1, $2, $3) 
		RETURNING id
	`

	var id int

	err := postgres.db.QueryRow(ctx, query, task.Title, task.Description, task.UserID).Scan(&id)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return id, fmt.Errorf("%w: task '%s' already exists", ErrDuplicate, task.Title)
		}

		return id, fmt.Errorf("failed to create task: %w", err)
	}

	return id, err

}

func (postgres *PgxRepository) UpdateTask(ctx context.Context, id int, update *model.RequestTask) error {
	const query = `
        UPDATE tasks
        SET
            title = $1,
            description = $2
        	user_id = $3
        WHERE id = $4
    `

	current, err := postgres.GetTask(ctx, id)
	if err != nil {
		return err
	}

	name := current.Title
	if update.Title != "" {
		name = update.Title
	}

	description := current.Description
	if update.Description != "" {
		description = update.Description
	}

	userID := current.UserID
	if update.UserID != 0 {
		userID = update.UserID
	}

	result, err := postgres.db.Exec(ctx, query, name, description, id, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (postgres *PgxRepository) CompleteTask(ctx context.Context, id int) error {
	const query = `
        UPDATE tasks
        SET
            completed = true
        WHERE id = $1
    `

	result, err := postgres.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (postgres *PgxRepository) DeleteTask(ctx context.Context, id int) error {
	const query = ` 
		DELETE FROM tasks 
		WHERE id = $1
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
