package app_jwt

import (
	"github.com/Kofandr/To-do_list/config"
	"github.com/Kofandr/To-do_list/internal/domain/model"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func GenerateTokens(user *model.User, cfg *config.Configuration) (*model.Tokens, error) {
	accessClaims := jwt.MapClaims{
		"user_id":  user.UserID,
		"username": user.Username,
		"type":     "access",
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
	}
	refreshClaims := jwt.MapClaims{
		"user_id":  user.UserID,
		"username": user.Username,
		"type":     "refresh",
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessString, err := accessToken.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return nil, err
	}
	refreshString, err := refreshToken.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &model.Tokens{AccessToken: accessString, RefreshToken: refreshString}, nil
}
