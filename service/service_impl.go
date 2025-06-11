package service

import (
	"bytes"
	"catering-admin-go/domain"
	"catering-admin-go/repository"
	"catering-admin-go/web"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	tx, err := svc.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	result, err := svc.repo.Login(ctx, tx, request)
	if err != nil {
		return nil, err
	}

	byteBody, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	bodyReader := bytes.NewBuffer(byteBody)

	byteResult, err := http.Post("http://localhost:8081/v1/auth", "application/json", bodyReader)
	if err != nil {
		return nil, err
	}

	byteBodyResult, err := io.ReadAll(byteResult.Body)
	if err != nil {
		return nil, err
	}
	var token web.Token

	json.Unmarshal(byteBodyResult, &token)

	response := &web.AdminResponse{
		Username:    result.Username,
		AccessToken: token.AccessToken,
	}

	return response, nil
}

func (svc *ServiceImpl) AddProduct(ctx context.Context, request *web.Request) (data *domain.Domain, err error) {
	tx, err := svc.db.Begin()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	id := uuid.New()
	date := time.Now()

	request.Id = id
	request.CreatedAt = &date
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	data, err = svc.repo.AddProduct(ctx, tx, (*domain.Domain)(request))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return data, nil

}

func (svc *ServiceImpl) GetProducts(ctx context.Context) (data []*domain.Domain, err error) {
	tx, err := svc.db.Begin()
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

	products, err := svc.repo.GetProducts(ctx, tx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return products, nil
}

func (svc *ServiceImpl) DeleteProduct(ctx context.Context, id string) error {
	tx, err := svc.db.Begin()
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = svc.repo.DeleteProduct(ctx, tx, id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (svc *ServiceImpl) UpdateProduct(ctx context.Context, request *web.Request, id string) (data *domain.Domain, err error) {
	tx, err := svc.db.Begin()
	if err != nil {
		return nil, err
	}

	date := time.Now()
	request.ModifiedAt = &date
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	data, err = svc.repo.UpdateProduct(ctx, tx, (*domain.Domain)(request), id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return data, nil
}
