package repository

import (
	"catering-admin-go/domain"
	"context"
	"database/sql"
)

type Repository interface {
	AddProduct(ctx context.Context, tx *sql.Tx, entity *domain.Domain) (*domain.Domain, error)
	GetProducts(ctx context.Context, tx *sql.Tx) ([]*domain.Domain, error)
	DeleteProduct(ctx context.Context, tx *sql.Tx, id string) error
	UpdateProduct(ctx context.Context, tx *sql.Tx, entity *domain.Domain, id string) (*domain.Domain, error)
}
