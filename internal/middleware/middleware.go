package middleware

import (
	"fmt"
	"github.com/Kofandr/To-do_list/config"
	"github.com/Kofandr/To-do_list/internal/appctx"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

func RequestLogger(logg *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqID := uuid.NewString()
			reqLog := logg.With("request_id", reqID)

			ctx := c.Request().Context()

			ctx = appctx.WithLogger(ctx, reqLog)

			c.SetRequest(c.Request().WithContext(ctx))

			start := time.Now()

			err := next(c)

			duration := time.Since(start)

			req := c.Request()
			res := c.Response()

			logFields := []any{
				"method", req.Method,
				"path", req.URL.Path,
				"ip", c.RealIP(),
				"status", res.Status,
				"duration", duration,
			}

			if err != nil {
				logFields = append(logFields, "err", err.Error())
				logg.Error("request failed", logFields...)
			} else {
				logg.Info("request handled", logFields...)
			}

			return err
		}
	}
}

func JWTAuth(cfg *config.Configuration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			logg := appctx.LoggerFromContext(ctx)

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				logg.Info("Missing authorization header")

				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing authorization header"})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				logg.Info("Invalid authorization header format")

				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid authorization header format"})
			}

			tokenString := parts[1]
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(cfg.JWTSecret), nil
			})
			if err != nil || !token.Valid {
				logg.Info("Invalid token")

				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				logg.Info("Invalid token claims")

				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token claims"})
			}

			if tokenType, ok := claims["type"].(string); !ok || tokenType != "access" {
				logg.Info("IInvalid token type, expected access")

				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token type, expected access"})
			}

			c.Set("userID", claims["user_id"])
			c.Set("username", claims["username"])

			return next(c)
		}
	}
}
