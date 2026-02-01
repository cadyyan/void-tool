package services

import (
	"context"
	"time"
)

type VoidPlayerService interface {
	GetAllPlayers(ctx context.Context) ([]Player, error)
}

type Player struct {
	AccountName string
	Experience  map[string]float64
	Levels      map[string]int
	CreatedOn   time.Time
}
