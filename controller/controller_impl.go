package controller

import (
	"catering-admin-go/domain"
	"catering-admin-go/helper"
	"catering-admin-go/service"
	"catering-admin-go/web"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ControllerImpl struct {
	svc service.Service
}

func NewControllerImpl(svc service.Service) Controller {
	return &ControllerImpl{
		svc: svc,
	}
}

func (ctrl *ControllerImpl) AddProduct(c *fiber.Ctx) error {
	var reqBody web.Request

	reqBody.Name = c.FormValue("name")
	reqBody.Description = c.FormValue("description")

	// konversi string ke int
	price, err := strconv.Atoi(c.FormValue("price"))
	if err != nil {
		return web.ErrorResponse(c, 400, "Invalid price", err.Error())
	}
	reqBody.Price = price

	stock, err := strconv.Atoi(c.FormValue("stock"))
	if err != nil {
		return web.ErrorResponse(c, 400, "Invalid stock", err.Error())
	}
	reqBody.Stock = stock

	// validasi
	if err := helper.ValidateStruct(reqBody); err != nil {
		return web.ErrorResponse(c, 400, "Validation failed", err.Error())
	}

	// lanjut ke service
	result, err := ctrl.svc.AddProduct(c.Context(), &reqBody)
	if err != nil {
		return web.ErrorResponse(c, 400, "Service error", err.Error())
	}

	return web.SuccessResponse[*domain.Domain](c, 201, "Created", result)
}

func (ctrl *ControllerImpl) GetProducts(c *fiber.Ctx) error {
	products, err := ctrl.svc.GetProducts(c.Context())
	if err != nil {
		return web.ErrorResponse(c, 400, "Error", err.Error())
	}

	return web.SuccessResponse[[]*domain.Domain](c, 200, "OK", products)
}
