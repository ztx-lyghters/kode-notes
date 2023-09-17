package core

type User struct {
	Id       int    `json:"-"`
	Username string `json:"username"`
	Password string `json:"password"`
}
