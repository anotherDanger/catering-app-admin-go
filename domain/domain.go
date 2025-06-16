package domain

import (
	"time"
)

type Domain struct {
	Id          string     `json:"id" validate:"required"`
	Name        string     `json:"name" validate:"required,min=5,max=50"`
	Description string     `json:"description" validate:"alphanum"`
	Stock       int        `json:"stock" validate:"required,number"`
	Price       int        `json:"price" validate:"required,number"`
	CreatedAt   *time.Time `json:"created_at" validate:"required"`
	ModifiedAt  *time.Time `json:"modified_at" validate:"required"`
}

type Orders struct {
	Id          string     `json:"id"`
	ProductId   string     `json:"product_id"`
	ProductName string     `json:"product_name"`
	Username    string     `json:"username"`
	Quantity    int        `json:"quantity"`
	Total       float64    `json:"total" validate:"required"`
	Status      string     `json:"status"`
	CreatedAt   *time.Time `json:"created_at" validate:"required"`
	ModifiedAt  *time.Time `json:"modified_at" validate:"required"`
}
