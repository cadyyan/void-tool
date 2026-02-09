package web

import (
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/cadyyan/void-tool/internal/logging"
	"github.com/cadyyan/void-tool/internal/services"
	"github.com/go-chi/chi/v5"
)

func HandlerPlayerPage(
	logger *slog.Logger,
	templateFS fs.FS,
	storageService services.StorageService,
) http.HandlerFunc {
	tmpl := template.Must(
		template.New("player_page.html").
			Funcs(DefaultMacros).
			ParseFS(templateFS, "templates/player_page.html"),
	)

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		username := chi.URLParam(r, "username")
		if username == "" {
			// TODO: proper 404 page
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Username not provided"))

			return
		}

		player, err := storageService.GetPlayerByUsername(ctx, username)
		if err != nil {
			logger.ErrorContext(ctx, "Unable to get user", logging.Err(err))

			// TODO: proper error handling
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Unable to get user by given username"))

			return
		}

		skills, err := storageService.GetPlayerSkills(ctx, player.Username)
		if err != nil {
			logger.ErrorContext(ctx, "Unable to get user skills", logging.Err(err))

			// TODO: proper error handling
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unable to get user skills"))

			return
		}

		// TODO: combat level

		totalExperience := 0.0
		totalLevel := 0
		for _, skill := range skills {
			totalExperience += skill.Experience
			totalLevel += skill.Level
		}

		w.WriteHeader(http.StatusOK)

		templateData := map[string]any{
			"Player":          player,
			"Skills":          skills,
			"SkillOrder":      skillOrder,
			"TotalExperience": totalExperience,
			"TotalLevel":      totalLevel,
		}
		if err := tmpl.Execute(w, templateData); err != nil {
			panic(err) // TODO: better handling
		}
	}
}
