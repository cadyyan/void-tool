package web

import (
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/cadyyan/void-tool/internal/logging"
	"github.com/cadyyan/void-tool/internal/services"
)

func HandlerHome(
	logger *slog.Logger,
	templateFS fs.FS,
	storageService services.StorageService,
) http.HandlerFunc {
	tmpl := template.Must(template.ParseFS(templateFS, "templates/home.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		logger.DebugContext(ctx, "Getting all players")

		players, err := storageService.GetAllPlayers(ctx)
		if err != nil {
			logger.ErrorContext(ctx, "Unable to get a list of all players", logging.Err(err))
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
		logger.DebugContext(ctx, "Found all users")

		logger.DebugContext(ctx, "Getting skills for each player")
		playerSkills := make(map[string]map[string]services.PlayerSkillRecord)
		for _, player := range players {
			skills, err := storageService.GetPlayerSkills(ctx, player.Username)
			if err != nil {
				logger.ErrorContext(ctx, "Unable to get player stats", logging.Err(err))
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			playerSkills[player.Username] = skills
		}
		logger.DebugContext(ctx, "Get all player skills")

		w.WriteHeader(http.StatusOK)

		templateData := map[string]any{
			"Players":      players,
			"PlayerSkills": playerSkills,
			"SkillOrder":   skillOrder,
		}
		if err := tmpl.Execute(w, templateData); err != nil {
			panic(err) // TODO: better handling
		}
	}
}

var skillOrder = []string{
	"Attack",
	"Defence",
	"Strength",
	"Constitution",
	"Ranged",
	"Prayer",
	"Magic",
	"Cooking",
	"Woodcutting",
	"Fletching",
	"Fishing",
	"Firemaking",
	"Crafting",
	"Smithing",
	"Mining",
	"Herblore",
	"Agility",
	"Thieving",
	"Slayer",
	"Farming",
	"Runecrafting",
	"Hunter",
	"Construction",
	"Summoning",
	"Dungeoneering",
}
