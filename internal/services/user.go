package services

import (
	"context"
)

type UserService interface {
	GetAllUsers(ctx context.Context) ([]Player, error)
}

type Player struct {
	AccountName string
	Experience  map[string]float64
	Levels      map[string]int
}
