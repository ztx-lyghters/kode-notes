package core

type Note struct {
	User_ID     uint   `json:"-"`
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
