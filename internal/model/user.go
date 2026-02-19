package model

type User struct {
	Id       int
	Name     string
	Username string
	Email    string
	Password string
	Profile  string
	Role     string
	Active   bool
}