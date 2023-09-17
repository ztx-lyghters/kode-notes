package core

type Note struct {
	User_ID     int    `json:"-"`
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
