
package generated

import (
    "context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"go-backend-scaffold/config"
	pb "go-backend-scaffold/proto" // 替换为实际路径
	"go-backend-scaffold/services/discovery"

	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
)

// UserServiceServer 是服务接口，用户需要实现这个接口
type UserServiceServer interface {
    
    GetUser(context.Context, *pb.GetUserRequest) (*pb.GetUserResponse, error)
    
    Login(context.Context, *pb.LoginRequest) (*pb.LoginResponse, error)
    
}

// server 是 gRPC 服务器的实现
type UserServiceImplement struct {
    pb.UnimplementedUserServiceServer
    Handler UserServiceServer
}


// GetUser 实现了 UserService 的 GetUser 方法
func (s *UserServiceImplement) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    return s.Handler.GetUser(ctx, req)
}

// Login 实现了 UserService 的 Login 方法
func (s *UserServiceImplement) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
    return s.Handler.Login(ctx, req)
}


// StartUserServiceServer 启动 gRPC 服务器
func StartUserServiceServer(ctx context.Context, port string, handler UserServiceServer) error {
	//设置监听端口
    lis, err := net.Listen("tcp", ":" + port)
    if err != nil {
        return fmt.Errorf("failed to listen: %v", err)
    }

	//创建grpc服务器
    grpcServer := grpc.NewServer()
    pb.RegisterUserServiceServer(grpcServer, &UserServiceImplement{Handler: handler})

	//创建健康检查服务
	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			log.Println("收到健康检查请求")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
		if err := http.ListenAndServe(":8000", nil); err != nil {
			log.Fatalf("无法启动健康检查 HTTP 服务器: %v", err)
		}
	}()

	//注册服务
	serverId := uuid.New().String()
	iport, _ := strconv.Atoi(port)
	err = discovery.RegisterService("UserService", serverId, iport)
	if err != nil {
		log.Printf("服务注册失败")
		panic(err)
	}
	log.Printf("服务 %s 注册成功", serverId)

	//注销服务
	consulConfig := api.DefaultConfig()
	consulConfig.Address = config.Consul.Address + ":8500"
	client, _ := api.NewClient(consulConfig)
	defer func() {
		err := discovery.DeregisterService(serverId, client)
		if err != nil {
			log.Printf("服务注销失败: %v", err)
		} else {
			log.Printf("服务 %s 已注销", serverId)
		}
	}()

	//关闭服务
	go func() {
		<-ctx.Done()
		log.Println("接收到关闭信号，正在优雅地停止 gRPC 服务器...")
		grpcServer.GracefulStop()
	}()

    log.Printf("gRPC server UserService is running on port :%s", port)
    return grpcServer.Serve(lis)
}
