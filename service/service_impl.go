package service

import (
	"catering-admin-go/domain"
	"catering-admin-go/repository"
	"catering-admin-go/web"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ServiceImpl struct {
	repo repository.Repository
	db   *sql.DB
}

func NewServiceImpl(repo repository.Repository, db *sql.DB) Service {
	return &ServiceImpl{
		repo: repo,
		db:   db,
	}
}

func (svc *ServiceImpl) AddProduct(ctx context.Context, request *web.Request) (*domain.Domain, error) {
	tx, err := svc.db.Begin()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	id := uuid.New()
	date := time.Now()

	request.Id = id
	request.CreatedAt = date

	data, err := svc.repo.AddProduct(ctx, tx, (*domain.Domain)(request))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	return data, nil

}

func (svc *ServiceImpl) GetProducts(ctx context.Context) ([]*domain.Domain, error) {
	tx, err := svc.db.Begin()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	products, err := svc.repo.GetProducts(ctx, tx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	return products, nil
}
