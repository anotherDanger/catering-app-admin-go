package repository

import (
	"catering-admin-go/domain"
	"catering-admin-go/logger"
	"context"
	"database/sql"
	"errors"
)

type RepositoryImpl struct{}

func NewRepositoryImpl() Repository {
	return &RepositoryImpl{}
}

func (repo *RepositoryImpl) Login(ctx context.Context, tx *sql.Tx, entity *domain.Admin) (*domain.Admin, error) {
	query := "SELECT id, username, password FROM admin WHERE username = ?"
	row := tx.QueryRowContext(ctx, query, entity.Username)

	var response domain.Admin
	err := row.Scan(&response.Id, &response.Username, &response.Password)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (repo *RepositoryImpl) AddProduct(ctx context.Context, tx *sql.Tx, entity *domain.Domain) (*domain.Domain, error) {
	query := "INSERT INTO products(id, name, description, stock, price, created_at) VALUES(?, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, query, entity.Id, entity.Name, entity.Description, entity.Stock, entity.Price, entity.CreatedAt)
	if err != nil {
		logger.GetLogger("repository-log").Log("add product", "error", err.Error())
		return nil, err
	}

	rowAff, err := result.RowsAffected()
	if err != nil || rowAff == 0 {
		logger.GetLogger("repository-log").Log("add product", "error", err.Error())
		return nil, err
	}

	return entity, nil
}

func (repo *RepositoryImpl) GetProducts(ctx context.Context, db *sql.DB) ([]*domain.Domain, error) {
	query := "SELECT id, name, description, stock, price, created_at, modified_at FROM products"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		logger.GetLogger("repository-log").Log("get product", "error", err.Error())
		return nil, err
	}
	defer rows.Close()

	var products []*domain.Domain
	for rows.Next() {
		var product domain.Domain
		var description sql.NullString

		err := rows.Scan(&product.Id, &product.Name, &description, &product.Stock, &product.Price, &product.CreatedAt, &product.ModifiedAt)
		if err != nil {
			logger.GetLogger("repository-log").Log("get product", "error", err.Error())
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
	query := "DELETE FROM products WHERE id = ?"
	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		logger.GetLogger("repository-log").Log("delete product", "error", err.Error())
		return err
	}

	rowAff, err := result.RowsAffected()
	if err != nil || rowAff == 0 {
		logger.GetLogger("repository-log").Log("delete product", "error", err.Error())
		return err
	}

	return nil
}

func (repo *RepositoryImpl) UpdateProduct(ctx context.Context, tx *sql.Tx, entity *domain.Domain, id string) (*domain.Domain, error) {
	query := "UPDATE products SET name = ?, description = ?, stock = ?, price = ?, modified_at = ? WHERE id = ?"
	result, err := tx.ExecContext(ctx, query, entity.Name, entity.Description, entity.Stock, entity.Price, entity.ModifiedAt, id)
	if err != nil {
		logger.GetLogger("repository-log").Log("update product", "error", err.Error())
		return nil, err
	}

	rowAff, err := result.RowsAffected()
	if err != nil {
		logger.GetLogger("repository-log").Log("update product", "error", err.Error())
		return nil, err
	}
	if rowAff == 0 {
		return nil, errors.New("no rows updated")
	}

	var product domain.Domain
	row := tx.QueryRowContext(ctx, "SELECT id, name, description, stock, price, created_at, modified_at FROM products WHERE id = ?", id)
	err = row.Scan(&product.Id, &product.Name, &product.Description, &product.Stock, &product.Price, &product.CreatedAt, &product.ModifiedAt)
	if err != nil {
		logger.GetLogger("repository-log").Log("update product", "error", err.Error())
		return nil, err
	}

	return &product, nil
}

func (repo *RepositoryImpl) GetOrders(ctx context.Context, db *sql.DB) ([]*domain.Orders, error) {
	query := "SELECT id, product_name, username, quantity, status FROM orders"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		logger.GetLogger("repository-log").Log("get orders", "error", err.Error())
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Orders
	for rows.Next() {
		var order domain.Orders
		err := rows.Scan(&order.Id, &order.ProductName, &order.Username, &order.Quantity, &order.Status)
		if err != nil {
			logger.GetLogger("repository-log").Log("get orders", "error", err.Error())
			return nil, err
		}
		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		logger.GetLogger("repository-log").Log("get orders", "error", err.Error())
		return nil, err
	}

	return orders, nil
}

func (repo *RepositoryImpl) UpdateOrder(ctx context.Context, tx *sql.Tx, entity *domain.Orders, id string) error {
	query := "UPDATE orders SET status = ? WHERE id = ?"
	result, err := tx.ExecContext(ctx, query, entity.Status, id)
	if err != nil {
		logger.GetLogger("repository-log").Log("update orders", "error", err.Error())
		return err
	}

	rowAff, err := result.RowsAffected()
	if err != nil || rowAff == 0 {
		logger.GetLogger("repository-log").Log("update orders", "error", err.Error())
		return err
	}

	return nil
}
