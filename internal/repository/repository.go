package repository

import (
	"context"
	"github.com/Kofandr/To-do_list/internal/domain/model"
)

type Repository interface {
	ToDoRepository
	UsersRepository
}

type ToDoRepository interface {
	GetTask(ctx context.Context, id int) (*model.Task, error)
	GetTasksUser(ctx context.Context, id int) (*[]model.Task, error)
	CreateTask(ctx context.Context, task *model.RequestTask) (int, error)
	UpdateTask(ctx context.Context, id int, update *model.RequestTask) error
	CompleteTask(ctx context.Context, id int) error
	DeleteTask(ctx context.Context, id int) error
}

type UsersRepository interface {
	CreateUser(ctx context.Context, user *model.NewUser) (int, error)
	GetUsers(ctx context.Context) (*[]model.User, error)
	DeleteUser(ctx context.Context, id int) error
	UserExists(ctx context.Context, id int) (bool, error)
	GetUsersByName(ctx context.Context, username string) (*model.User, error)
	CreateTwoFACode(ctx context.Context, code *model.TwoFACode) error
	GetTwoFACode(ctx context.Context, userID int, forLogin bool) (*model.TwoFACode, error)
	DeleteTwoFACode(ctx context.Context, userID int, forLogin bool) error
	UpdateTelegramConfirmed(ctx context.Context, userID int, confirmed bool) error
	GetUserByTelegramUsername(ctx context.Context, tgUsername string) (*model.User, error)
	BindTelegramChat(ctx context.Context, chatID int64, linkCode string) error
	VerifyTelegramCode(ctx context.Context, userID int, code string) (bool, error)
	AssignLinkCode(ctx context.Context, userID int) (string, error)
	SetTelegramCode(ctx context.Context, userID int, code string) error
}
