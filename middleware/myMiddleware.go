package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func MyMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	splitHeader := strings.SplitN(authHeader, " ", 2)
	if len(splitHeader) != 2 || splitHeader[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	tokenString := splitHeader[1]
	c.Locals("token", tokenString)

	return c.Next()
}
