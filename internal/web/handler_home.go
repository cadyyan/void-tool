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
	voidPlayerService services.VoidPlayerService,
) http.HandlerFunc {
	tmpl := template.Must(template.ParseFS(templateFS, "templates/home.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		logger.DebugContext(ctx, "Getting all players")

		players, err := voidPlayerService.GetAllPlayers(ctx)
		if err != nil {
			logger.ErrorContext(ctx, "Unable to get a list of all players", logging.Err(err))
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
		logger.DebugContext(ctx, "Found all users")

		w.WriteHeader(http.StatusOK)

		templateData := map[string]any{
			"Players": players,
			"Skills":  skillOrder,
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
