package service

import (
	"catering-admin-go/domain"
	"catering-admin-go/helper"
	"catering-admin-go/logger"
	"catering-admin-go/repository"
	"catering-admin-go/web"
	"context"
	"database/sql"

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

func (svc *ServiceImpl) Login(ctx context.Context, request *domain.Admin) (*web.AdminResponse, error) {

	result, err := svc.repo.Login(ctx, svc.db, request)
	if err != nil {
		return nil, err
	}

	response := &web.AdminResponse{
		Username: result.Username,
	}

	return response, nil
}

func (svc *ServiceImpl) AddProduct(ctx context.Context, request *web.Request) (data *domain.Domain, err error) {
	tx, err := svc.db.Begin()
	if err != nil {
		logger.GetLogger("service-log").Log("add product", "error", err.Error())
		return nil, err
	}

	id := uuid.New()
	date := time.Now()

	request.Id = id
	request.CreatedAt = &date
	defer helper.WithTransaction(tx, &err)

	data, err = svc.repo.AddProduct(ctx, tx, (*domain.Domain)(request))
	if err != nil {
		logger.GetLogger("service-log").Log("add product", "error", err.Error())
		return nil, err
	}

	return data, nil

}

func (svc *ServiceImpl) GetProducts(ctx context.Context) (data []*domain.Domain, err error) {

	products, err := svc.repo.GetProducts(ctx, svc.db)
	if err != nil {
		logger.GetLogger("service-log").Log("get product", "error", err.Error())
		return nil, err
	}

	return products, nil
}

func (svc *ServiceImpl) DeleteProduct(ctx context.Context, id string) error {
	tx, err := svc.db.Begin()
	if err != nil {
		logger.GetLogger("service-log").Log("delete product", "error", err.Error())
		return err
	}

	defer helper.WithTransaction(tx, &err)

	err = svc.repo.DeleteProduct(ctx, tx, id)
	if err != nil {
		logger.GetLogger("service-log").Log("delete product", "error", err.Error())
		return err
	}

	return nil
}

func (svc *ServiceImpl) UpdateProduct(ctx context.Context, request *web.Request, id string) (data *domain.Domain, err error) {
	tx, err := svc.db.Begin()
	if err != nil {
		logger.GetLogger("service-log").Log("update product", "error", err.Error())
		return nil, err
	}

	date := time.Now()
	request.ModifiedAt = &date
	defer helper.WithTransaction(tx, &err)
	data, err = svc.repo.UpdateProduct(ctx, tx, (*domain.Domain)(request), id)
	if err != nil {
		logger.GetLogger("service-log").Log("update product", "error", err.Error())
		return nil, err
	}

	return data, nil
}

func (svc *ServiceImpl) GetOrders(ctx context.Context) (orders []*domain.Orders, err error) {

	orders, err = svc.repo.GetOrders(ctx, svc.db)
	if err != nil {
		logger.GetLogger("service-log").Log("get orders", "error", err.Error())
		return nil, err
	}

	return orders, nil
}

func (svc *ServiceImpl) UpdateOrder(ctx context.Context, entity *domain.Orders, id string) (err error) {
	tx, err := svc.db.Begin()
	if err != nil {
		logger.GetLogger("service-log").Log("update order", "error", err.Error())
		return err
	}

	defer helper.WithTransaction(tx, &err)

	err = svc.repo.UpdateOrder(ctx, tx, entity, id)
	if err != nil {
		logger.GetLogger("service-log").Log("update order", "error", err.Error())
		return err
	}

	return nil
}
