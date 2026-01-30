package web

import (
	"embed"
	"log/slog"
	"net/http"

	"github.com/cadyyan/void-tool/internal/configuration"
	"github.com/cadyyan/void-tool/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
)

//go:embed templates/*.html
var templateFS embed.FS

func NewRouter(
	logger *slog.Logger,
	config configuration.Configuration,
	userService services.UserService,
) *chi.Mux {
	router := chi.NewRouter()

	router.Use(
		// The logger needs to be the first middleware
		httplog.RequestLogger(config.Logging.BuildAccessLogger()),
		middleware.Recoverer,
		middleware.RealIP,
		middleware.RedirectSlashes,
		middleware.CleanPath,
		middleware.Timeout(config.HTTP.RequestTimeout),
		middleware.ContentCharset("UTF-8", "Latin-1", ""),
	)

	// TODO: rate limits

	router.Get("/api/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	router.Get("/", HandlerHome(logger, templateFS, userService))

	return router
}
