package domain

import "github.com/google/uuid"

type Admin struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username" validate:"required"`
	Password string    `json:"password" validate:"required"`
}
