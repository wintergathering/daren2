package daren2

import "context"

type Dare struct {
	Title   string `json:"title"`
	Text    string `json:"text"`
	Seen    bool   `json:"seen"`
	AddedBy string `json:"addedBy"`
}

type DareService interface {
	CreateDare(ctx context.Context, dare Dare) error
	GetRandomDare(ctx context.Context) (Dare, error)
	GetAllDares(ctx context.Context) ([]Dare, error)
}
