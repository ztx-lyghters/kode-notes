package core

type User struct {
	Id       uint   `json:"-"`
	Username string `json:"username"`
	Password string `json:"password"`
}
