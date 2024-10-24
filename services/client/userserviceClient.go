
package client

import (
	"context"
	pb "go-backend-scaffold/proto"
	"go-backend-scaffold/services/discovery"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


func GetUser(req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	service := discovery.GetService("UserService")
	serviceAddress := service.Address
	servicePort := service.Port

	conn, err := grpc.NewClient(serviceAddress+":"+strconv.Itoa(servicePort), grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	service := discovery.GetService("UserService")
	serviceAddress := service.Address
	servicePort := service.Port

	conn, err := grpc.NewClient(serviceAddress+":"+strconv.Itoa(servicePort), grpc.WithTransportCredentials(insecure.NewCredentials()))
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

