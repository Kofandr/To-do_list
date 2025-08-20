package middleware

import (
	"github.com/Kofandr/To-do_list/internal/appctx"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"log/slog"
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
