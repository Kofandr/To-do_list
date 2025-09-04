package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/Kofandr/To-do_list/internal/domain/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
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
		SELECT user_id, username, password, telegram_chat_id  
		FROM users 
		WHERE username = $1
		`

	var user model.User
	err := postgres.db.QueryRow(ctx, query, username).Scan(&user.UserID, &user.Username, &user.Password, &user.TelegramChatID)

	return &user, err
}

func (postgres *PgxRepository) CreateTwoFACode(ctx context.Context, code *model.TwoFACode) error {
	const query = `
		INSERT INTO twofa_codes (user_id, code_hash, expires_at, for_login)
		VALUES ($1, $2, $3, $4)
	`

	_, err := postgres.db.Exec(ctx, query, code.UserID, code.CodeHash, code.ExpiresAt, code.ForLogin)
	if err != nil {
		return fmt.Errorf("failed to create twofa code: %w", err)
	}

	return nil
}

func (postgres *PgxRepository) GetTwoFACode(ctx context.Context, userID int, forLogin bool) (*model.TwoFACode, error) {
	const query = `
		SELECT user_id, code_hash, expires_at, for_login
		FROM twofa_codes
		WHERE user_id = $1 AND for_login = $2
		ORDER BY created_at DESC
		LIMIT 1
	`

	var twofaCode model.TwoFACode
	err := postgres.db.QueryRow(ctx, query, userID, forLogin).Scan(
		&twofaCode.UserID,
		&twofaCode.CodeHash,
		&twofaCode.ExpiresAt,
		&twofaCode.ForLogin,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get twofa code: %w", err)
	}

	return &twofaCode, nil
}

func (postgres *PgxRepository) DeleteTwoFACode(ctx context.Context, userID int, forLogin bool) error {
	const query = `
		DELETE FROM twofa_codes
		WHERE user_id = $1 AND for_login = $2
	`

	res, err := postgres.db.Exec(ctx, query, userID, forLogin)
	if err != nil {
		return fmt.Errorf("failed to delete twofa code: %w", err)
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (postgres *PgxRepository) UpdateTelegramConfirmed(ctx context.Context, userID int, confirmed bool) error {
	const query = `
		UPDATE users
		SET telegram_confirmed = $1
		WHERE user_id = $2
	`

	res, err := postgres.db.Exec(ctx, query, confirmed, userID)
	if err != nil {
		return fmt.Errorf("failed to update telegram confirmed: %w", err)
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (postgres *PgxRepository) GetUserByTelegramUsername(ctx context.Context, tgUsername string) (*model.User, error) {
	const query = `
		SELECT user_id, username, password, telegram_username, telegram_confirmed
		FROM users
		WHERE telegram_username = $1
	`

	var user model.User
	err := postgres.db.QueryRow(ctx, query, tgUsername).Scan(
		&user.UserID,
		&user.Username,
		&user.Password,
		&user.TelegramUsername,
		&user.TelegramConfirmed,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get user by telegram username: %w", err)
	}

	return &user, nil
}

func (postgres *PgxRepository) BindTelegramChat(ctx context.Context, chatID int64, linkCode string) error {
	const query = `
		UPDATE users
		SET telegram_chat_id = $1, link_code = NULL
		WHERE link_code = $2;
	`

	cmdTag, err := postgres.db.Exec(ctx, query, chatID, linkCode)
	if err != nil {
		return fmt.Errorf("BindTelegramChat exec failed: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (postgres *PgxRepository) VerifyTelegramCode(ctx context.Context, userID int, code string) (bool, error) {
	const query = `
		SELECT telegram_code, code_expires_at
		FROM users
		WHERE user_id = $1;
	`

	var dbCode string
	var expiresAt time.Time

	err := postgres.db.QueryRow(ctx, query, userID).Scan(&dbCode, &expiresAt)
	if err != nil {
		return false, fmt.Errorf("VerifyTelegramCode query failed: %w", err)
	}

	if dbCode != code || time.Now().After(expiresAt) {
		return false, nil
	}

	const deleteQuery = `
		UPDATE users SET telegram_code = NULL WHERE user_id = $1;
	`

	if _, err := postgres.db.Exec(ctx, deleteQuery, userID); err != nil {
		return false, fmt.Errorf("VerifyTelegramCode cleanup failed: %w", err)
	}

	return true, nil
}

func (postgres *PgxRepository) AssignLinkCode(ctx context.Context, userID int) (string, error) {
	const query = `
		UPDATE users SET link_code = $1 WHERE user_id = $2;
	`

	linkCode := uuid.NewString()

	_, err := postgres.db.Exec(ctx, query, linkCode, userID)
	if err != nil {
		return "", fmt.Errorf("AssignLinkCode exec failed: %w", err)
	}

	return linkCode, nil
}

func (postgres *PgxRepository) SetTelegramCode(ctx context.Context, userID int, code string) error {
	const query = `
		UPDATE users
		SET telegram_code = $1,
		    code_expires_at = now() + interval '5 minutes'
		WHERE user_id = $2;
	`

	_, err := postgres.db.Exec(ctx, query, code, userID)
	if err != nil {
		return fmt.Errorf("SetTelegramCode exec failed: %w", err)
	}

	return nil
}
