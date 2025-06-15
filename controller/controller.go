package controller

import "github.com/gofiber/fiber/v2"

type Controller interface {
	Login(c *fiber.Ctx) error
	AddProduct(c *fiber.Ctx) error
	GetProducts(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
	GetOrders(c *fiber.Ctx) error
	UpdateOrder(c *fiber.Ctx) error
}
