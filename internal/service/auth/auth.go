package auth

import (
	"context"
	"fmt"
	"github.com/Kofandr/To-do_list/config"
	"github.com/Kofandr/To-do_list/internal/domain/model"
	"github.com/Kofandr/To-do_list/internal/repository"
	"github.com/Kofandr/To-do_list/internal/service/app_jwt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db  repository.Repository
	cfg config.Configuration
}

func New(db repository.Repository, cfg config.Configuration) *Service {
	return &Service{db: db, cfg: cfg}
}

func (service *Service) Register(user *model.NewUser, ctx context.Context) (int, error) {
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

func (service *Service) Login(user *model.NewUser, ctx context.Context) (*model.Tokens, error) {
	dbUser, err := service.db.GetUsersByName(ctx, user.Username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	tokens, err := app_jwt.GenerateTokens(dbUser, service.cfg)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}
func (service *Service) Refresh(token *model.Tokens, ctx context.Context) (*model.Tokens, error) {
	jwtToken, err := jwt.Parse(token.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(service.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if !jwtToken.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	if tokenType, ok := claims["type"].(string); !ok || tokenType != "refresh" {
		return nil, fmt.Errorf("invalid token type, expected refresh")
	}

	userName, ok := claims["username"].(string)
	if !ok {
		return nil, fmt.Errorf("username not found in token")
	}

	user, err := service.db.GetUsersByName(ctx, userName)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	newTokens, err := app_jwt.GenerateTokens(user, service.cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return newTokens, nil
}
