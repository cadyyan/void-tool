package configuration

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/go-chi/httplog/v2"
	"github.com/kelseyhightower/envconfig"
	"github.com/lmittmann/tint"
)

type LoggingConfiguration struct {
	Level LogLevel `default:"warn"`
	Color bool     `default:"true"`
}

func (config LoggingConfiguration) BuildLogger() *slog.Logger {
	handler := tint.NewHandler(os.Stderr, &tint.Options{
		AddSource:   true,
		Level:       slog.Level(config.Level),
		ReplaceAttr: nil,
		TimeFormat:  time.RFC3339,
		NoColor:     !config.Color,
	})

	return slog.New(handler)
}

func (config LoggingConfiguration) BuildAccessLogger() *httplog.Logger {
	return httplog.NewLogger("http", httplog.Options{
		JSON:                 false,
		LogLevel:             slog.Level(config.Level),
		LevelFieldName:       "level",
		Concise:              true,
		HideRequestHeaders:   []string{},
		RequestHeaders:       false,
		ResponseHeaders:      false,
		MessageFieldName:     "message",
		TimeFieldName:        "timestamp",
		TimeFieldFormat:      time.RFC3339,
		Tags:                 map[string]string{},
		ReplaceAttrsOverride: nil,
		QuietDownRoutes:      []string{"/api/healthcheck"},
		QuietDownPeriod:      1 * time.Minute,
		Writer:               nil,
		Trace:                nil,
	})
}

type LogLevel slog.Level

var (
	LogLevelDebug LogLevel = LogLevel(slog.LevelDebug)
	LogLevelInfo  LogLevel = LogLevel(slog.LevelInfo)
	LogLevelWarn  LogLevel = LogLevel(slog.LevelWarn)
	LogLevelError LogLevel = LogLevel(slog.LevelError)
)

var _ envconfig.Decoder = (*LogLevel)(nil)

func (l *LogLevel) Decode(value string) error {
	var s slog.Level
	if err := s.UnmarshalText([]byte(value)); err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}

	*l = LogLevel(s)

	return nil
}
