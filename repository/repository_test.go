package repository

import (
	"catering-admin-go/domain"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetProducts(t *testing.T) {
	now := time.Now()
	id := "123e4567-e89b-12d3-a456-426614174000"

	tests := []struct {
		name           string
		setupMock      func(mock sqlmock.Sqlmock)
		expectedErr    bool
		expectedResult []*domain.Domain
	}{
		{
			name: "Test GetProducts Success",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := mock.NewRows([]string{
					"Id", "Name", "Description", "Stock", "Price", "CreatedAt", "ModifiedAt",
				}).AddRow(
					id,
					"Product 1",
					"1st Product",
					10,
					1000,
					now,
					now,
				)

				mock.ExpectQuery("(?i)select .* from products").WillReturnRows(rows)
			},
			expectedErr: false,
			expectedResult: []*domain.Domain{
				{
					Id:          id,
					Name:        "Product 1",
					Description: "1st Product",
					Stock:       10,
					Price:       1000,
					CreatedAt:   &now,
					ModifiedAt:  &now,
				},
			},
		},
		{
			name: "1 column missing",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := mock.NewRows([]string{
					"Name", "Description", "Stock", "Price", "CreatedAt", "ModifiedAt",
				}).AddRow("Product 1", "1st Product", 10, 1000, now, now)
				mock.ExpectQuery("(?i)select \\* from products").WillReturnRows(rows)
			},
			expectedErr: true,
			expectedResult: []*domain.Domain{
				{
					Name:        "Product 1",
					Description: "1st Product",
					Stock:       10,
					Price:       1000,
					CreatedAt:   &now,
					ModifiedAt:  &now,
				},
			},
		},
		{
			name: "Failed",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("(?i)select \\* from products").WillReturnError(errors.New("failed to get products"))
			},
			expectedErr:    true,
			expectedResult: nil,
		},
		{
			name: "empty result",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := mock.NewRows([]string{
					"Id", "Name", "Description", "Stock", "Price", "CreatedAt", "ModifiedAt",
				})
				mock.ExpectQuery("(?i)select .* from products").WillReturnRows(rows)
			},
			expectedErr:    true,
			expectedResult: nil,
		},
		{
			name: "scan error due to type mismatch",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := mock.NewRows([]string{
					"Id", "Name", "Description", "Stock", "Price", "CreatedAt", "ModifiedAt",
				}).AddRow(
					"wrong-type", // should be UUID
					123,          // should be string
					"desc",
					"invalid-int",
					"invalid-float",
					time.Now(),
					time.Now(),
				)
				mock.ExpectQuery("(?i)select .* from products").WillReturnRows(rows)
			},
			expectedErr:    true,
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}

			tt.setupMock(mock)

			repo := NewRepositoryImpl()

			result, err := repo.GetProducts(context.Background(), db)

			if tt.expectedErr {
				if tt.name == "1 column missing" {
					assert.Empty(t, result)
				}

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult[0].Id, result[0].Id)
			assert.Equal(t, tt.expectedResult[0].Name, result[0].Name)
			assert.Equal(t, tt.expectedResult[0].Description, result[0].Description)
			assert.Equal(t, tt.expectedResult[0].Stock, result[0].Stock)
			assert.Equal(t, tt.expectedResult[0].Price, result[0].Price)
			assert.WithinDuration(t, *tt.expectedResult[0].CreatedAt, *result[0].CreatedAt, time.Second)
			assert.WithinDuration(t, *tt.expectedResult[0].ModifiedAt, *result[0].ModifiedAt, time.Second)
		})
	}

}

