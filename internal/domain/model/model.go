package model

import "time"

type NewUser struct {
	Username         string `json:"username" validate:"required,min=1"`
	Password         string `json:"password" validate:"required,min=8"`
	TelegramUsername string `json:"telegram_username" validate:"required,min=1"`
}

type User struct {
	UserID            int    `json:"user_id"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	TelegramUsername  string `json:"telegram_username"`
	TelegramConfirmed bool   `json:"telegram_confirmed"`
}
type RequestTask struct {
	Title       string `json:"title" validate:"required,min=1"`
	Description string `json:"description" validate:"required,min=1"`
	UserID      int    `json:"user_id" validate:"required,min=1"`
}

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      int    `json:"user_id"`
	Completed   bool   `json:"completed"`
}

type Tokens struct {
	AccessToken  string `json:"access_token" validate:"required,min=1"`
	RefreshToken string `json:"refresh_token" validate:"required,min=1"`
}

type TwoFACode struct {
	UserID    int       `json:"user_id"`
	CodeHash  string    `json:"-"`
	ExpiresAt time.Time `json:"-"`
	ForLogin  bool      `json:"-"`
}

type VerifyCodeRequest struct {
	Code string `json:"code" validate:"required,len=6"`
}

type TelegramCodeChatID struct {
	ChatID int64  `json:"chat_id"`
	Code   string `json:"code"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
