package service

import (
	"catering-admin-go/domain"
	"catering-admin-go/web"
	"context"
)

type Service interface {
	AddProduct(ctx context.Context, request *web.Request) (*domain.Domain, error)
}
