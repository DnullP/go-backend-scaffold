package client

import (
	"context"
	pb "go-backend-scaffold/proto"
	"go-backend-scaffold/services/discovery"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetUser(req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	serviceAddress := discovery.GetService("UserService")
	println(serviceAddress)

	conn, err := grpc.NewClient(serviceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.GetUser(ctx, req)
	if err != nil {
		panic(err)
	}
	return res, nil
}

func Login(req *pb.LoginRequest) (*pb.LoginResponse, error) {
	serviceAddress := discovery.GetService("UserService")

	conn, err := grpc.NewClient(serviceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.Login(ctx, req)
	if err != nil {
		panic(err)
	}
	return res, nil
}
