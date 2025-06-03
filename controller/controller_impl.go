package controller

import (
	"catering-admin-go/domain"
	"catering-admin-go/service"
	"catering-admin-go/web"
	"net/http"

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
		return err
	}

	result, err := ctrl.svc.AddProduct(c.Context(), reqBody)
	if err != nil {
		c.JSON(&web.Response[*domain.Domain]{
			Code:   http.StatusBadRequest,
			Status: http.StatusText(400),
			Data:   nil,
		})
		c.Status(400)
	}

	c.JSON(&web.Response[*domain.Domain]{
		Code:   http.StatusCreated,
		Status: http.StatusText(http.StatusCreated),
		Data:   result,
	})
	c.Status(201)
	return nil
}
