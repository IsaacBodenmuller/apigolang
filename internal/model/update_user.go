package model

type UpdateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Profile  string `json:"profile"`
	Role     string `json:"role"`
}