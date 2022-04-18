package models

type Dare struct {
	ID    string `json:"id" binding:"required"`
	Title string `json:"title" binding:"required"`
	Text  string `json:"text" binding:"required"`
}
