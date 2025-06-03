package controller

import "github.com/gofiber/fiber/v2"

type Controller interface {
	AddProduct(c *fiber.Ctx) error
	GetProducts(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
}
