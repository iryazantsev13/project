package service

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/exp/slog"
)

// awaitServiceTermination - helper функция по ожиданию работы сервиса
// Мониторит сигналы и ошибки от приложения, и помогает остановить приложение корректно
func AwaitTermination(service interface{ Wait() }, gracefulShutdown func(), errs chan error) {
	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigs:
		var interruptTimeout = 5 * time.Second
		switch sig {
		case syscall.SIGINT:
			interruptTimeout = 30 * time.Second
			slog.Info("SIGINT received. Graceful shutdown.", slog.Duration("timeout", interruptTimeout))
			gracefulShutdown()
		case syscall.SIGTERM:
			slog.Info("SIGTERM received. Trying to stop gracefully.", slog.Duration("timeout", interruptTimeout))
			gracefulShutdown()
		default:
			slog.Info("Unexpected signal received. Quiting.", "signal", sig)
			os.Exit(1)
		}

		select {
		case <-time.After(interruptTimeout):
			slog.Info("Interrupt timeout exceeded")
			os.Exit(1)
		case sig := <-sigs:
			slog.Info("Another signal received. Quiting.", "signal", sig)
			os.Exit(1)
		case err := <-errs:
			if err != nil {
				slog.Info("Service interrupted", "error", err)
				os.Exit(1)
			}
		}

	case err := <-errs:
		switch err {
		case nil:
			slog.Info("Service successfully finished it's work")
		case err:
			const awaitTimeout = 3 * time.Second
			slog.Error("Service failed. Awaiting started tasks.", "error", err, slog.Duration("timeout", awaitTimeout))
			gracefulShutdown()
			time.AfterFunc(awaitTimeout, func() {
				slog.Info("Service tasks timeout exceeded.")
				os.Exit(1)
			})
			service.Wait()
			slog.Info("Service run failed. Service graceful shutdown successfully finished")
			os.Exit(1)
		}
	}
}
