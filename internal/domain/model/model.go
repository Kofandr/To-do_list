package model

type NewUser struct {
	Username string `json:"username" validate:"required,min=1"`
	Password string `json:"password" validate:"required,min=8"`
}

type User struct {
	Username string `json:"username" validate:"required,min=1"`
	UserID   int    `json:"user_id" validate:"required,min=1"`
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
