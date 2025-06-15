package controller

import (
	"catering-admin-go/domain"
	"catering-admin-go/helper"
	"catering-admin-go/logger"
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

func (ctrl *ControllerImpl) Login(c *fiber.Ctx) error {
	var reqBody domain.Admin

	reqBody.Username = c.FormValue("username")
	reqBody.Password = c.FormValue("password")

	result, err := ctrl.svc.Login(c.Context(), &reqBody)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(result)
}

func (ctrl *ControllerImpl) AddProduct(c *fiber.Ctx) error {
	var reqBody web.Request

	reqBody.Name = c.FormValue("name")
	reqBody.Description = c.FormValue("description")

	price, err := strconv.Atoi(c.FormValue("price"))
	if err != nil {
		logger.GetLogger("controller-log").Log("controller", "error", err.Error())
		return web.ErrorResponse(c, 400, "Invalid price", err.Error())
	}
	reqBody.Price = price

	stock, err := strconv.Atoi(c.FormValue("stock"))
	if err != nil {
		logger.GetLogger("controller-log").Log("controller", "error", err.Error())
		return web.ErrorResponse(c, 400, "Invalid stock", err.Error())
	}
	reqBody.Stock = stock

	if err := helper.ValidateStruct(reqBody); err != nil {
		logger.GetLogger("controller-log").Log("controller", "error", err.Error())
		return web.ErrorResponse(c, 400, "Validation failed", err.Error())
	}

	result, err := ctrl.svc.AddProduct(c.Context(), &reqBody)
	if err != nil {
		logger.GetLogger("controller-log").Log("controller", "error", err.Error())
		return web.ErrorResponse(c, 400, "Service error", err.Error())
	}

	logger.GetLogger("controller-log").Log("controller", "error", "Product added successfully")
	return web.SuccessResponse[*domain.Domain](c, 201, "Created", result)
}

func (ctrl *ControllerImpl) GetProducts(c *fiber.Ctx) error {

	products, err := ctrl.svc.GetProducts(c.Context())
	if err != nil {
		logger.GetLogger("controller-log").Log("Controller", "error", err.Error())
		return web.ErrorResponse(c, 400, "Error", err.Error())
	}

	return web.SuccessResponse[[]*domain.Domain](c, 200, "OK", products)
}

func (ctrl *ControllerImpl) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	err := ctrl.svc.DeleteProduct(c.Context(), id)
	if err != nil {
		logger.GetLogger("controller-log").Log("Controller", "error", err.Error())
		return web.ErrorResponse(c, 400, "Error", err.Error())
	}

	return c.SendStatus(204)

}

func (ctrl *ControllerImpl) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	name := c.FormValue("name")
	description := c.FormValue("description")
	stockStr := c.FormValue("stock")
	priceStr := c.FormValue("price")

	stock, err := strconv.Atoi(stockStr)
	if err != nil {
		logger.GetLogger("controller-log").Log("Controller", "error", err.Error())
		return web.ErrorResponse(c, 400, "Invalid stock", err.Error())
	}

	price, err := strconv.Atoi(priceStr)
	if err != nil {
		logger.GetLogger("controller-log").Log("Controller", "error", err.Error())
		return web.ErrorResponse(c, 400, "Invalid price", err.Error())
	}

	reqBody := &web.Request{
		Name:        name,
		Description: description,
		Stock:       stock,
		Price:       price,
	}

	response, err := ctrl.svc.UpdateProduct(c.Context(), reqBody, id)
	if err != nil {
		logger.GetLogger("controller-log").Log("Controller", "error", err.Error())
		return web.ErrorResponse(c, 400, "Error", err.Error())
	}

	return web.SuccessResponse(c, 200, "OK", response)
}

func (ctrl *ControllerImpl) GetOrders(c *fiber.Ctx) error {
	orders, err := ctrl.svc.GetOrders(c.Context())
	if err != nil {

		logger.GetLogger("controller-log").Log("Controller getOrders", "error", err.Error())
		return web.ErrorResponse(c, 400, "Error", err.Error())
	}

	return web.SuccessResponse[[]*domain.Orders](c, 200, "OK", orders)
}

func (ctrl *ControllerImpl) UpdateOrder(c *fiber.Ctx) error {
	var reqBody domain.Orders
	id := c.Params("id")
	err := c.BodyParser(&reqBody)
	if err != nil {
		logger.GetLogger("controller-log").Log("Controller update order", "error", err.Error())
		return web.ErrorResponse(c, 400, "Error", err.Error())
	}

	err = ctrl.svc.UpdateOrder(c.Context(), &reqBody, id)
	if err != nil {
		logger.GetLogger("controller-log").Log("Controller update order", "error", err.Error())
		return web.ErrorResponse(c, 400, "Error", err.Error())
	}

	return web.SuccessResponse[interface{}](c, 200, "OK", nil)
}
