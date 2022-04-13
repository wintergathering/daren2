package models

type Dare struct {
	Title string `json:"title" binding:"required"`
	Text  string `json:"text" binding:"required"`
}
