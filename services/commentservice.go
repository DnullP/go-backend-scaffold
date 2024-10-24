
package services

import (
	"context"
	pb "go-backend-scaffold/proto"
	"go-backend-scaffold/services/generated"
)

type CommentService struct {
	generated.CommentServiceServer
}


func (u CommentService) GetCommentList(context.Context, *pb.GetCommentListRequest) (*pb.GetCommentListResponse, error) {
	//complete this function
	return &pb.GetCommentListResponse{}, nil
}


func StartCommentServiceServer(ctx context.Context, port string) error {
	return generated.StartCommentServiceServer(ctx, port, &CommentService{})
}
