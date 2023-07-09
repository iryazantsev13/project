package app

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"stub-service/internal/settings"
	"stub-service/pkg/config"
	"stub-service/pkg/logger"
	"stub-service/pkg/service"

	"golang.org/x/exp/slog"
)

// Run ...
func Run() {
	settings := getSettings()
	if err := logger.Configure(settings.Log.Level, settings.Log.View); err != nil {
		log.Fatal(fmt.Sprintf("logger.Configure error: %v", err))
	}

	serviceName := NewService(settings)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errs := make(chan error)
	go runService(ctx, serviceName, errs)
	slog.Info("Service successfully runned")

	// ожидаем сигнал или error от service.Run()
	service.AwaitTermination(serviceName, cancel, errs)
	slog.Info("Service successfully finished")
}

// runService - функция запускаемая в горутине. Мониторит ошибки из запускаемого сервиса
func runService(ctx context.Context, service *Service, errs chan error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errs <- service.Run(ctx)
}

// NewService - Создание объекта Service
func NewService(settings *settings.Settings) *Service {
	return &Service{
		settings: settings,
	}
}

// getSettings - получение структуры конфигурации из файла указанного в качестве
// аргументов командной строки. Если файл недоступен/несуществует или
// структура не соответствует конфигурации - завершение работы приложения
func getSettings() *settings.Settings {
	s := &settings.Settings{}
	if err := config.LoadConfig(getConfigPathFromArgs(), s); err != nil {
		log.Fatal(err)
	}
	return s
}

// getConfigPathFromArgs - Функция парсит аргументы заданные в командной строке
// и получает путь к файлу конфигурации иначе завершает программу с кодом 1
func getConfigPathFromArgs() string {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <path_to_config_file>\n", os.Args[0])
		flag.PrintDefaults()
	}
	var configPath string
	flag.Parse()
	if configPath = flag.Arg(0); configPath == "" {
		flag.Usage()
		os.Exit(1)
	}
	return configPath
}
