package web

import (
	"time"

	"github.com/google/uuid"
)

type Request struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name" validate:"alpha,required,min=5,max=50"`
	Description string    `json:"description" validate:"omitempty,alphanum"`
	Stock       int       `json:"stock" validate:"required,number"`
	Price       int       `json:"price" validate:"required,number"`
	CreatedAt   time.Time `json:"created_at"`
}
