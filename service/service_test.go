package service

import (
	"catering-admin-go/domain"
	"catering-admin-go/repository/mocks"
	"catering-admin-go/web"
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetProducts(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(sqlmock.Sqlmock, *mocks.Repository)
		expectedErr    bool
		expectedResult []*domain.Domain
		checkResult    func(t *testing.T, result []*domain.Domain, err error)
	}{
		{
			name: "Success",
			setupMock: func(dbmock sqlmock.Sqlmock, repo *mocks.Repository) {
				response := []*domain.Domain{
					{
						Id:          uuid.New(),
						Name:        "Product 1",
						Description: "1st Product",
						Price:       100,
						Stock:       10,
					},
					{
						Id:          uuid.New(),
						Name:        "Product 2",
						Description: "2nd Product",
						Price:       100,
						Stock:       10,
					},
				}
				dbmock.ExpectBegin()
				repo.On("GetProducts", mock.Anything, mock.Anything).Return(response, nil)
				dbmock.ExpectCommit()
			},
			expectedErr: false,
			checkResult: func(t *testing.T, result []*domain.Domain, err error) {
				assert.NoError(t, err)
				assert.Len(t, result, 2)
				assert.Equal(t, "Product 1", result[0].Name)
				assert.Equal(t, "Product 2", result[1].Name)
			},
		},
		{
			name: "Failed",
			setupMock: func(dbmock sqlmock.Sqlmock, repo *mocks.Repository) {
				dbmock.ExpectBegin()
				repo.On("GetProducts", mock.Anything, mock.Anything).Return(nil, errors.New("cannot get products"))
			},
			expectedErr: true,
			checkResult: func(t *testing.T, result []*domain.Domain, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
			},
		},
		{
			name: "NotFound",
			setupMock: func(dbmock sqlmock.Sqlmock, repo *mocks.Repository) {
				response := []*domain.Domain{}
				dbmock.ExpectBegin()
				repo.On("GetProducts", mock.Anything, mock.Anything).Return(response, nil)
				dbmock.ExpectCommit()
			},
			expectedErr: false,
			checkResult: func(t *testing.T, result []*domain.Domain, err error) {
				assert.NoError(t, err)
				assert.Empty(t, result)
			},
		},
		{
			name: "ErrTransaction",
			setupMock: func(dbmock sqlmock.Sqlmock, repo *mocks.Repository) {
				dbmock.ExpectBegin().WillReturnError(errors.New("cannot start transaction"))
			},
			expectedErr: true,
			checkResult: func(t *testing.T, result []*domain.Domain, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
			},
		},
		{
			name: "ErrTransactionCommit",
			setupMock: func(dbmock sqlmock.Sqlmock, repo *mocks.Repository) {
				dbmock.ExpectBegin()
				repo.On("GetProducts", mock.Anything, mock.Anything).Return(nil, errors.New("repository error"))
			},
			expectedErr: true,
			checkResult: func(t *testing.T, result []*domain.Domain, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
			},
		},
		{
			name: "ErrTransactionRollback",
			setupMock: func(dbmock sqlmock.Sqlmock, repo *mocks.Repository) {
				dbmock.ExpectBegin()
				repo.On("GetProducts", mock.Anything, mock.Anything).Return(nil, errors.New("repository error"))
			},
			expectedErr: true,
			checkResult: func(t *testing.T, result []*domain.Domain, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, dbmock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()

			var repo *mocks.Repository
			if tt.name != "ErrTransaction" {
				repo = mocks.NewRepository(t)
			}

			tt.setupMock(dbmock, repo)

			svc := NewServiceImpl(repo, db)
			products, err := svc.GetProducts(context.Background())

			tt.checkResult(t, products, err)

			if repo != nil {
				repo.AssertExpectations(t)
			}
			assert.NoError(t, dbmock.ExpectationsWereMet())
		})
	}
}

