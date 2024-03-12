package daren2

import (
	"context"
	"errors"
)

var ErrNoDare = errors.New("no dares available")

type Dare struct {
	UUID    string `json:"uuid"`
	Title   string `json:"title"`
	Text    string `json:"text"`
	Seen    bool   `json:"seen"`
	AddedBy string `json:"addedBy"`
}

type DareService interface {
	CreateDare(ctx context.Context, dare *Dare) error
	GetRandomDare(ctx context.Context) (*Dare, error)
	GetAllDares(ctx context.Context) ([]*Dare, error)
}
