package auth

import (
	"context"
	"github.com/Kofandr/To-do_list/internal/domain/model"
	"golang.org/x/crypto/bcrypt"
)

func (service *AuthService) Register(user *model.NewUser, ctx context.Context) (int, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	user.Password = string(hash)

	id, err := service.db.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (service *AuthService) Login(user *model.NewUser, ctx context.Context) {
	dbUser := service.db.

}
