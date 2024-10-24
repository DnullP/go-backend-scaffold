
package services

import (
	"context"
	pb "go-backend-scaffold/proto"
	"go-backend-scaffold/services/generated"
)

type UserService struct {
	generated.UserServiceServer
}


func (u UserService) GetUser(context.Context, *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	//complete this function
	return &pb.GetUserResponse{}, nil
}

func (u UserService) Login(context.Context, *pb.LoginRequest) (*pb.LoginResponse, error) {
	//complete this function
	return &pb.LoginResponse{}, nil
}


func StartUserServiceServer(ctx context.Context, port string) error {
	return generated.StartUserServiceServer(ctx, port, &UserService{})
}
