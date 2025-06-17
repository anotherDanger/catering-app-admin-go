package controller

import (
	"catering-admin-go/domain"
	"catering-admin-go/helper"
	"catering-admin-go/service"
	"catering-admin-go/web"
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ControllerImpl struct {
	svc service.Service
}

func NewControllerImpl(svc service.Service) Controller {
	return &ControllerImpl{svc: svc}
}

func (ctrl *ControllerImpl) Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	var reqBody domain.Admin
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Request data is invalid."})
	}
	if err := helper.ValidateStruct(reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Please fill all required fields correctly."})
	}
	result, err := ctrl.svc.Login(ctx, &reqBody)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Incorrect username or password."})
	}
	return c.Status(fiber.StatusOK).JSON(result)
}

func (ctrl *ControllerImpl) AddProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	var reqBody web.Request
	reqBody.Name = c.FormValue("name")
	reqBody.Description = c.FormValue("description")
	price, err := strconv.Atoi(c.FormValue("price"))
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Price must be a valid number.", "")
	}
	reqBody.Price = price
	stock, err := strconv.Atoi(c.FormValue("stock"))
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Stock must be a valid number.", "")
	}
	reqBody.Stock = stock
	if err := helper.ValidateStruct(reqBody); err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Please complete all required product fields.", "")
	}
	result, err := ctrl.svc.AddProduct(ctx, &reqBody)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusInternalServerError, "Unable to add product. Please try again later.", "")
	}
	return web.SuccessResponse[*domain.Domain](c, fiber.StatusCreated, "Product successfully added.", result)
}

func (ctrl *ControllerImpl) GetProducts(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	products, err := ctrl.svc.GetProducts(ctx)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to load products. Please try again later.", "")
	}
	return web.SuccessResponse[[]*domain.Domain](c, fiber.StatusOK, "Products loaded successfully.", products)
}

func (ctrl *ControllerImpl) DeleteProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	id := c.Params("id")
	err := ctrl.svc.DeleteProduct(ctx, id)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusInternalServerError, "Unable to delete product. Please try again later.", "")
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (ctrl *ControllerImpl) UpdateProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	id := c.Params("id")
	name := c.FormValue("name")
	description := c.FormValue("description")
	stockStr := c.FormValue("stock")
	priceStr := c.FormValue("price")
	stock, err := strconv.Atoi(stockStr)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Stock must be a valid number.", "")
	}
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Price must be a valid number.", "")
	}
	reqBody := &web.Request{
		Name:        name,
		Description: description,
		Stock:       stock,
		Price:       price,
	}
	response, err := ctrl.svc.UpdateProduct(ctx, reqBody, id)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update product. Please try again later.", "")
	}
	return web.SuccessResponse(c, fiber.StatusOK, "Product successfully updated.", response)
}

func (ctrl *ControllerImpl) GetOrders(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	orders, err := ctrl.svc.GetOrders(ctx)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to load orders. Please try again later.", "")
	}
	return web.SuccessResponse[[]*domain.Orders](c, fiber.StatusOK, "Orders loaded successfully.", orders)
}

func (ctrl *ControllerImpl) UpdateOrder(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	var reqBody domain.Orders
	err := c.BodyParser(&reqBody)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update order. Please try again later.", "")
	}
	id := c.Params("id")

	if err := ctrl.svc.UpdateOrder(ctx, &reqBody, id); err != nil {
		return web.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update order. Please try again later.", "")
	}
	return web.SuccessResponse[interface{}](c, fiber.StatusOK, "Order successfully updated.", nil)
}

func (ctrl *ControllerImpl) DeleteOrder(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	id := c.Params("id")
	if err := ctrl.svc.DeleteOrder(ctx, id); err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete order", "")
	}
	return web.SuccessResponse[interface{}](c, fiber.StatusNoContent, "Order successfully deleted", nil)
}
