package repository

import (
	"catering-admin-go/domain"
	"catering-admin-go/repository/mocks"
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetProducts(t *testing.T) {
	now := time.Now()
	id := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	tests := []struct {
		name           string
		mockSetup      func(repo *mocks.Repository)
		expectedErr    bool
		expectedResult []*domain.Domain
	}{
		{
			name: "Success",
			mockSetup: func(repo *mocks.Repository) {
				repo.On("GetProducts", mock.Anything, mock.Anything).
					Return([]*domain.Domain{
						{
							Id:          id,
							Name:        "Product 1",
							Description: "1st Product",
							Stock:       10,
							Price:       1000,
							CreatedAt:   &now,
							ModifiedAt:  &now,
						},
					}, nil)
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
			name: "Failed",
			mockSetup: func(repo *mocks.Repository) {
				repo.On("GetProducts", mock.Anything, mock.Anything).
					Return(nil, errors.New("failed to get products"))
			},
			expectedErr:    true,
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewRepository(t)
			tt.mockSetup(repo)

			result, err := repo.GetProducts(context.Background(), nil)

			if tt.expectedErr {
				require.Error(t, err)
				require.Nil(t, result)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestAddProduct(t *testing.T) {
	createdAt := time.Now()
	id := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	tests := []struct {
		name           string
		input          *domain.Domain
		mockSetup      func(repo *mocks.Repository)
		expectedErr    bool
		expectedResult *domain.Domain
	}{
		{
			name: "Success",
			input: &domain.Domain{
				Id:          id,
				Name:        "Product 1",
				Description: "1st Product",
				Stock:       100,
				Price:       2000,
				CreatedAt:   &createdAt,
			},
			mockSetup: func(repo *mocks.Repository) {
				repo.On("AddProduct", mock.Anything, mock.Anything, mock.AnythingOfType("*domain.Domain")).
					Return(&domain.Domain{
						Id:          id,
						Name:        "Product 1",
						Description: "1st Product",
						Stock:       100,
						Price:       2000,
						CreatedAt:   &createdAt,
					}, nil)
			},
			expectedErr:    false,
			expectedResult: &domain.Domain{Id: id, Name: "Product 1", Description: "1st Product", Stock: 100, Price: 2000, CreatedAt: &createdAt},
		},
		{
			name: "Failed",
			input: &domain.Domain{
				Id:          id,
				Name:        "",
				Description: "1st Product",
				Stock:       100,
				Price:       2000,
				CreatedAt:   &createdAt,
			},
			mockSetup: func(repo *mocks.Repository) {
				repo.On("AddProduct", mock.Anything, mock.Anything, mock.AnythingOfType("*domain.Domain")).
					Return(nil, errors.New("field name cannot empty"))
			},
			expectedErr:    true,
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewRepository(t)
			tt.mockSetup(repo)

			result, err := repo.AddProduct(context.Background(), nil, tt.input)

			if tt.expectedErr {
				require.Error(t, err)
				require.Nil(t, result)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	modifiedAt := time.Now().Add(2 * time.Hour)
	id := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	tests := []struct {
		name           string
		input          *domain.Domain
		idStr          string
		mockSetup      func(repo *mocks.Repository)
		expectedErr    bool
		expectedResult *domain.Domain
	}{
		{
			name: "Success",
			input: &domain.Domain{
				Name:        "Product 2",
				Description: "2nd Product",
				Stock:       10,
				Price:       1000,
				ModifiedAt:  &modifiedAt,
			},
			idStr: id.String(),
			mockSetup: func(repo *mocks.Repository) {
				repo.On("UpdateProduct", mock.Anything, mock.Anything, mock.AnythingOfType("*domain.Domain"), id.String()).
					Return(&domain.Domain{
						Id:          id,
						Name:        "Product 2",
						Description: "2nd Product",
						Stock:       10,
						Price:       1000,
						ModifiedAt:  &modifiedAt,
					}, nil)
			},
			expectedErr: false,
			expectedResult: &domain.Domain{
				Id:          id,
				Name:        "Product 2",
				Description: "2nd Product",
				Stock:       10,
				Price:       1000,
				ModifiedAt:  &modifiedAt,
			},
		},
		{
			name: "Failed",
			input: &domain.Domain{
				Name:        "Product 2",
				Description: "2nd Product",
				Stock:       10,
				Price:       1000,
				ModifiedAt:  &modifiedAt,
			},
			idStr: "non-existing-id",
			mockSetup: func(repo *mocks.Repository) {
				repo.On("UpdateProduct", mock.Anything, mock.Anything, mock.AnythingOfType("*domain.Domain"), "non-existing-id").
					Return(nil, errors.New("product not found"))
			},
			expectedErr:    true,
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewRepository(t)
			tt.mockSetup(repo)

			result, err := repo.UpdateProduct(context.Background(), nil, tt.input, tt.idStr)

			if tt.expectedErr {
				require.Error(t, err)
				require.Nil(t, result)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	id := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	tests := []struct {
		name        string
		idStr       string
		mockSetup   func(repo *mocks.Repository)
		expectedErr bool
	}{
		{
			name:  "Success",
			idStr: id.String(),
			mockSetup: func(repo *mocks.Repository) {
				repo.On("DeleteProduct", mock.Anything, mock.Anything, id.String()).
					Return(nil)
			},
			expectedErr: false,
		},
		{
			name:  "Failed",
			idStr: id.String(),
			mockSetup: func(repo *mocks.Repository) {
				repo.On("DeleteProduct", mock.Anything, mock.Anything, id.String()).
					Return(errors.New("id not found"))
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewRepository(t)
			tt.mockSetup(repo)

			err := repo.DeleteProduct(context.Background(), nil, tt.idStr)

			if tt.expectedErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestLogin(t *testing.T) {
	id := uuid.MustParse("e7b8a9d4-3f5a-4c82-b7e2-2c3f49b0e9c1")
	username := "Admin"
	password := "hashed_pass"

	tests := []struct {
		name        string
		username    string
		password    string
		mockSetup   func(repo *mocks.Repository)
		expectedErr bool
		expectedRes *domain.Admin
	}{
		{
			name:     "Success",
			username: username,
			password: password,
			mockSetup: func(repo *mocks.Repository) {
				repo.On("Login", mock.Anything, mock.Anything, mock.AnythingOfType("*domain.Admin")).
					Return(&domain.Admin{
						Id:       id,
						Username: username,
						Password: password,
					}, nil)
			},
			expectedErr: false,
			expectedRes: &domain.Admin{
				Id:       id,
				Username: username,
				Password: password,
			},
		},
		{
			name:     "Failed",
			username: username,
			password: "wrong_pass",
			mockSetup: func(repo *mocks.Repository) {
				repo.On("Login", mock.Anything, mock.Anything, mock.AnythingOfType("*domain.Admin")).
					Return(nil, sql.ErrNoRows)
			},
			expectedErr: true,
			expectedRes: nil,
		},
		{
			name: "Connection failed",
			mockSetup: func(repo *mocks.Repository) {
				repo.On("Login", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, sql.ErrConnDone)
			},
			expectedErr: true,
			expectedRes: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()

			repo := mocks.NewRepository(t)
			tt.mockSetup(repo)

			res, err := repo.Login(context.Background(), db, tt.expectedRes)

			if tt.expectedErr {
				require.Error(t, err)
				require.Nil(t, res)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedRes, res)
		})
	}
}
