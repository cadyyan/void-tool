package bgtasks

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/cadyyan/void-tool/internal/logging"
	"github.com/cadyyan/void-tool/internal/services"
	"github.com/google/uuid"
)

func ScrapePlayerSkills(
	ctx context.Context,
	logger *slog.Logger,
	storageService services.StorageService,
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
			storageService,
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
	storageService services.StorageService,
	player services.VoidPlayer,
) error {
	playerRecord, err := storageService.GetOrCreatePlayerByUsername(
		ctx,
		services.GetOrCreatePlayerByUsernameParams{
			Username:  player.AccountName,
			CreatedOn: player.CreatedOn,
		},
	)
	if err != nil {
		logger.ErrorContext(ctx, "Unable to get player", logging.Err(err))

		return fmt.Errorf("unable to get player: %w", err)
	}

	skillUpdate := make(map[string]services.PlayerSkillRecord)
	for skillName, level := range player.Levels {
		experience := player.Experience[skillName]

		skillUpdate[skillName] = services.PlayerSkillRecord{
			Level:      level,
			Experience: experience,
		}
	}

	today := time.Now().UTC()
	err = storageService.RecordPlayerSkills(
		ctx,
		services.RecordPlayerSkillsParams{
			PlayerID: playerRecord.ID,
			Date:     today,
			Skills:   skillUpdate,
		},
	)
	if err != nil {
		logger.ErrorContext(ctx, "Unable to record player skills", logging.Err(err))

		return fmt.Errorf("unable to record player skills: %w", err)
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
