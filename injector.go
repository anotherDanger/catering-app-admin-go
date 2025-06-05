//go:build wireinject
// +build wireinject

package main

import (
	"catering-admin-go/controller"
	"catering-admin-go/helper"
	"catering-admin-go/repository"
	"catering-admin-go/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

var ServerSet = wire.NewSet(
	repository.NewRepositoryImpl,
	service.NewServiceImpl,
	controller.NewControllerImpl,
	helper.NewDb,
	NewServer,
)

func InitServer() (*fiber.App, func(), error) {
	wire.Build(ServerSet)
	return nil, nil, nil
}
