package controller

import "github.com/gofiber/fiber/v2"

type Controller interface {
	AddProduct(c *fiber.Ctx) error
}
