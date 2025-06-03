package controller

import (
	"catering-admin-go/domain"
	"catering-admin-go/helper"
	"catering-admin-go/service"
	"catering-admin-go/web"

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

	var reqBody *web.Request
	err := c.BodyParser(&reqBody)
	if err != nil {
		return web.ErrorResponse(c, 400, "Bad Request", err.Error())
	}
	err = helper.ValidateStruct(reqBody)
	if err != nil {
		return web.ErrorResponse(c, 400, "Bad Request", err.Error())
	}

	result, err := ctrl.svc.AddProduct(c.Context(), reqBody)
	if err != nil {
		return web.ErrorResponse(c, 400, "Bad Request", err.Error())
	}

	web.SuccessResponse[*domain.Domain](c, 201, "OK", result)

	return nil
}
