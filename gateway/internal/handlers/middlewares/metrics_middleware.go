package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"time"
	"xlink/common/prometheus"
)

func MetricsMiddlewareFiber() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start).Seconds()

		prometheus.RecordRequest(string(c.Request().Header.Method()), c.Path(), duration)

		return err
	}
}