func TestAddProduct(t *testing.T) {
	created_at := time.Now()
	id := "123e4567-e89b-12d3-a456-426614174000"
	name := "Product 1"
	description := "1st Product"
	stock := 100
	price := 2000

	tests := []struct {
		name           string
		setupMock      func(mock sqlmock.Sqlmock)
		expectedErr    bool
		expectedResult *domain.Domain
	}{
		{
			name: "Success",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("(?i)insert\\s+into\\s+products\\s*\\(\\s*id\\s*,\\s*name\\s*,\\s*description\\s*,\\s*stock\\s*,\\s*price\\s*,\\s*created_at\\s*\\)\\s*values\\s*\\(\\s*\\?\\s*,\\s*\\?\\s*,\\s*\\?\\s*,\\s*\\?\\s*,\\s*\\?\\s*,\\s*\\?\\s*\\)").
					WithArgs(id, name, description, stock, price, created_at).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: false,
			expectedResult: &domain.Domain{
				Id:          id,
				Name:        name,
				Description: description,
				Stock:       stock,
				Price:       price,
				CreatedAt:   &created_at,
			},
		},
		{
			name: "1 column missing except description",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("(?i)insert\\s+into\\s+products\\s*\\(\\s*id\\s*,\\s*name\\s*,\\s*description\\s*,\\s*stock\\s*,\\s*price\\s*,\\s*created_at\\s*\\)\\s*values\\s*\\(\\s*\\?\\s*,\\s*\\?\\s*,\\s*\\?\\s*,\\s*\\?\\s*,\\s*\\?\\s*,\\s*\\?\\s*\\)").
					WithArgs(id, "", description, stock, price, created_at).
					WillReturnError(errors.New("field name cannot empty"))
			},
			expectedErr: true,
			expectedResult: &domain.Domain{
				Id:          id,
				Name:        "",
				Description: description,
				Stock:       stock,
				Price:       price,
				CreatedAt:   &created_at,
			},
		},
		{
			name: "Transaction failed",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("cannot start transaction"))
			},
			expectedErr:    true,
			expectedResult: nil,
		},
		{
			name: "Connection failed",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("connection failed"))
			},
			expectedErr:    true,
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)

			tt.setupMock(mock)

			tx, err := db.Begin()
			if err != nil {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			repo := NewRepositoryImpl()
			result, err := repo.AddProduct(context.Background(), tx, tt.expectedResult)

			if tt.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.Id, result.Id)
				assert.Equal(t, tt.expectedResult.Name, result.Name)
				assert.Equal(t, tt.expectedResult.Description, result.Description)
				assert.Equal(t, tt.expectedResult.Stock, result.Stock)
				assert.Equal(t, tt.expectedResult.Price, result.Price)
				assert.WithinDuration(t, *tt.expectedResult.CreatedAt, *result.CreatedAt, time.Second)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	modified_at := time.Now().Add(2 * time.Hour)
	id := "123e4567-e89b-12d3-a456-426614174000"
	name := "Product 2"
	description := "2nd Product"
	stock := 10
	price := 1000

	tests := []struct {
		name           string
		setupMock      func(mock sqlmock.Sqlmock)
		inputEntity    *domain.Domain
		expectedErr    bool
		expectedResult *domain.Domain
	}{
		{
			name: "Success",
			inputEntity: &domain.Domain{
				Name:        name,
				Description: description,
				Stock:       stock,
				Price:       price,
				ModifiedAt:  &modified_at,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`(?i)^update\s+products\s+set\s+name\s*=\s*\?,\s*description\s*=\s*\?,\s*stock\s*=\s*\?,\s*price\s*=\s*\?,\s*modified_at\s*=\s*\?\s+where\s+id\s*=\s*\?\s*$`).
					WithArgs(name, description, stock, price, modified_at, id).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectQuery(`(?i)^select id, name, description, stock, price, created_at, modified_at from products where id = \?$`).
					WithArgs(id).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "stock", "price", "created_at", "modified_at"}).
						AddRow(id, name, description, stock, price, time.Now(), modified_at))
			},
			expectedErr: false,
			expectedResult: &domain.Domain{
				Id:          id,
				Name:        name,
				Description: description,
				Stock:       stock,
				Price:       price,
				ModifiedAt:  &modified_at,
			},
		},
		{
			name: "1 column missing",
			inputEntity: &domain.Domain{
				Name:        name,
				Description: description,
				Stock:       stock,
				Price:       price,
				ModifiedAt:  &modified_at,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`(?i)^update\s+products\s+set\s+name\s*=\s*\?,\s*description\s*=\s*\?,\s*stock\s*=\s*\?,\s*price\s*=\s*\?,\s*modified_at\s*=\s*\?\s+where\s+id\s*=\s*\?\s*$`).
					WithArgs(name, description, stock, price, modified_at, id).
					WillReturnError(errors.New("1 column missing"))
			},
			expectedErr:    true,
			expectedResult: nil,
		},
		{
			name: "Failed",
			inputEntity: &domain.Domain{
				Name:        name,
				Description: description,
				Stock:       stock,
				Price:       price,
				ModifiedAt:  &modified_at,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`(?i)^update\s+products\s+set\s+name\s*=\s*\?,\s*description\s*=\s*\?,\s*stock\s*=\s*\?,\s*price\s*=\s*\?,\s*modified_at\s*=\s*\?\s+where\s+id\s*=\s*\?\s*$`).
					WithArgs(name, description, stock, price, modified_at, id).
					WillReturnError(errors.New("failed to update product"))
			},
			expectedErr:    true,
			expectedResult: nil,
		},
		{
			name: "Tx failed",
			inputEntity: &domain.Domain{
				Name:        name,
				Description: description,
				Stock:       stock,
				Price:       price,
				ModifiedAt:  &modified_at,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("transaction failed"))
			},
			expectedErr:    true,
			expectedResult: nil,
		},
		{
			name: "Connection failed",
			inputEntity: &domain.Domain{
				Name:        name,
				Description: description,
				Stock:       stock,
				Price:       price,
				ModifiedAt:  &modified_at,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("connection failed"))
			},
			expectedErr:    true,
			expectedResult: nil,
		},
		{
			name: "update but no rows affected",
			inputEntity: &domain.Domain{
				Name:        name,
				Description: description,
				Stock:       stock,
				Price:       price,
				ModifiedAt:  &modified_at,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`(?i)^update\s+products`).
					WithArgs(name, description, stock, price, modified_at, id).
					WillReturnResult(sqlmock.NewResult(0, 0)) // No rows affected
			},
			expectedErr:    true,
			expectedResult: nil,
		},
		{
			name: "select after update failed",
			inputEntity: &domain.Domain{
				Name:        name,
				Description: description,
				Stock:       stock,
				Price:       price,
				ModifiedAt:  &modified_at,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`(?i)^update\s+products`).
					WithArgs(name, description, stock, price, modified_at, id).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectQuery(`(?i)^select id, name, description, stock, price, created_at, modified_at from products where id = \?$`).
					WithArgs(id).
					WillReturnError(errors.New("select failed"))
			},
			expectedErr:    true,
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)

			tt.setupMock(mock)

			tx, err := db.Begin()
			if err != nil {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			repo := NewRepositoryImpl()

			result, err := repo.UpdateProduct(context.Background(), tx, tt.inputEntity, id)

			if tt.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, id, result.Id)
				assert.Equal(t, tt.inputEntity.Name, result.Name)
				assert.Equal(t, tt.inputEntity.Description, result.Description)
				assert.Equal(t, tt.inputEntity.Stock, result.Stock)
				assert.Equal(t, tt.inputEntity.Price, result.Price)
				assert.WithinDuration(t, *tt.inputEntity.ModifiedAt, *result.ModifiedAt, time.Second)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetOrders(t *testing.T) {
	createdAt := time.Now()
	modifiedAt := time.Now()

	tests := []struct {
		name           string
		setupMock      func(mock sqlmock.Sqlmock)
		expectedErr    bool
		expectedResult []*domain.Orders
	}{
		{
			name: "success get orders",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "product_id", "product_name", "username", "quantity", "total", "status", "created_at", "modified_at",
				}).AddRow(
					"1", "101", "ProductA", "user1", 2, 100.0, "pending", createdAt, modifiedAt,
				)

				mock.ExpectQuery("SELECT \\* FROM orders").WillReturnRows(rows)
			},
			expectedErr: false,
			expectedResult: []*domain.Orders{
				{
					Id:          "1",
					ProductId:   "101",
					ProductName: "ProductA",
					Username:    "user1",
					Quantity:    2,
					Total:       100.0,
					Status:      "pending",
					CreatedAt:   &createdAt,
					ModifiedAt:  &modifiedAt,
				},
			},
		},
		{
			name: "order not found",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "product_id", "product_name", "username", "quantity", "total", "status", "created_at", "modified_at",
				})
				mock.ExpectQuery("SELECT \\* FROM orders").WillReturnRows(rows)
			},
			expectedErr:    true,
			expectedResult: nil,
		},
		{
			name: "data corrupted on scan",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "product_id", "product_name", "username", "quantity", "total", "status", "created_at", "modified_at",
				}).AddRow("1", "101", "ProductA", "user1", "invalid", "total", "done", createdAt, modifiedAt)

				mock.ExpectQuery("SELECT \\* FROM orders").WillReturnRows(rows)
			},
			expectedErr:    true,
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.setupMock(mock)

			repo := NewRepositoryImpl()
			result, err := repo.GetOrders(context.Background(), db)

			if tt.expectedErr {
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, len(tt.expectedResult), len(result))
				assert.Equal(t, tt.expectedResult[0].Id, result[0].Id)
				assert.Equal(t, tt.expectedResult[0].ProductId, result[0].ProductId)
				assert.Equal(t, tt.expectedResult[0].ProductName, result[0].ProductName)
				assert.Equal(t, tt.expectedResult[0].Username, result[0].Username)
				assert.Equal(t, tt.expectedResult[0].Quantity, result[0].Quantity)
				assert.Equal(t, tt.expectedResult[0].Total, result[0].Total)
				assert.Equal(t, tt.expectedResult[0].Status, result[0].Status)
				assert.WithinDuration(t, *tt.expectedResult[0].CreatedAt, *result[0].CreatedAt, time.Second)
				assert.WithinDuration(t, *tt.expectedResult[0].ModifiedAt, *result[0].ModifiedAt, time.Second)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdateOrder(t *testing.T) {
	id := "1"
	status := "done"

	tests := []struct {
		name           string
		setupMock      func(mock sqlmock.Sqlmock)
		expectedErr    bool
		expectedResult *domain.Orders
	}{
		{
			name: "success update order",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE orders SET status = \\? WHERE id = \\?").
					WithArgs(status, id).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedErr: false,
			expectedResult: &domain.Orders{
				Id:     id,
				Status: status,
			},
		},
		{
			name: "update failed",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE orders SET status = \\? WHERE id = \\?").
					WithArgs(status, id).
					WillReturnError(errors.New("update error"))
			},
			expectedErr:    true,
			expectedResult: nil,
		},
		{
			name: "transaction begin failed",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("begin tx failed"))
			},
			expectedErr:    true,
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.setupMock(mock)

			tx, err := db.Begin()
			if tt.expectedErr && tt.name == "transaction begin failed" {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			repo := NewRepositoryImpl()
			order := &domain.Orders{
				Id:     id,
				Status: status,
			}
			err = repo.UpdateOrder(context.Background(), tx, order, id)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
