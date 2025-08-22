package model

type NewUser struct {
	Username string `json:"username" validate:"required,min=1"`
}

type User struct {
	Username string `json:"username" validate:"required,min=1"`
	UserID   int    `json:"user_id" validate:"required,min=1"`
}
type NewTask struct {
	Title       string
	Description string
	UserID      int
}

type Task struct {
	ID          int
	Title       string
	Description string
	UserID      int
	Completed   bool
}
