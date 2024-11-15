package client

import (
	"context"
	"go-backend-scaffold/services/generated/discovery"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func clientRequest(ctx context.Context, serviceName string) *grpc.ClientConn {
	serviceAddress := discovery.GetService(serviceName)

	conn, err := grpc.NewClient(serviceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		panic(err)
	}
	go func() {
		<-ctx.Done()
		conn.Close()
	}()
	return conn
}
