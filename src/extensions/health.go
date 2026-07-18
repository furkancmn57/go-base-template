package extensions

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type healthStatus struct {
	Status   string `json:"status"`
	Postgres string `json:"postgres"`
	Redis    string `json:"redis"`
}

// RegisterHealth mounts GET /health on the given router.
func RegisterHealth(router fiber.Router, db *gorm.DB, redisClient *redis.Client) {
	router.Get("/health", func(c *fiber.Ctx) error {
		result := healthStatus{Status: "ok", Postgres: "ok", Redis: "ok"}
		healthy := true

		if sqlDB, err := db.DB(); err != nil || sqlDB.Ping() != nil {
			result.Postgres = "down"
			healthy = false
		}

		if err := redisClient.Ping(c.Context()).Err(); err != nil {
			result.Redis = "down"
			healthy = false
		}

		if !healthy {
			result.Status = "degraded"
			return c.Status(fiber.StatusServiceUnavailable).JSON(result)
		}

		return c.Status(fiber.StatusOK).JSON(result)
	})
}
