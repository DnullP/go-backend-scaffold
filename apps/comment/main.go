package main

import (
	"context"
	"go-backend-scaffold/init_service"
	"go-backend-scaffold/services"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	init_service.InitServiceManage(ctx)
	services.StartCommentServiceServer(ctx, "5008")
}
