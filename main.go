package main

import (
	"catering-admin-go/controller"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewServer(handler controller.Controller) *fiber.App {
	app := fiber.New()
	app.Use(cors.New())

	app.Post("/v1/products", handler.AddProduct)
	app.Get("/v1/products", handler.GetProducts)
	app.Delete("/v1/products/:id", handler.DeleteProduct)
	app.Put("/v1/products/:id", handler.UpdateProduct)

	return app
}

func main() {
	app, cleanup, err := InitServer()
	if err != nil {
		panic(err)
	}

	defer cleanup()

	app.Listen(":8080")
}
