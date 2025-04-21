package middlewares

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"xlink/common/logger"
)

func LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		contextWithLogger, err := logger.New(c.Context())
		if err != nil {
			ctx, _ := logger.New(context.Background())
			logger.GetLoggerFromCtx(ctx).Error(ctx, "couldn't assign logger to context", zap.Error(err))
		} else {
			c.SetUserContext(contextWithLogger)
		}

		return c.Next()
	}
}
