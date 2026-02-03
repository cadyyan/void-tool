package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cadyyan/void-tool/internal/database/sqlitedb"
	"github.com/google/uuid"
)

type StorageSQLiteService struct {
	db      *sql.DB
	queries *sqlitedb.Queries
}

var _ StorageService = (*StorageSQLiteService)(nil)

func NewStorageSQLiteService(
	db *sql.DB,
	queries *sqlitedb.Queries,
) *StorageSQLiteService {
	return &StorageSQLiteService{
		db:      db,
		queries: queries,
	}
}

func (service *StorageSQLiteService) CreatePlayer(
	ctx context.Context,
	params CreatePlayerParams,
) (Player, error) {
	id := uuid.New().String()

	record, err := service.queries.CreatePlayer(ctx, sqlitedb.CreatePlayerParams{
		ID:        id,
		Username:  params.Username,
		CreatedOn: params.CreatedOn.Format(time.RFC3339),
	})
	if err != nil {
		return Player{}, fmt.Errorf("unable to create player record in SQLite: %w", err)
	}

	player, err := playerSQLiteRecordToPlayer(record)
	if err != nil {
		return Player{}, err
	}

	return player, nil
}

func (service *StorageSQLiteService) GetAllPlayers(
	ctx context.Context,
) ([]Player, error) {
	records, err := service.queries.GetAllPlayers(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get all players from SQLite: %w", err)
	}

	players := make([]Player, len(records))
	for index, record := range records {
		player, err := playerSQLiteRecordToPlayer(record)
		if err != nil {
			return nil, err
		}

		players[index] = player
	}

	return players, nil
}

func (service *StorageSQLiteService) GetPlayerByUsername(
	ctx context.Context,
	username string,
) (Player, error) {
	record, err := service.queries.GetPlayerByName(
		ctx,
		username,
	)
	if err != nil {
		return Player{}, fmt.Errorf("unable to get player by username from SQLite: %w", err)
	}

	player, err := playerSQLiteRecordToPlayer(record)
	if err != nil {
		return Player{}, err
	}

	return player, nil
}

func (service *StorageSQLiteService) GetOrCreatePlayerByUsername(
	ctx context.Context,
	params GetOrCreatePlayerByUsernameParams,
) (Player, error) {
	err := service.queries.CreatePlayerIfNotExist(
		ctx,
		sqlitedb.CreatePlayerIfNotExistParams{
			Username:  params.Username,
			CreatedOn: params.CreatedOn.Format(time.RFC3339),
		},
	)
	if err != nil {
		return Player{}, fmt.Errorf("unable to ensure that player exists in SQLite: %w", err)
	}

	return service.GetPlayerByUsername(ctx, params.Username)
}

func (service *StorageSQLiteService) RecordPlayerSkills(
	ctx context.Context,
	params RecordPlayerSkillsParams,
) error {
	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("unable to start SQLite transaction to record player skills: %w", err)
	}
	defer tx.Rollback()

	date := params.Date.Format(time.DateOnly)

	queriesWithTx := service.queries.WithTx(tx)
	for name, skill := range params.Skills {
		err := queriesWithTx.RecordPlayerSkill(ctx, sqlitedb.RecordPlayerSkillParams{
			PlayerID:   params.PlayerID,
			Day:        date,
			Name:       name,
			Experience: float64(skill.Experience),
			Level:      int64(skill.Level),
		})
		if err != nil {
			return fmt.Errorf("unable to record player skill update to SQLite: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("unable to ")
	}

	return nil
}

func (service *StorageSQLiteService) GetPlayerSkills(
	ctx context.Context,
	username string,
) (map[string]PlayerSkillRecord, error) {
	records, err := service.queries.GetAllPlayerSkillsByPlayerName(
		ctx,
		username,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to get player skills from SQLite: %w", err)
	}

	skills := make(map[string]PlayerSkillRecord)
	for _, record := range records {
		skills[record.Name] = PlayerSkillRecord{
			Experience: record.Experience,
			Level:      int(record.Level),
		}
	}

	return skills, nil
}

func (service *StorageSQLiteService) GetHighscoresForSkill(
	ctx context.Context,
	skill string,
) ([]HighscoreSkillRecord, error) {
	records, err := service.queries.GetHighscoresForSkill(ctx, skill)
	if err != nil {
		return nil, fmt.Errorf("unable to get highscores from SQLite: %w", err)
	}

	highscores := make([]HighscoreSkillRecord, len(records))
	for index, record := range highscores {
		highscores[index] = HighscoreSkillRecord{
			PlayerID:   record.PlayerID,
			Username:   record.Username,
			Experience: record.Experience,
			Level:      record.Level,
		}
	}

	return highscores, nil
}

func playerSQLiteRecordToPlayer(dbRecord sqlitedb.Player) (Player, error) {
	createdOn, err := time.Parse(time.RFC3339, dbRecord.CreatedOn)
	if err != nil {
		return Player{}, fmt.Errorf("unable to parse player created on timestamp from SQLite: %w", err)
	}

	return Player{
		ID:        dbRecord.ID,
		Username:  dbRecord.Username,
		CreatedOn: createdOn,
	}, nil
}
