package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/cadyyan/void-tool/internal/configuration"
	"github.com/cadyyan/void-tool/internal/services"
	"github.com/cadyyan/void-tool/internal/web"
)

type Server struct {
	http *http.Server
}

func NewServer(
	logger *slog.Logger,
	config configuration.Configuration,
	userService services.UserService,
) *Server {
	return &Server{
		http: &http.Server{
			Addr:              config.HTTP.BindAddress(),
			Handler:           web.NewRouter(logger, config, userService),
			ReadHeaderTimeout: readHeaderTimeout,
			IdleTimeout:       idleTimeout,
			// TODO: error logger
		},
	}
}

func (server *Server) Start(ctx context.Context) error {
	err := server.http.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("unable to start HTTP server: %w", err)
	}

	return nil
}

func (server *Server) Stop(ctx context.Context) error {
	if err := server.http.Shutdown(ctx); err != nil {
		return fmt.Errorf("unable to cleanly shutdown HTTP server: %w", err)
	}

	return nil
}

const (
	readHeaderTimeout = 10 * time.Second
	idleTimeout       = 1 * time.Minute
)
