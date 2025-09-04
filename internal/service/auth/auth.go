package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/Kofandr/To-do_list/config"
	"github.com/Kofandr/To-do_list/internal/botclient"
	"github.com/Kofandr/To-do_list/internal/domain/model"
	"github.com/Kofandr/To-do_list/internal/repository"
	"github.com/Kofandr/To-do_list/internal/service/app_jwt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"math/big"
)

type Service struct {
	db     repository.Repository
	cfg    *config.Configuration
	BotURL string
}

func New(db repository.Repository, cfg *config.Configuration, BotURL string) *Service {
	return &Service{db: db, cfg: cfg, BotURL: BotURL}
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

func (service *Service) Login(user *model.NewUser, ctx context.Context, logg *slog.Logger) (*model.LoginResult, error) {
	dbUser, err := service.db.GetUsersByName(ctx, user.Username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if dbUser.TelegramChatID == 0 {
		linkCode, err := service.db.AssignLinkCode(ctx, dbUser.UserID)
		if err != nil {
			return nil, err
		}

		return &model.LoginResult{
			LinkCode: linkCode,
			UserID:   dbUser.UserID,
		}, nil
	}

	code, err := Generate2FACode()
	if err != nil {
		return nil, err
	}

	err = service.db.SetTelegramCode(ctx, dbUser.UserID, code)
	if err != nil {
		return nil, err
	}

	go botclient.SendCode(service.BotURL, dbUser.TelegramChatID, code, logg)

	return &model.LoginResult{
		Message:     "2FA code sent to Telegram",
		Requires2FA: true,
		UserID:      dbUser.UserID,
	}, nil
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

func (service *Service) Verify2FA(req *model.Verify2FARequest, ctx context.Context, logg *slog.Logger) (*model.Tokens, error) {
	user, err := service.db.GetUsersByName(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	logg = logg.With("user", user.Username)

	valid, err := service.db.VerifyTelegramCode(ctx, user.UserID, req.Code)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, fmt.Errorf("fverify 2FA: db error")
	}
	logg = logg.With("user", user.Username)

	newTokens, err := app_jwt.GenerateTokens(user, service.cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return newTokens, nil
}

func Generate2FACode() (string, error) {
	max := big.NewInt(1000000) // [0..999999]
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", fmt.Errorf("generate 2FA code: %w", err)
	}

	return fmt.Sprintf("%06d", n.Int64()), nil
}
