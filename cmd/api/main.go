package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	connection "example.com/golang-test-task-api-hitalent/internal/infrastructure/postgres"
	repo "example.com/golang-test-task-api-hitalent/internal/repository/postgres"
	service "example.com/golang-test-task-api-hitalent/internal/service"
	router "example.com/golang-test-task-api-hitalent/internal/transport"
	handlers "example.com/golang-test-task-api-hitalent/internal/transport/handlers"
	swagger "example.com/golang-test-task-api-hitalent/internal/transport/swagger"
)

type config struct {
	httpAddr  string
	dbConnect string
}

func main() {
	// Создаем Логгер
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Получаем конфигурацию подключений
	cfg := loadConfig()
	/*
		cfg := config{
			httpAddr:  ":8080",
			dbConnect: "postgres://postgres:postgres1@localhost:5432/postgres?sslmode=disable",
		}
		//*/

	// Подписываемся на сигналы syscall.SIGINT (завершение процесса CTRL+C), syscall.SIGTERM (Завершение, например, по kill pid)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Получаем открытое подключение к БД
	conn, err := connection.NewConnection(ctx, cfg.dbConnect)
	if err != nil {
		logger.Error("Error during get database connection", "error", err)
		os.Exit(1)
	}

	// Закрываем подключение в отложенной функции
	defer func() {
		logger.Info("Closing database connection")
		db, err := conn.DB()
		if err != nil {
			logger.Error("Got error during get database", "error", err)
			os.Exit(1)
		}
		if db != nil {
			if err := db.Close(); err != nil {
				logger.Error("Error during closing database connection", "error", err)
				os.Exit(1)
			}
		}
	}()

	departmentRepository := repo.NewRepository(conn)
	service := service.NewService(departmentRepository)
	departmentHandler := handlers.NewHandler(service)
	swaggerHandler := swagger.NewHandler()
	router := router.NewRouter(departmentHandler, swaggerHandler)

	// Сам сервер
	server := &http.Server{
		Addr:              cfg.httpAddr,
		Handler:           router,
		ReadHeaderTimeout: 15 * time.Second,
	}

	// Запускаем сервер в отдельной горутине
	go func(logger *slog.Logger, cfr config) {
		logger.Info("Server started!", "HTTP_ADDR", cfg.httpAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("ListenAndServe", "error", err)
			os.Exit(1)
		}
	}(logger, cfg)

	// Обрабатываем graceful shutdown
	<-ctx.Done()
	logger.Info("Got the completion signal")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		panic("The server didn`t have time to stop")
	}
	logger.Info("The server is stopped")
}

func loadConfig() config {
	cfg := config{
		httpAddr:  os.Getenv("HTTP_ADDR"),
		dbConnect: os.Getenv("DB_CONNECT"),
	}

	if cfg.httpAddr == "" {
		panic(fmt.Errorf("HTTP_ADDR not received"))
	}
	if cfg.dbConnect == "" {
		panic(fmt.Errorf("DB_CONNECT not received"))
	}

	return cfg
}
