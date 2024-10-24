package main

import (
	"context"
	"go-backend-scaffold/config"
	"go-backend-scaffold/services"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	config.LoadConfig()
	services.StartUserServiceServer(ctx, "5004")
}
