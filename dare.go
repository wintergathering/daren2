package daren

import (
	"errors"
	"time"
)

var ErrNoDare = errors.New("no dares available")

type Dare struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	AddedBy   string    `json:"addedBy"`
	Seen      bool      `json:"seen"`
	CreatedAt time.Time `json:"createdAt"`
}

type DareService interface {
	CreateDare(d *Dare) (int, error)
	GetDareByID(id int) (*Dare, error)
	GetRandomDare() (*Dare, error)
	GetAllDares() ([]*Dare, error)
}
