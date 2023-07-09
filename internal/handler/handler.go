package handler

import (
	"context"
	"stub-service/internal/settings"
	"time"

	"golang.org/x/exp/slog"
)

// Yекоторая структура для работы нашего сервиса (пример)
type ServiceSettings struct {
	DebugMode   bool  `yaml:"debug_mode"`
	WeekDays    []int `yaml:"week_days"`
	PackageSize int   `yaml:"package_size"`
}

type Handler struct {
	ID       int
	settings *ServiceSettings
	log      *slog.Logger
	stop     func()
	// должен быть объект storage в котором или БД или что то еще
}

// isRunning мониторит получаемые приложением сигналы, для безопасного останова
func isRunning(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return false
	default:
	}
	return true
}

// Run ...
func (h *Handler) Run(ctx context.Context) error {
	defer func() {
		h.stop()
		h.log.Debug("Stopped")
	}()
	for isRunning(ctx) {
		h.log.Debug("Something payload")
		time.Sleep(time.Second * 3)
	}
	return nil
}

// New - создание объекта Handler
func New(conf *settings.Settings, onWaitDone func(), goroutineId int) *Handler {
	return &Handler{
		ID:       goroutineId,
		settings: (*ServiceSettings)(&conf.ServiceSettings),
		log:      slog.With(slog.Group("Handler", slog.Int("id", goroutineId))),
		stop:     onWaitDone,
	}
}
