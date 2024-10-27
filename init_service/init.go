package init_service

import (
	"context"
	"go-backend-scaffold/config"
)

func InitServiceManage(ctx context.Context) {
	config.LoadConfig()
}
