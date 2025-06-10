package repository

import (
	"catering-admin-go/domain"
	"context"
	"database/sql"
	"errors"
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

func (repo *RepositoryImpl) GetProducts(ctx context.Context, tx *sql.Tx) ([]*domain.Domain, error) {
	query := "select * from products"
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var products []*domain.Domain
	defer rows.Close()
	for rows.Next() {
		var product domain.Domain
		var description sql.NullString
		err := rows.Scan(&product.Id, &product.Name, &description, &product.Stock, &product.Price, &product.CreatedAt, &product.ModifiedAt)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		if description.Valid {
			product.Description = description.String
		}

		products = append(products, &product)
	}

	return products, nil
}

func (repo *RepositoryImpl) DeleteProduct(ctx context.Context, tx *sql.Tx, id string) error {
	query := "delete from products where id = ?"
	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowAff, _ := result.RowsAffected()
	if rowAff == 0 {
		return fmt.Errorf("no product found with id %s", id)
	}

	return nil
}

func (repo *RepositoryImpl) UpdateProduct(ctx context.Context, tx *sql.Tx, entity *domain.Domain, id string) (*domain.Domain, error) {
	query := "update products set name = ?, description = ?, stock = ?, price = ?, modified_at = ? where id = ?"

	result, err := tx.ExecContext(ctx, query, entity.Name, entity.Description, entity.Stock, entity.Price, entity.ModifiedAt, id)
	if err != nil {
		return nil, err
	}

	rowAff, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowAff == 0 {
		return nil, errors.New("no rows updated")
	}

	var product domain.Domain

	row := tx.QueryRowContext(ctx, "select id, name, description, stock, price, created_at, modified_at from products where id = ?", id)
	err = row.Scan(&product.Id, &product.Name, &product.Description, &product.Stock, &product.Price, &product.CreatedAt, &product.ModifiedAt)
	if err != nil {
		return nil, err
	}

	return &product, nil
}
