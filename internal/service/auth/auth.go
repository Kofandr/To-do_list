package auth

import "github.com/Kofandr/To-do_list/internal/repository"

type AuthService struct {
	db repository.Repository
}

func New(db repository.Repository) *AuthService {
	return &AuthService{db: db}
}
