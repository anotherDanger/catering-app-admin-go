package service

import (
	"catering-admin-go/domain"
	"catering-admin-go/web"
	"context"
)

type Service interface {
	Login(ctx context.Context, request *domain.Admin) (*web.AdminResponse, error)
	AddProduct(ctx context.Context, request *web.Request) (*domain.Domain, error)
	GetProducts(ctx context.Context) ([]*domain.Domain, error)
	DeleteProduct(ctx context.Context, id string) error
	UpdateProduct(ctx context.Context, request *web.Request, id string) (*domain.Domain, error)
}
