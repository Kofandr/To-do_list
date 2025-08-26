package model

type NewUser struct {
	Username string `json:"username" validate:"required,min=1"`
	Password string `json:"password" validate:"required,min=8"`
}

type User struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
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
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
