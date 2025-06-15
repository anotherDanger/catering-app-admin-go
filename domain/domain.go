package domain

import (
	"time"

	"github.com/google/uuid"
)

type Domain struct {
	Id          uuid.UUID  `json:"id" validate:"required,uuid"`
	Name        string     `json:"name" validate:"required,min=5,max=50"`
	Description string     `json:"description" validate:"alphanum"`
	Stock       int        `json:"stock" validate:"required,number"`
	Price       int        `json:"price" validate:"required,number"`
	CreatedAt   *time.Time `json:"created_at" validate:"required"`
	ModifiedAt  *time.Time `json:"modified_at" validate:"required"`
}

type Orders struct {
	Id          uuid.UUID  `json:"id"`
	ProductName string     `json:"product_name"`
	Username    string     `json:"username"`
	Quantity    int        `json:"quantity"`
	Total       float64    `json:"total" validate:"required"`
	Status      string     `json:"status"`
	CreatedAt   *time.Time `json:"created_at" validate:"required"`
}
