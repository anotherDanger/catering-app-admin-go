package domain

type Admin struct {
	Id       string `json:"username" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
