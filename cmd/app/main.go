package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"crud_ef/internal/adapter/http"
	"crud_ef/internal/adapter/repository/postgres"
	"crud_ef/internal/config"
	"crud_ef/internal/db"
	"crud_ef/internal/usecase/subscription"

	_ "crud_ef/docs"
)

// @title Subscriptions API
// @version 1.0
// @description REST API для управления подписками и расчёта сумм за период.
// @BasePath /
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx := context.Background()

	pg, err := db.New(ctx, cfg)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	defer pg.Close()

	repo := postgres.NewSubscriptionRepo(pg.Pool)
	svc := subscription.NewService(repo)

	srv := http.New(cfg, svc)

	errCh := make(chan error, 1)
	go func() { errCh <- srv.Run() }()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		log.Fatalf("server: %v", err)
	case <-sigCh:
		log.Println("shutting down")
	}
}
