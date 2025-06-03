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
