
package client

import (
	"context"
	pb "go-backend-scaffold/proto"
	"go-backend-scaffold/services/discovery"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


func GetCommentList(req *pb.GetCommentListRequest) (*pb.GetCommentListResponse, error) {
	serviceAddress := discovery.GetService("CommentService")

	conn, err := grpc.NewClient(serviceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := pb.NewCommentServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.GetCommentList(ctx, req)
	if err != nil {
		panic(err)
	}
	return res, nil
}

