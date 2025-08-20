package model

type NewUser struct {
	Name string
}

type User struct {
	Name string
	ID   int
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
