package dto

type CreateItem struct {
	Name  string `json:"name" validate:"required"`
	Count int    `json:"count" validate:"required"`
}

type CreateUser struct {
	Name string `json:"name" validate:"required"`
	Role string `json:"role" validate:"required"`
}

type UpdateItem struct {
	ID     int
	UserID int     `json:"user_id" validate:"required"`
	Name   *string `json:"name"`
	Count  *int    `json:"count"`
}
