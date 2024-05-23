package daren2

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrNoDare = errors.New("no dares available")

type Dare struct {
	UUID    string `json:"uuid"`
	Title   string `json:"title" validate:"required"`
	Text    string `json:"text" validate:"required"`
	Seen    bool   `json:"seen"`
	AddedBy string `json:"addedBy"`
}

func NewDare(title string, text string, addedBy string) *Dare {
	id := uuid.New().String()
	return &Dare{
		UUID:    id,
		Title:   title,
		Text:    text,
		Seen:    false,
		AddedBy: addedBy,
	}
}

type DareService interface {
	CreateDare(ctx context.Context, dare *Dare) error
	GetRandomDare(ctx context.Context) (*Dare, error)
	GetAllDares(ctx context.Context) ([]*Dare, error)
}
