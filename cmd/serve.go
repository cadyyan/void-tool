package cmd

import (
	"context"
	"fmt"

	"github.com/cadyyan/void-tool/internal"
	"github.com/cadyyan/void-tool/internal/configuration"
	"github.com/cadyyan/void-tool/internal/logging"
	"github.com/cadyyan/void-tool/internal/services"
	"github.com/spf13/cobra"
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

			logger.DebugContext(ctx, "Setting up file based user service")
			userService := services.NewUserFileService(config.RS.DataDirFS())

			logger.DebugContext(ctx, "Starting server")
			server := internal.NewServer(logger, config, userService)

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