func TestAddProduct(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(dbmock sqlmock.Sqlmock, repo *mocks.Repository)
		expectedErr    bool
		expectedResult *domain.Domain
	}{
		{
			name: "Success",
			mockSetup: func(dbmock sqlmock.Sqlmock, repo *mocks.Repository) {
				response := &domain.Domain{
					Id:          uuid.New(),
					Name:        "Product 1",
					Description: "1st Product",
					Price:       100,
					Stock:       10,
				}

				dbmock.ExpectBegin()
				repo.On("AddProduct", mock.Anything, mock.Anything, mock.Anything).Return(response, nil)
				dbmock.ExpectCommit()
			},
			expectedErr: false,
			expectedResult: &domain.Domain{
				Name:        "Product 1",
				Description: "1st Product",
				Price:       100,
				Stock:       10,
			},
		},
		{
			name: "Failed",
			mockSetup: func(dbmock sqlmock.Sqlmock, repo *mocks.Repository) {
				dbmock.ExpectBegin()
				repo.On("AddProduct", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("field cannot be empty"))
				dbmock.ExpectRollback()
			},
			expectedErr:    true,
			expectedResult: nil,
		},
		{
			name: "Transaction Failed",
			mockSetup: func(dbmock sqlmock.Sqlmock, repo *mocks.Repository) {
				dbmock.ExpectBegin().WillReturnError(errors.New("transaction error"))
			},
			expectedErr:    true,
			expectedResult: nil,
		},
		{
			name: "Connection Failed",
			mockSetup: func(dbmock sqlmock.Sqlmock, repo *mocks.Repository) {
				dbmock.ExpectBegin().WillReturnError(errors.New("connection error"))
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
			defer db.Close()

			repo := mocks.NewRepository(t)
			tt.mockSetup(mock, repo)

			svc := NewServiceImpl(repo, db)
			result, err := svc.AddProduct(context.Background(), &web.Request{
				Id:          uuid.New(),
				Name:        "Product 1",
				Description: "1st Product",
				Price:       100,
				Stock:       10,
			})

			if tt.expectedErr && err == nil {
				t.Error("expected error, but got none")
			} else if !tt.expectedErr && err != nil {
				t.Errorf("expected no error, but got: %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}

			if tt.expectedResult != nil {
				assert.Equal(t, tt.expectedResult.Name, result.Name, "Name mismatch")
				assert.Equal(t, tt.expectedResult.Description, result.Description, "Description mismatch")
				assert.Equal(t, tt.expectedResult.Price, result.Price, "Price mismatch")
				assert.Equal(t, tt.expectedResult.Stock, result.Stock, "Stock mismatch")
			}
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(sqlmock sqlmock.Sqlmock, repo *mocks.Repository)
		expectedErr    bool
		expectedResult *domain.Domain
	}{
		{
			name: "Success",
			mockSetup: func(sqlmock sqlmock.Sqlmock, repo *mocks.Repository) {
				response := &domain.Domain{
					Name:        "Product 1",
					Description: "1st Product",
					Price:       2000,
					Stock:       100,
				}
				sqlmock.ExpectBegin()
				repo.On("UpdateProduct", mock.Anything, mock.Anything, mock.MatchedBy(func(p *domain.Domain) bool {
					return p.Name == "Product 1" &&
						p.Description == "1st Product" &&
						p.Price == 2000 &&
						p.Stock == 100
				}), mock.Anything).Return(response, nil)
				sqlmock.ExpectCommit()

			},
			expectedErr: false,
			expectedResult: &domain.Domain{
				Name:        "Product 1",
				Description: "1st Product",
				Price:       2000,
				Stock:       100,
			},
		},
		{
			name: "Failed",
			mockSetup: func(sqlmock sqlmock.Sqlmock, repo *mocks.Repository) {
				sqlmock.ExpectBegin()
				repo.On("UpdateProduct", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("failed update product"))
				sqlmock.ExpectRollback()
			},
			expectedErr:    true,
			expectedResult: nil,
		},
		{
			name: "Transaction failed",
			mockSetup: func(sqlmock sqlmock.Sqlmock, repo *mocks.Repository) {
				sqlmock.ExpectBegin().WillReturnError(errors.New("transaction failed"))

			},
			expectedErr:    true,
			expectedResult: nil,
		},
		{
			name: "Connection failed",
			mockSetup: func(sqlmock sqlmock.Sqlmock, repo *mocks.Repository) {
				sqlmock.ExpectBegin().WillReturnError(errors.New("connection failed"))

			},
			expectedErr:    true,
			expectedResult: nil,
		},
		{
			name: "Id not found",
			mockSetup: func(sqlmock sqlmock.Sqlmock, repo *mocks.Repository) {
				sqlmock.ExpectBegin()

				repo.On("UpdateProduct", mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("id not found"))
				sqlmock.ExpectRollback()
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
			defer db.Close()

			repo := mocks.NewRepository(t)
			tt.mockSetup(mock, repo)

			id := "8154bf2e-2723-4149-b366-01998b2b3f00"

			svc := NewServiceImpl(repo, db)
			result, err := svc.UpdateProduct(context.Background(), &web.Request{

				Name:        "Product 1",
				Description: "1st Product",
				Price:       2000,
				Stock:       100,
			}, id)

			if tt.expectedErr && err != nil {
				assert.Error(t, err, "Success error test")
			} else if !tt.expectedErr && err != nil {
				t.Errorf("expected no error, but got: %v", err)
			}

			if tt.expectedResult != nil {
				assert.Equal(t, tt.expectedResult.Name, result.Name)
				assert.Equal(t, tt.expectedResult.Description, result.Description)
				assert.Equal(t, tt.expectedResult.Price, result.Price)
				assert.Equal(t, tt.expectedResult.Stock, result.Stock)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	tests := []struct {
		name        string
		mockSetup   func(sqlmock sqlmock.Sqlmock, repo *mocks.Repository)
		expectedErr bool
	}{
		{
			name: "Success",
			mockSetup: func(sqlmock sqlmock.Sqlmock, repo *mocks.Repository) {
				sqlmock.ExpectBegin()
				repo.On("DeleteProduct", mock.Anything, mock.Anything, mock.AnythingOfType("string")).Return(nil)
				sqlmock.ExpectCommit()
			},
			expectedErr: false,
		},
		{
			name: "Failed",
			mockSetup: func(sqlmock sqlmock.Sqlmock, repo *mocks.Repository) {
				sqlmock.ExpectBegin()
				repo.On("DeleteProduct", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("delete product failed"))
				sqlmock.ExpectRollback()
			},
			expectedErr: true,
		},
		{
			name: "Id not found",
			mockSetup: func(sqlmock sqlmock.Sqlmock, repo *mocks.Repository) {
				sqlmock.ExpectBegin()
				repo.On("DeleteProduct", mock.Anything, mock.Anything, mock.AnythingOfType("string")).Return(errors.New("id not found"))
				sqlmock.ExpectRollback()
			},
			expectedErr: true,
		},
		{
			name: "Transaction failed",
			mockSetup: func(sqlmock sqlmock.Sqlmock, repo *mocks.Repository) {
				sqlmock.ExpectBegin().WillReturnError(errors.New("transaction failed"))
			},
			expectedErr: true,
		},
		{
			name: "Connection failed",
			mockSetup: func(sqlmock sqlmock.Sqlmock, repo *mocks.Repository) {
				sqlmock.ExpectBegin().WillReturnError(errors.New("Connection failed"))
			},
			expectedErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}

			repo := mocks.NewRepository(t)
			tt.mockSetup(mock, repo)

			id := uuid.New()
			svc := NewServiceImpl(repo, db)
			err = svc.DeleteProduct(context.Background(), id.String())

			if tt.expectedErr {
				assert.Error(t, err, "error happen")
				if tt.name == "Id_not_found" {
					assert.Equal(t, errors.New("delete product failed"), err)
				} else if tt.name == "Transaction_failed" {
					assert.Equal(t, errors.New("transaction failed"), err)
				} else if tt.name == "Connection_failed" {
					assert.Equal(t, errors.New("connection failed"), err)
				}

			}

			if !tt.expectedErr {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}
