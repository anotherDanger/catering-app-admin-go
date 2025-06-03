package domain

import (
	"time"
)

type Domain struct {
	Id          string    `json:"id" validate:"required,uuid"`
	Name        string    `json:"name" validate:"alpha,required,string,min=5,max=50"`
	Description string    `json:"description" validate:"alphanum"`
	Stock       int       `json:"stock" validate:"required,number"`
	Price       int       `json:"price" validate:"required,number"`
	CreatedAt   time.Time `json:"created_at" validate:"required"`
}
