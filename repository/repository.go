package repository

import (
	"catering-admin-go/domain"
	"context"
	"database/sql"
)

type Repository interface {
	Login(ctx context.Context, tx *sql.Tx, entity *domain.Admin) (*domain.Admin, error)
	AddProduct(ctx context.Context, tx *sql.Tx, entity *domain.Domain) (*domain.Domain, error)
	GetProducts(ctx context.Context, tx *sql.Tx) ([]*domain.Domain, error)
	DeleteProduct(ctx context.Context, tx *sql.Tx, id string) error
	UpdateProduct(ctx context.Context, tx *sql.Tx, entity *domain.Domain, id string) (*domain.Domain, error)
}
