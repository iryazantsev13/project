package logger

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/exp/slog"
)

const (
	LevelTrace  = slog.Level(-8)
	LevelNotice = slog.Level(2)
	LevelFatal  = slog.Level(12)
)

// Configure - настраивает логирование через модуль slog,
// создает handler (view: json/text) с требуемым уровнем логирования
// (logLevel: trace/debug/info/notice/warning/error/fatal) или
// возвращает ошибку, если уровень логирования или способ вывода неизвестен
// trace, notice, fatal - самописные уровни логирования
func Configure(logLevel, view string) error {
	level, err := parseLevel(logLevel)
	if err != nil {
		return err
	}

	var LevelNames = map[slog.Leveler]string{
		LevelTrace:  "TRACE",
		LevelNotice: "NOTICE",
		LevelFatal:  "FATAL",
	}

	opts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel, exists := LevelNames[level]
				if !exists {
					levelLabel = level.String()
				}

				a.Value = slog.StringValue(levelLabel)
			}

			return a
		},
	}

	switch strings.ToLower(view) {
	case "text":
		configureTextLogger(opts)
	case "json":
		configureJSONLogger(opts)
	default:
		return errors.New("unknown view output param (need: text/json)")
	}
	return nil
}

// configureTextLogger - Устанавливает текстовое отображение лога
func configureTextLogger(opts *slog.HandlerOptions) {
	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

// configureTextLogger - Устанавливает отображение лога в формате JSON
func configureJSONLogger(opts *slog.HandlerOptions) {
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

// parseLevel - определяет уровень логирования или возвращает
// ошибку, если уровень логирования неизвестен
func parseLevel(level string) (slog.Leveler, error) {
	switch strings.ToLower(level) {
	case "info":
		return slog.LevelInfo, nil
	case "debug":
		return slog.LevelDebug, nil
	case "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	case "trace":
		return LevelTrace, nil
	default:
		return nil, errors.New(fmt.Sprintf("unknown level Log - %v ", level))
	}
}
