package repository

import (
	"catering-admin-go/domain"
	"context"
	"database/sql"
	"fmt"
)

type RepositoryImpl struct{}

func NewRepositoryImpl() Repository {
	return &RepositoryImpl{}
}

func (repo *RepositoryImpl) AddProduct(ctx context.Context, tx *sql.Tx, entity *domain.Domain) (*domain.Domain, error) {
	query := "insert into products(id, name, description, stock, price, created_at) values(?, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, query, entity.Id, entity.Name, entity.Description, entity.Stock, entity.Price, entity.CreatedAt)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	rowAff, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if rowAff == 0 {
		fmt.Println(err)
		return nil, err
	}

	response := &domain.Domain{
		Id:          entity.Id,
		Name:        entity.Name,
		Description: entity.Description,
		Stock:       entity.Stock,
		Price:       entity.Price,
		CreatedAt:   entity.CreatedAt,
	}

	return response, nil

}
