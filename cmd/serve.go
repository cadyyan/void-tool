package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/cadyyan/void-tool/database/sqlite"
	"github.com/cadyyan/void-tool/internal"
	"github.com/cadyyan/void-tool/internal/configuration"
	"github.com/cadyyan/void-tool/internal/database/sqlitedb"
	"github.com/cadyyan/void-tool/internal/logging"
	"github.com/cadyyan/void-tool/internal/services"
	"github.com/golang-migrate/migrate/v4"
	migrateSQLiteDriver "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/spf13/cobra"
	_ "modernc.org/sqlite"
)

func newServeCommand() *cobra.Command {
	command := cobra.Command{
		Use:   "serve",
		Short: "Run the HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			config, err := configuration.NewConfigurationFromEnv()
			if err != nil {
				return fmt.Errorf("unable to prepare server: %w", err)
			}

			logger := config.Logging.BuildLogger()
			logger.DebugContext(ctx, "Configuration loaded")

			logger.DebugContext(ctx, "Opening connection to SQLite database")

			sqliteConnection, err := config.SQLite.Connect()
			if err != nil {
				logger.ErrorContext(ctx, "Unable to connect to SQLite database", logging.Err(err))

				return fmt.Errorf("unable to connect to SQLite database: %w", err)
			}
			defer sqliteConnection.Close()

			err = runSQLiteMigrations(
				ctx,
				logger,
				sqliteConnection,
			)
			if err != nil {
				return err
			}

			playerQueries := sqlitedb.New(sqliteConnection)

			logger.DebugContext(ctx, "Setting up file based user service")

			voidPlayerService := services.NewVoidPlayerFileService(config.RS.DataDirFS())

			logger.DebugContext(ctx, "Setting up server")

			server, err := internal.NewServer(
				logger,
				config,
				sqliteConnection,
				playerQueries,
				voidPlayerService,
			)
			if err != nil {
				return fmt.Errorf("unable to setup server: %w", err)
			}

			logger.DebugContext(ctx, "Starting server")
			serverErrChan := make(chan error, 1)
			go func(errChan chan<- error) {
				if err := server.Start(ctx); err != nil {
					logger.ErrorContext(ctx, "Error starting server", logging.Err(err))

					errChan <- fmt.Errorf("unable to start server: %w", err)
				}
			}(serverErrChan)

			select {
			case err := <-serverErrChan:
				logger.ErrorContext(ctx, "Server crashed", logging.Err(err))

			case <-ctx.Done():
				cancelCtx, done := context.WithTimeout(ctx, config.HTTP.ShutdownTimeout)
				defer done()

				if err := server.Stop(cancelCtx); err != nil {
					logger.ErrorContext(
						cancelCtx,
						"Unable to cleanly shutdown server",
						logging.Err(err),
					)

					return fmt.Errorf("unable to cleanly shutdown server: %w", err)
				}
			}

			return nil
		},
	}

	return &command
}

func runSQLiteMigrations(
	ctx context.Context,
	logger *slog.Logger,
	dbConn *sql.DB,
) error {
	logger.DebugContext(ctx, "Applying database migrations")

	sourceDriver, err := iofs.New(sqlite.MigrationsFS, "migrations")
	if err != nil {
		logger.ErrorContext(ctx, "Unable to get migrations", logging.Err(err))

		return fmt.Errorf("unable to get migrations: %w", err)
	}

	dbDriver, err := migrateSQLiteDriver.WithInstance(dbConn, &migrateSQLiteDriver.Config{})
	if err != nil {
		logger.ErrorContext(ctx, "Unable to create migration database driver", logging.Err(err))

		return fmt.Errorf("unable to create migration database driver: %w", err)
	}

	migrator, err := migrate.NewWithInstance(
		"iofs",
		sourceDriver,
		"sqlite",
		dbDriver,
	)
	if err != nil {
		logger.ErrorContext(
			ctx,
			"Unable to setup SQLite migrator",
			logging.Err(err),
		)

		return fmt.Errorf("unable to setup SQLite migrator: %w", err)
	}

	err = migrator.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.ErrorContext(
			ctx,
			"Unable to apply migrations to SQLite",
			logging.Err(err),
		)

		return fmt.Errorf("unable to apply migrations to SQLite: %w", err)
	}

	return nil
}
