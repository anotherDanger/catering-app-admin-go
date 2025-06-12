package main

import (
	"catering-admin-go/controller"
	"catering-admin-go/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewServer(handler controller.Controller) *fiber.App {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE",
	}))

	app.Post("/v1/login", handler.Login)

	protectedRoute := app.Group("/api")
	protectedRoute.Use(middleware.MyMiddleware)
	protectedRoute.Post("/v1/products", handler.AddProduct)
	protectedRoute.Get("/v1/products", handler.GetProducts)
	protectedRoute.Delete("/v1/products/:id", handler.DeleteProduct)
	protectedRoute.Put("/v1/products/:id", handler.UpdateProduct)

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
