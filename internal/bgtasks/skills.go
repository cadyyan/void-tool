package bgtasks

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/cadyyan/void-tool/internal/database/sqlitedb"
	"github.com/cadyyan/void-tool/internal/logging"
	"github.com/cadyyan/void-tool/internal/services"
	"github.com/google/uuid"
)

func ScrapePlayerSkills(
	ctx context.Context,
	logger *slog.Logger,
	sqliteDB *sql.DB,
	playerQueries *sqlitedb.Queries,
	playerService services.VoidPlayerService,
) {
	// TODO: job timeout?
	traceID := must(uuid.NewV7()).String()

	logger = logger.With(slog.String("traceId", traceID))

	logger.DebugContext(ctx, "Fetching player stats")
	defer logger.DebugContext(ctx, "Finished fetching player stats")

	players, err := playerService.GetAllPlayers(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "Unable to fetch players", logging.Err(err))

		return
	}

	for _, player := range players {
		err := recordPlayerSkills(
			ctx,
			logger.With("playerName", player.AccountName),
			sqliteDB,
			playerQueries,
			player,
		)
		if err != nil {
			continue
		}
	}
}

func recordPlayerSkills(
	ctx context.Context,
	logger *slog.Logger,
	sqliteDB *sql.DB,
	playerQueries *sqlitedb.Queries,
	player services.Player,
) error {
	tx, err := sqliteDB.BeginTx(ctx, nil)
	if err != nil {
		logger.ErrorContext(ctx, "Unable to start transaction", logging.Err(err))

		return fmt.Errorf("unable to start transaction: %w", err)
	}
	defer tx.Rollback()

	playerQueriesWithTx := playerQueries.WithTx(tx)

	err = playerQueriesWithTx.CreatePlayerIfNotExist(
		ctx,
		sqlitedb.CreatePlayerIfNotExistParams{
			ID:        uuid.New().String(),
			Username:  player.AccountName,
			CreatedOn: player.CreatedOn.Format(time.DateOnly),
		},
	)
	if err != nil {
		logger.ErrorContext(ctx, "Unable to create player record", logging.Err(err))

		return fmt.Errorf("unable to create player record: %w", err)
	}

	playerRecord, err := playerQueriesWithTx.GetPlayerByName(ctx, player.AccountName)
	if err != nil {
		logger.ErrorContext(ctx, "Unable to find player by account name", logging.Err(err))

		return fmt.Errorf("unable to find player by account name: %w", err)
	}

	today := time.Now().UTC().Format(time.DateOnly)

	for skillName, level := range player.Levels {
		experience := player.Experience[skillName]

		err := playerQueriesWithTx.RecordPlayerSkill(
			ctx,
			sqlitedb.RecordPlayerSkillParams{
				PlayerID:   playerRecord.ID,
				Name:       skillName,
				Day:        today,
				Level:      int64(level),
				Experience: experience,
			},
		)
		if err != nil {
			logger.ErrorContext(ctx, "Unable to record player skills", logging.Err(err))

			return fmt.Errorf("unable to record player skills: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		logger.ErrorContext(ctx, "Unable to commit player skills update", logging.Err(err))

		return fmt.Errorf("unable to commit player skills update: %w", err)
	}

	return nil
}

//nolint:ireturn // This is a bug in the linter
func must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}

	return value
}
