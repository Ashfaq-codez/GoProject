package middleware

import (
	"time"
	"user-api/internal/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Inject Request ID
		reqID := uuid.New().String()
		c.Set("X-Request-ID", reqID)

		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		logger.Log.Info("Incoming Request",
			zap.String("request_id", reqID),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.String("duration", duration.String()),
		)

		return err
	}
}