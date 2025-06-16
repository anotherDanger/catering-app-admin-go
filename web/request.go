package web

import (
	"time"
)

type Request struct {
	Id          string     `json:"id"`
	Name        string     `json:"name" validate:"required,min=5,max=50"`
	Description string     `json:"description" validate:"omitempty,alphanum"`
	Stock       int        `json:"stock" validate:"required,number"`
	Price       int        `json:"price" validate:"required,number"`
	CreatedAt   *time.Time `json:"created_at"`
	ModifiedAt  *time.Time `json:"modified_at"`
}
