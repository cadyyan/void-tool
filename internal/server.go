package internal

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/cadyyan/void-tool/internal/bgtasks"
	"github.com/cadyyan/void-tool/internal/configuration"
	"github.com/cadyyan/void-tool/internal/database/sqlitedb"
	"github.com/cadyyan/void-tool/internal/services"
	"github.com/cadyyan/void-tool/internal/web"
	"github.com/go-co-op/gocron/v2"
)

type Server struct {
	http *http.Server
	cron gocron.Scheduler
}

func NewServer(
	logger *slog.Logger,
	config configuration.Configuration,
	sqliteDB *sql.DB,
	playerQueries *sqlitedb.Queries,
	voidPlayerService services.VoidPlayerService,
) (*Server, error) {
	cron, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("unable to create background task scheduler: %w", err)
	}

	_, err = cron.NewJob(
		gocron.DurationJob(config.RS.PollFrequency),
		gocron.NewTask(
			bgtasks.ScrapePlayerSkills,
			logger.WithGroup("background--ingestSkills"),
			sqliteDB,
			playerQueries,
			voidPlayerService,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to schedule highscores polling task: %w", err)
	}

	return &Server{
		http: &http.Server{
			Addr:              config.HTTP.BindAddress(),
			Handler:           web.NewRouter(logger, config, voidPlayerService),
			ReadHeaderTimeout: readHeaderTimeout,
			IdleTimeout:       idleTimeout,
			// TODO: error logger
		},
		cron: cron,
	}, nil
}

func (server *Server) Start(ctx context.Context) error {
	server.cron.Start()

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

	// TODO: run these in parallel?
	err := server.cron.Shutdown()
	if err != nil {
		return fmt.Errorf("unable to cleanly shutdown background tasks: %w", err)
	}

	return nil
}

const (
	readHeaderTimeout = 10 * time.Second
	idleTimeout       = 1 * time.Minute
)
