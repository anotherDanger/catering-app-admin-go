package main

import (
	"catering-admin-go/controller"
	"catering-admin-go/helper"
	"catering-admin-go/repository"
	"catering-admin-go/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

	db, cleanup, _ := helper.NewDb()
	repo := repository.NewRepositoryImpl()
	svc := service.NewServiceImpl(repo, db)
	ctrl := controller.NewControllerImpl(svc)

	defer cleanup()

	app.Post("/v1/products", ctrl.AddProduct)
	app.Get("/v1/products", ctrl.GetProducts)
	app.Delete("/v1/products/:id", ctrl.DeleteProduct)
	app.Put("/v1/products/:id", ctrl.UpdateProduct)

	app.Listen(":8080")

}
