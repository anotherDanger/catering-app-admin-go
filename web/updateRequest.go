package web

import (
	"time"
)

type UpdateProductRequest struct {
	Name        *string    `json:"name" validate:"omitempty,alpha,min=5,max=50"`
	Description *string    `json:"description" validate:"omitempty,alphanum"`
	Stock       *int       `json:"stock" validate:"omitempty,number"`
	Price       *int       `json:"price" validate:"omitempty,number"`
	CreatedAt   *time.Time `json:"created_at"`
	ModifiedAt  *time.Time `json:"modified_at"`
}
