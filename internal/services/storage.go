package services

import (
	"context"
	"time"
)

type StorageService interface {
	CreatePlayer(ctx context.Context, params CreatePlayerParams) (Player, error)
	GetAllPlayers(ctx context.Context) ([]Player, error)
	GetPlayerByUsername(ctx context.Context, username string) (Player, error)
	GetOrCreatePlayerByUsername(
		ctx context.Context,
		params GetOrCreatePlayerByUsernameParams,
	) (Player, error)
	RecordPlayerSkills(ctx context.Context, params RecordPlayerSkillsParams) error
	GetPlayerSkills(ctx context.Context, username string) (map[string]PlayerSkillRecord, error)
	GetHighscoresForSkill(ctx context.Context, skill string) ([]HighscoreSkillRecord, error)
}

// TODO: parameter validation

type CreatePlayerParams struct {
	Username  string
	CreatedOn time.Time
}

type GetOrCreatePlayerByUsernameParams struct {
	Username  string
	CreatedOn time.Time
}

type RecordPlayerSkillsParams struct {
	PlayerID string
	Skills   map[string]PlayerSkillRecord
	Date     time.Time
}

type PlayerSkillRecord struct {
	Level      int
	Experience float64
}

type HighscoreSkillRecord struct {
	PlayerID   string
	Username   string
	Level      int
	Experience float64
}

type Player struct {
	ID        string
	Username  string
	CreatedOn time.Time
}
