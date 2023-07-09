package app

import (
	"context"
	"runtime"
	"stub-service/internal/handler"
	"stub-service/internal/settings"
	"sync"

	"github.com/pkg/errors"
)

type Service struct {
	settings *settings.Settings
	wg       sync.WaitGroup
}

// Run - запускает запускает воркеры, после чего блочиться до тех пор пока
// не вернется ошибка или все воркеры не завершаться успехом
func (s *Service) Run(ctx context.Context) error {
	// канал результатов исполнения handler
	runRes := make(chan poolRunResult, 1)
	// запускаем потоки
	for i := 0; i < s.settings.Service.WorkersCount; i++ {
		s.wg.Add(1)
		handler := handler.New(s.settings, s.wg.Done, runtime.NumGoroutine())
		go func() {
			err := handler.Run(ctx)
			runRes <- poolRunResult{handler.ID, err}
		}()
	}

	// мониторим ошибки полученные от работы наших хэндлеров
	for i := 0; i < s.settings.Service.WorkersCount; i++ {
		res := <-runRes
		if res.Err != nil {
			return errors.WithMessagef(res.Err, " handler_id %d failed", res.ID)
		}
	}
	return nil
}

// Wait - блокироует завершение приложения до тех пор пока все треды не будут
// остановлены
func (s *Service) Wait() {
	s.wg.Wait()
}

type poolRunResult struct {
	ID  int
	Err error
}
